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

func RecruitmentTypeValidation(fl validator.FieldLevel) bool {
	status := fl.Field().String()
	if status == "" {
		return true
	}
	switch entity.RecruitmentTypeEnum(status) {
	case entity.RecruitmentTypeEnumMT, entity.RecruitmentTypeEnumPH, entity.RecruitmentTypeEnumNS:
		return true
	default:
		return false
	}
}

func MaritalStatusValidation(fl validator.FieldLevel) bool {
	status := fl.Field().String()
	if status == "" {
		return true
	}
	switch entity.MaritalStatusEnum(status) {
	case entity.MaritalStatusEnumSingle, entity.MaritalStatusEnumMarried, entity.MaritalStatusEnumDivorced, entity.MaritalStatusEnumWidowed:
		return true
	default:
		return false
	}
}

func MinimumEducationValidation(fl validator.FieldLevel) bool {
	status := fl.Field().String()
	if status == "" {
		return true
	}
	switch entity.EducationEnum(status) {
	case entity.EducationEnumSD, entity.EducationEnumSMP, entity.EducationEnumSMA, entity.EducationEnumD3, entity.EducationEnumS1, entity.EducationEnumS2, entity.EducationEnumS3:
		return true
	default:
		return false
	}
}

func MPRequestStatusValidation(fl validator.FieldLevel) bool {
	status := fl.Field().String()
	if status == "" {
		return true
	}
	switch entity.MPRequestStatus(status) {
	case entity.MPRequestStatusDraft, entity.MPRequestStatusSubmitted, entity.MPRequestStatusRejected, entity.MPRequestStatusApproved:
		return true
	default:
		return false
	}
}

func MPRequestTypeEnumValidation(fl validator.FieldLevel) bool {
	status := fl.Field().String()
	if status == "" {
		return true
	}
	switch entity.MPRequestTypeEnum(status) {
	case entity.MPRequestTypeEnumOnBudget, entity.MPRequestTypeEnumOffBudget:
		return true
	default:
		return false
	}
}

func RecruitmentTypeEnumValidation(fl validator.FieldLevel) bool {
	status := fl.Field().String()
	if status == "" {
		return true
	}
	switch entity.RecruitmentTypeEnum(status) {
	case entity.RecruitmentTypeEnumMT, entity.RecruitmentTypeEnumPH, entity.RecruitmentTypeEnumNS:
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
