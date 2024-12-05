package handler

import (
	"net/http"
	"strconv"

	"github.com/IlhamSetiaji/julong-manpower-be/internal/config"
	"github.com/IlhamSetiaji/julong-manpower-be/internal/http/request"
	"github.com/IlhamSetiaji/julong-manpower-be/internal/usecase"
	"github.com/IlhamSetiaji/julong-manpower-be/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
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
	page, err := strconv.Atoi(ctx.Query("page"))
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(ctx.Query("pageSize"))
	if err != nil || pageSize < 1 {
		pageSize = 10
	}

	search := ctx.Query("search")
	if search == "" {
		search = ""
	}

	req := request.FindAllPaginatedMPPPeriodRequest{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	resp, err := h.UseCase.FindAllPaginated(req)
	if err != nil {
		h.Log.Errorf("[MPPPeriodHandler.FindAllPaginated] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "error", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "find all paginated success", resp)
}

func (h *MPPPeriodHandler) FindById(ctx *gin.Context) {
	id := ctx.Param("id")

	req := request.FindByIdMPPPeriodRequest{
		ID: uuid.MustParse(id),
	}

	resp, err := h.UseCase.FindById(req)
	if err != nil {
		h.Log.Errorf("[MPPPeriodHandler.FindById] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "error", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "find by id success", resp.MPPPeriod)
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

	utils.SuccessResponse(ctx, http.StatusCreated, "mpp period created successfully", resp)
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

	utils.SuccessResponse(ctx, http.StatusOK, "mpp period updated successfully", resp)
}

func (h *MPPPeriodHandler) Delete(ctx *gin.Context) {
	id := ctx.Param("id")

	req := request.DeleteMPPPeriodRequest{
		ID: uuid.MustParse(id),
	}

	err := h.UseCase.Delete(req)
	if err != nil {
		h.Log.Errorf("[MPPPeriodHandler.Delete] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "error", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "mpp period deleted successfully", nil)
}
