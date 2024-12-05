package repository

import (
	"errors"

	"github.com/IlhamSetiaji/julong-manpower-be/internal/entity"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type IMPPlanningRepository interface {
	FindAllHeadersPaginated(page int, pageSize int, search string) (*[]entity.MPPlanningHeader, int64, error)
	FindHeaderById(id int) (*entity.MPPlanningHeader, error)
	CreateHeader(mppHeader *entity.MPPlanningHeader) (*entity.MPPlanningHeader, error)
	UpdateHeader(mppHeader *entity.MPPlanningHeader) (*entity.MPPlanningHeader, error)
	DeleteHeader(id int) error
	FindAllLinesByHeaderId(headerId int) (*[]entity.MPPlanningLine, error)
	FindLineById(id int) (*entity.MPPlanningLine, error)
	CreateLine(mppLine *entity.MPPlanningLine) (*entity.MPPlanningLine, error)
	UpdateLine(mppLine *entity.MPPlanningLine) (*entity.MPPlanningLine, error)
	DeleteLine(id int) error
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

func (r *MPPlanningRepository) FindAllHeadersPaginated(page int, pageSize int, search string) (*[]entity.MPPlanningHeader, int64, error) {
	var mppHeaders []entity.MPPlanningHeader
	var total int64

	query := r.DB.Model(&entity.MPPlanningHeader{})

	if search != "" {
		query = query.Where("name LIKE ?", "%"+search+"%")
	}

	if err := query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&mppHeaders).Error; err != nil {
		r.Log.Errorf("[MPPlanningRepository.FindAllHeadersPaginated] " + err.Error())
		return nil, 0, errors.New("[MPPlanningRepository.FindAllHeadersPaginated] " + err.Error())
	}

	if err := query.Count(&total).Error; err != nil {
		r.Log.Errorf("[MPPlanningRepository.FindAllHeadersPaginated] " + err.Error())
		return nil, 0, errors.New("[MPPlanningRepository.FindAllHeadersPaginated] " + err.Error())
	}

	return &mppHeaders, total, nil
}

func (r *MPPlanningRepository) FindHeaderById(id int) (*entity.MPPlanningHeader, error) {
	var mppHeader entity.MPPlanningHeader

	if err := r.DB.Where("id = ?", id).First(&mppHeader).Error; err != nil {
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

func (r *MPPlanningRepository) DeleteHeader(id int) error {
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

func (r *MPPlanningRepository) FindAllLinesByHeaderId(headerId int) (*[]entity.MPPlanningLine, error) {
	var mppLines []entity.MPPlanningLine

	if err := r.DB.Where("mpp_header_id = ?", headerId).Find(&mppLines).Error; err != nil {
		r.Log.Errorf("[MPPlanningRepository.FindAllLinesByHeaderId] " + err.Error())
		return nil, errors.New("[MPPlanningRepository.FindAllLinesByHeaderId] " + err.Error())
	}

	return &mppLines, nil
}

func (r *MPPlanningRepository) FindLineById(id int) (*entity.MPPlanningLine, error) {
	var mppLine entity.MPPlanningLine

	if err := r.DB.Where("id = ?", id).First(&mppLine).Error; err != nil {
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

	if err := tx.Save(mppLine).Error; err != nil {
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

func (r *MPPlanningRepository) DeleteLine(id int) error {
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
