package repository

import (
	"errors"

	"github.com/IlhamSetiaji/julong-manpower-be/internal/config"
	"github.com/IlhamSetiaji/julong-manpower-be/internal/entity"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type IRequestCategoryRepository interface {
	FindAll() (*[]entity.RequestCategory, error)
	FindById(id uuid.UUID) (*entity.RequestCategory, error)
}

type RequestCategoryRepository struct {
	Log *logrus.Logger
	DB  *gorm.DB
}

func NewRequestCategoryRepository(log *logrus.Logger, db *gorm.DB) IRequestCategoryRepository {
	return &RequestCategoryRepository{Log: log, DB: db}
}

func (r *RequestCategoryRepository) FindAll() (*[]entity.RequestCategory, error) {
	var requestCategories []entity.RequestCategory
	if err := r.DB.Find(&requestCategories).Error; err != nil {
		r.Log.Errorf("[RequestCategoryRepository.FindAll] %s", err.Error())
		return nil, errors.New("[RequestCategoryRepository.FindAll] Internal server error")
	}
	return &requestCategories, nil
}

func (r *RequestCategoryRepository) FindById(id uuid.UUID) (*entity.RequestCategory, error) {
	var requestCategory entity.RequestCategory
	if err := r.DB.First(&requestCategory, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			r.Log.Errorf("[MPPlanningRepository.FindHeaderById] Request Category with ID %s not found", id)
			return nil, errors.New("[MPPlanningRepository.FindHeaderById] Request Category not found")
		} else {
			r.Log.Errorf("[MPPlanningRepository.FindHeaderById] %s", err.Error())
			return nil, errors.New("[MPPlanningRepository.FindHeaderById] Internal server error")
		}
	}
	return &requestCategory, nil
}

func RequestCategoryRepositoryFactory(log *logrus.Logger) IRequestCategoryRepository {
	db := config.NewDatabase()
	return NewRequestCategoryRepository(log, db)
}
