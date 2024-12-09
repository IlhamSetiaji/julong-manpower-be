package handler

import (
	"net/http"

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

type IRequestCategoryHandler interface {
	FindAll(ctx *gin.Context)
	FindById(ctx *gin.Context)
}

type RequestCategoryHandler struct {
	Log      *logrus.Logger
	Viper    *viper.Viper
	Uc       usecase.IRequestCategoryUseCase
	Validate *validator.Validate
}

func NewRequestCategoryHandler(log *logrus.Logger, viper *viper.Viper, uc usecase.IRequestCategoryUseCase, validate *validator.Validate) IRequestCategoryHandler {
	return &RequestCategoryHandler{Log: log, Viper: viper, Uc: uc, Validate: validate}
}

func (h *RequestCategoryHandler) FindAll(ctx *gin.Context) {
	requestCategories, err := h.Uc.FindAll()
	if err != nil {
		h.Log.Errorf("[RequestCategoryHandler.FindAll] error: %s", err.Error())
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Internal Server Error", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Find All Request Categories Success", requestCategories)
}

func (h *RequestCategoryHandler) FindById(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Bad Request", "ID is required")
		return
	}

	req := request.FindByIdRequestCategoryRequest{ID: uuid.MustParse(id)}
	if err := h.Validate.Struct(req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Bad Request", err.Error())
		return
	}

	requestCategory, err := h.Uc.FindById(&req)
	if err != nil {
		h.Log.Errorf("[RequestCategoryHandler.FindById] error: %s", err.Error())
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Internal Server Error", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Find Request Category By ID Success", requestCategory)
}

func RequestCategoryHandlerFactory(log *logrus.Logger, viper *viper.Viper) IRequestCategoryHandler {
	uc := usecase.RequestCategoryUseCaseFactory(log)
	validate := config.NewValidator(viper)
	return NewRequestCategoryHandler(log, viper, uc, validate)
}
