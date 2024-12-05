package usecase

import (
	"errors"
	"time"

	"github.com/IlhamSetiaji/julong-manpower-be/internal/entity"
	"github.com/IlhamSetiaji/julong-manpower-be/internal/http/request"
	"github.com/IlhamSetiaji/julong-manpower-be/internal/http/response"
	"github.com/IlhamSetiaji/julong-manpower-be/internal/repository"
	"github.com/sirupsen/logrus"
)

type IMPPPeriodUseCase interface {
	FindAllPaginated(request request.FindAllPaginatedMPPPeriodRequest) (*response.FindAllPaginatedMPPPeriodResponse, error)
	FindById(request request.FindByIdMPPPeriodRequest) (*response.FindByIdMPPPeriodResponse, error)
	Create(request request.CreateMPPPeriodRequest) (*response.CreateMPPPeriodResponse, error)
	Update(request request.UpdateMPPPeriodRequest) (*response.UpdateMPPPeriodResponse, error)
	Delete(request request.DeleteMPPPeriodRequest) error
}

type MPPPeriodUseCase struct {
	Log                 *logrus.Logger
	MPPPeriodRepository repository.IMPPPeriodRepository
}

func NewMPPPeriodUseCase(log *logrus.Logger, repo repository.IMPPPeriodRepository) IMPPPeriodUseCase {
	return &MPPPeriodUseCase{
		Log:                 log,
		MPPPeriodRepository: repo,
	}
}

func (uc *MPPPeriodUseCase) FindAllPaginated(req request.FindAllPaginatedMPPPeriodRequest) (*response.FindAllPaginatedMPPPeriodResponse, error) {
	mppPeriods, total, err := uc.MPPPeriodRepository.FindAllPaginated(req.Page, req.PageSize, req.Search)
	if err != nil {
		uc.Log.Errorf("[MPPPeriodUseCase.FindAllPaginated] " + err.Error())
		return nil, err
	}

	return &response.FindAllPaginatedMPPPeriodResponse{
		MPPPeriods: mppPeriods,
		Total:      total,
	}, nil
}

func (uc *MPPPeriodUseCase) FindById(req request.FindByIdMPPPeriodRequest) (*response.FindByIdMPPPeriodResponse, error) {
	mppPeriod, err := uc.MPPPeriodRepository.FindById(req.ID)
	if err != nil {
		uc.Log.Errorf("[MPPPeriodUseCase.FindById] " + err.Error())
		return nil, err
	}

	return &response.FindByIdMPPPeriodResponse{
		MPPPeriod: mppPeriod,
	}, nil
}

func (uc *MPPPeriodUseCase) Create(req request.CreateMPPPeriodRequest) (*response.CreateMPPPeriodResponse, error) {
	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		uc.Log.Errorf("[MPPPeriodUseCase.Create] " + err.Error())
		return nil, err
	}

	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		uc.Log.Errorf("[MPPPeriodUseCase.Create] " + err.Error())
		return nil, err
	}

	periodExist, err := uc.MPPPeriodRepository.FindByCurrentDateAndStatus(entity.MPPeriodStatusOpen)
	if err != nil {
		uc.Log.Errorf("[MPPPeriodUseCase.Create] " + err.Error())
		return nil, err
	}

	if periodExist != nil {
		uc.Log.Errorf("[MPPPeriodUseCase.Create] " + "MPP Period already exist")
		return nil, errors.New("MPP Period already exist")
	}

	mppPeriodEntity := &entity.MPPPeriod{
		Title:     req.Title,
		StartDate: startDate,
		EndDate:   endDate,
		Status:    req.Status,
	}

	mppPeriod, err := uc.MPPPeriodRepository.Create(mppPeriodEntity)
	if err != nil {
		uc.Log.Errorf("[MPPPeriodUseCase.Create] " + err.Error())
		return nil, err
	}

	return &response.CreateMPPPeriodResponse{
		ID:        mppPeriod.ID,
		Title:     mppPeriod.Title,
		StartDate: mppPeriod.StartDate.Format("2006-01-02"),
		EndDate:   mppPeriod.EndDate.Format("2006-01-02"),
		Status:    mppPeriod.Status,
	}, nil
}

func (uc *MPPPeriodUseCase) Update(req request.UpdateMPPPeriodRequest) (*response.UpdateMPPPeriodResponse, error) {
	exist, err := uc.MPPPeriodRepository.FindById(req.ID)
	if err != nil {
		uc.Log.Errorf("[MPPPeriodUseCase.Update] " + err.Error())
		return nil, err
	}

	if exist == nil {
		uc.Log.Errorf("[MPPPeriodUseCase.Update] " + "MPP Period not found")
		return nil, errors.New("MPP Period not found")
	}

	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		uc.Log.Errorf("[MPPPeriodUseCase.Create] " + err.Error())
		return nil, err
	}

	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		uc.Log.Errorf("[MPPPeriodUseCase.Create] " + err.Error())
		return nil, err
	}

	mppPeriodEntity := &entity.MPPPeriod{
		ID:        req.ID,
		Title:     req.Title,
		StartDate: startDate,
		EndDate:   endDate,
		Status:    req.Status,
	}

	mppPeriod, err := uc.MPPPeriodRepository.Update(mppPeriodEntity)
	if err != nil {
		uc.Log.Errorf("[MPPPeriodUseCase.Update] " + err.Error())
		return nil, err
	}

	return &response.UpdateMPPPeriodResponse{
		ID:        mppPeriod.ID,
		Title:     mppPeriod.Title,
		StartDate: mppPeriod.StartDate.Format("2006-01-02"),
		EndDate:   mppPeriod.EndDate.Format("2006-01-02"),
		Status:    mppPeriod.Status,
	}, nil
}

func (uc *MPPPeriodUseCase) Delete(req request.DeleteMPPPeriodRequest) error {
	exist, err := uc.MPPPeriodRepository.FindById(req.ID)
	if err != nil {
		uc.Log.Errorf("[MPPPeriodUseCase.Delete] " + err.Error())
		return err
	}

	if exist == nil {
		uc.Log.Errorf("[MPPPeriodUseCase.Delete] " + "MPP Period not found")
		return errors.New("MPP Period not found")
	}

	err = uc.MPPPeriodRepository.Delete(req.ID)
	if err != nil {
		uc.Log.Errorf("[MPPPeriodUseCase.Delete] " + err.Error())
		return err
	}

	return nil
}

func MPPPeriodUseCaseFactory(log *logrus.Logger) IMPPPeriodUseCase {
	repo := repository.MPPPeriodRepositoryFactory(log)
	return NewMPPPeriodUseCase(log, repo)
}
