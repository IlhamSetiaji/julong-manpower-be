package repository

import (
	"github.com/IlhamSetiaji/julong-manpower-be/internal/config"
	"github.com/IlhamSetiaji/julong-manpower-be/internal/entity"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type IRequestMajorRepository interface {
	Create(requestMajor *entity.RequestMajor) (*entity.RequestMajor, error)
}

type RequestMajorRepository struct {
	Log *logrus.Logger
	DB  *gorm.DB
}

func NewRequestMajorRepository(log *logrus.Logger, db *gorm.DB) IRequestMajorRepository {
	return &RequestMajorRepository{Log: log, DB: db}
}

func (r *RequestMajorRepository) Create(requestMajor *entity.RequestMajor) (*entity.RequestMajor, error) {
	tx := r.DB.Begin()

	if err := tx.Create(requestMajor).Error; err != nil {
		tx.Rollback()
		r.Log.Errorf("[RequestMajorRepository.Create] error when create request major: %v", err)
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		r.Log.Errorf("[RequestMajorRepository.Create] error when commit transaction: %v", err)
		return nil, err
	}

	return requestMajor, nil
}

func RequestMajorFactory(log *logrus.Logger) IRequestMajorRepository {
	db := config.NewDatabase()
	return NewRequestMajorRepository(log, db)
}
