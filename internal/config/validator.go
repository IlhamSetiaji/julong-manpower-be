package config

import (
	"github.com/IlhamSetiaji/julong-manpower-be/internal/http/request"
	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

func NewValidator(viper *viper.Viper) *validator.Validate {
	validate := validator.New()
	validate.RegisterValidation("MPPPeriodStatusValidation", request.MPPPeriodStatusValidation)
	return validate
}
