package messaging

import (
	"errors"
	"log"

	"github.com/IlhamSetiaji/julong-manpower-be/internal/http/request"
	"github.com/IlhamSetiaji/julong-manpower-be/internal/http/response"
	"github.com/IlhamSetiaji/julong-manpower-be/utils"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type IMPRequestMessage interface {
	SendCloneMPR(id uuid.UUID) (*string, error)
}

type MPRequestMessage struct {
	Log *logrus.Logger
}

func NewMPRequestMessage(log *logrus.Logger) IMPRequestMessage {
	return &MPRequestMessage{
		Log: log,
	}
}

func MPRequestMessageFactory(log *logrus.Logger) IMPRequestMessage {
	return NewMPRequestMessage(log)
}

func (m *MPRequestMessage) SendCloneMPR(id uuid.UUID) (*string, error) {
	payload := map[string]interface{}{
		"mpr_clone_id": id.String(),
	}

	docMsg := &request.RabbitMQRequest{
		ID:          uuid.New().String(),
		MessageType: "clone_mp_request",
		MessageData: payload,
		ReplyTo:     "julong_manpower",
	}

	log.Printf("INFO: document message: %v", docMsg)

	// create channel and add to rchans with uid
	rchan := make(chan response.RabbitMQResponse)
	utils.Rchans[docMsg.ID] = rchan

	// publish rabbit message
	msg := utils.RabbitMsg{
		QueueName: "julong_recruitment",
		Message:   *docMsg,
	}
	utils.Pchan <- msg

	// wait for reply
	resp, err := waitReply(docMsg.ID, rchan)
	if err != nil {
		return nil, err
	}

	log.Printf("INFO: response: %v", resp.MessageData["mpr"])

	if errMsg, ok := resp.MessageData["error"].(string); ok && errMsg != "" {
		return nil, errors.New("[SendCloneMPR] " + errMsg)
	}

	messsage := "Success clone MPR"

	return &messsage, nil
}
