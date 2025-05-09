package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/IlhamSetiaji/julong-manpower-be/internal/config"
	"github.com/IlhamSetiaji/julong-manpower-be/internal/entity"
	"github.com/IlhamSetiaji/julong-manpower-be/internal/http/helper"
	"github.com/IlhamSetiaji/julong-manpower-be/internal/http/request"
	"github.com/IlhamSetiaji/julong-manpower-be/internal/http/service"
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
	FindByCurrentDateAndStatus(ctx *gin.Context)
	FindByStatus(ctx *gin.Context)
	UpdateStatusByDate(ctx *gin.Context)
}

type MPPPeriodHandler struct {
	Log                 *logrus.Logger
	Viper               *viper.Viper
	UseCase             usecase.IMPPPeriodUseCase
	Validate            *validator.Validate
	NotificationService service.INotificationService
	UserHelper          helper.IUserHelper
}

func NewMPPPeriodHandler(
	log *logrus.Logger,
	viper *viper.Viper,
	useCase usecase.IMPPPeriodUseCase,
	validate *validator.Validate,
	notificationService service.INotificationService,
	userHelper helper.IUserHelper,
) IMPPPeriodHander {
	return &MPPPeriodHandler{
		Log:                 log,
		Viper:               viper,
		UseCase:             useCase,
		Validate:            validate,
		NotificationService: notificationService,
		UserHelper:          userHelper,
	}
}

func MPPPeriodHandlerFactory(log *logrus.Logger, viper *viper.Viper) IMPPPeriodHander {
	usecase := usecase.MPPPeriodUseCaseFactory(log)
	validate := config.NewValidator(viper)
	notificationService := service.NotificationServiceFactory(viper, log)
	userHelper := helper.UserHelperFactory(log)
	return NewMPPPeriodHandler(log, viper, usecase, validate, notificationService, userHelper)
}

func (h *MPPPeriodHandler) FindAllPaginated(ctx *gin.Context) {
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

	// user, err := middleware.GetUser(ctx, h.Log)
	// if err != nil {
	// 	h.Log.Errorf("Error when getting user: %v", err)
	// 	utils.ErrorResponse(ctx, 500, "error", err.Error())
	// 	return
	// }
	// if user == nil {
	// 	h.Log.Errorf("User not found")
	// 	utils.ErrorResponse(ctx, 404, "error", "User not found")
	// 	return
	// }
	// userUUID, err := h.UserHelper.GetUserId(user)
	// if err != nil {
	// 	h.Log.Errorf("Error when getting user id: %v", err)
	// 	utils.ErrorResponse(ctx, 500, "error", err.Error())
	// 	return
	// }

	// err = h.NotificationService.CreatePeriodNotification(userUUID.String())
	// if err != nil {
	// 	h.Log.Errorf("Error when creating notification: %v", err)
	// 	utils.ErrorResponse(ctx, 500, "error", err.Error())
	// 	return
	// }

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

func (h *MPPPeriodHandler) FindByCurrentDateAndStatus(ctx *gin.Context) {
	status := ctx.Query("status")

	if status == "" {
		status = "open"
	}

	// req := request.FindByCurrentDateAndStatusMPPPeriodRequest{
	// 	Status: entity.MPPPeriodStatus(status),
	// }

	resp, err := h.UseCase.FindByStatus(entity.MPPPeriodStatus(status))
	if err != nil {
		h.Log.Errorf("[MPPPeriodHandler.FindByCurrentDateAndStatus] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "error", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "find by current date and status success", resp.MPPPeriod)
}

func (h *MPPPeriodHandler) FindByStatus(ctx *gin.Context) {
	status := ctx.Query("status")

	if status == "" {
		h.Log.Errorf("[MPPPeriodHandler.FindByStatus] " + "status is required")
		utils.ErrorResponse(ctx, http.StatusBadRequest, "error", "status is required")
		return
	}

	resp, err := h.UseCase.FindByStatus(entity.MPPPeriodStatus(status))
	if err != nil {
		h.Log.Errorf("[MPPPeriodHandler.FindByStatus] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "error", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "find by status success", resp)
}

func (h *MPPPeriodHandler) UpdateStatusByDate(ctx *gin.Context) {
	date := ctx.Query("date")

	if date == "" {
		h.Log.Errorf("[MPPPeriodHandler.UpdateStatusByDate] " + "date is required")
		utils.ErrorResponse(ctx, http.StatusBadRequest, "error", "date is required")
		return
	}

	status := ctx.Query("status")
	if status == "" {
		h.Log.Errorf("[MPPPeriodHandler.UpdateStatusByDate] " + "status is required")
		utils.ErrorResponse(ctx, http.StatusBadRequest, "error", "status is required")
		return
	}
	parsedDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		h.Log.Errorf("[MPPPeriodHandler.UpdateStatusByDate] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusBadRequest, "error", "invalid date format")
		return
	}
	if status == "open" {
		err = h.UseCase.UpdateStatusToOpenByDate(parsedDate)
		if err != nil {
			h.Log.Errorf("[MPPPeriodHandler.UpdateStatusByDate] " + err.Error())
			utils.ErrorResponse(ctx, http.StatusInternalServerError, "error", err.Error())
			return
		}
	} else if status == "close" {
		err := h.UseCase.UpdateStatusToCloseByDate(parsedDate)
		if err != nil {
			h.Log.Errorf("[MPPPeriodHandler.UpdateStatusByDate] " + err.Error())
			utils.ErrorResponse(ctx, http.StatusInternalServerError, "error", err.Error())
			return
		}
	} else {
		h.Log.Errorf("[MPPPeriodHandler.UpdateStatusByDate] " + "status is invalid")
		utils.ErrorResponse(ctx, http.StatusBadRequest, "error", "status is invalid")
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "update status by date success", nil)
}
