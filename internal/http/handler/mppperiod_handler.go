package handler

import (
	"net/http"

	"github.com/IlhamSetiaji/julong-manpower-be/internal/config"
	"github.com/IlhamSetiaji/julong-manpower-be/internal/http/request"
	"github.com/IlhamSetiaji/julong-manpower-be/internal/usecase"
	"github.com/IlhamSetiaji/julong-manpower-be/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type IMPPPeriodHander interface {
	FindAllPaginated(ctx *gin.Context)
	FindById(ctx *gin.Context)
	Create(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
}

type MPPPeriodHandler struct {
	Log      *logrus.Logger
	Viper    *viper.Viper
	UseCase  usecase.IMPPPeriodUseCase
	Validate *validator.Validate
}

func NewMPPPeriodHandler(log *logrus.Logger, viper *viper.Viper, useCase usecase.IMPPPeriodUseCase, validate *validator.Validate) IMPPPeriodHander {
	return &MPPPeriodHandler{
		Log:      log,
		Viper:    viper,
		UseCase:  useCase,
		Validate: validate,
	}
}

func MPPPeriodHandlerFactory(log *logrus.Logger, viper *viper.Viper) IMPPPeriodHander {
	usecase := usecase.MPPPeriodUseCaseFactory(log)
	validate := config.NewValidator(viper)
	return NewMPPPeriodHandler(log, viper, usecase, validate)
}

func (h *MPPPeriodHandler) FindAllPaginated(ctx *gin.Context) {
	var req request.FindAllPaginatedMPPPeriodRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		h.Log.Errorf("[MPPPeriodHandler.FindAllPaginated] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusBadRequest, "error", err.Error())
		return
	}

	resp, err := h.UseCase.FindAllPaginated(req)
	if err != nil {
		h.Log.Errorf("[MPPPeriodHandler.FindAllPaginated] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "error", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "success", resp)
}

func (h *MPPPeriodHandler) FindById(ctx *gin.Context) {
	var req request.FindByIdMPPPeriodRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		h.Log.Errorf("[MPPPeriodHandler.FindById] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusBadRequest, "error", err.Error())
		return
	}

	resp, err := h.UseCase.FindById(req)
	if err != nil {
		h.Log.Errorf("[MPPPeriodHandler.FindById] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "error", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "success", resp)
}

func (h *MPPPeriodHandler) Create(ctx *gin.Context) {
	var req request.CreateMPPPeriodRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.Log.Errorf("[MPPPeriodHandler.Create] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusBadRequest, "error", err.Error())
		return
	}

	if err := h.Validate.Struct(req); err != nil {
		h.Log.Errorf("[MPPPeriodHandler.Create] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusBadRequest, "error", err.Error())
		return
	}

	resp, err := h.UseCase.Create(req)
	if err != nil {
		h.Log.Errorf("[MPPPeriodHandler.Create] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "error", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusCreated, "success", resp)
}

func (h *MPPPeriodHandler) Update(ctx *gin.Context) {
	var req request.UpdateMPPPeriodRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.Log.Errorf("[MPPPeriodHandler.Update] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusBadRequest, "error", err.Error())
		return
	}

	if err := h.Validate.Struct(req); err != nil {
		h.Log.Errorf("[MPPPeriodHandler.Update] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusBadRequest, "error", err.Error())
		return
	}

	resp, err := h.UseCase.Update(req)
	if err != nil {
		h.Log.Errorf("[MPPPeriodHandler.Update] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "error", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "success", resp)
}

func (h *MPPPeriodHandler) Delete(ctx *gin.Context) {
	var req request.DeleteMPPPeriodRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		h.Log.Errorf("[MPPPeriodHandler.Delete] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusBadRequest, "error", err.Error())
		return
	}

	err := h.UseCase.Delete(req)
	if err != nil {
		h.Log.Errorf("[MPPPeriodHandler.Delete] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "error", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "success", nil)
}
