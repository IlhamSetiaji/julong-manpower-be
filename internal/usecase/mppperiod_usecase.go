package usecase

import (
	"github.com/IlhamSetiaji/julong-manpower-be/internal/repository"
	"github.com/IlhamSetiaji/julong-manpower-be/internal/http/request"
	"github.com/IlhamSetiaji/julong-manpower-be/internal/http/response"
	"github.com/sirupsen/logrus"
)

type IMPPPeriodUseCase interface {
	FindAllPaginated(request request.FindAllPaginatedMPPPeriodRequest) (*response.FindAllPaginatedMPPPeriodResponse, error)
	FindById(request request.FindByIdMPPPeriodRequest) (*response.FindByIdMPPPeriodResponse, error)
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

func MPPPeriodUseCaseFactory(log *logrus.Logger) IMPPPeriodUseCase {
	repo := repository.MPPPeriodRepositoryFactory(log)
	return NewMPPPeriodUseCase(log, repo)
}
