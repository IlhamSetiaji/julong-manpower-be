package handler

import (
	"net/http"

	"github.com/IlhamSetiaji/julong-manpower-be/internal/config"
	"github.com/IlhamSetiaji/julong-manpower-be/internal/entity"
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
	FindByStatus(c *gin.Context)
	FindById(c *gin.Context)
	FindDocumentByID(c *gin.Context)
	FindByCurrentDocumentDateAndStatus(c *gin.Context)
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

func (h *BatchHandler) FindByStatus(c *gin.Context) {
	status := c.Param("status")
	approvalStatus := entity.BatchHeaderApprovalStatus(status)
	batch, err := h.UseCase.FindByStatus(approvalStatus)
	if err != nil {
		h.Log.Error(err)
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to find batch by status", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Batch found", batch)
}

func (h *BatchHandler) FindById(c *gin.Context) {
	id := c.Param("id")
	batch, err := h.UseCase.FindById(id)
	if err != nil {
		h.Log.Error(err)
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to find batch by id", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Batch found", batch)
}

func (h *BatchHandler) FindDocumentByID(c *gin.Context) {
	id := c.Param("id")
	batch, err := h.UseCase.FindDocumentByID(id)
	if err != nil {
		h.Log.Error(err)
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to find document by id", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Document found", batch)
}

func (h *BatchHandler) FindByCurrentDocumentDateAndStatus(c *gin.Context) {
	status := c.Param("status")
	approvalStatus := entity.BatchHeaderApprovalStatus(status)
	batch, err := h.UseCase.FindByCurrentDocumentDateAndStatus(approvalStatus)
	if err != nil {
		h.Log.Error(err)
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to find batch by current document date and status", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Batch found", batch)
}

func BatchHandlerFactory(log *logrus.Logger, viper *viper.Viper) IBatchHandler {
	validate := config.NewValidator(viper)
	userHelper := helper.NewUserHelper(log)
	useCase := usecase.BatchUsecaseFactory(viper, log)
	return NewBatchHandler(log, viper, useCase, validate, userHelper)
}
