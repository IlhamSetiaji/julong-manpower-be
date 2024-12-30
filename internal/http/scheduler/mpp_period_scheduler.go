package scheduler

import (
	"time"

	"github.com/IlhamSetiaji/julong-manpower-be/internal/usecase"
	"github.com/sirupsen/logrus"
)

type IMPPPeriodScheduler interface {
	UpdateStatusToOpenByDate() error
	UpdateStatusToCloseByDate() error
}

type MPPPeriodScheduler struct {
	Log              *logrus.Logger
	MPPPeriodUseCase usecase.IMPPPeriodUseCase
}

func NewMPPPeriodScheduler(log *logrus.Logger, uc usecase.IMPPPeriodUseCase) IMPPPeriodScheduler {
	return &MPPPeriodScheduler{
		Log:              log,
		MPPPeriodUseCase: uc,
	}
}

func (s *MPPPeriodScheduler) UpdateStatusToOpenByDate() error {
	dateNow := time.Now()

	err := s.MPPPeriodUseCase.UpdateStatusToOpenByDate(dateNow)
	if err != nil {
		s.Log.Errorf("[MPPPeriodScheduler.UpdateStatusToOpenByDate] " + err.Error())
		return err
	}

	return nil
}

func (s *MPPPeriodScheduler) UpdateStatusToCloseByDate() error {
	dateNow := time.Now()

	err := s.MPPPeriodUseCase.UpdateStatusToCloseByDate(dateNow)
	if err != nil {
		s.Log.Errorf("[MPPPeriodScheduler.UpdateStatusToCloseByDate] " + err.Error())
		return err
	}

	return nil
}

func MPPPeriodSchedulerFactory(log *logrus.Logger) IMPPPeriodScheduler {
	uc := usecase.MPPPeriodUseCaseFactory(log)
	return NewMPPPeriodScheduler(log, uc)
}
