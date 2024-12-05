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
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type IMPPlanningHandler interface {
	FindAllHeadersPaginated(ctx *gin.Context)
	FindById(ctx *gin.Context)
	Create(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
	FindAllLinesByHeaderIdPaginated(ctx *gin.Context)
	FindLineById(ctx *gin.Context)
	CreateLine(ctx *gin.Context)
	UpdateLine(ctx *gin.Context)
	DeleteLine(ctx *gin.Context)
}

type MPPlanningHandler struct {
	Log      *logrus.Logger
	Viper    *viper.Viper
	UseCase  usecase.IMPPlanningUseCase
	Validate *validator.Validate
}

func NewMPPlanningHandler(log *logrus.Logger, viper *viper.Viper, useCase usecase.IMPPlanningUseCase, validate *validator.Validate) IMPPlanningHandler {
	return &MPPlanningHandler{
		Log:      log,
		Viper:    viper,
		UseCase:  useCase,
		Validate: validate,
	}
}

func MPPlanningHandlerFactory(log *logrus.Logger, viper *viper.Viper) IMPPlanningHandler {
	usecase := usecase.MPPlanningUseCaseFactory(log)
	validate := config.NewValidator(viper)
	return NewMPPlanningHandler(log, viper, usecase, validate)
}

func (h *MPPlanningHandler) FindAllHeadersPaginated(ctx *gin.Context) {
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

	req := request.FindAllHeadersPaginatedMPPlanningRequest{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	resp, err := h.UseCase.FindAllHeadersPaginated(&req)
	if err != nil {
		h.Log.Errorf("[MPPlanningHandler.FindAllHeadersPaginated] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "error", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "find all headers paginated success", resp)
}

func (h *MPPlanningHandler) FindById(ctx *gin.Context) {
	id := ctx.Param("id")

	if id == "" {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "error", "id is required")
		return
	}

	req := request.FindHeaderByIdMPPlanningRequest{
		ID: id,
	}

	resp, err := h.UseCase.FindById(&req)
	if err != nil {
		h.Log.Errorf("[MPPlanningHandler.FindById] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "error", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "find by id success", resp.MPPlanningHeader)
}

func (h *MPPlanningHandler) Create(ctx *gin.Context) {
	var req request.CreateHeaderMPPlanningRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.Log.Errorf("[MPPlanningHandler.Create] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusBadRequest, "error", err.Error())
		return
	}

	if err := h.Validate.Struct(req); err != nil {
		h.Log.Errorf("[MPPlanningHandler.Create] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusBadRequest, "error", err.Error())
		return
	}

	resp, err := h.UseCase.Create(&req)
	if err != nil {
		h.Log.Errorf("[MPPlanningHandler.Create] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "error", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusCreated, "create success", resp)
}

func (h *MPPlanningHandler) Update(ctx *gin.Context) {
	var req request.UpdateHeaderMPPlanningRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.Log.Errorf("[MPPlanningHandler.Update] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusBadRequest, "error", err.Error())
		return
	}

	if err := h.Validate.Struct(req); err != nil {
		h.Log.Errorf("[MPPlanningHandler.Update] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusBadRequest, "error", err.Error())
		return
	}

	resp, err := h.UseCase.Update(&req)
	if err != nil {
		h.Log.Errorf("[MPPlanningHandler.Update] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "error", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "update success", resp)
}

func (h *MPPlanningHandler) Delete(ctx *gin.Context) {
	id := ctx.Param("id")

	if id == "" {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "error", "id is required")
		return
	}

	req := request.DeleteHeaderMPPlanningRequest{
		ID: id,
	}

	err := h.UseCase.Delete(&req)
	if err != nil {
		h.Log.Errorf("[MPPlanningHandler.Delete] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "error", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "delete success", nil)
}

func (h *MPPlanningHandler) FindAllLinesByHeaderIdPaginated(ctx *gin.Context) {
	headerId := ctx.Param("header_id")

	if headerId == "" {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "error", "header_id is required")
		return
	}

	page, err := strconv.Atoi(ctx.Query("page"))
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(ctx.Query("pageSize"))
	if err != nil || pageSize < 1 {
		pageSize = 10
	}

	req := request.FindAllLinesByHeaderIdPaginatedMPPlanningLineRequest{
		HeaderID: headerId,
		Page:     page,
		PageSize: pageSize,
	}

	resp, err := h.UseCase.FindAllLinesByHeaderIdPaginated(&req)
	if err != nil {
		h.Log.Errorf("[MPPlanningHandler.FindAllLinesByHeaderIdPaginated] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "error", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "find all lines by header id paginated success", resp)
}

func (h *MPPlanningHandler) FindLineById(ctx *gin.Context) {
	id := ctx.Param("id")

	if id == "" {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "error", "id is required")
		return
	}

	req := request.FindLineByIdMPPlanningLineRequest{
		ID: id,
	}

	resp, err := h.UseCase.FindLineById(&req)
	if err != nil {
		h.Log.Errorf("[MPPlanningHandler.FindLineById] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "error", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "find line by id success", resp.MPPlanningLine)
}

func (h *MPPlanningHandler) CreateLine(ctx *gin.Context) {
	var req request.CreateLineMPPlanningLineRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.Log.Errorf("[MPPlanningHandler.CreateLine] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusBadRequest, "error", err.Error())
		return
	}

	if err := h.Validate.Struct(req); err != nil {
		h.Log.Errorf("[MPPlanningHandler.CreateLine] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusBadRequest, "error", err.Error())
		return
	}

	resp, err := h.UseCase.CreateLine(&req)
	if err != nil {
		h.Log.Errorf("[MPPlanningHandler.CreateLine] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "error", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusCreated, "create line success", resp)
}

func (h *MPPlanningHandler) UpdateLine(ctx *gin.Context) {
	var req request.UpdateLineMPPlanningLineRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.Log.Errorf("[MPPlanningHandler.UpdateLine] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusBadRequest, "error", err.Error())
		return
	}

	if err := h.Validate.Struct(req); err != nil {
		h.Log.Errorf("[MPPlanningHandler.UpdateLine] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusBadRequest, "error", err.Error())
		return
	}

	resp, err := h.UseCase.UpdateLine(&req)
	if err != nil {
		h.Log.Errorf("[MPPlanningHandler.UpdateLine] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "error", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "update line success", resp)
}

func (h *MPPlanningHandler) DeleteLine(ctx *gin.Context) {
	id := ctx.Param("id")

	if id == "" {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "error", "id is required")
		return
	}

	req := request.DeleteLineMPPlanningLineRequest{
		ID: id,
	}

	err := h.UseCase.DeleteLine(&req)
	if err != nil {
		h.Log.Errorf("[MPPlanningHandler.DeleteLine] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "error", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "delete line success", nil)
}