package usecase

import (
	"github.com/IlhamSetiaji/julong-manpower-be/internal/http/request"
	"github.com/IlhamSetiaji/julong-manpower-be/internal/http/response"
	"github.com/IlhamSetiaji/julong-manpower-be/internal/repository"
	"github.com/sirupsen/logrus"
)

type IRequestCategoryUseCase interface {
	FindAll() (*[]response.RequestCategoryResponse, error)
	FindById(request *request.FindByIdRequestCategoryRequest) (*response.RequestCategoryResponse, error)
}

type RequestCategoryUseCase struct {
	Log  *logrus.Logger
	Repo repository.IRequestCategoryRepository
}

func NewRequestCategoryUseCase(log *logrus.Logger, repo repository.IRequestCategoryRepository) IRequestCategoryUseCase {
	return &RequestCategoryUseCase{Log: log, Repo: repo}
}

func (u *RequestCategoryUseCase) FindAll() (*[]response.RequestCategoryResponse, error) {
	requestCategories, err := u.Repo.FindAll()
	if err != nil {
		return nil, err
	}

	var responseCategories []response.RequestCategoryResponse
	for _, requestCategory := range *requestCategories {
		responseCategories = append(responseCategories, response.RequestCategoryResponse{
			ID:            requestCategory.ID,
			Name:          requestCategory.Name,
			IsReplacement: requestCategory.IsReplacement,
			CreatedAt:     requestCategory.CreatedAt,
			UpdatedAt:     requestCategory.UpdatedAt,
		})
	}

	return &responseCategories, nil
}

func (u *RequestCategoryUseCase) FindById(req *request.FindByIdRequestCategoryRequest) (*response.RequestCategoryResponse, error) {
	requestCategory, err := u.Repo.FindById(req.ID)
	if err != nil {
		return nil, err
	}

	responseCategory := response.RequestCategoryResponse{
		ID:            requestCategory.ID,
		Name:          requestCategory.Name,
		IsReplacement: requestCategory.IsReplacement,
		CreatedAt:     requestCategory.CreatedAt,
		UpdatedAt:     requestCategory.UpdatedAt,
	}

	return &responseCategory, nil
}

func RequestCategoryUseCaseFactory(log *logrus.Logger) IRequestCategoryUseCase {
	repo := repository.RequestCategoryRepositoryFactory(log)
	return NewRequestCategoryUseCase(log, repo)
}
