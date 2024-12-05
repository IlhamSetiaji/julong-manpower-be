package messaging

import (
	"errors"
	"log"
	"time"

	// "github.com/IlhamSetiaji/go-rabbitmq-utils/rabbitmq"
	"github.com/IlhamSetiaji/julong-manpower-be/internal/http/request"
	"github.com/IlhamSetiaji/julong-manpower-be/internal/http/response"
	jobResponse "github.com/IlhamSetiaji/julong-manpower-be/internal/http/response"
	"github.com/IlhamSetiaji/julong-manpower-be/utils"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type IJobPlafonMessage interface {
	SendCheckJobExistMessage(request request.CheckJobExistMessageRequest) (*jobResponse.CheckJobExistMessageResponse, error)
}

type JobPlafonMessage struct {
	Log *logrus.Logger
}

func NewJobPlafonMessage(log *logrus.Logger) IJobPlafonMessage {
	return &JobPlafonMessage{
		Log: log,
	}
}

func (m *JobPlafonMessage) SendCheckJobExistMessage(req request.CheckJobExistMessageRequest) (*jobResponse.CheckJobExistMessageResponse, error) {
	payload := map[string]interface{}{
		"job_id": req.ID,
	}

	docMsg := &request.RabbitMQRequest{
		ID:          uuid.New().String(),
		MessageType: "check_job_exist",
		MessageData: payload,
		ReplyTo:     "julong_manpower",
	}

	log.Printf("INFO: document message: %v", docMsg)

	// create channel and add to rchans with uid
	rchan := make(chan response.RabbitMQResponse)
	utils.Rchans[docMsg.ID] = rchan

	// publish rabbit message
	msg := utils.RabbitMsg{
		QueueName: "julong_sso",
		Message:   *docMsg,
	}
	utils.Pchan <- msg

	// wait for reply
	resp, err := waitReply(docMsg.ID, rchan)
	if err != nil {
		return nil, err
	}

	log.Printf("INFO: response: %v", resp)

	return &jobResponse.CheckJobExistMessageResponse{
		JobID: uuid.MustParse(resp.MessageData["job_id"].(string)),
		Exist: resp.MessageData["exists"].(bool),
	}, nil
}

func waitReply(id string, rchan chan response.RabbitMQResponse) (response.RabbitMQResponse, error) {
	for {
		select {
		case docReply := <-rchan:
			// responses received
			log.Printf("INFO: received reply: %v uid: %s", docReply, id)

			delete(utils.Rchans, id)
			return docReply, nil
		case <-time.After(10 * time.Second):
			// timeout
			log.Printf("ERROR: request timeout uid: %s", id)

			// remove channel from rchans
			delete(utils.Rchans, id)
			return response.RabbitMQResponse{}, errors.New("request timeout")
		}
	}
}

func JobPlafonMessageFactory(log *logrus.Logger) IJobPlafonMessage {
	return NewJobPlafonMessage(log)
}
