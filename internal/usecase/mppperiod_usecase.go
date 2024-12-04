package usecase

import (
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
	mppPeriodEntity := &entity.MPPPeriod{
		Title:     req.Title,
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
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
	mppPeriodEntity := &entity.MPPPeriod{
		ID:        req.ID,
		Title:     req.Title,
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
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
	err := uc.MPPPeriodRepository.Delete(req.ID)
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
