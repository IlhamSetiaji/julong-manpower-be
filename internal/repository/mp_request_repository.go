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

type IMPRequestRepository interface {
	GetHighestDocumentNumberByDate(date string) (int, error)
	Create(mpRequestHeader *entity.MPRequestHeader) (*entity.MPRequestHeader, error)
	FindAllPaginated(page int, pageSize int, search string, filter map[string]interface{}) ([]entity.MPRequestHeader, int64, error)
	FindAll() ([]entity.MPRequestHeader, error)
	GetHeadersByDocumentDate(documentDate string) ([]entity.MPRequestHeader, error)
	GetHeadersByCreatedAt(createdAt string) ([]entity.MPRequestHeader, error)
	GetRequestApprovalHistoryByHeaderId(headerID uuid.UUID, status entity.MPRequestApprovalHistoryStatus) ([]entity.MPRequestApprovalHistory, error)
	FindById(id uuid.UUID) (*entity.MPRequestHeader, error)
	FindByIDOnly(id uuid.UUID) (*entity.MPRequestHeader, error)
	Update(mpRequestHeader *entity.MPRequestHeader) (*entity.MPRequestHeader, error)
	UpdateStatusHeader(id uuid.UUID, status string, approvedBy string, approvalHistory *entity.MPRequestApprovalHistory) error
	StoreAttachmentToApprovalHistory(mppApprovalHistory *entity.MPRequestApprovalHistory, attachment entity.ManpowerAttachment) (*entity.MPRequestApprovalHistory, error)
	DeleteHeader(id uuid.UUID) error
	CountTotalApprovalHistoryByStatus(mpHeaderID uuid.UUID, status entity.MPRequestApprovalHistoryStatus) (int64, error)
	FindByKeys(keys map[string]interface{}) (*entity.MPRequestHeader, error)
	FindAllByMajorIds(majorIds []string) ([]entity.MPRequestHeader, error)
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

func (r *MPRequestRepository) FindAllByMajorIds(majorIds []string) ([]entity.MPRequestHeader, error) {
	var mpRequestHeaders []entity.MPRequestHeader

	r.Log.Infof("[MPRequestRepository.FindAllByMajorIds] majorIds: %v", majorIds)
	if err := r.DB.Preload("RequestCategory").
		Preload("RequestMajors.Major").
		Preload("MPPlanningHeader").
		Joins("JOIN request_majors ON request_majors.mp_request_header_id = mp_request_headers.id").
		Where("request_majors.major_id IN ?", majorIds).
		Find(&mpRequestHeaders).Error; err != nil {
		r.Log.Errorf("[MPRequestRepository.FindAllByMajorIds] error when query mp request headers: %v", err)
		return nil, errors.New("[MPRequestRepository.FindAllByMajorIds] error when query mp request headers " + err.Error())
	}

	return mpRequestHeaders, nil
}

func (r *MPRequestRepository) GetHighestDocumentNumberByDate(date string) (int, error) {
	var maxNumber int
	err := r.DB.Raw(`
			SELECT COALESCE(MAX(CAST(SUBSTRING(document_number FROM '[0-9]+$') AS INTEGER)), 0)
			FROM mp_request_headers
			WHERE DATE(created_at) = ?
	`, date).Scan(&maxNumber).Error
	if err != nil {
		r.Log.Errorf("[MPPlanningRepository.GetHighestDocumentNumberByDate] error when querying max document number: %v", err)
		return 0, err
	}
	return maxNumber, nil
}

func (r *MPRequestRepository) FindByIDOnly(id uuid.UUID) (*entity.MPRequestHeader, error) {
	var mpRequestHeader entity.MPRequestHeader

	if err := r.DB.Preload("MPPPeriod").Preload("RequestCategory").Preload("RequestMajors.Major").First(&mpRequestHeader, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			r.Log.Errorf("[MPRequestRepository.FindByIDOnly] mp request header with id %s not found", id)
			return nil, nil
		} else {
			r.Log.Errorf("[MPRequestRepository.FindByIDOnly] error when query mp request header: %v", err)
		}

		return nil, errors.New("[MPRequestRepository.FindByIDOnly] error when query mp request header " + err.Error())
	}

	return &mpRequestHeader, nil
}

func (r *MPRequestRepository) GetHeadersByDocumentDate(documentDate string) ([]entity.MPRequestHeader, error) {
	var mpRequestHeaders []entity.MPRequestHeader

	if err := r.DB.Preload("RequestCategory").Preload("RequestMajors.Major").Preload("MPPlanningHeader").Where("document_date = ?", documentDate).Find(&mpRequestHeaders).Error; err != nil {
		r.Log.Errorf("[MPRequestRepository.GetHeadersByDocumentDate] error when query mp request headers: %v", err)
		return nil, errors.New("[MPRequestRepository.GetHeadersByDocumentDate] error when query mp request headers " + err.Error())
	}

	return mpRequestHeaders, nil
}

func (r *MPRequestRepository) GetRequestApprovalHistoryByHeaderId(headerID uuid.UUID, status entity.MPRequestApprovalHistoryStatus) ([]entity.MPRequestApprovalHistory, error) {
	var mpRequestApprovalHistories []entity.MPRequestApprovalHistory
	var whereStatus string = ""
	if status != "" {
		whereStatus = "status = '" + string(status) + "'"
	}

	if err := r.DB.Where("mp_request_header_id = ?", headerID).Where(whereStatus).Find(&mpRequestApprovalHistories).Error; err != nil {
		r.Log.Errorf("[MPRequestRepository.GetRequestApprovalHistoryByHeaderId] error when query mp request approval histories: %v", err)
		return nil, errors.New("[MPRequestRepository.GetRequestApprovalHistoryByHeaderId] error when query mp request approval histories " + err.Error())
	}

	return mpRequestApprovalHistories, nil
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

func (r *MPRequestRepository) GetHeadersByCreatedAt(createdAt string) ([]entity.MPRequestHeader, error) {
	var mpRequestHeaders []entity.MPRequestHeader

	var formatTimeCreatedAt time.Time
	var err error

	formats := []string{
		"2006-01-02",
		time.RFC3339,
	}

	for _, format := range formats {
		formatTimeCreatedAt, err = time.Parse(format, createdAt)
		if err == nil {
			break
		}
	}

	if err != nil {
		r.Log.Errorf("[MPPlanningRepository.GetHeadersByCreatedAt] " + err.Error())
		return nil, errors.New("[MPPlanningRepository.GetHeadersByCreatedAt] " + err.Error())
	}

	// Use only the date part for comparison
	dateOnly := formatTimeCreatedAt.Format("2006-01-02")

	if err := r.DB.Where("DATE(created_at) = ?", dateOnly).Find(&mpRequestHeaders).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			r.Log.Warnf("[MPPlanningRepository.GetHeadersByCreatedAt] no records found")
			return nil, nil
		} else {
			r.Log.Errorf("[MPPlanningRepository.GetHeadersByCreatedAt] " + err.Error())
			return nil, errors.New("[MPPlanningRepository.GetHeadersByCreatedAt] " + err.Error())
		}
	}

	return mpRequestHeaders, nil
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

	var mprHeader entity.MPRequestHeader

	if err := tx.First(&mprHeader, id).Error; err != nil {
		tx.Rollback()
		r.Log.Errorf("[MPRequestRepository.DeleteHeader] error when query mp request header: %v", err)
		return errors.New("[MPRequestRepository.DeleteHeader] error when query mp request header " + err.Error())
	}

	if err := tx.Where("id = ?", id).Delete(&mprHeader).Error; err != nil {
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

	if err := r.DB.Preload("MPRequestApprovalHistories").Preload("MPPPeriod").Preload("RequestCategory").Preload("RequestMajors.Major").Preload("MPPlanningHeader.MPPlanningLines").First(&mpRequestHeader, id).Error; err != nil {
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
	var isAdmin bool = false
	// var includedIDs []string = []string{}

	query := r.DB.Preload("MPPPeriod").Preload("RequestCategory").Preload("RequestMajors.Major").Model(&entity.MPRequestHeader{})

	if search != "" {
		query = query.Where("document_number LIKE ?", "%"+search+"%")
	}

	if filter != nil {
		// if _, ok := filter["requestor_id"]; ok {
		// 	query = query.Where("requestor_id = ?", filter["requestor_id"])
		// }
		// if _, ok := filter["department_head"]; ok {
		// 	if filter["department_head"] == "NULL" {
		// 		query = query.Where("department_head IS NULL")
		// 	} else {
		// 		query = query.Where("department_head IS NOT NULL")
		// 	}
		// }
		// if _, ok := filter["vp_gm_director"]; ok {
		// 	if filter["vp_gm_director"] == "NULL" {
		// 		query = query.Where("vp_gm_director IS NULL")
		// 	} else {
		// 		query = query.Where("vp_gm_director IS NOT NULL")
		// 	}
		// }
		// if _, ok := filter["ceo"]; ok {
		// 	if filter["ceo"] == "NULL" {
		// 		query = query.Where("ceo IS NULL")
		// 	} else {
		// 		query = query.Where("ceo IS NOT NULL")
		// 	}
		// }
		// if _, ok := filter["hrd_ho_unit"]; ok {
		// 	if filter["hrd_ho_unit"] == "NULL" {
		// 		query = query.Where("hrd_ho_unit IS NULL")
		// 	} else {
		// 		query = query.Where("hrd_ho_unit IS NOT NULL")
		// 	}
		// }
		if _, ok := filter["approver_type"]; ok {
			switch filter["approver_type"] {
			case "requestor":
				query = query.Where("requestor_id = ?", filter["requestor_id"])
			case "department_head":
				r.Log.Info("Included IDs: ", filter["included_ids"])
				query = query.Where("department_head IS NULL OR department_head IS NOT NULL").Where("for_organization_structure_id IN (?)", filter["included_ids"]).Or("requestor_id = ?", filter["requestor_id"])
			case "vp_gm_director":
				query = query.Where("vp_gm_director IS NULL OR vp_gm_director IS NOT NULL").Where("organization_id = ?", filter["organization_id"]).Where("mp_request_type = ?", entity.MPRequestTypeEnumOffBudget).Or("requestor_id = ?", filter["requestor_id"])
			case "ceo":
				query = query.Where("ceo IS NULL OR ceo IS NOT NULL").Where("mp_request_type = ?", entity.MPRequestTypeEnumOffBudget)
			case "hrd_ho_unit":
				query = query.Where("hrd_ho_unit IS NULL OR hrd_ho_unit IS NOT NULL").Where("status IN (?)", []entity.MPRequestStatus{entity.MPRequestStatusApproved, entity.MPRequestStatusCompleted}).Or("requestor_id = ?", filter["requestor_id"])
			default:
				query = query.Where("requestor_id = ?", filter["requestor_id"])
			}
		}
		if _, ok := filter["is_admin"]; ok {
			isAdmin = true
		}
		if status, ok := filter["status"]; ok {
			query = query.Where("status = ?", status)
		}
		if !isAdmin {
			if filter["included_ids"] != nil {
				// includedIDs = filter["included_ids"].([]string)
				// query = query.Where("id IN (?)", includedIDs)
			}
		}
	}

	if err := query.Count(&total).Error; err != nil {
		r.Log.Errorf("[MPRequestRepository.FindAllPaginated] error when count mp request headers: %v", err)
		return nil, 0, errors.New("[MPRequestRepository.FindAllPaginated] error when count mp request headers " + err.Error())
	}

	if err := query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&mpRequestHeaders).Error; err != nil {
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
			r.Log.Info("Ini bukan reject")
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
			r.Log.Info("Ini reject")
			if err := tx.Model(&entity.MPRequestHeader{}).Where("id = ?", id).Select("Status", "DepartmentHead", "VpGmDirector", "CEO", "HrdHoUnit").Updates(&entity.MPRequestHeader{
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

func (r *MPRequestRepository) FindByKeys(keys map[string]interface{}) (*entity.MPRequestHeader, error) {
	var mpRequestHeader entity.MPRequestHeader

	if err := r.DB.Where(keys).First(&mpRequestHeader).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			r.Log.Errorf("[MPRequestRepository.FindByKeys] mp request header not found")
			return nil, nil
		} else {
			r.Log.Errorf("[MPRequestRepository.FindByKeys] error when query mp request header: %v", err)
		}
	}

	return &mpRequestHeader, nil
}

func MPRequestRepositoryFactory(log *logrus.Logger) IMPRequestRepository {
	db := config.NewDatabase()
	return NewMPRequestRepository(log, db)
}
