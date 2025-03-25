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

type IGradeMessage interface {
	SendFindByIDMessage(id string) (*response.GradeResponse, error)
}

type GradeMessage struct {
	Log *logrus.Logger
}

func NewGradeMessage(log *logrus.Logger) IGradeMessage {
	return &GradeMessage{
		Log: log,
	}
}

func (m *GradeMessage) SendFindByIDMessage(id string) (*response.GradeResponse, error) {
	payload := map[string]interface{}{
		"id": id,
	}

	docMsg := &request.RabbitMQRequest{
		ID:          uuid.New().String(),
		MessageType: "find_grade_by_id",
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
		return nil, errors.New("[GradeMessage.SendFindByIDMessage] " + errMsg)
	}

	grade := &response.GradeResponse{
		ID:           resp.MessageData["id"].(string),
		Name:         resp.MessageData["name"].(string),
		JobLevelID:   resp.MessageData["job_level_id"].(string),
		JobLevelName: resp.MessageData["job_level_name"].(string),
	}

	return grade, nil
}

func GradeMessageFactory(log *logrus.Logger) IGradeMessage {
	return NewGradeMessage(log)
}
