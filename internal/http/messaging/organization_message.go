package messaging

import (
	"errors"
	"log"

	"github.com/IlhamSetiaji/julong-manpower-be/internal/http/request"
	"github.com/IlhamSetiaji/julong-manpower-be/internal/http/response"
	orgResponse "github.com/IlhamSetiaji/julong-manpower-be/internal/http/response"
	"github.com/IlhamSetiaji/julong-manpower-be/utils"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type IOrganizationMessage interface {
	SendFindOrganizationByIDMessage(request request.SendFindOrganizationByIDMessageRequest) (*orgResponse.SendFindOrganizationByIDMessageResponse, error)
	SendFindOrganizationLocationByIDMessage(request request.SendFindOrganizationLocationByIDMessageRequest) (*orgResponse.SendFindOrganizationLocationByIDMessageResponse, error)
}

type OrganizationMessage struct {
	Log *logrus.Logger
}

func NewOrganizationMessage(log *logrus.Logger) IOrganizationMessage {
	return &OrganizationMessage{
		Log: log,
	}
}

func (m *OrganizationMessage) SendFindOrganizationByIDMessage(req request.SendFindOrganizationByIDMessageRequest) (*orgResponse.SendFindOrganizationByIDMessageResponse, error) {
	payload := map[string]interface{}{
		"organization_id": req.ID,
	}

	docMsg := &request.RabbitMQRequest{
		ID:          uuid.New().String(),
		MessageType: "find_organization_by_id",
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
		return nil, errors.New("[SendFindOrganizationByIDMessage] " + errMsg)
	}

	return &orgResponse.SendFindOrganizationByIDMessageResponse{
		OrganizationID: resp.MessageData["organization_id"].(string),
		Name:           resp.MessageData["name"].(string),
	}, nil
}

func (m *OrganizationMessage) SendFindOrganizationLocationByIDMessage(req request.SendFindOrganizationLocationByIDMessageRequest) (*orgResponse.SendFindOrganizationLocationByIDMessageResponse, error) {
	payload := map[string]interface{}{
		"organization_location_id": req.ID,
	}

	docMsg := &request.RabbitMQRequest{
		ID:          uuid.New().String(),
		MessageType: "find_organization_location_by_id",
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
		return nil, errors.New("[SendFindOrganizationLocationByIDMessage] " + errMsg)
	}

	return &orgResponse.SendFindOrganizationLocationByIDMessageResponse{
		OrganizationLocationID: resp.MessageData["organization_location_id"].(string),
		Name:                   resp.MessageData["name"].(string),
	}, nil
}

func OrganizationMessageFactory(log *logrus.Logger) IOrganizationMessage {
	return NewOrganizationMessage(log)
}
