package messaging

import (
	"log"

	"github.com/IlhamSetiaji/julong-manpower-be/internal/http/dto"
	"github.com/IlhamSetiaji/julong-manpower-be/internal/http/request"
	"github.com/IlhamSetiaji/julong-manpower-be/internal/http/response"
	"github.com/IlhamSetiaji/julong-manpower-be/utils"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type IJobMessage interface {
	SendFindJobDataByIdMessage(request request.SendFindJobByIDMessageRequest) (*response.JobResponse, error)
}

type JobMessage struct {
	Log *logrus.Logger
}

func NewJobMessage(log *logrus.Logger) IJobMessage {
	return &JobMessage{
		Log: log,
	}
}

func (m *JobMessage) SendFindJobDataByIdMessage(req request.SendFindJobByIDMessageRequest) (*response.JobResponse, error) {
	payload := map[string]interface{}{
		"job_id": req.ID,
	}

	docMsg := &request.RabbitMQRequest{
		ID:          uuid.New().String(),
		MessageType: "find_job_data_by_id",
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

	log.Printf("INFO: response: %v", resp.MessageData["job"])

	jobData := resp.MessageData["job"].(map[string]interface{})
	return dto.ConvertInterfaceToJobResponse(jobData), nil
}

func JobMessageFactory(log *logrus.Logger) IJobMessage {
	return NewJobMessage(log)
}
