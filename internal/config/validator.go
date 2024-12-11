package config

import (
	"github.com/IlhamSetiaji/julong-manpower-be/internal/http/request"
	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

func NewValidator(viper *viper.Viper) *validator.Validate {
	validate := validator.New()
	validate.RegisterValidation("MPPPeriodStatusValidation", request.MPPPeriodStatusValidation)
	validate.RegisterValidation("MPPlaningStatusValidation", request.MPPlaningStatusValidation)
	validate.RegisterValidation("RecruitmentTypeValidation", request.RecruitmentTypeValidation)
	validate.RegisterValidation("MaritalStatusValidation", request.MaritalStatusValidation)
	validate.RegisterValidation("MinimumEducationValidation", request.MinimumEducationValidation)
	validate.RegisterValidation("MPRequestStatusValidation", request.MPRequestStatusValidation)
	validate.RegisterValidation("MPRequestTypeEnumValidation", request.MPRequestTypeEnumValidation)
	validate.RegisterValidation("RecruitmentTypeEnumValidation", request.RecruitmentTypeEnumValidation)
	validate.RegisterValidation("date_today_or_later", request.ValidateDateMoreThanEqualToday)
	return validate
}
