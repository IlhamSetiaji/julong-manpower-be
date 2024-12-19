package repository

import (
	"errors"

	"github.com/IlhamSetiaji/julong-manpower-be/internal/config"
	"github.com/IlhamSetiaji/julong-manpower-be/internal/entity"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type IRequestMajorRepository interface {
	Create(requestMajor *entity.RequestMajor) (*entity.RequestMajor, error)
	Delete(id uuid.UUID) error
	DeleteByMPRequestHeaderID(mpRequestHeaderID uuid.UUID) error
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

	if err := r.DB.Preload("Major").First(requestMajor, requestMajor.ID).Error; err != nil {
		r.Log.Errorf("[RequestMajorRepository.Create] error when preloading associations: %v", err)
		return nil, errors.New("[RequestMajorRepository.Create] error when preloading associations " + err.Error())
	}

	return requestMajor, nil
}

func (r *RequestMajorRepository) Delete(id uuid.UUID) error {
	tx := r.DB.Begin()

	if err := tx.Where("id = ?", id).Delete(&entity.RequestMajor{}).Error; err != nil {
		tx.Rollback()
		r.Log.Errorf("[RequestMajorRepository.Delete] error when delete request major: %v", err)
		return err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		r.Log.Errorf("[RequestMajorRepository.Delete] error when commit transaction: %v", err)
		return err
	}

	return nil
}

func (r *RequestMajorRepository) DeleteByMPRequestHeaderID(mpRequestHeaderID uuid.UUID) error {
	tx := r.DB.Begin()

	if err := tx.Where("mp_request_header_id = ?", mpRequestHeaderID).Delete(&entity.RequestMajor{}).Error; err != nil {
		tx.Rollback()
		r.Log.Errorf("[RequestMajorRepository.DeleteByMPRequestHeaderID] error when delete request major by mp request header id: %v", err)
		return err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		r.Log.Errorf("[RequestMajorRepository.DeleteByMPRequestHeaderID] error when commit transaction: %v", err)
		return err
	}

	return nil
}

func RequestMajorRepositoryFactory(log *logrus.Logger) IRequestMajorRepository {
	db := config.NewDatabase()
	return NewRequestMajorRepository(log, db)
}
