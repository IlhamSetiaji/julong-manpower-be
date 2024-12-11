package request

import (
	"time"

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
	switch entity.EducationLevelEnum(status) {
	case entity.EducationLevelEnumD3, entity.EducationLevelEnumD4, entity.EducationLevelEnumSMA,
		entity.EducationLevelEnumBachelor, entity.EducationLevelEnumDoctoral, entity.EducationLevelEnumMaster, entity.EducationLevelEnumD1, entity.EducationLevelEnumD2,
		entity.EducationLevelEnumSD, entity.EducationLevelEnumSMP, entity.EducationLevelEnumTK:
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

func MPPlanningApprovalHistoryLevelValidation(fl validator.FieldLevel) bool {
	status := fl.Field().String()
	if status == "" {
		return true
	}
	switch entity.MPPlanningApprovalHistoryLevel(status) {
	case entity.MPPlanningApprovalHistoryLevelManager, entity.MPPlanningApprovalHistoryLevelRecruitment:
		return true
	default:
		return false
	}
}

func ValidateDateMoreThanEqualToday(fl validator.FieldLevel) bool {
	startDateStr := fl.Field().String()
	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		return false
	}

	today := time.Now().Truncate(24 * time.Hour)
	return !startDate.Before(today)
}

type RabbitMQRequest struct {
	ID          string                 `json:"id"`
	MessageType string                 `json:"message_type"`
	MessageData map[string]interface{} `json:"message_data"`
	ReplyTo     string                 `json:"reply_to"`
}
