package repository

import (
	"github.com/IlhamSetiaji/julong-manpower-be/internal/config"
	"github.com/IlhamSetiaji/julong-manpower-be/internal/entity"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type IBatchRepository interface {
	CreateBatchHeaderAndLines(batchHeader *entity.BatchHeader, batchLines []entity.BatchLine) (*entity.BatchHeader, error)
	InsertLinesByBatchHeaderID(batchHeaderID string, batchLines []entity.BatchLine) error
	DeleteLinesNotInBatchLines(batchHeaderID string, batchLines []entity.BatchLine) error
	FindByStatus(status entity.BatchHeaderApprovalStatus) (*entity.BatchHeader, error)
	FindById(id string) (*entity.BatchHeader, error)
	GetHeadersByDocumentDate(documentDate string) ([]entity.BatchHeader, error)
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

func (r *BatchRepository) InsertLinesByBatchHeaderID(batchHeaderID string, batchLines []entity.BatchLine) error {
	tx := r.DB.Begin()
	for i := range batchLines {
		batchLines[i].BatchHeaderID = uuid.MustParse(batchHeaderID)
		if err := tx.Create(&batchLines[i]).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (r *BatchRepository) DeleteLinesNotInBatchLines(batchHeaderID string, batchLines []entity.BatchLine) error {
	tx := r.DB.Begin()
	var batchLineIDs []uuid.UUID
	for _, bl := range batchLines {
		batchLineIDs = append(batchLineIDs, bl.ID)
	}

	if err := tx.Where("batch_header_id = ? AND id NOT IN ?", batchHeaderID, batchLineIDs).Delete(&entity.BatchLine{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}

	return nil
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
	if err := r.DB.Preload("BatchLines.MPPlanningHeader.MPPPeriod").Preload("BatchLines.MPPlanningHeader.MPPlanningLines").Where("id = ?", id).First(&batchHeader).Error; err != nil {
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

func (r *BatchRepository) GetHeadersByDocumentDate(documentDate string) ([]entity.BatchHeader, error) {
	var batchHeaders []entity.BatchHeader
	if err := r.DB.Where("document_date = ?", documentDate).Find(&batchHeaders).Error; err != nil {
		r.Log.Error(err)
		return nil, err
	}

	return batchHeaders, nil
}

func BatchRepositoryFactory(log *logrus.Logger) IBatchRepository {
	db := config.NewDatabase()
	return NewBatchRepository(log, db)
}
