package repository

import (
	"github.com/IlhamSetiaji/julong-manpower-be/internal/config"
	"github.com/IlhamSetiaji/julong-manpower-be/internal/entity"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type IBatchRepository interface {
	CreateBatchHeaderAndLines(batchHeader *entity.BatchHeader, batchLines []entity.BatchLine) (*entity.BatchHeader, error)
	FindByStatus(status entity.BatchHeaderApprovalStatus) (*entity.BatchHeader, error)
	FindById(id string) (*entity.BatchHeader, error)
}

type BatchRepository struct {
	Log *logrus.Logger
	DB  *gorm.DB
}

func NewBatchRepository(log *logrus.Logger, db *gorm.DB) IBatchRepository {
	return &BatchRepository{
		Log: log,
		DB:  db,
	}
}

func (r *BatchRepository) CreateBatchHeaderAndLines(batchHeader *entity.BatchHeader, batchLines []entity.BatchLine) (*entity.BatchHeader, error) {
	tx := r.DB.Begin()
	if batchHeader.Status == "" {
		batchHeader.Status = entity.BatchHeaderApprovalStatusNeedApproval
	}
	if err := tx.Create(batchHeader).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	for i := range batchLines {
		batchLines[i].BatchHeaderID = batchHeader.ID
		if err := tx.Create(&batchLines[i]).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Fetch the batch header with the batch lines
	if err := r.DB.Preload("BatchLines.MPPlanningHeader.MPPPeriod").Preload("BatchLines.MPPlanningHeader.MPPlanningLines").First(batchHeader).Error; err != nil {
		return nil, err
	}

	return batchHeader, nil
}

func (r *BatchRepository) FindByStatus(status entity.BatchHeaderApprovalStatus) (*entity.BatchHeader, error) {
	var batchHeader entity.BatchHeader
	if err := r.DB.Where("status = ?", status).First(&batchHeader).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			r.Log.Warnf("Batch header with status %s not found", status)
			return nil, nil
		} else {
			r.Log.Error(err)
			return nil, err
		}
	}

	return &batchHeader, nil
}

func (r *BatchRepository) FindById(id string) (*entity.BatchHeader, error) {
	var batchHeader entity.BatchHeader
	if err := r.DB.Where("id = ?", id).First(&batchHeader).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			r.Log.Warnf("Batch header with id %s not found", id)
			return nil, nil
		} else {
			r.Log.Error(err)
			return nil, err
		}
	}

	return &batchHeader, nil
}

func BatchRepositoryFactory(log *logrus.Logger) IBatchRepository {
	db := config.NewDatabase()
	return NewBatchRepository(log, db)
}
