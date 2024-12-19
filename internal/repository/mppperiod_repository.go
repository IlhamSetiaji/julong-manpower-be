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

type IMPPPeriodRepository interface {
	FindAllPaginated(page int, pageSize int, search string) (*[]entity.MPPPeriod, int64, error)
	FindById(id uuid.UUID) (*entity.MPPPeriod, error)
	Create(mppPeriod *entity.MPPPeriod) (*entity.MPPPeriod, error)
	Update(mppPeriod *entity.MPPPeriod) (*entity.MPPPeriod, error)
	Delete(id uuid.UUID) error
	FindByCurrentDateAndStatus(status entity.MPPPeriodStatus) (*entity.MPPPeriod, error)
	FindByStatus(status entity.MPPPeriodStatus) (*entity.MPPPeriod, error)
}

type MPPPeriodRepository struct {
	Log *logrus.Logger
	DB  *gorm.DB
}

func NewMPPPeriodRepository(log *logrus.Logger, db *gorm.DB) IMPPPeriodRepository {
	return &MPPPeriodRepository{
		Log: log,
		DB:  db,
	}
}

func (r *MPPPeriodRepository) FindByStatus(status entity.MPPPeriodStatus) (*entity.MPPPeriod, error) {
	var mppPeriod entity.MPPPeriod

	err := r.DB.Where("status = ?", status).First(&mppPeriod).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			r.Log.Warn("[MPPPeriodRepository.FindByStatus] User not found")
			return nil, nil
		} else {
			r.Log.Error("[MPPPeriodRepository.FindByStatus] " + err.Error())
			return nil, errors.New("[UserRepository.FindByEmail] " + err.Error())
		}
	}

	return &mppPeriod, nil
}

func (r *MPPPeriodRepository) FindAllPaginated(page int, pageSize int, search string) (*[]entity.MPPPeriod, int64, error) {
	var mppPeriods []entity.MPPPeriod
	var total int64

	query := r.DB.Model(&entity.MPPPeriod{})

	if search != "" {
		query = query.Where("name LIKE ?", "%"+search+"%")
	}

	if err := query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&mppPeriods).Error; err != nil {
		r.Log.Errorf("[MPPPeriodRepository.FindAllPaginated] " + err.Error())
		return nil, 0, errors.New("[MPPPeriodRepository.FindAllPaginated] " + err.Error())
	}

	if err := query.Count(&total).Error; err != nil {
		r.Log.Errorf("[MPPPeriodRepository.FindAllPaginated] " + err.Error())
		return nil, 0, errors.New("[MPPPeriodRepository.FindAllPaginated] " + err.Error())
	}

	return &mppPeriods, total, nil
}

func (r *MPPPeriodRepository) FindById(id uuid.UUID) (*entity.MPPPeriod, error) {
	var mppPeriod entity.MPPPeriod

	if err := r.DB.Where("id = ?", id).First(&mppPeriod).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			r.Log.Warn("[MPPPeriodRepository.FindById] Job not found")
			return nil, nil
		} else {
			r.Log.Errorf("[MPPPeriodRepository.FindById] " + err.Error())
			return nil, err
		}
	}

	return &mppPeriod, nil
}

func MPPPeriodRepositoryFactory(log *logrus.Logger) IMPPPeriodRepository {
	db := config.NewDatabase()
	return NewMPPPeriodRepository(log, db)
}

func (r *MPPPeriodRepository) Create(mppPeriod *entity.MPPPeriod) (*entity.MPPPeriod, error) {
	tx := r.DB.Begin()
	if tx.Error != nil {
		r.Log.Errorf("[MPPPeriodRepository.Create] " + tx.Error.Error())
		return nil, errors.New("[MPPPeriodRepository.Create] " + tx.Error.Error())
	}

	if mppPeriod.Status != "draft" {
		dateNow := time.Now().Format("2006-01-02")
		if dateNow < mppPeriod.StartDate.Format("2006-01-02") {
			mppPeriod.Status = entity.MPPeriodStatusNotOpen
		}
	}

	if err := tx.Create(mppPeriod).Error; err != nil {
		tx.Rollback()
		r.Log.Errorf("[MPPPeriodRepository.Create] " + err.Error())
		return nil, errors.New("[MPPPeriodRepository.Create] " + err.Error())
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		r.Log.Errorf("[MPPPeriodRepository.Create] " + err.Error())
		return nil, errors.New("[MPPPeriodRepository.Create] " + err.Error())
	}

	return mppPeriod, nil
}

func (r *MPPPeriodRepository) Update(mppPeriod *entity.MPPPeriod) (*entity.MPPPeriod, error) {
	tx := r.DB.Begin()
	if tx.Error != nil {
		r.Log.Errorf("[MPPPeriodRepository.Update] " + tx.Error.Error())
		return nil, errors.New("[MPPPeriodRepository.Update] " + tx.Error.Error())
	}

	if mppPeriod.Status != "draft" {
		dateNow := time.Now().Format("2006-01-02")
		if dateNow < mppPeriod.StartDate.Format("2006-01-02") {
			mppPeriod.Status = entity.MPPeriodStatusNotOpen
		}
	}

	if err := tx.Model(mppPeriod).Where("id = ?", mppPeriod.ID).Updates(mppPeriod).Error; err != nil {
		tx.Rollback()
		r.Log.Errorf("[MPPPeriodRepository.Update] " + err.Error())
		return nil, errors.New("[MPPPeriodRepository.Update] " + err.Error())
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		r.Log.Errorf("[MPPPeriodRepository.Update] " + err.Error())
		return nil, errors.New("[MPPPeriodRepository.Update] " + err.Error())
	}

	return mppPeriod, nil
}

func (r *MPPPeriodRepository) Delete(id uuid.UUID) error {
	tx := r.DB.Begin()
	if tx.Error != nil {
		r.Log.Errorf("[MPPPeriodRepository.Delete] " + tx.Error.Error())
		return errors.New("[MPPPeriodRepository.Delete] " + tx.Error.Error())
	}

	if err := tx.Where("id = ?", id).Delete(&entity.MPPPeriod{}).Error; err != nil {
		tx.Rollback()
		r.Log.Errorf("[MPPPeriodRepository.Delete] " + err.Error())
		return errors.New("[MPPPeriodRepository.Delete] " + err.Error())
	}

	if err := tx.Commit().Error; err != nil {
		r.Log.Errorf("[MPPPeriodRepository.Delete] " + err.Error())
		return errors.New("[MPPPeriodRepository.Delete] " + err.Error())
	}

	return nil
}

func (r *MPPPeriodRepository) FindByCurrentDateAndStatus(status entity.MPPPeriodStatus) (*entity.MPPPeriod, error) {
	var mppPeriod entity.MPPPeriod

	err := r.DB.Where("status = ? AND start_date <= ? AND end_date >= ?", status, time.Now().Format("2006-01-02"), time.Now().Format("2006-01-02")).First(&mppPeriod).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			r.Log.Warn("[MPPPeriodRepository.FindByCurrentDateAndStatus] User not found")
			return nil, nil
		} else {
			r.Log.Error("[MPPPeriodRepository.FindByCurrentDateAndStatus] " + err.Error())
			return nil, errors.New("[UserRepository.FindByEmail] " + err.Error())
		}
	}

	return &mppPeriod, nil
}
