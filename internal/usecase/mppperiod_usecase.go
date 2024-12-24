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
	FindByStatus(status entity.MPPPeriodStatus) (*response.FindByCurrentDateAndStatusMPPPeriodResponse, error)
	Create(request request.CreateMPPPeriodRequest) (*response.CreateMPPPeriodResponse, error)
	Update(request request.UpdateMPPPeriodRequest) (*response.UpdateMPPPeriodResponse, error)
	Delete(request request.DeleteMPPPeriodRequest) error
	FindByCurrentDateAndStatus(request request.FindByCurrentDateAndStatusMPPPeriodRequest) (*response.FindByCurrentDateAndStatusMPPPeriodResponse, error)
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

func (uc *MPPPeriodUseCase) FindByStatus(status entity.MPPPeriodStatus) (*response.FindByCurrentDateAndStatusMPPPeriodResponse, error) {
	mppPeriod, err := uc.MPPPeriodRepository.FindByStatus(status)
	if err != nil {
		uc.Log.Errorf("[MPPPeriodUseCase.FindByStatus] " + err.Error())
		return nil, err
	}

	return &response.FindByCurrentDateAndStatusMPPPeriodResponse{
		MPPPeriod: mppPeriod,
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

	budgetStartDate, err := time.Parse("2006-01-02", req.BudgetStartDate)
	if err != nil {
		uc.Log.Errorf("[MPPPeriodUseCase.Create] " + err.Error())
		return nil, err
	}

	budgetEndDate, err := time.Parse("2006-01-02", req.BudgetEndDate)
	if err != nil {
		uc.Log.Errorf("[MPPPeriodUseCase.Create] " + err.Error())
		return nil, err
	}

	// periodExist, err := uc.MPPPeriodRepository.FindByCurrentDateAndStatus(entity.MPPeriodStatusOpen)
	// if err != nil {
	// 	uc.Log.Errorf("[MPPPeriodUseCase.Create] " + err.Error())
	// 	return nil, err
	// }
	periodExist, err := uc.MPPPeriodRepository.FindByStatus(entity.MPPeriodStatusOpen)
	if err != nil {
		uc.Log.Errorf("[MPPPeriodUseCase.Create] " + err.Error())
		return nil, err
	}

	if periodExist != nil {
		uc.Log.Errorf("[MPPPeriodUseCase.Create] " + "MPP Period already exist")
		return nil, errors.New("MPP Period already exist")
	}

	mppPeriodEntity := &entity.MPPPeriod{
		Title:           req.Title,
		StartDate:       startDate,
		EndDate:         endDate,
		BudgetStartDate: budgetStartDate,
		BudgetEndDate:   budgetEndDate,
		Status:          req.Status,
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
		uc.Log.Errorf("[MPPPeriodUseCase.Update] " + err.Error())
		return nil, err
	}

	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		uc.Log.Errorf("[MPPPeriodUseCase.Update] " + err.Error())
		return nil, err
	}

	if req.Status != entity.MPPPeriodStatusDraft {
		// periodExist, err := uc.MPPPeriodRepository.FindByCurrentDateAndStatus(entity.MPPeriodStatusOpen)
		// if err != nil {
		// 	uc.Log.Errorf("[MPPPeriodUseCase.Update] " + err.Error())
		// 	return nil, err
		// }
		periodExist, err := uc.MPPPeriodRepository.FindByStatus(entity.MPPeriodStatusOpen)
		if err != nil {
			uc.Log.Errorf("[MPPPeriodUseCase.Update] " + err.Error())
			return nil, err
		}

		if periodExist != nil {
			uc.Log.Errorf("[MPPPeriodUseCase.Update] " + "MPP Period already exist")
			return nil, errors.New("MPP Period already exist")
		}
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

func (uc *MPPPeriodUseCase) FindByCurrentDateAndStatus(req request.FindByCurrentDateAndStatusMPPPeriodRequest) (*response.FindByCurrentDateAndStatusMPPPeriodResponse, error) {
	mppPeriod, err := uc.MPPPeriodRepository.FindByCurrentDateAndStatus(req.Status)
	if err != nil {
		uc.Log.Errorf("[MPPPeriodUseCase.FindByCurrentDateAndStatus] " + err.Error())
		return nil, err
	}

	return &response.FindByCurrentDateAndStatusMPPPeriodResponse{
		MPPPeriod: mppPeriod,
	}, nil
}

func MPPPeriodUseCaseFactory(log *logrus.Logger) IMPPPeriodUseCase {
	repo := repository.MPPPeriodRepositoryFactory(log)
	return NewMPPPeriodUseCase(log, repo)
}
