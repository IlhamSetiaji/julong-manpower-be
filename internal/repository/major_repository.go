package repository

import (
	"errors"

	"github.com/IlhamSetiaji/julong-manpower-be/internal/config"
	"github.com/IlhamSetiaji/julong-manpower-be/internal/entity"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type IMajorRepository interface {
	FindAll() (*[]entity.Major, error)
	FindById(id uuid.UUID) (*entity.Major, error)
}

type MajorRepository struct {
	Log *logrus.Logger
	DB  *gorm.DB
}

func NewMajorRepository(log *logrus.Logger, db *gorm.DB) *MajorRepository {
	return &MajorRepository{Log: log, DB: db}
}

func (r *MajorRepository) FindAll() (*[]entity.Major, error) {
	var majors []entity.Major
	err := r.DB.Find(&majors).Error
	if err != nil {
		return nil, err
	}
	return &majors, nil
}

func (r *MajorRepository) FindById(id uuid.UUID) (*entity.Major, error) {
	var major entity.Major
	if err := r.DB.First(&major, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			r.Log.Errorf("[MajorRepository.FindById] Major with ID %s not found", id)
			return nil, nil
		} else {
			r.Log.Errorf("[MajorRepository.FindById] %s", err.Error())
			return nil, errors.New("[MajorRepository.FindById] Internal server error")
		}
	}

	return &major, nil
}

func MajorRepositoryFactory(log *logrus.Logger) *MajorRepository {
	db := config.NewDatabase()
	return NewMajorRepository(log, db)
}
