package rabbitmq

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/IlhamSetiaji/julong-manpower-be/internal/http/messaging"
	"github.com/IlhamSetiaji/julong-manpower-be/internal/http/request"
	"github.com/IlhamSetiaji/julong-manpower-be/internal/http/response"
	"github.com/IlhamSetiaji/julong-manpower-be/internal/usecase"
	"github.com/IlhamSetiaji/julong-manpower-be/utils"
	"github.com/google/uuid"
	"github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func InitConsumer(viper *viper.Viper, log *logrus.Logger) {
	// conn
	conn, err := amqp091.Dial(viper.GetString("rabbitmq.url"))
	if err != nil {
		log.Printf("ERROR: fail init consumer: %s", err.Error())
		os.Exit(1)
	}

	log.Printf("INFO: done init consumer conn")

	// create channel
	amqpChannel, err := conn.Channel()
	if err != nil {
		log.Printf("ERROR: fail create channel: %s", err.Error())
		os.Exit(1)
	}

	// create queue
	queue, err := amqpChannel.QueueDeclare(
		"julong_manpower", // channelname
		true,              // durable
		false,             // delete when unused
		false,             // exclusive
		false,             // no-wait
		nil,               // arguments
	)
	if err != nil {
		log.Printf("ERROR: fail create queue: %s", err.Error())
		os.Exit(1)
	}

	// channel
	msgChannel, err := amqpChannel.Consume(
		queue.Name, // queue
		"",         // consumer
		false,      // auto-ack
		false,      // exclusive
		false,      // no-local
		false,      // no-wait
		nil,        // args
	)
	if err != nil {
		log.Printf("ERROR: fail create channel: %s", err.Error())
		os.Exit(1)
	}

	// consume
	for {
		select {
		case msg := <-msgChannel:
			// unmarshal
			docRply := &response.RabbitMQResponse{}
			docMsg := &request.RabbitMQRequest{}
			err = json.Unmarshal(msg.Body, docRply)
			if err != nil {
				log.Printf("ERROR: fail unmarshl: %s", msg.Body)
				continue
			}
			log.Printf("INFO: received docRply: %v", docRply)

			err = json.Unmarshal(msg.Body, docMsg)
			if err != nil {
				log.Printf("ERROR: fail unmarshl: %s", msg.Body)
				continue
			}
			log.Printf("INFO: received docMsg: %v", docMsg)

			// ack for message
			err = msg.Ack(true)
			if err != nil {
				log.Printf("ERROR: fail to ack: %s", err.Error())
			}

			// find waiting channel(with uid) and forward the reply to it
			if rchan, ok := utils.Rchans[docRply.ID]; ok {
				rchan <- *docRply
			}

			handleMsg(docMsg, log, viper)
		}
	}
}

func handleMsg(docMsg *request.RabbitMQRequest, log *logrus.Logger, viper *viper.Viper) {
	// switch case
	var msgData map[string]interface{}

	switch docMsg.MessageType {
	case "reply":
		log.Printf("INFO: received reply message")
		return
	case "find_job_plafon_by_job_id":
		jobID, ok := docMsg.MessageData["job_id"].(string)
		if !ok {
			log.Printf("Invalid request format: missing 'job_id'")
			msgData = map[string]interface{}{
				"error": errors.New("missing 'job_id'").Error(),
			}
			break
		}

		messageFactory := messaging.JobPlafonMessageFactory(log)
		message, err := messageFactory.FindJobPlafonByJobIDMessage(uuid.MustParse(jobID))

		if err != nil {
			log.Printf("Failed to execute message: %v", err)
			msgData = map[string]interface{}{
				"error": err.Error(),
			}
			break
		}

		msgData = map[string]interface{}{
			"id":     message.ID,
			"plafon": message.Plafon,
		}
	case "find_mp_request_header_by_id":
		mpRequestHeaderID, ok := docMsg.MessageData["mp_request_header_id"].(string)
		if !ok {
			log.Printf("Invalid request format: missing 'mp_request_header_id'")
			msgData = map[string]interface{}{
				"error": errors.New("missing 'mp_request_header_id'").Error(),
			}
			break
		}

		uc := usecase.MPRequestUseCaseFactory(viper, log)
		mpRequestHeader, err := uc.FindByIDOnly(uuid.MustParse(mpRequestHeaderID))
		if err != nil {
			log.Printf("Failed to execute usecase: %v", err)
			msgData = map[string]interface{}{
				"error": err.Error(),
			}
			break
		}

		msgData = map[string]interface{}{
			"mp_request_header": mpRequestHeader,
		}
	default:
		log.Printf("Unknown message type, please recheck your type: %s", docMsg.MessageType)

		msgData = map[string]interface{}{
			"error": errors.New("unknown message type").Error(),
		}
	}
	// reply
	reply := response.RabbitMQResponse{
		ID: docMsg.ID,
		// MessageType: docMsg.MessageType,
		MessageType: "reply",
		MessageData: msgData,
	}
	msg := RabbitMsg{
		QueueName: docMsg.ReplyTo,
		Reply:     reply,
	}
	rchan <- msg
}
