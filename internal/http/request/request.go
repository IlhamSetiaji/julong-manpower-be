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
