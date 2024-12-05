package request

import (
	"github.com/IlhamSetiaji/julong-manpower-be/internal/entity"
	"github.com/go-playground/validator/v10"
)

func MPPPeriodStatusValidation(fl validator.FieldLevel) bool {
	status := fl.Field().String()
	if status == "" {
		return true
	}
	switch entity.MPPPeriodStatus(status) {
	case entity.MPPeriodStatusOpen, entity.MPPeriodStatusComplete, entity.MPPeriodStatusClose:
		return true
	default:
		return false
	}
}

func MPPlaningStatusValidation(fl validator.FieldLevel) bool {
	status := fl.Field().String()
	if status == "" {
		return true
	}
	switch entity.MPPlaningStatus(status) {
	case entity.MPPlaningStatusDraft, entity.MPPlaningStatusReject, entity.MPPlaningStatusSubmit, entity.MPPlaningStatusComplete:
		return true
	default:
		return false
	}
}

type RabbitMQRequest struct {
	ID          string                 `json:"id"`
	MessageType string                 `json:"message_type"`
	MessageData map[string]interface{} `json:"message_data"`
	ReplyTo     string                 `json:"reply_to"`
}
