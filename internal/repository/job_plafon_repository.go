package repository

import (
	"errors"

	"github.com/IlhamSetiaji/julong-manpower-be/internal/config"
	"github.com/IlhamSetiaji/julong-manpower-be/internal/entity"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type IJobPlafonRepository interface {
	FindAllPaginated(page int, pageSize int, search string, filter map[string]interface{}) (*[]entity.JobPlafon, int64, error)
	FindById(id uuid.UUID) (*entity.JobPlafon, error)
	FindByJobId(jobId uuid.UUID) (*entity.JobPlafon, error)
	Create(jobPlafon *entity.JobPlafon) (*entity.JobPlafon, error)
	Update(jobPlafon *entity.JobPlafon) (*entity.JobPlafon, error)
	Delete(id uuid.UUID) error
}

type JobPlafonRepository struct {
	Log *logrus.Logger
	DB  *gorm.DB
}

func NewJobPlafonRepository(log *logrus.Logger, db *gorm.DB) IJobPlafonRepository {
	return &JobPlafonRepository{
		Log: log,
		DB:  db,
	}
}

func (r *JobPlafonRepository) FindAllPaginated(page int, pageSize int, search string, filter map[string]interface{}) (*[]entity.JobPlafon, int64, error) {
	var jobPlafons []entity.JobPlafon
	var total int64
	var jobIDs []string

	query := r.DB.Model(&entity.JobPlafon{})

	if filter != nil {
		if _, ok := filter["job_ids"]; ok {
			jobIDs = filter["job_ids"].([]string)
			query = query.Where("job_id IN (?)", jobIDs)
		}
	}

	// if search != "" {
	// 	query = query.Where("name LIKE ?", "%"+search+"%")
	// }

	if err := query.Count(&total).Error; err != nil {
		r.Log.Errorf("[JobPlafonRepository.FindAllPaginated] " + err.Error())
		return nil, 0, err
	}

	if err := query.Order("plafon DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&jobPlafons).Error; err != nil {
		r.Log.Errorf("[JobPlafonRepository.FindAllPaginated] " + err.Error())
		return nil, 0, err
	}

	return &jobPlafons, total, nil
}

func (r *JobPlafonRepository) FindById(id uuid.UUID) (*entity.JobPlafon, error) {
	var jobPlafon entity.JobPlafon

	if err := r.DB.Where("id = ?", id).First(&jobPlafon).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			r.Log.Warn("[JobPlafonRepository.FindById] Job not found")
			return nil, nil
		} else {
			r.Log.Errorf("[JobPlafonRepository.FindById] " + err.Error())
			return nil, err
		}
	}

	return &jobPlafon, nil
}

func (r *JobPlafonRepository) FindByJobId(jobId uuid.UUID) (*entity.JobPlafon, error) {
	var jobPlafon entity.JobPlafon

	if err := r.DB.Where("job_id = ?", jobId).First(&jobPlafon).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			r.Log.Warn("[JobPlafonRepository.FindByJobId] Job not found")
			return nil, nil
		} else {
			r.Log.Errorf("[JobPlafonRepository.FindByJobId] " + err.Error())
			return nil, err
		}
	}

	return &jobPlafon, nil
}

func (r *JobPlafonRepository) Create(jobPlafon *entity.JobPlafon) (*entity.JobPlafon, error) {
	tx := r.DB.Begin()

	if tx.Error != nil {
		r.Log.Errorf("[JobPlafonRepository.Create] " + tx.Error.Error())
		return nil, errors.New("[JobPlafonRepository.Create] " + tx.Error.Error())
	}

	if err := tx.Create(jobPlafon).Error; err != nil {
		tx.Rollback()
		r.Log.Errorf("[JobPlafonRepository.Create] " + err.Error())
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		r.Log.Errorf("[JobPlafonRepository.Create] " + err.Error())
		return nil, errors.New("[JobPlafonRepository.Create] " + err.Error())
	}

	return jobPlafon, nil
}

func (r *JobPlafonRepository) Update(jobPlafon *entity.JobPlafon) (*entity.JobPlafon, error) {
	tx := r.DB.Begin()

	if tx.Error != nil {
		r.Log.Errorf("[JobPlafonRepository.Update] " + tx.Error.Error())
		return nil, errors.New("[JobPlafonRepository.Update] " + tx.Error.Error())
	}

	if err := tx.Model(&entity.JobPlafon{}).Where("id = ?", jobPlafon.ID).Updates(jobPlafon).Error; err != nil {
		tx.Rollback()
		r.Log.Errorf("[JobPlafonRepository.Update] " + err.Error())
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		r.Log.Errorf("[JobPlafonRepository.Update] " + err.Error())
		return nil, errors.New("[JobPlafonRepository.Update] " + err.Error())
	}

	return jobPlafon, nil
}

func (r *JobPlafonRepository) Delete(id uuid.UUID) error {
	tx := r.DB.Begin()

	if tx.Error != nil {
		r.Log.Errorf("[JobPlafonRepository.Delete] " + tx.Error.Error())
		return errors.New("[JobPlafonRepository.Delete] " + tx.Error.Error())
	}

	if err := tx.Where("id = ?", id).Delete(&entity.JobPlafon{}).Error; err != nil {
		tx.Rollback()
		r.Log.Errorf("[JobPlafonRepository.Delete] " + err.Error())
		return errors.New("[JobPlafonRepository.Delete] " + err.Error())
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		r.Log.Errorf("[JobPlafonRepository.Delete] " + err.Error())
		return errors.New("[JobPlafonRepository.Delete] " + err.Error())
	}

	return nil
}

func JobPlafonRepositoryFactory(log *logrus.Logger) IJobPlafonRepository {
	db := config.NewDatabase()
	return NewJobPlafonRepository(log, db)
}
