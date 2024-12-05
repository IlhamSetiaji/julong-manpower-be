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
	FindAllPaginated(page int, pageSize int, search string) (*[]entity.JobPlafon, int64, error)
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

func (r *JobPlafonRepository) FindAllPaginated(page int, pageSize int, search string) (*[]entity.JobPlafon, int64, error) {
	var jobPlafons []entity.JobPlafon
	var total int64

	query := r.DB.Model(&entity.JobPlafon{})

	if search != "" {
		query = query.Where("name LIKE ?", "%"+search+"%")
	}

	if err := query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&jobPlafons).Error; err != nil {
		r.Log.Errorf("[JobPlafonRepository.FindAllPaginated] " + err.Error())
		return nil, 0, err
	}

	if err := query.Count(&total).Error; err != nil {
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
	if err := r.DB.Create(jobPlafon).Error; err != nil {
		r.Log.Errorf("[JobPlafonRepository.Create] " + err.Error())
		return nil, err
	}

	return jobPlafon, nil
}

func (r *JobPlafonRepository) Update(jobPlafon *entity.JobPlafon) (*entity.JobPlafon, error) {
	if err := r.DB.Save(jobPlafon).Error; err != nil {
		r.Log.Errorf("[JobPlafonRepository.Update] " + err.Error())
		return nil, err
	}

	return jobPlafon, nil
}

func (r *JobPlafonRepository) Delete(id uuid.UUID) error {
	if err := r.DB.Where("id = ?", id).Delete(&entity.JobPlafon{}).Error; err != nil {
		r.Log.Errorf("[JobPlafonRepository.Delete] " + err.Error())
		return err
	}

	return nil
}

func JobPlafonRepositoryFactory(log *logrus.Logger) IJobPlafonRepository {
	db := config.NewDatabase()
	return NewJobPlafonRepository(log, db)
}
