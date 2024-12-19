package repository

import (
	"errors"

	"github.com/IlhamSetiaji/julong-manpower-be/internal/config"
	"github.com/IlhamSetiaji/julong-manpower-be/internal/entity"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type IMPPlanningRepository interface {
	FindAllHeadersPaginated(page int, pageSize int, search string, approverType string, orgLocationId string, orgId string, status entity.MPPlaningStatus) (*[]entity.MPPlanningHeader, int64, error)
	FindAllHeadersByRequestorIDPaginated(requestorID uuid.UUID, page int, pageSize int, search string) (*[]entity.MPPlanningHeader, int64, error)
	FindAllHeaders() (*[]entity.MPPlanningHeader, error)
	FindHeaderById(id uuid.UUID) (*entity.MPPlanningHeader, error)
	FindHeaderByOrganizationLocationID(organizationLocationID uuid.UUID) (*entity.MPPlanningHeader, error)
	FindHeaderByOrganizationLocationIDAndStatus(organizationLocationID uuid.UUID, status entity.MPPlaningStatus) (*entity.MPPlanningHeader, error)
	GetHeadersByStatus(status entity.MPPlaningStatus) (*[]entity.MPPlanningHeader, error)
	UpdateStatusHeader(id uuid.UUID, status string, approvedBy string, approvalHistory *entity.MPPlanningApprovalHistory) error
	GetHeadersByDocumentDate(documentDate string) (*[]entity.MPPlanningHeader, error)
	CreateHeader(mppHeader *entity.MPPlanningHeader) (*entity.MPPlanningHeader, error)
	UpdateHeader(mppHeader *entity.MPPlanningHeader) (*entity.MPPlanningHeader, error)
	StoreAttachmentToHeader(mppHeader *entity.MPPlanningHeader, attachment entity.ManpowerAttachment) (*entity.MPPlanningHeader, error)
	StoreAttachmentToApprovalHistory(mppApprovalHistory *entity.MPPlanningApprovalHistory, attachment entity.ManpowerAttachment) (*entity.MPPlanningApprovalHistory, error)
	StoreAttachmentsToApprovalHistory(mppApprovalHistories *entity.MPPlanningApprovalHistory, attachments []entity.ManpowerAttachment) (*entity.MPPlanningApprovalHistory, error)
	GetPlanningApprovalHistoryByHeaderId(headerId uuid.UUID) (*[]entity.MPPlanningApprovalHistory, error)
	GetPlanningApprovalHistoryAttachmentsByApprovalHistoryId(approvalHistoryId uuid.UUID) (*[]entity.ManpowerAttachment, error)
	DeleteAttachmentFromHeader(mppHeader *entity.MPPlanningHeader, attachmentID uuid.UUID) (*entity.MPPlanningHeader, error)
	DeleteHeader(id uuid.UUID) error
	FindHeaderByMPPPeriodId(mppPeriodId uuid.UUID) (*entity.MPPlanningHeader, error)
	FindAllLinesByHeaderIdPaginated(headerId uuid.UUID, page int, pageSize int, search string) (*[]entity.MPPlanningLine, int64, error)
	FindLineById(id uuid.UUID) (*entity.MPPlanningLine, error)
	CreateLine(mppLine *entity.MPPlanningLine) (*entity.MPPlanningLine, error)
	UpdateLine(mppLine *entity.MPPlanningLine) (*entity.MPPlanningLine, error)
	DeleteLine(id uuid.UUID) error
	DeleteLineIfNotInIDs(ids []uuid.UUID) error
	DeleteAllLinesByHeaderID(headerID uuid.UUID) error
	DeleteLinesByIDs(ids []uuid.UUID) error
}

type MPPlanningRepository struct {
	Log *logrus.Logger
	DB  *gorm.DB
}

func NewMPPlanningRepository(log *logrus.Logger, db *gorm.DB) IMPPlanningRepository {
	return &MPPlanningRepository{
		Log: log,
		DB:  db,
	}
}

func (r *MPPlanningRepository) FindAllHeaders() (*[]entity.MPPlanningHeader, error) {
	var mppHeaders []entity.MPPlanningHeader

	if err := r.DB.Preload("MPPlanningLines").Preload("MPPlanningLines.MPRequestLine").Preload("MPPPeriod").Find(&mppHeaders).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			r.Log.Errorf("[MPPlanningRepository.FindAllHeaders] " + err.Error())
			return nil, nil
		} else {
			r.Log.Errorf("[MPPlanningRepository.FindAllHeaders] " + err.Error())
			return nil, errors.New("[MPPlanningRepository.FindAllHeaders] " + err.Error())
		}
	}

	return &mppHeaders, nil
}

func (r *MPPlanningRepository) FindAllHeadersPaginated(page int, pageSize int, search string, approverType string, orgLocationId string, orgId string, status entity.MPPlaningStatus) (*[]entity.MPPlanningHeader, int64, error) {
	var mppHeaders []entity.MPPlanningHeader
	var total int64

	query := r.DB.Model(&entity.MPPlanningHeader{}).Preload("MPPlanningLines").Preload("MPPPeriod")

	if search != "" {
		query = query.Where("document_number LIKE ?", "%"+search+"%")
	}

	if approverType != "" {
		if approverType == "manager" {
			r.Log.Info("Approver Type: Manager")
			query = query.Where("approver_manager_id IS NOT NULL")
		} else {
			query = query.Where("approver_recruitment_id IS NOT NULL")
		}
	}

	if orgLocationId != "" {
		query = query.Where("organization_location_id = ?", orgLocationId)
	}

	if orgId != "" {
		query = query.Where("organization_id = ?", orgId)
	}

	if status != "" {
		query = query.Where("status = ?", status)
	}

	countQuery := query.Session(&gorm.Session{})
	if err := countQuery.Count(&total).Error; err != nil {
		r.Log.Errorf("[MPPlanningRepository.FindAllHeadersPaginated - count  side] " + err.Error())
		return nil, 0, errors.New("[MPPlanningRepository.FindAllHeadersPaginated - count side] " + err.Error())
	}

	if err := query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&mppHeaders).Error; err != nil {
		r.Log.Errorf("[MPPlanningRepository.FindAllHeadersPaginated - pagination side] " + err.Error())
		return nil, 0, errors.New("[MPPlanningRepository.FindAllHeadersPaginated - pagination side] " + err.Error())
	}

	return &mppHeaders, total, nil
}

func (r *MPPlanningRepository) FindAllHeadersByRequestorIDPaginated(requestorID uuid.UUID, page int, pageSize int, search string) (*[]entity.MPPlanningHeader, int64, error) {
	var mppHeaders []entity.MPPlanningHeader

	query := r.DB.Model(&entity.MPPlanningHeader{}).Preload("MPPlanningLines").Preload("MPPPeriod").Where("requestor_id = ?", requestorID)

	if search != "" {
		query = query.Where("name LIKE ?", "%"+search+"%")
	}

	if err := query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&mppHeaders).Error; err != nil {
		r.Log.Errorf("[MPPlanningRepository.FindAllHeadersByRequestorIDPaginated] " + err.Error())
		return nil, 0, errors.New("[MPPlanningRepository.FindAllHeadersByRequestorIDPaginated] " + err.Error())
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		r.Log.Errorf("[MPPlanningRepository.FindAllHeadersByRequestorIDPaginated] " + err.Error())
		return nil, 0, errors.New("[MPPlanningRepository.FindAllHeadersByRequestorIDPaginated] " + err.Error())
	}

	return &mppHeaders, total, nil
}

func (r *MPPlanningRepository) FindHeaderById(id uuid.UUID) (*entity.MPPlanningHeader, error) {
	var mppHeader entity.MPPlanningHeader

	if err := r.DB.Preload("MPPlanningLines").Preload("MPPPeriod").Preload("ManpowerAttachments").Where("id = ?", id).First(&mppHeader).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			r.Log.Errorf("[MPPlanningRepository.FindHeaderById] " + err.Error())
			return nil, nil
		} else {
			r.Log.Errorf("[MPPlanningRepository.FindHeaderById] " + err.Error())
			return nil, errors.New("[MPPlanningRepository.FindHeaderById] " + err.Error())
		}
	}

	return &mppHeader, nil
}

func (r *MPPlanningRepository) FindHeaderByOrganizationLocationID(organizationLocationID uuid.UUID) (*entity.MPPlanningHeader, error) {
	var mppHeader entity.MPPlanningHeader
	statuses := []entity.MPPlaningStatus{entity.MPPlaningStatusApproved, entity.MPPlaningStatusReject}

	if err := r.DB.Preload("MPPlanningLines").Preload("MPPPeriod").Preload("ManpowerAttachments").Where("organization_location_id = ?", organizationLocationID).Where("status IN ?", statuses).First(&mppHeader).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			r.Log.Errorf("[MPPlanningRepository.FindHeaderByOrganizationLocationID] " + err.Error())
			return nil, nil
		} else {
			r.Log.Errorf("[MPPlanningRepository.FindHeaderByOrganizationLocationID] " + err.Error())
			return nil, errors.New("[MPPlanningRepository.FindHeaderByOrganizationLocationID] " + err.Error())
		}
	}

	return &mppHeader, nil
}

func (r *MPPlanningRepository) FindHeaderByOrganizationLocationIDAndStatus(organizationLocationID uuid.UUID, status entity.MPPlaningStatus) (*entity.MPPlanningHeader, error) {
	var mppHeader entity.MPPlanningHeader

	if err := r.DB.Preload("MPPlanningLines").Preload("MPPPeriod").Preload("ManpowerAttachments").Where("organization_location_id = ?", organizationLocationID).Where("status = ?", status).First(&mppHeader).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			r.Log.Errorf("[MPPlanningRepository.FindHeaderByOrganizationLocationIDAndStatus] " + err.Error())
			return nil, nil
		} else {
			r.Log.Errorf("[MPPlanningRepository.FindHeaderByOrganizationLocationIDAndStatus] " + err.Error())
			return nil, errors.New("[MPPlanningRepository.FindHeaderByOrganizationLocationIDAndStatus] " + err.Error())
		}
	}

	return &mppHeader, nil
}

func (r *MPPlanningRepository) GetHeadersByStatus(status entity.MPPlaningStatus) (*[]entity.MPPlanningHeader, error) {
	var mppHeaders []entity.MPPlanningHeader

	if err := r.DB.Preload("MPPlanningLines").Preload("MPPPeriod").Preload("ManpowerAttachments").Where("status = ?", status).Find(&mppHeaders).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			r.Log.Errorf("[MPPlanningRepository.GetHeadersByStatus] " + err.Error())
			return nil, nil
		} else {
			r.Log.Errorf("[MPPlanningRepository.GetHeadersByStatus] " + err.Error())
			return nil, errors.New("[MPPlanningRepository.GetHeadersByStatus] " + err.Error())
		}
	}

	return &mppHeaders, nil
}

func (r *MPPlanningRepository) GetHeadersByDocumentDate(documentDate string) (*[]entity.MPPlanningHeader, error) {
	var mppHeaders []entity.MPPlanningHeader

	if err := r.DB.Preload("MPPlanningLines").Preload("MPPPeriod").Preload("ManpowerAttachments").Where("document_date = ?", documentDate).Find(&mppHeaders).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			r.Log.Errorf("[MPPlanningRepository.GetHeadersByDocumentDate] " + err.Error())
			return nil, nil
		} else {
			r.Log.Errorf("[MPPlanningRepository.GetHeadersByDocumentDate] " + err.Error())
			return nil, errors.New("[MPPlanningRepository.GetHeadersByDocumentDate] " + err.Error())
		}
	}

	return &mppHeaders, nil
}

func (r *MPPlanningRepository) UpdateStatusHeader(id uuid.UUID, status string, approvedBy string, approvalHistory *entity.MPPlanningApprovalHistory) error {
	tx := r.DB.Begin()

	if tx.Error != nil {
		r.Log.Errorf("[MPPlanningRepository.UpdateStatusHeader] " + tx.Error.Error())
		return errors.New("[MPPlanningRepository.UpdateStatusHeader] " + tx.Error.Error())
	}

	if approvalHistory != nil {
		if approvalHistory.Status != entity.MPPlanningApprovalHistoryStatusRejected {
			if approvalHistory.Level == string(entity.MPPlanningApprovalHistoryLevelHRDUnit) {
				if err := tx.Model(&entity.MPPlanningHeader{}).Where("id = ?", id).Updates(&entity.MPPlanningHeader{
					Status: entity.MPPlaningStatus(status),
					// ApprovedBy:        approvedBy,
					ApproverManagerID: &approvalHistory.ApproverID,
					NotesManager:      approvalHistory.Notes,
				}).Error; err != nil {
					tx.Rollback()
					r.Log.Errorf("[MPPlanningRepository.UpdateStatusHeader] " + err.Error())
					return errors.New("[MPPlanningRepository.UpdateStatusHeader] " + err.Error())
				}
			} else if approvalHistory.Level == string(entity.MPPlanningApprovalHistoryLevelDirekturUnit) {
				if err := tx.Model(&entity.MPPlanningHeader{}).Where("id = ?", id).Updates(&entity.MPPlanningHeader{
					Status: entity.MPPlaningStatus(status),
					// ApprovedBy:            approvedBy,
					ApproverRecruitmentID: &approvalHistory.ApproverID,
					NotesRecruitment:      approvalHistory.Notes,
				}).Error; err != nil {
					tx.Rollback()
					r.Log.Errorf("[MPPlanningRepository.UpdateStatusHeader] " + err.Error())
					return errors.New("[MPPlanningRepository.UpdateStatusHeader] " + err.Error())
				}
			} else if approvalHistory.Level == string(entity.MPPlanningApprovalHistoryLevelCEO) {
				if err := tx.Model(&entity.MPPlanningHeader{}).Where("id = ?", id).Updates(&entity.MPPlanningHeader{
					Status:     entity.MPPlaningStatus(status),
					ApprovedBy: approvedBy,
				}).Error; err != nil {
					tx.Rollback()
					r.Log.Errorf("[MPPlanningRepository.UpdateStatusHeader] " + err.Error())
					return errors.New("[MPPlanningRepository.UpdateStatusHeader] " + err.Error())
				}
			} else {
				if err := tx.Model(&entity.MPPlanningHeader{}).Where("id = ?", id).Updates(&entity.MPPlanningHeader{
					Status: entity.MPPlaningStatus(status),
				}).Error; err != nil {
					tx.Rollback()
					r.Log.Errorf("[MPPlanningRepository.UpdateStatusHeader] " + err.Error())
					return errors.New("[MPPlanningRepository.UpdateStatusHeader] " + err.Error())
				}
			}
		} else {
			if err := tx.Model(&entity.MPPlanningHeader{}).Where("id = ?", id).Updates(&entity.MPPlanningHeader{
				Status: entity.MPPlaningStatus(status),
			}).Error; err != nil {
				tx.Rollback()
				r.Log.Errorf("[MPPlanningRepository.UpdateStatusHeader] " + err.Error())
				return errors.New("[MPPlanningRepository.UpdateStatusHeader] " + err.Error())
			}
		}
	} else {
		if err := tx.Model(&entity.MPPlanningHeader{}).Where("id = ?", id).Updates(&entity.MPPlanningHeader{
			Status: entity.MPPlaningStatus(status),
		}).Error; err != nil {
			tx.Rollback()
			r.Log.Errorf("[MPPlanningRepository.UpdateStatusHeader] " + err.Error())
			return errors.New("[MPPlanningRepository.UpdateStatusHeader] " + err.Error())
		}
	}

	if approvalHistory != nil {
		if err := tx.Create(approvalHistory).Error; err != nil {
			tx.Rollback()
			r.Log.Errorf("[MPPlanningRepository.UpdateStatusHeader] " + err.Error())
			return errors.New("[MPPlanningRepository.UpdateStatusHeader] " + err.Error())
		}
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		r.Log.Errorf("[MPPlanningRepository.UpdateStatusHeader] " + err.Error())
		return errors.New("[MPPlanningRepository.UpdateStatusHeader] " + err.Error())
	}

	return nil
}

func (r *MPPlanningRepository) CreateHeader(mppHeader *entity.MPPlanningHeader) (*entity.MPPlanningHeader, error) {
	tx := r.DB.Begin()

	if tx.Error != nil {
		r.Log.Errorf("[MPPlanningRepository.CreateHeader] " + tx.Error.Error())
		return nil, errors.New("[MPPlanningRepository.CreateHeader] " + tx.Error.Error())
	}

	if err := tx.Create(mppHeader).Error; err != nil {
		tx.Rollback()
		r.Log.Errorf("[MPPlanningRepository.CreateHeader] " + err.Error())
		return nil, errors.New("[MPPlanningRepository.CreateHeader] " + err.Error())
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		r.Log.Errorf("[MPPlanningRepository.CreateHeader] " + err.Error())
		return nil, errors.New("[MPPlanningRepository.CreateHeader] " + err.Error())
	}

	return mppHeader, nil
}

func (r *MPPlanningRepository) UpdateHeader(mppHeader *entity.MPPlanningHeader) (*entity.MPPlanningHeader, error) {
	tx := r.DB.Begin()

	if tx.Error != nil {
		r.Log.Errorf("[MPPlanningRepository.UpdateHeader] " + tx.Error.Error())
		return nil, errors.New("[MPPlanningRepository.UpdateHeader] " + tx.Error.Error())
	}

	if err := tx.Save(mppHeader).Error; err != nil {
		tx.Rollback()
		r.Log.Errorf("[MPPlanningRepository.UpdateHeader] " + err.Error())
		return nil, errors.New("[MPPlanningRepository.UpdateHeader] " + err.Error())
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		r.Log.Errorf("[MPPlanningRepository.UpdateHeader] " + err.Error())
		return nil, errors.New("[MPPlanningRepository.UpdateHeader] " + err.Error())
	}

	return mppHeader, nil
}

func (r *MPPlanningRepository) StoreAttachmentToHeader(mppHeader *entity.MPPlanningHeader, attachment entity.ManpowerAttachment) (*entity.MPPlanningHeader, error) {
	tx := r.DB.Begin()

	if tx.Error != nil {
		r.Log.Errorf("[MPPlanningRepository.StoreAttachmentToHeader] " + tx.Error.Error())
		return nil, errors.New("[MPPlanningRepository.StoreAttachmentToHeader] " + tx.Error.Error())
	}

	if err := tx.Model(mppHeader).Association("ManpowerAttachments").Append(&attachment); err != nil {
		tx.Rollback()
		r.Log.Errorf("[MPPlanningRepository.StoreAttachmentToHeader] " + err.Error())
		return nil, errors.New("[MPPlanningRepository.StoreAttachmentToHeader] " + err.Error())
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		r.Log.Errorf("[MPPlanningRepository.StoreAttachmentToHeader] " + err.Error())
		return nil, errors.New("[MPPlanningRepository.StoreAttachmentToHeader] " + err.Error())
	}

	return mppHeader, nil
}

func (r *MPPlanningRepository) StoreAttachmentToApprovalHistory(mppApprovalHistory *entity.MPPlanningApprovalHistory, attachment entity.ManpowerAttachment) (*entity.MPPlanningApprovalHistory, error) {
	tx := r.DB.Begin()

	if tx.Error != nil {
		r.Log.Errorf("[MPPlanningRepository.StoreAttachmentToApprovalHistory] " + tx.Error.Error())
		return nil, errors.New("[MPPlanningRepository.StoreAttachmentToApprovalHistory] " + tx.Error.Error())
	}

	if err := tx.Create(&attachment).Error; err != nil {
		tx.Rollback()
		r.Log.Errorf("[MPPlanningRepository.StoreAttachmentToApprovalHistory] " + err.Error())
		return nil, errors.New("[MPPlanningRepository.StoreAttachmentToApprovalHistory] " + err.Error())
	}

	// if err := tx.Model(mppApprovalHistory).Association("ManpowerAttachments").Append(&attachment); err != nil {
	// 	tx.Rollback()
	// 	r.Log.Errorf("[MPPlanningRepository.StoreAttachmentToApprovalHistory] " + err.Error())
	// 	return nil, errors.New("[MPPlanningRepository.StoreAttachmentToApprovalHistory] " + err.Error())
	// }

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		r.Log.Errorf("[MPPlanningRepository.StoreAttachmentToApprovalHistory] " + err.Error())
		return nil, errors.New("[MPPlanningRepository.StoreAttachmentToApprovalHistory] " + err.Error())
	}

	return mppApprovalHistory, nil
}

func (r *MPPlanningRepository) StoreAttachmentsToApprovalHistory(mppApprovalHistories *entity.MPPlanningApprovalHistory, attachments []entity.ManpowerAttachment) (*entity.MPPlanningApprovalHistory, error) {
	tx := r.DB.Begin()

	if tx.Error != nil {
		r.Log.Errorf("[MPPlanningRepository.StoreAttachmentsToApprovalHistory] " + tx.Error.Error())
		return nil, errors.New("[MPPlanningRepository.StoreAttachmentsToApprovalHistory] " + tx.Error.Error())
	}

	if err := tx.Model(mppApprovalHistories).Association("ManpowerAttachments").Append(&attachments); err != nil {
		tx.Rollback()
		r.Log.Errorf("[MPPlanningRepository.StoreAttachmentsToApprovalHistory] " + err.Error())
		return nil, errors.New("[MPPlanningRepository.StoreAttachmentsToApprovalHistory] " + err.Error())
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		r.Log.Errorf("[MPPlanningRepository.StoreAttachmentsToApprovalHistory] " + err.Error())
		return nil, errors.New("[MPPlanningRepository.StoreAttachmentsToApprovalHistory] " + err.Error())
	}

	return mppApprovalHistories, nil
}

func (r *MPPlanningRepository) GetPlanningApprovalHistoryByHeaderId(headerId uuid.UUID) (*[]entity.MPPlanningApprovalHistory, error) {
	var mppApprovalHistories []entity.MPPlanningApprovalHistory

	if err := r.DB.Preload("ManpowerAttachments").Where("mp_planning_header_id = ?", headerId).Find(&mppApprovalHistories).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			r.Log.Errorf("[MPPlanningRepository.GetPlanningApprovalHistoryByHeaderId] " + err.Error())
			return nil, nil
		} else {
			r.Log.Errorf("[MPPlanningRepository.GetPlanningApprovalHistoryByHeaderId] " + err.Error())
			return nil, errors.New("[MPPlanningRepository.GetPlanningApprovalHistoryByHeaderId] " + err.Error())
		}
	}

	return &mppApprovalHistories, nil
}

func (r *MPPlanningRepository) GetPlanningApprovalHistoryAttachmentsByApprovalHistoryId(approvalHistoryId uuid.UUID) (*[]entity.ManpowerAttachment, error) {
	var attachments []entity.ManpowerAttachment

	if err := r.DB.Where("owner_id = ? AND owner_type = ?", approvalHistoryId, "mp_planning_approval_histories").Find(&attachments).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			r.Log.Errorf("[MPPlanningRepository.GetPlanningApprovalHistoryAttachmentsByApprovalHistoryId] " + err.Error())
			return nil, nil
		} else {
			r.Log.Errorf("[MPPlanningRepository.GetPlanningApprovalHistoryAttachmentsByApprovalHistoryId] " + err.Error())
			return nil, errors.New("[MPPlanningRepository.GetPlanningApprovalHistoryAttachmentsByApprovalHistoryId] " + err.Error())
		}
	}

	return &attachments, nil
}

func (r *MPPlanningRepository) DeleteAttachmentFromHeader(mppHeader *entity.MPPlanningHeader, attachmentID uuid.UUID) (*entity.MPPlanningHeader, error) {
	tx := r.DB.Begin()

	if tx.Error != nil {
		r.Log.Errorf("[MPPlanningRepository.DeleteAttachmentFromHeader] " + tx.Error.Error())
		return nil, errors.New("[MPPlanningRepository.DeleteAttachmentFromHeader] " + tx.Error.Error())
	}

	if err := tx.Model(mppHeader).Association("ManpowerAttachments").Delete(&entity.ManpowerAttachment{ID: attachmentID}); err != nil {
		tx.Rollback()
		r.Log.Errorf("[MPPlanningRepository.DeleteAttachmentFromHeader] " + err.Error())
		return nil, errors.New("[MPPlanningRepository.DeleteAttachmentFromHeader] " + err.Error())
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		r.Log.Errorf("[MPPlanningRepository.DeleteAttachmentFromHeader] " + err.Error())
		return nil, errors.New("[MPPlanningRepository.DeleteAttachmentFromHeader] " + err.Error())
	}

	return mppHeader, nil
}

func (r *MPPlanningRepository) DeleteHeader(id uuid.UUID) error {
	tx := r.DB.Begin()

	if tx.Error != nil {
		r.Log.Errorf("[MPPlanningRepository.DeleteHeader] " + tx.Error.Error())
		return errors.New("[MPPlanningRepository.DeleteHeader] " + tx.Error.Error())
	}

	if err := tx.Where("id = ?", id).Delete(&entity.MPPlanningHeader{}).Error; err != nil {
		tx.Rollback()
		r.Log.Errorf("[MPPlanningRepository.DeleteHeader] " + err.Error())
		return errors.New("[MPPlanningRepository.DeleteHeader] " + err.Error())
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		r.Log.Errorf("[MPPlanningRepository.DeleteHeader] " + err.Error())
		return errors.New("[MPPlanningRepository.DeleteHeader] " + err.Error())
	}

	return nil
}

func (r *MPPlanningRepository) FindHeaderByMPPPeriodId(mppPeriodId uuid.UUID) (*entity.MPPlanningHeader, error) {
	var mppHeader entity.MPPlanningHeader

	if err := r.DB.Preload("MPPlanningLines").Preload("MPPPeriod").Where("mpp_period_id = ?", mppPeriodId).First(&mppHeader).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			r.Log.Errorf("[MPPlanningRepository.FindHeaderByMPPPeriodId] " + err.Error())
			return nil, nil
		} else {
			r.Log.Errorf("[MPPlanningRepository.FindHeaderByMPPPeriodId] " + err.Error())
			return nil, errors.New("[MPPlanningRepository.FindHeaderByMPPPeriodId] " + err.Error())
		}
	}

	return &mppHeader, nil
}

func (r *MPPlanningRepository) FindAllLinesByHeaderIdPaginated(headerId uuid.UUID, page int, pageSize int, search string) (*[]entity.MPPlanningLine, int64, error) {
	var mppLines []entity.MPPlanningLine

	query := r.DB.Model(&entity.MPPlanningLine{}).Preload("MPPlanningHeader").Where("mp_planning_header_id = ?", headerId)

	if search != "" {
		query = query.Where("name LIKE ?", "%"+search+"%")
	}

	if err := query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&mppLines).Error; err != nil {
		r.Log.Errorf("[MPPlanningRepository.FindAllLinesByHeaderIdPaginated] " + err.Error())
		return nil, 0, errors.New("[MPPlanningRepository.FindAllLinesByHeaderIdPaginated] " + err.Error())
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		r.Log.Errorf("[MPPlanningRepository.FindAllLinesByHeaderIdPaginated] " + err.Error())
		return nil, 0, errors.New("[MPPlanningRepository.FindAllLinesByHeaderIdPaginated] " + err.Error())
	}

	return &mppLines, total, nil
}

func (r *MPPlanningRepository) FindLineById(id uuid.UUID) (*entity.MPPlanningLine, error) {
	var mppLine entity.MPPlanningLine

	if err := r.DB.Preload("MPPlanningHeader").Where("id = ?", id).First(&mppLine).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			r.Log.Errorf("[MPPlanningRepository.FindLineById] " + err.Error())
			return nil, nil
		} else {
			r.Log.Errorf("[MPPlanningRepository.FindLineById] " + err.Error())
			return nil, errors.New("[MPPlanningRepository.FindLineById] " + err.Error())
		}
	}

	return &mppLine, nil
}

func (r *MPPlanningRepository) CreateLine(mppLine *entity.MPPlanningLine) (*entity.MPPlanningLine, error) {
	tx := r.DB.Begin()

	if tx.Error != nil {
		r.Log.Errorf("[MPPlanningRepository.CreateLine] " + tx.Error.Error())
		return nil, errors.New("[MPPlanningRepository.CreateLine] " + tx.Error.Error())
	}

	if err := tx.Create(mppLine).Error; err != nil {
		tx.Rollback()
		r.Log.Errorf("[MPPlanningRepository.CreateLine] " + err.Error())
		return nil, errors.New("[MPPlanningRepository.CreateLine] " + err.Error())
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		r.Log.Errorf("[MPPlanningRepository.CreateLine] " + err.Error())
		return nil, errors.New("[MPPlanningRepository.CreateLine] " + err.Error())
	}

	return mppLine, nil
}

func (r *MPPlanningRepository) UpdateLine(mppLine *entity.MPPlanningLine) (*entity.MPPlanningLine, error) {
	tx := r.DB.Begin()

	if tx.Error != nil {
		r.Log.Errorf("[MPPlanningRepository.UpdateLine] " + tx.Error.Error())
		return nil, errors.New("[MPPlanningRepository.UpdateLine] " + tx.Error.Error())
	}

	if err := tx.Model(mppLine).Where("id = ?", mppLine.ID).Updates(mppLine).Error; err != nil {
		tx.Rollback()
		r.Log.Errorf("[MPPlanningRepository.UpdateLine] " + err.Error())
		return nil, errors.New("[MPPlanningRepository.UpdateLine] " + err.Error())
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		r.Log.Errorf("[MPPlanningRepository.UpdateLine] " + err.Error())
		return nil, errors.New("[MPPlanningRepository.UpdateLine] " + err.Error())
	}

	return mppLine, nil
}

func (r *MPPlanningRepository) DeleteLine(id uuid.UUID) error {
	tx := r.DB.Begin()

	if tx.Error != nil {
		r.Log.Errorf("[MPPlanningRepository.DeleteLine] " + tx.Error.Error())
		return errors.New("[MPPlanningRepository.DeleteLine] " + tx.Error.Error())
	}

	if err := tx.Where("id = ?", id).Delete(&entity.MPPlanningLine{}).Error; err != nil {
		tx.Rollback()
		r.Log.Errorf("[MPPlanningRepository.DeleteLine] " + err.Error())
		return errors.New("[MPPlanningRepository.DeleteLine] " + err.Error())
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		r.Log.Errorf("[MPPlanningRepository.DeleteLine] " + err.Error())
		return errors.New("[MPPlanningRepository.DeleteLine] " + err.Error())
	}

	return nil
}

func (r *MPPlanningRepository) DeleteLineIfNotInIDs(ids []uuid.UUID) error {
	tx := r.DB.Begin()

	if tx.Error != nil {
		r.Log.Errorf("[MPPlanningRepository.DeleteLineIfNotInIDs] " + tx.Error.Error())
		return errors.New("[MPPlanningRepository.DeleteLineIfNotInIDs] " + tx.Error.Error())
	}

	if err := tx.Where("id NOT IN ?", ids).Delete(&entity.MPPlanningLine{}).Error; err != nil {
		tx.Rollback()
		r.Log.Errorf("[MPPlanningRepository.DeleteLineIfNotInIDs] " + err.Error())
		return errors.New("[MPPlanningRepository.DeleteLineIfNotInIDs] " + err.Error())
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		r.Log.Errorf("[MPPlanningRepository.DeleteLineIfNotInIDs] " + err.Error())
		return errors.New("[MPPlanningRepository.DeleteLineIfNotInIDs] " + err.Error())
	}

	return nil
}

func (r *MPPlanningRepository) DeleteAllLinesByHeaderID(headerID uuid.UUID) error {
	tx := r.DB.Begin()

	if tx.Error != nil {
		r.Log.Errorf("[MPPlanningRepository.DeleteAllLinesByHeaderID] " + tx.Error.Error())
		return errors.New("[MPPlanningRepository.DeleteAllLinesByHeaderID] " + tx.Error.Error())
	}

	if err := tx.Where("mp_planning_header_id = ?", headerID).Delete(&entity.MPPlanningLine{}).Error; err != nil {
		tx.Rollback()
		r.Log.Errorf("[MPPlanningRepository.DeleteAllLinesByHeaderID] " + err.Error())
		return errors.New("[MPPlanningRepository.DeleteAllLinesByHeaderID] " + err.Error())
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		r.Log.Errorf("[MPPlanningRepository.DeleteAllLinesByHeaderID] " + err.Error())
		return errors.New("[MPPlanningRepository.DeleteAllLinesByHeaderID] " + err.Error())
	}

	r.Log.Infof("[MPPlanningRepository.DeleteAllLinesByHeaderID] Successfully deleted all lines by header ID")

	return nil
}

func (r *MPPlanningRepository) DeleteLinesByIDs(ids []uuid.UUID) error {
	tx := r.DB.Begin()

	if tx.Error != nil {
		r.Log.Errorf("[MPPlanningRepository.DeleteLinesByIDs] " + tx.Error.Error())
		return errors.New("[MPPlanningRepository.DeleteLinesByIDs] " + tx.Error.Error())
	}

	if err := tx.Where("id IN ?", ids).Delete(&entity.MPPlanningLine{}).Error; err != nil {
		tx.Rollback()
		r.Log.Errorf("[MPPlanningRepository.DeleteLinesByIDs] " + err.Error())
		return errors.New("[MPPlanningRepository.DeleteLinesByIDs] " + err.Error())
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		r.Log.Errorf("[MPPlanningRepository.DeleteLinesByIDs] " + err.Error())
		return errors.New("[MPPlanningRepository.DeleteLinesByIDs] " + err.Error())
	}

	r.Log.Infof("[MPPlanningRepository.DeleteLinesByIDs] Successfully deleted lines by IDs")

	return nil
}

func MPPlanningRepositoryFactory(log *logrus.Logger) IMPPlanningRepository {
	db := config.NewDatabase()
	return NewMPPlanningRepository(log, db)
}
