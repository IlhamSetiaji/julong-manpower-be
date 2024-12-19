package repository

import (
	"errors"

	"github.com/IlhamSetiaji/julong-manpower-be/internal/config"
	"github.com/IlhamSetiaji/julong-manpower-be/internal/entity"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type IMPRequestRepository interface {
	Create(mpRequestHeader *entity.MPRequestHeader) (*entity.MPRequestHeader, error)
	FindAllPaginated(page int, pageSize int, search string, filter map[string]interface{}) ([]entity.MPRequestHeader, int64, error)
	FindAll() ([]entity.MPRequestHeader, error)
	GetHeadersByDocumentDate(documentDate string) ([]entity.MPRequestHeader, error)
	FindById(id uuid.UUID) (*entity.MPRequestHeader, error)
	Update(mpRequestHeader *entity.MPRequestHeader) (*entity.MPRequestHeader, error)
	UpdateStatusHeader(id uuid.UUID, status string, approvedBy string, approvalHistory *entity.MPRequestApprovalHistory) error
	StoreAttachmentToApprovalHistory(mppApprovalHistory *entity.MPRequestApprovalHistory, attachment entity.ManpowerAttachment) (*entity.MPRequestApprovalHistory, error)
	DeleteHeader(id uuid.UUID) error
	CountTotalApprovalHistoryByStatus(mpHeaderID uuid.UUID, status entity.MPRequestApprovalHistoryStatus) (int64, error)
}

type MPRequestRepository struct {
	Log *logrus.Logger
	DB  *gorm.DB
}

func NewMPRequestRepository(log *logrus.Logger, db *gorm.DB) IMPRequestRepository {
	return &MPRequestRepository{Log: log, DB: db}
}

func (r *MPRequestRepository) FindAll() ([]entity.MPRequestHeader, error) {
	var mpRequestHeaders []entity.MPRequestHeader

	if err := r.DB.Preload("RequestCategory").Preload("RequestMajors.Major").Preload("MPPlanningHeader").Find(&mpRequestHeaders).Error; err != nil {
		r.Log.Errorf("[MPRequestRepository.FindAll] error when query mp request headers: %v", err)
		return nil, errors.New("[MPRequestRepository.FindAll] error when query mp request headers " + err.Error())
	}

	return mpRequestHeaders, nil
}

func (r *MPRequestRepository) GetHeadersByDocumentDate(documentDate string) ([]entity.MPRequestHeader, error) {
	var mpRequestHeaders []entity.MPRequestHeader

	if err := r.DB.Preload("RequestCategory").Preload("RequestMajors.Major").Preload("MPPlanningHeader").Where("document_date = ?", documentDate).Find(&mpRequestHeaders).Error; err != nil {
		r.Log.Errorf("[MPRequestRepository.GetHeadersByDocumentDate] error when query mp request headers: %v", err)
		return nil, errors.New("[MPRequestRepository.GetHeadersByDocumentDate] error when query mp request headers " + err.Error())
	}

	return mpRequestHeaders, nil
}

func (r *MPRequestRepository) CountTotalApprovalHistoryByStatus(mpHeaderID uuid.UUID, status entity.MPRequestApprovalHistoryStatus) (int64, error) {
	var total int64

	if err := r.DB.Model(&entity.MPRequestApprovalHistory{}).Where("mp_request_header_id = ? AND status = ?", mpHeaderID, status).Count(&total).Error; err != nil {
		r.Log.Errorf("[MPRequestRepository.CountTotalApprovalHistoryByStatus] error when count total approval history: %v", err)
		return 0, errors.New("[MPRequestRepository.CountTotalApprovalHistoryByStatus] error when count total approval history " + err.Error())
	}

	return total, nil
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

	if err := r.DB.Preload("RequestCategory").Preload("RequestMajors.Major").Preload("MPPlanningHeader").First(mpRequestHeader, mpRequestHeader.ID).Error; err != nil {
		r.Log.Errorf("[MPRequestRepository.Create] error when preloading associations: %v", err)
		return nil, errors.New("[MPRequestRepository.Create] error when preloading associations " + err.Error())
	}

	return mpRequestHeader, nil
}

func (r *MPRequestRepository) Update(mpRequestHeader *entity.MPRequestHeader) (*entity.MPRequestHeader, error) {
	tx := r.DB.Begin()

	if err := tx.Where("id = ?", mpRequestHeader.ID).Updates(&mpRequestHeader).Error; err != nil {
		tx.Rollback()
		r.Log.Errorf("[MPRequestRepository.Update] error when update mp request header: %v", err)
		return nil, errors.New("[MPRequestRepository.Update] error when update mp request header " + err.Error())
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		r.Log.Errorf("[MPRequestRepository.Update] error when commit transaction: %v", err)
		return nil, errors.New("[MPRequestRepository.Update] error when commit transaction " + err.Error())
	}

	if err := r.DB.Preload("RequestCategory").Preload("RequestMajors.Major").Preload("MPPlanningHeader").First(mpRequestHeader, mpRequestHeader.ID).Error; err != nil {
		r.Log.Errorf("[MPRequestRepository.Update] error when preloading associations: %v", err)
		return nil, errors.New("[MPRequestRepository.Update] error when preloading associations " + err.Error())
	}

	return mpRequestHeader, nil
}

func (r *MPRequestRepository) DeleteHeader(id uuid.UUID) error {
	tx := r.DB.Begin()

	if tx.Error != nil {
		r.Log.Errorf("[MPRequestRepository.DeleteHeader] " + tx.Error.Error())
		return errors.New("[MPRequestRepository.DeleteHeader] " + tx.Error.Error())
	}

	if err := tx.Where("id = ?", id).Delete(&entity.MPRequestHeader{}).Error; err != nil {
		tx.Rollback()
		r.Log.Errorf("[MPRequestRepository.DeleteHeader] error when delete mp request header: %v", err)
		return errors.New("[MPRequestRepository.DeleteHeader] error when delete mp request header " + err.Error())
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		r.Log.Errorf("[MPRequestRepository.DeleteHeader] error when commit transaction: %v", err)
		return errors.New("[MPRequestRepository.DeleteHeader] error when commit transaction " + err.Error())
	}

	return nil
}

func (r *MPRequestRepository) FindById(id uuid.UUID) (*entity.MPRequestHeader, error) {
	var mpRequestHeader entity.MPRequestHeader

	if err := r.DB.Preload("MPPPeriod").Preload("RequestCategory").Preload("RequestMajors.Major").Preload("MPPlanningHeader.MPPlanningLines").First(&mpRequestHeader, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			r.Log.Errorf("[MPRequestRepository.FindById] mp request header with id %s not found", id)
			return nil, nil
		} else {
			r.Log.Errorf("[MPRequestRepository.FindById] error when query mp request header: %v", err)
		}
	}

	return &mpRequestHeader, nil
}

func (r *MPRequestRepository) StoreAttachmentToApprovalHistory(mpApprovalHistory *entity.MPRequestApprovalHistory, attachment entity.ManpowerAttachment) (*entity.MPRequestApprovalHistory, error) {
	tx := r.DB.Begin()

	if tx.Error != nil {
		r.Log.Errorf("[MPRequestRepository.StoreAttachmentToApprovalHistory] " + tx.Error.Error())
		return nil, errors.New("[MPRequestRepository.StoreAttachmentToApprovalHistory] " + tx.Error.Error())
	}

	if err := tx.Create(&attachment).Error; err != nil {
		tx.Rollback()
		r.Log.Errorf("[MPRequestRepository.StoreAttachmentToApprovalHistory] error when create attachment: %v", err)
		return nil, errors.New("[MPRequestRepository.StoreAttachmentToApprovalHistory] error when create attachment " + err.Error())
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		r.Log.Errorf("[MPRequestRepository.StoreAttachmentToApprovalHistory] error when commit transaction: %v", err)
		return nil, errors.New("[MPRequestRepository.StoreAttachmentToApprovalHistory] error when commit transaction " + err.Error())
	}

	return mpApprovalHistory, nil
}

func (r *MPRequestRepository) FindAllPaginated(page int, pageSize int, search string, filter map[string]interface{}) ([]entity.MPRequestHeader, int64, error) {
	var mpRequestHeaders []entity.MPRequestHeader
	var total int64

	query := r.DB.Preload("MPPPeriod").Model(&entity.MPRequestHeader{})

	if search != "" {
		query = query.Where("document_number LIKE ?", "%"+search+"%")
	}

	if filter != nil {
		if _, ok := filter["department_head"]; ok {
			query = query.Where("department_head IS NOT NULL")
		}
		if _, ok := filter["vp_gm_director"]; ok {
			query = query.Where("vp_gm_director IS NOT NULL")
		}
		if _, ok := filter["ceo"]; ok {
			query = query.Where("ceo IS NOT NULL")
		}
		if _, ok := filter["hrd_ho_unit"]; ok {
			query = query.Where("hrd_ho_unit IS NOT NULL")
		}
		if status, ok := filter["status"]; ok {
			query = query.Where("status = ?", status)
		}
	}

	if err := query.Count(&total).Error; err != nil {
		r.Log.Errorf("[MPRequestRepository.FindAllPaginated] error when count mp request headers: %v", err)
		return nil, 0, errors.New("[MPRequestRepository.FindAllPaginated] error when count mp request headers " + err.Error())
	}

	if err := query.Preload("RequestCategory").Preload("RequestMajors.Major").Preload("MPPlanningHeader").Offset((page - 1) * pageSize).Limit(pageSize).Find(&mpRequestHeaders).Error; err != nil {
		r.Log.Errorf("[MPRequestRepository.FindAllPaginated] error when find mp request headers: %v", err)
		return nil, 0, errors.New("[MPRequestRepository.FindAllPaginated] error when find mp request headers " + err.Error())
	}

	return mpRequestHeaders, total, nil
}

func (r *MPRequestRepository) UpdateStatusHeader(id uuid.UUID, status string, approvedBy string, approvalHistory *entity.MPRequestApprovalHistory) error {
	tx := r.DB.Begin()

	if tx.Error != nil {
		r.Log.Errorf("[MPPlanningRepository.UpdateStatusHeader] " + tx.Error.Error())
		return errors.New("[MPPlanningRepository.UpdateStatusHeader] " + tx.Error.Error())
	}

	var approvedByPtr *uuid.UUID

	if approvedBy != "" {
		approvedByUUID, err := uuid.Parse(approvedBy)
		if err != nil {
			tx.Rollback()
			r.Log.Errorf("[MPRequestRepository.UpdateStatusHeader] error when parse approved by: %v", err)
			return errors.New("[MPRequestRepository.UpdateStatusHeader] error when parse approved by " + err.Error())
		}
		approvedByPtr = &approvedByUUID
	}

	if approvalHistory != nil {
		if approvalHistory.Status != entity.MPRequestApprovalHistoryStatusRejected {
			if approvalHistory.Level == string(entity.MPRequestApprovalHistoryLevelCEO) {
				if err := tx.Model(&entity.MPRequestHeader{}).Where("id = ?", id).Updates(&entity.MPRequestHeader{
					Status: entity.MPRequestStatus(status),
					CEO:    approvedByPtr,
				}).Error; err != nil {
					tx.Rollback()
					r.Log.Errorf("[MPRequestRepository.UpdateStatusHeader] error when update mp request header: %v", err)
					return errors.New("[MPRequestRepository.UpdateStatusHeader] error when update mp request header " + err.Error())
				}
			} else if approvalHistory.Level == string(entity.MPRequestApprovalHistoryLevelVP) {
				if err := tx.Model(&entity.MPRequestHeader{}).Where("id = ?", id).Updates(&entity.MPRequestHeader{
					Status:       entity.MPRequestStatus(status),
					VpGmDirector: approvedByPtr,
				}).Error; err != nil {
					tx.Rollback()
					r.Log.Errorf("[MPRequestRepository.UpdateStatusHeader] error when update mp request header: %v", err)
					return errors.New("[MPRequestRepository.UpdateStatusHeader] error when update mp request header " + err.Error())
				}
			} else if approvalHistory.Level == string(entity.MPRequestApprovalHistoryLevelHeadDept) {
				if err := tx.Model(&entity.MPRequestHeader{}).Where("id = ?", id).Updates(&entity.MPRequestHeader{
					Status:         entity.MPRequestStatus(status),
					DepartmentHead: approvedByPtr,
				}).Error; err != nil {
					tx.Rollback()
					r.Log.Errorf("[MPRequestRepository.UpdateStatusHeader] error when update mp request header: %v", err)
					return errors.New("[MPRequestRepository.UpdateStatusHeader] error when update mp request header " + err.Error())
				}
			} else if approvalHistory.Level == string(entity.MPRequestApprovalHistoryLevelStaff) {
				if err := tx.Model(&entity.MPRequestHeader{}).Where("id = ?", id).Updates(&entity.MPRequestHeader{
					Status: entity.MPRequestStatus(status),
				}).Error; err != nil {
					tx.Rollback()
					r.Log.Errorf("[MPRequestRepository.UpdateStatusHeader] error when update mp request header: %v", err)
					return errors.New("[MPRequestRepository.UpdateStatusHeader] error when update mp request header " + err.Error())
				}
			} else if approvalHistory.Level == string(entity.MPPRequestApprovalHistoryLevelHRDHO) {
				if err := tx.Model(&entity.MPRequestHeader{}).Where("id = ?", id).Updates(&entity.MPRequestHeader{
					Status:    entity.MPRequestStatus(status),
					HrdHoUnit: approvedByPtr,
				}).Error; err != nil {
					tx.Rollback()
					r.Log.Errorf("[MPRequestRepository.UpdateStatusHeader] error when update mp request header: %v", err)
					return errors.New("[MPRequestRepository.UpdateStatusHeader] error when update mp request header " + err.Error())
				}
			}
		} else {
			if err := tx.Model(&entity.MPRequestHeader{}).Where("id = ?", id).Updates(&entity.MPRequestHeader{
				Status:         entity.MPRequestStatus(status),
				DepartmentHead: nil,
				VpGmDirector:   nil,
				CEO:            nil,
				HrdHoUnit:      nil,
			}).Error; err != nil {
				tx.Rollback()
				r.Log.Errorf("[MPRequestRepository.UpdateStatusHeader] error when update mp request header: %v", err)
				return errors.New("[MPRequestRepository.UpdateStatusHeader] error when update mp request header " + err.Error())
			}
		}
	} else {
		if err := tx.Model(&entity.MPRequestHeader{}).Where("id = ?", id).Updates(&entity.MPRequestHeader{
			Status: entity.MPRequestStatus(status),
		}).Error; err != nil {
			tx.Rollback()
			r.Log.Errorf("[MPRequestRepository.UpdateStatusHeader] error when update mp request header: %v", err)
			return errors.New("[MPRequestRepository.UpdateStatusHeader] error when update mp request header " + err.Error())
		}
	}

	if approvalHistory != nil {
		if err := tx.Create(approvalHistory).Error; err != nil {
			tx.Rollback()
			r.Log.Errorf("[MPRequestRepository.UpdateStatusHeader] error when create mp request approval history: %v", err)
			return errors.New("[MPRequestRepository.UpdateStatusHeader] error when create mp request approval history " + err.Error())
		}
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		r.Log.Errorf("[MPRequestRepository.UpdateStatusHeader] error when commit transaction: %v", err)
		return errors.New("[MPRequestRepository.UpdateStatusHeader] error when commit transaction " + err.Error())
	}

	return nil
}

func MPRequestRepositoryFactory(log *logrus.Logger) IMPRequestRepository {
	db := config.NewDatabase()
	return NewMPRequestRepository(log, db)
}
