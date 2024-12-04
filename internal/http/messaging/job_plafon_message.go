package messaging

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/IlhamSetiaji/go-rabbitmq-utils/rabbitmq"
	"github.com/IlhamSetiaji/julong-manpower-be/internal/http/request"
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

func (m *JobPlafonMessage) SendCheckJobExistMessage(request request.CheckJobExistMessageRequest) (*jobResponse.CheckJobExistMessageResponse, error) {
	payload := map[string]interface{}{
		"message_type": "check_job_exists_request",
		"job_id":       request.ID,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		m.Log.Errorf("Failed to serialize message: %v", err)
		return nil, err
	}
	err = rabbitmq.PublishMessage("", "julong_queue_request", body)
	if err != nil {
		return nil, err
	}

	select {
	case response := <-utils.ResponseChannel:
		if request.ID == response["user_id"] && response["exists"].(bool) {
			m.Log.Printf("User validated. Proceeding with employee creation.")
			return &jobResponse.CheckJobExistMessageResponse{
				JobID: uuid.MustParse(request.ID),
				Exist: response["exists"].(bool),
			}, nil
		}
		m.Log.Errorf("User does not exist. Aborting.")
		return nil, errors.New("user does not exist")
	case <-time.After(5 * time.Second):
		m.Log.Errorf("Validation response timeout. Aborting.")
		return nil, errors.New("validation response timeout")
	}
}

func JobPlafonMessageFactory(log *logrus.Logger) IJobPlafonMessage {
	return NewJobPlafonMessage(log)
}
