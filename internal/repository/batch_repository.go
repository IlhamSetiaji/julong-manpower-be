package repository

import (
	"github.com/IlhamSetiaji/julong-manpower-be/internal/config"
	"github.com/IlhamSetiaji/julong-manpower-be/internal/entity"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type IBatchRepository interface {
	CreateBatchHeaderAndLines(batchHeader *entity.BatchHeader, batchLines []entity.BatchLine) (*entity.BatchHeader, error)
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

	tx.Commit()
	return batchHeader, nil
}

func BatchRepositoryFactory(log *logrus.Logger) IBatchRepository {
	db := config.NewDatabase()
	return NewBatchRepository(log, db)
}
