package repository

import (
	"errors"

	"github.com/IlhamSetiaji/julong-manpower-be/internal/config"
	"github.com/IlhamSetiaji/julong-manpower-be/internal/entity"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type IMPPPeriodRepository interface {
	FindAllPaginated(page int, pageSize int, search string) (*[]entity.MPPPeriod, int64, error)
	FindById(id uuid.UUID) (*entity.MPPPeriod, error)
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
		r.Log.Errorf("[MPPPeriodRepository.FindById] " + err.Error())
		return nil, errors.New("[MPPPeriodRepository.FindById] " + err.Error())
	}

	return &mppPeriod, nil
}

func MPPPeriodRepositoryFactory(log *logrus.Logger) IMPPPeriodRepository {
	db := config.NewDatabase()
	return NewMPPPeriodRepository(log, db)
}
