package handler

import (
	"net/http"

	"github.com/IlhamSetiaji/julong-manpower-be/internal/http/helper"
	"github.com/IlhamSetiaji/julong-manpower-be/internal/http/request"
	"github.com/IlhamSetiaji/julong-manpower-be/internal/usecase"
	"github.com/IlhamSetiaji/julong-manpower-be/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type IBatchHandler interface {
	CreateBatchHeaderAndLines(c *gin.Context)
}

type BatchHandler struct {
	Log        *logrus.Logger
	Viper      *viper.Viper
	UseCase    usecase.IBatchUsecase
	Validate   *validator.Validate
	UserHelper helper.IUserHelper
}

func NewBatchHandler(log *logrus.Logger, viper *viper.Viper, useCase usecase.IBatchUsecase, validate *validator.Validate, userHelper helper.IUserHelper) IBatchHandler {
	return &BatchHandler{
		Log:        log,
		Viper:      viper,
		UseCase:    useCase,
		Validate:   validate,
		UserHelper: userHelper,
	}
}

func (h *BatchHandler) CreateBatchHeaderAndLines(c *gin.Context) {
	var req request.CreateBatchHeaderAndLinesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.Log.Error(err)
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	if err := h.Validate.Struct(req); err != nil {
		h.Log.Error(err)
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	batchHeader, err := h.UseCase.CreateBatchHeaderAndLines(&req)
	if err != nil {
		h.Log.Error(err)
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create batch header and lines", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Batch header and lines created", batchHeader)
}

func BatchHandlerFactory(log *logrus.Logger, viper *viper.Viper) IBatchHandler {
	validate := validator.New()
	userHelper := helper.NewUserHelper(log)
	useCase := usecase.BatchUsecaseFactory(viper, log)
	return NewBatchHandler(log, viper, useCase, validate, userHelper)
}
