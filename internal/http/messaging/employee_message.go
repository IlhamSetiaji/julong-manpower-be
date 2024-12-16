package messaging

import (
	"errors"
	"log"

	"github.com/IlhamSetiaji/julong-manpower-be/internal/http/dto"
	"github.com/IlhamSetiaji/julong-manpower-be/internal/http/request"
	"github.com/IlhamSetiaji/julong-manpower-be/internal/http/response"
	"github.com/IlhamSetiaji/julong-manpower-be/utils"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type IEmployeeMessage interface {
	SendFindEmployeeByIDMessage(req request.SendFindEmployeeByIDMessageRequest) (*response.EmployeeResponse, error)
}

type EmployeeMessage struct {
	Log *logrus.Logger
}

func NewEmployeeMessage(log *logrus.Logger) IEmployeeMessage {
	return &EmployeeMessage{
		Log: log,
	}
}

func (m *EmployeeMessage) SendFindEmployeeByIDMessage(req request.SendFindEmployeeByIDMessageRequest) (*response.EmployeeResponse, error) {
	payload := map[string]interface{}{
		"employee_id": req.ID,
	}

	docMsg := &request.RabbitMQRequest{
		ID:          uuid.New().String(),
		MessageType: "find_employee_by_id",
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

	if errMsg, ok := resp.MessageData["error"]; ok {
		return nil, errors.New("[EmployeeMessage.SendFindEmployeeByIDMessage] " + errMsg.(string))
	}

	employeeData := resp.MessageData["employee"].(map[string]interface{})
	employee := dto.ConvertInterfaceToEmployeeResponse(employeeData)

	return employee, nil
}

func EmployeeMessageFactory(log *logrus.Logger) IEmployeeMessage {
	return NewEmployeeMessage(log)
}
