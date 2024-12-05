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

type IUserMessage interface {
	SendFindUserByIDMessage(request request.SendFindUserByIDMessageRequest) (*response.SendFindUserByIDResponse, error)
}

type UserMessage struct {
	Log *logrus.Logger
}

func NewUserMessage(log *logrus.Logger) IUserMessage {
	return &UserMessage{
		Log: log,
	}
}

func (m *UserMessage) SendFindUserByIDMessage(req request.SendFindUserByIDMessageRequest) (*response.SendFindUserByIDResponse, error) {
	payload := map[string]interface{}{
		"user_id": req.ID,
	}

	docMsg := &request.RabbitMQRequest{
		ID:          uuid.New().String(),
		MessageType: "find_user_by_id",
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

	if errMsg, ok := resp.MessageData["error"].(string); ok && errMsg != "" {
		return nil, errors.New(errMsg)
	}

	return &response.SendFindUserByIDResponse{
		ID:   resp.MessageData["user_id"].(string),
		Name: resp.MessageData["name"].(string),
	}, nil
}

func UserMessageFactory(log *logrus.Logger) IUserMessage {
	return NewUserMessage(log)
}
