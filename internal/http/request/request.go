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
	case entity.MPPeriodStatusOpen, entity.MPPeriodStatusComplete, entity.MPPeriodStatusClose, entity.MPPPeriodStatusDraft, entity.MPPeriodStatusNotOpen:
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
	case entity.MPPlaningStatusDraft, entity.MPPlaningStatusReject, entity.MPPlaningStatusSubmit, entity.MPPlaningStatusComplete, entity.MPPlaningStatusApproved, entity.MPPlanningStatusInProgress, entity.MPPlaningStatusNeedApproval:
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
	case entity.MaritalStatusEnumSingle, entity.MaritalStatusEnumMarried, entity.MaritalStatusEnumDivorced, entity.MaritalStatusEnumWidowed, entity.MaritalStatusEnumAny:
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
	case entity.MPRequestStatusDraft, entity.MPRequestStatusSubmitted, entity.MPRequestStatusRejected, entity.MPRequestStatusApproved, entity.MPRequestStatusNeedApproval, entity.MPRequestStatusCompleted, entity.MPRequestStatusInProgress:
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
	case entity.MPPlanningApprovalHistoryLevelHRDUnit, entity.MPPlanningApprovalHistoryLevelDirekturUnit, entity.MPPlanningApprovalHistoryLevelRecruitment, entity.MPPlanningApprovalHistoryLevelCEO:
		return true
	default:
		return false
	}
}

func BatchHeaderApprovalStatusValidation(fl validator.FieldLevel) bool {
	status := fl.Field().String()
	if status == "" {
		return true
	}
	switch entity.BatchHeaderApprovalStatus(status) {
	case entity.BatchHeaderApprovalStatusApproved, entity.BatchHeaderApprovalStatusRejected, entity.BatchHeaderApprovalStatusNeedApproval, entity.BatchHeaderApprovalStatusCompleted:
		return true
	default:
		return false
	}
}

func MPRequestApprovalHistoryStatusValidation(fl validator.FieldLevel) bool {
	status := fl.Field().String()
	if status == "" {
		return true
	}
	switch entity.MPRequestApprovalHistoryStatus(status) {
	case entity.MPRequestApprovalHistoryStatusApproved, entity.MPRequestApprovalHistoryStatusRejected, entity.MPRequestApprovalHistoryStatusNeedApproval, entity.MPRequestApprovalHistoryStatusCompleted:
		return true
	default:
		return false
	}
}

func MPRequestApprovalHistoryLevelValidation(fl validator.FieldLevel) bool {
	status := fl.Field().String()
	if status == "" {
		return true
	}
	switch entity.MPRequestApprovalHistoryLevel(status) {
	case entity.MPRequestApprovalHistoryLevelStaff, entity.MPRequestApprovalHistoryLevelHeadDept, entity.MPRequestApprovalHistoryLevelVP, entity.MPRequestApprovalHistoryLevelCEO, entity.MPPRequestApprovalHistoryLevelHRDHO:
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

func MessageValidation(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email"
	case "datetime":
		return "Invalid datetime format"
	case "date_today_or_later":
		return "Date must be today or later"
	case "uuid":
		return "Invalid UUID"
	case "MaritalStatusValidation":
		return "Invalid marital status"
	case "MinimumEducationValidation":
		return "Invalid minimum education"
	case "MPRequestStatusValidation":
		return "Invalid MP request status"
	case "MPRequestTypeEnumValidation":
		return "Invalid MP request type"
	case "RecruitmentTypeEnumValidation":
		return "Invalid recruitment type"
	case "MPPlanningApprovalHistoryLevelValidation":
		return "Invalid MP planning approval history level"
	case "BatchHeaderApprovalStatusValidation":
		return "Invalid batch header approval status"
	case "MPRequestApprovalHistoryStatusValidation":
		return "Invalid MP request approval history status"
	case "MPRequestApprovalHistoryLevelValidation":
		return "Invalid MP request approval history level"
	case "dive":
		return "Invalid array"
	}
	return fe.Error() // default error
}

type RabbitMQRequest struct {
	ID          string                 `json:"id"`
	MessageType string                 `json:"message_type"`
	MessageData map[string]interface{} `json:"message_data"`
	ReplyTo     string                 `json:"reply_to"`
}
