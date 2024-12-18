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
	FindAllPaginated(page int, pageSize int, search string, filter map[string]interface{}) ([]entity.MPRequestHeader, int64, error)
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

	if err := r.DB.Preload("RequestCategory").Preload("RequestMajors.Major").Preload("MPPlanningHeader").First(mpRequestHeader, mpRequestHeader.ID).Error; err != nil {
		r.Log.Errorf("[MPRequestRepository.Create] error when preloading associations: %v", err)
		return nil, errors.New("[MPRequestRepository.Create] error when preloading associations " + err.Error())
	}

	return mpRequestHeader, nil
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

func MPRequestRepositoryFactory(log *logrus.Logger) IMPRequestRepository {
	db := config.NewDatabase()
	return NewMPRequestRepository(log, db)
}
