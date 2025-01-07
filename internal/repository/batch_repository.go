package repository

import (
	"errors"
	"time"

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
	FindByStatus(status entity.BatchHeaderApprovalStatus, approverType string, orgID string) (*entity.BatchHeader, error)
	FindById(id string) (*entity.BatchHeader, error)
	FindByNeedApproval(approverType string, orgID string) (*entity.BatchHeader, error)
	GetHeadersByDocumentDate(documentDate string) ([]entity.BatchHeader, error)
	FindByCurrentDocumentDateAndStatus(status entity.BatchHeaderApprovalStatus) (*entity.BatchHeader, error)
	UpdateStatusBatchHeader(batchHeader *entity.BatchHeader, status entity.BatchHeaderApprovalStatus, approvedBy string, approverName string) error
	GetBatchHeadersByStatus(status entity.BatchHeaderApprovalStatus) ([]entity.BatchHeader, error)
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

func (r *BatchRepository) GetBatchHeadersByStatus(status entity.BatchHeaderApprovalStatus) ([]entity.BatchHeader, error) {
	var batchHeaders []entity.BatchHeader
	if err := r.DB.Preload("BatchLines.MPPlanningHeader.MPPPeriod").Preload("BatchLines.MPPlanningHeader.MPPlanningLines").Where("status = ?", status).Find(&batchHeaders).Error; err != nil {
		r.Log.Error(err)
		return nil, err
	}

	return batchHeaders, nil
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

func (r *BatchRepository) FindByStatus(status entity.BatchHeaderApprovalStatus, approverType string, orgID string) (*entity.BatchHeader, error) {
	var batchHeader entity.BatchHeader
	var whereApproverType string
	var whereOrgID string
	if approverType == "" || approverType == "CEO" {
		whereApproverType = "approver_type = 'CEO'"
	} else {
		whereApproverType = "approver_type = 'DIRECTOR'"
	}

	if orgID != "" {
		if approverType != "" && approverType == "DIRECTOR" {
			whereOrgID = "organization_id = '" + orgID + "'"
		}
	}
	if err := r.DB.Preload("BatchLines.MPPlanningHeader.MPPPeriod").Preload("BatchLines.MPPlanningHeader.MPPlanningLines").Where("status = ?", status).Where(whereOrgID).Where(whereApproverType).First(&batchHeader).Error; err != nil {
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

func (r *BatchRepository) FindByNeedApproval(approverType string, orgID string) (*entity.BatchHeader, error) {
	var batchHeader entity.BatchHeader
	var whereApproverType string
	var whereOrgID string
	if approverType == "" || approverType == "CEO" {
		whereApproverType = "approver_type = 'CEO'"
	} else {
		whereApproverType = "approver_type = 'DIRECTOR'"
	}
	if orgID != "" {
		if approverType != "" && approverType == "DIRECTOR" {
			whereOrgID = "organization_id = '" + orgID + "'"
		}
	}
	if err := r.DB.Preload("BatchLines.MPPlanningHeader.MPPPeriod").Preload("BatchLines.MPPlanningHeader.MPPlanningLines").Where("status = ?", entity.BatchHeaderApprovalStatusNeedApproval).Where(whereApproverType).Where(whereOrgID).First(&batchHeader).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			r.Log.Warnf("Batch header with status %s not found", entity.BatchHeaderApprovalStatusNeedApproval)
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

func (r *BatchRepository) FindByCurrentDocumentDateAndStatus(status entity.BatchHeaderApprovalStatus) (*entity.BatchHeader, error) {
	var batchHeader entity.BatchHeader
	dateNow := time.Now()
	if err := r.DB.Preload("BatchLines.MPPlanningHeader.MPPPeriod").Preload("BatchLines.MPPlanningHeader.MPPlanningLines").Where("status = ? AND document_date = ?", status, dateNow.Format("2006-01-02")).First(&batchHeader).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			r.Log.Warnf("Batch header with status %s and document date %s not found", status, dateNow.Format("2006-01-02"))
			return nil, nil
		} else {
			r.Log.Error(err)
			return nil, err
		}
	}

	return &batchHeader, nil
}

func (r *BatchRepository) UpdateStatusBatchHeader(batchHeader *entity.BatchHeader, status entity.BatchHeaderApprovalStatus, approvedBy string, approverName string) error {
	tx := r.DB.Begin()

	var approvedByPtr uuid.UUID

	if approvedBy != "" {
		approvedByPtr = uuid.MustParse(approvedBy)
	}

	// loop through the batch lines and update the status
	for _, bl := range batchHeader.BatchLines {
		var mppPeriodCompleted entity.MPPPeriod
		if err := r.DB.Where("status = ?", entity.MPPeriodStatusComplete).First(&mppPeriodCompleted).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				r.Log.Warnf("MPP Period with status %s not found", entity.MPPeriodStatusComplete)
			} else {
				r.Log.Error(err)
				return err
			}
		}

		if status != entity.BatchHeaderApprovalStatusCompleted {
			var approvalHistory *entity.MPPlanningApprovalHistory

			approvalHistory = &entity.MPPlanningApprovalHistory{
				MPPlanningHeaderID: bl.MPPlanningHeaderID,
				ApproverID:         uuid.MustParse(approvedBy),
				ApproverName:       approverName,
				Notes:              "",
				Level:              string(entity.MPPlanningApprovalHistoryLevelCEO),
				Status:             entity.MPPlanningApprovalHistoryStatus(status),
			}

			if approvalHistory != nil {
				if approvalHistory.Status != entity.MPPlanningApprovalHistoryStatusRejected {
					if approvalHistory.Level == string(entity.MPPlanningApprovalHistoryLevelHRDUnit) {
						if err := tx.Model(&entity.MPPlanningHeader{}).Where("id = ?", bl.MPPlanningHeaderID).Updates(&entity.MPPlanningHeader{
							Status: entity.MPPlaningStatus(status),
							// ApprovedBy:        approvedBy,
							ApproverManagerID: &approvalHistory.ApproverID,
							NotesManager:      approvalHistory.Notes,
						}).Error; err != nil {
							tx.Rollback()
							r.Log.Errorf("[BatchRepository.UpdateStatusBatchHeader] " + err.Error())
							return errors.New("[BatchRepository.UpdateStatusBatchHeader] " + err.Error())
						}
					} else if approvalHistory.Level == string(entity.MPPlanningApprovalHistoryLevelDirekturUnit) {
						if err := tx.Model(&entity.MPPlanningHeader{}).Where("id = ?", bl.MPPlanningHeaderID).Updates(&entity.MPPlanningHeader{
							Status: entity.MPPlaningStatus(status),
							// ApprovedBy:            approvedBy,
							ApproverRecruitmentID: &approvalHistory.ApproverID,
							NotesRecruitment:      approvalHistory.Notes,
						}).Error; err != nil {
							tx.Rollback()
							r.Log.Errorf("[BatchRepository.UpdateStatusBatchHeader] " + err.Error())
							return errors.New("[BatchRepository.UpdateStatusBatchHeader] " + err.Error())
						}
					} else if approvalHistory.Level == string(entity.MPPlanningApprovalHistoryLevelCEO) {
						if err := tx.Model(&entity.MPPlanningHeader{}).Where("id = ?", bl.MPPlanningHeaderID).Updates(&entity.MPPlanningHeader{
							Status:     entity.MPPlaningStatus(status),
							ApprovedBy: approvedBy,
						}).Error; err != nil {
							tx.Rollback()
							r.Log.Errorf("[BatchRepository.UpdateStatusBatchHeader] " + err.Error())
							return errors.New("[BatchRepository.UpdateStatusBatchHeader] " + err.Error())
						}
					} else {
						if err := tx.Model(&entity.MPPlanningHeader{}).Where("id = ?", bl.MPPlanningHeaderID).Updates(&entity.MPPlanningHeader{
							Status: entity.MPPlaningStatus(status),
						}).Error; err != nil {
							tx.Rollback()
							r.Log.Errorf("[BatchRepository.UpdateStatusBatchHeader] " + err.Error())
							return errors.New("[BatchRepository.UpdateStatusBatchHeader] " + err.Error())
						}
					}
				} else {
					if err := tx.Model(&entity.MPPlanningHeader{}).Where("id = ?", bl.MPPlanningHeaderID).Updates(&entity.MPPlanningHeader{
						Status: entity.MPPlaningStatus(status),
					}).Error; err != nil {
						tx.Rollback()
						r.Log.Errorf("[BatchRepository.UpdateStatusBatchHeader] " + err.Error())
						return errors.New("[BatchRepository.UpdateStatusBatchHeader] " + err.Error())
					}
				}
			} else {
				if err := tx.Model(&entity.MPPlanningHeader{}).Where("id = ?", bl.MPPlanningHeaderID).Updates(&entity.MPPlanningHeader{
					Status: entity.MPPlaningStatus(status),
				}).Error; err != nil {
					tx.Rollback()
					r.Log.Errorf("[BatchRepository.UpdateStatusBatchHeader] " + err.Error())
					return errors.New("[BatchRepository.UpdateStatusBatchHeader] " + err.Error())
				}
			}

			if approvalHistory != nil {
				if err := tx.Create(approvalHistory).Error; err != nil {
					tx.Rollback()
					r.Log.Errorf("[BatchRepository.UpdateStatusBatchHeader] " + err.Error())
					return errors.New("[BatchRepository.UpdateStatusBatchHeader] " + err.Error())
				}
			}

			if err := tx.Model(&entity.BatchHeader{}).Where("id = ?", batchHeader.ID).Updates(&entity.BatchHeader{
				Status:       status,
				ApproverID:   &approvedByPtr,
				ApproverName: approverName,
			}).Error; err != nil {
				tx.Rollback()
				r.Log.Errorf("[BatchRepository.UpdateStatusBatchHeader] " + err.Error())
				return err
			}
		} else {
			if mppPeriodCompleted.ID != uuid.Nil {
				r.Log.Warnf("MPP Period with status %s found: %v", entity.MPPeriodStatusComplete, mppPeriodCompleted)
				return errors.New("MPP Period with status complete found")
			}
			if err := tx.Model(&entity.MPPlanningHeader{}).Where("id = ?", bl.MPPlanningHeaderID).Updates(&entity.MPPlanningHeader{
				Status: entity.MPPlaningStatus(status),
			}).Error; err != nil {
				tx.Rollback()
				r.Log.Errorf("[BatchRepository.UpdateStatusBatchHeader] " + err.Error())
				return errors.New("[BatchRepository.UpdateStatusBatchHeader] " + err.Error())
			}

			mppPeriod := bl.MPPlanningHeader.MPPPeriod
			if err := tx.Model(&entity.MPPPeriod{}).Where("id = ?", mppPeriod.ID).Updates(&entity.MPPPeriod{
				Status: entity.MPPeriodStatusComplete,
			}).Error; err != nil {
				tx.Rollback()
				r.Log.Errorf("[BatchRepository.UpdateStatusBatchHeader] " + err.Error())
				return errors.New("[BatchRepository.UpdateStatusBatchHeader] " + err.Error())
			}

			if err := tx.Model(&entity.BatchHeader{}).Where("id = ?", batchHeader.ID).Updates(&entity.BatchHeader{
				Status: status,
				// ApproverID:   &approvedByPtr,
				// ApproverName: approverName,
			}).Error; err != nil {
				tx.Rollback()
				r.Log.Errorf("[BatchRepository.UpdateStatusBatchHeader] " + err.Error())
				return err
			}

		}
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		r.Log.Errorf("[BatchRepository.UpdateStatusBatchHeader] " + err.Error())
		return err
	}

	return nil
}

func BatchRepositoryFactory(log *logrus.Logger) IBatchRepository {
	db := config.NewDatabase()
	return NewBatchRepository(log, db)
}
