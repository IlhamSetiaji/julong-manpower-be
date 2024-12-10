package repository

import (
	"errors"

	"github.com/IlhamSetiaji/julong-manpower-be/internal/config"
	"github.com/IlhamSetiaji/julong-manpower-be/internal/entity"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type IMPRequestRepository interface {
	Create(mpRequestHeader *entity.MPRequestHeader) (*entity.MPRequestHeader, error)
}

type MPRequestRepository struct {
	Log *logrus.Logger
	DB  *gorm.DB
}

func NewMPRequestRepository(log *logrus.Logger, db *gorm.DB) IMPRequestRepository {
	return &MPRequestRepository{Log: log, DB: db}
}

func (r *MPRequestRepository) Create(mpRequestHeader *entity.MPRequestHeader) (*entity.MPRequestHeader, error) {
	tx := r.DB.Begin()

	if err := tx.Create(mpRequestHeader).Error; err != nil {
		tx.Rollback()
		r.Log.Errorf("[MPRequestRepository.Create] error when create mp request header: %v", err)
		return nil, errors.New("[MPRequestRepository.Create] error when create mp request header " + err.Error())
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		r.Log.Errorf("[MPRequestRepository.Create] error when commit transaction: %v", err)
		return nil, errors.New("[MPRequestRepository.Create] error when commit transaction " + err.Error())
	}

	return mpRequestHeader, nil
}

func MPRequestRepositoryFactory(log *logrus.Logger) IMPRequestRepository {
	db := config.NewDatabase()
	return NewMPRequestRepository(log, db)
}
