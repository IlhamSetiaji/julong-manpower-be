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

type IJobPlafonHandler interface {
	FindAllPaginated(ctx *gin.Context)
	FindById(ctx *gin.Context)
	FindByJobId(ctx *gin.Context)
	Create(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
	SyncJobPlafon(ctx *gin.Context)
}

type JobPlafonHandler struct {
	Log      *logrus.Logger
	Viper    *viper.Viper
	UseCase  usecase.IJobPlafonUseCase
	Validate *validator.Validate
}

func NewJobPlafonHandler(log *logrus.Logger, viper *viper.Viper, useCase usecase.IJobPlafonUseCase, validate *validator.Validate) IJobPlafonHandler {
	return &JobPlafonHandler{
		Log:      log,
		Viper:    viper,
		UseCase:  useCase,
		Validate: validate,
	}
}

func JobPlafonHandlerFactory(log *logrus.Logger, viper *viper.Viper) IJobPlafonHandler {
	usecase := usecase.JobPlafonUseCaseFactory(log)
	validate := config.NewValidator(viper)
	return NewJobPlafonHandler(log, viper, usecase, validate)
}

func (h *JobPlafonHandler) SyncJobPlafon(ctx *gin.Context) {
	err := h.UseCase.SyncJobPlafon()
	if err != nil {
		h.Log.Errorf("[JobPlafonHandler.SyncJobPlafon] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "error", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "job plafon synced successfully", nil)
}

func (h *JobPlafonHandler) FindAllPaginated(ctx *gin.Context) {
	page, err := strconv.Atoi(ctx.Query("page"))
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(ctx.Query("page_size"))
	if err != nil || pageSize < 1 {
		pageSize = 10
	}

	search := ctx.Query("search")
	if search == "" {
		search = ""
	}

	resp, err := h.UseCase.FindAllPaginated(&request.FindAllPaginatedJobPlafonRequest{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	})
	if err != nil {
		h.Log.Errorf("[JobPlafonHandler.FindAllPaginated] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "error", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "find all paginated success", resp)
}

func (h *JobPlafonHandler) FindById(ctx *gin.Context) {
	id := ctx.Param("id")

	resp, err := h.UseCase.FindById(&request.FindByIdJobPlafonRequest{
		ID: id,
	})
	if err != nil {
		h.Log.Errorf("[JobPlafonHandler.FindById] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "error", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "find by id success", resp.JobPlafon)
}

func (h *JobPlafonHandler) FindByJobId(ctx *gin.Context) {
	jobId := ctx.Param("job_id")

	if jobId == "" {
		h.Log.Errorf("[JobPlafonHandler.FindByJobId] job id is required")
		utils.ErrorResponse(ctx, http.StatusBadRequest, "error", "job id is required")
		return
	}

	resp, err := h.UseCase.FindByJobId(&request.FindByJobIdJobPlafonRequest{
		JobID: jobId,
	})

	if err != nil {
		h.Log.Errorf("[JobPlafonHandler.FindByJobId] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "error", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "find by job id success", resp.JobPlafon)
}

func (h *JobPlafonHandler) Create(ctx *gin.Context) {
	var req request.CreateJobPlafonRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.Log.Errorf("[JobPlafonHandler.Create] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusBadRequest, "error", err.Error())
		return
	}

	if err := h.Validate.Struct(req); err != nil {
		h.Log.Errorf("[JobPlafonHandler.Create] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusBadRequest, "error", err.Error())
		return
	}

	resp, err := h.UseCase.Create(&req)
	if err != nil {
		h.Log.Errorf("[JobPlafonHandler.Create] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "error", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusCreated, "job plafon created successfully", resp.JobPlafon)
}

func (h *JobPlafonHandler) Update(ctx *gin.Context) {
	var req request.UpdateJobPlafonRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.Log.Errorf("[JobPlafonHandler.Update] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusBadRequest, "error", err.Error())
		return
	}

	if err := h.Validate.Struct(req); err != nil {
		h.Log.Errorf("[JobPlafonHandler.Update] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusBadRequest, "error", err.Error())
		return
	}

	resp, err := h.UseCase.Update(&req)
	if err != nil {
		h.Log.Errorf("[JobPlafonHandler.Update] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "error", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "job plafon updated successfully", resp.JobPlafon)
}

func (h *JobPlafonHandler) Delete(ctx *gin.Context) {
	id := ctx.Param("id")

	err := h.UseCase.Delete(&request.DeleteJobPlafonRequest{
		ID: id,
	})
	if err != nil {
		h.Log.Errorf("[JobPlafonHandler.Delete] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "error", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "job plafon deleted successfully", nil)
}
