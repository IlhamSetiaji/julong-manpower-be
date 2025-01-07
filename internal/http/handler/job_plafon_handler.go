package handler

import (
	"net/http"
	"strconv"

	"github.com/IlhamSetiaji/julong-manpower-be/internal/config"
	"github.com/IlhamSetiaji/julong-manpower-be/internal/http/helper"
	"github.com/IlhamSetiaji/julong-manpower-be/internal/http/middleware"
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
	Log        *logrus.Logger
	Viper      *viper.Viper
	UseCase    usecase.IJobPlafonUseCase
	Validate   *validator.Validate
	UserHelper helper.IUserHelper
}

func NewJobPlafonHandler(log *logrus.Logger, viper *viper.Viper, useCase usecase.IJobPlafonUseCase, validate *validator.Validate, userHelper helper.IUserHelper) IJobPlafonHandler {
	return &JobPlafonHandler{
		Log:        log,
		Viper:      viper,
		UseCase:    useCase,
		Validate:   validate,
		UserHelper: userHelper,
	}
}

func JobPlafonHandlerFactory(log *logrus.Logger, viper *viper.Viper) IJobPlafonHandler {
	usecase := usecase.JobPlafonUseCaseFactory(log)
	validate := config.NewValidator(viper)
	userHelper := helper.UserHelperFactory(log)
	return NewJobPlafonHandler(log, viper, usecase, validate, userHelper)
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

	user, err := middleware.GetUser(ctx, h.Log)
	if err != nil {
		h.Log.Errorf("Error when getting user: %v", err)
		utils.ErrorResponse(ctx, 500, "error", err.Error())
		return
	}
	if user == nil {
		h.Log.Errorf("User not found")
		utils.ErrorResponse(ctx, 404, "error", "User not found")
		return
	}
	requestorUUID, err := h.UserHelper.GetEmployeeId(user)
	if err != nil {
		h.Log.Errorf("Error when getting employee id: %v", err)
		utils.ErrorResponse(ctx, 500, "error", err.Error())
		return
	}
	orgUUID, err := h.UserHelper.GetOrganizationID(user)
	if err != nil {
		h.Log.Errorf("Error when getting organization id: %v", err)
		utils.ErrorResponse(ctx, 500, "error", err.Error())
		return
	}

	resp, err := h.UseCase.FindAllPaginated(&request.FindAllPaginatedJobPlafonRequest{
		Page:           page,
		PageSize:       pageSize,
		Search:         search,
		RequestorID:    requestorUUID.String(),
		OrganizationID: orgUUID.String(),
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
