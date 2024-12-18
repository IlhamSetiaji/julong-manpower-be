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

type IMPRequestHandler interface {
	Create(ctx *gin.Context)
	FindAllPaginated(ctx *gin.Context)
}

type MPRequestHandler struct {
	Log      *logrus.Logger
	Viper    *viper.Viper
	UseCase  usecase.IMPRequestUseCase
	Validate *validator.Validate
}

func NewMPRequestHandler(
	log *logrus.Logger,
	viper *viper.Viper,
	useCase usecase.IMPRequestUseCase,
	validate *validator.Validate,
) IMPRequestHandler {
	return &MPRequestHandler{
		Log:      log,
		Viper:    viper,
		UseCase:  useCase,
		Validate: validate,
	}
}

func MPRequestHandlerFactory(log *logrus.Logger, viper *viper.Viper) IMPRequestHandler {
	useCase := usecase.MPRequestUseCaseFactory(viper, log)
	validate := config.NewValidator(viper)
	return NewMPRequestHandler(log, viper, useCase, validate)
}

func (h *MPRequestHandler) FindAllPaginated(ctx *gin.Context) {
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

	filter := make(map[string]interface{})

	status := ctx.Query("status")
	if status != "" {
		filter["status"] = status
	}

	departmentHead := ctx.Query("department_head")
	if departmentHead != "" {
		filter["department_head"] = departmentHead
	}

	vpGmDirector := ctx.Query("vp_gm_director")
	if vpGmDirector != "" {
		filter["vp_gm_director"] = vpGmDirector
	}

	ceo := ctx.Query("ceo")
	if ceo != "" {
		filter["ceo"] = ceo
	}

	hrdHoUnit := ctx.Query("hrd_ho_unit")
	if hrdHoUnit != "" {
		filter["hrd_ho_unit"] = hrdHoUnit
	}

	res, err := h.UseCase.FindAllPaginated(page, pageSize, search, filter)
	if err != nil {
		h.Log.Errorf("[MPRequestHandler.FindAllPaginated] error when find all paginated: %v", err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to find all paginated", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "MP Request Headers found", res)
}

func (h *MPRequestHandler) Create(ctx *gin.Context) {
	var req request.CreateMPRequestHeaderRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.Log.Errorf("[MPRequestHandler.Create] error when bind json: %v", err)
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	if err := h.Validate.Struct(req); err != nil {
		h.Log.Errorf("[MPRequestHandler.Create] error when validate request: %v", err)
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	res, err := h.UseCase.Create(&req)
	if err != nil {
		h.Log.Errorf("[MPRequestHandler.Create] error when create mp request header: %v", err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to create mp request header", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusCreated, "MP Request Header created", res)
}
