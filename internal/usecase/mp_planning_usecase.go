package usecase

import (
	"errors"
	"time"

	"github.com/IlhamSetiaji/julong-manpower-be/internal/entity"
	"github.com/IlhamSetiaji/julong-manpower-be/internal/http/request"
	"github.com/IlhamSetiaji/julong-manpower-be/internal/http/response"
	"github.com/IlhamSetiaji/julong-manpower-be/internal/repository"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type IMPPlanningUseCase interface {
	FindAllHeadersPaginated(request *request.FindAllHeadersPaginatedMPPlanningRequest) (*response.FindAllHeadersPaginatedMPPlanningResponse, error)
	FindById(request *request.FindHeaderByIdMPPlanningRequest) (*response.FindByIdMPPlanningResponse, error)
	Create(request *request.CreateHeaderMPPlanningRequest) (*response.CreateMPPlanningResponse, error)
	Update(request *request.UpdateHeaderMPPlanningRequest) (*response.UpdateMPPlanningResponse, error)
	Delete(request *request.DeleteHeaderMPPlanningRequest) error
	FindAllLinesByHeaderIdPaginated(request *request.FindAllLinesByHeaderIdPaginatedMPPlanningLineRequest) (*response.FindAllLinesByHeaderIdPaginatedMPPlanningLineResponse, error)
	FindLineById(request *request.FindLineByIdMPPlanningLineRequest) (*response.FindByIdMPPlanningLineResponse, error)
	CreateLine(request *request.CreateLineMPPlanningLineRequest) (*response.CreateMPPlanningLineResponse, error)
	UpdateLine(request *request.UpdateLineMPPlanningLineRequest) (*response.UpdateMPPlanningLineResponse, error)
	DeleteLine(request *request.DeleteLineMPPlanningLineRequest) error
}

type MPPlanningUseCase struct {
	Log                  *logrus.Logger
	MPPlanningRepository repository.IMPPlanningRepository
}

func NewMPPlanningUseCase(log *logrus.Logger, repo repository.IMPPlanningRepository) IMPPlanningUseCase {
	return &MPPlanningUseCase{
		Log:                  log,
		MPPlanningRepository: repo,
	}
}

func (uc *MPPlanningUseCase) FindAllHeadersPaginated(req *request.FindAllHeadersPaginatedMPPlanningRequest) (*response.FindAllHeadersPaginatedMPPlanningResponse, error) {
	mpPlanningHeaders, total, err := uc.MPPlanningRepository.FindAllHeadersPaginated(req.Page, req.PageSize, req.Search)
	if err != nil {
		uc.Log.Errorf("[MPPlanningUseCase.FindAllHeadersPaginated] " + err.Error())
		return nil, err
	}

	return &response.FindAllHeadersPaginatedMPPlanningResponse{
		MPPlanningHeaders: mpPlanningHeaders,
		Total:             total,
	}, nil
}

func (uc *MPPlanningUseCase) FindById(req *request.FindHeaderByIdMPPlanningRequest) (*response.FindByIdMPPlanningResponse, error) {
	mpPlanningHeader, err := uc.MPPlanningRepository.FindHeaderById(uuid.MustParse(req.ID))
	if err != nil {
		uc.Log.Errorf("[MPPlanningUseCase.FindById] " + err.Error())
		return nil, err
	}

	if mpPlanningHeader == nil {
		uc.Log.Errorf("[MPPlanningUseCase.FindById] MP Planning Header not found")
		return nil, errors.New("MP Planning Header not found")
	}

	return &response.FindByIdMPPlanningResponse{
		MPPlanningHeader: mpPlanningHeader,
	}, nil
}

func (uc *MPPlanningUseCase) Create(req *request.CreateHeaderMPPlanningRequest) (*response.CreateMPPlanningResponse, error) {
	documentDate, err := time.Parse(time.RFC3339, req.DocumentDate)
	if err != nil {
		uc.Log.Errorf("[MPPlanningUseCase.Create] " + err.Error())
		return nil, err
	}

	mpPlanningHeader, err := uc.MPPlanningRepository.CreateHeader(&entity.MPPlanningHeader{
		MPPPeriodID:       req.MPPPeriodID,
		OrganizationID:    &req.OrganizationID,
		EmpOrganizationID: &req.EmpOrganizationID,
		DocumentNumber:    req.DocumentNumber,
		DocumentDate:      documentDate,
		Notes:             req.Notes,
		TotalRecruit:      req.TotalRecruit,
		TotalPromote:      req.TotalPromote,
		Status:            req.Status,
		RecommendedBy:     req.RecommendedBy,
		ApprovedBy:        req.ApprovedBy,
		RequestorID:       &req.RequestorID,
		NotesAttach:       req.NotesAttach,
	})
	if err != nil {
		uc.Log.Errorf("[MPPlanningUseCase.Create] " + err.Error())
		return nil, err
	}

	return &response.CreateMPPlanningResponse{
		ID:                mpPlanningHeader.ID.String(),
		MPPPeriodID:       mpPlanningHeader.MPPPeriodID.String(),
		OrganizationID:    mpPlanningHeader.OrganizationID.String(),
		EmpOrganizationID: mpPlanningHeader.EmpOrganizationID.String(),
		DocumentNumber:    mpPlanningHeader.DocumentNumber,
		DocumentDate:      mpPlanningHeader.DocumentDate,
		Notes:             mpPlanningHeader.Notes,
		TotalRecruit:      mpPlanningHeader.TotalRecruit,
		TotalPromote:      mpPlanningHeader.TotalPromote,
		Status:            mpPlanningHeader.Status,
		RecommendedBy:     mpPlanningHeader.RecommendedBy,
		ApprovedBy:        mpPlanningHeader.ApprovedBy,
		RequestorID:       mpPlanningHeader.RequestorID.String(),
		NotesAttach:       mpPlanningHeader.NotesAttach,
		CreatedAt:         mpPlanningHeader.CreatedAt,
		UpdatedAt:         mpPlanningHeader.UpdatedAt,
	}, nil
}

func (uc *MPPlanningUseCase) Update(req *request.UpdateHeaderMPPlanningRequest) (*response.UpdateMPPlanningResponse, error) {
	documentDate, err := time.Parse(time.RFC3339, req.DocumentDate)
	if err != nil {
		uc.Log.Errorf("[MPPlanningUseCase.Update] " + err.Error())
		return nil, err
	}

	mpPlanningHeader, err := uc.MPPlanningRepository.UpdateHeader(&entity.MPPlanningHeader{
		ID:                req.ID,
		MPPPeriodID:       req.MPPPeriodID,
		OrganizationID:    &req.OrganizationID,
		EmpOrganizationID: &req.EmpOrganizationID,
		DocumentNumber:    req.DocumentNumber,
		DocumentDate:      documentDate,
		Notes:             req.Notes,
		TotalRecruit:      req.TotalRecruit,
		TotalPromote:      req.TotalPromote,
		Status:            req.Status,
		RecommendedBy:     req.RecommendedBy,
		ApprovedBy:        req.ApprovedBy,
		RequestorID:       &req.RequestorID,
		NotesAttach:       req.NotesAttach,
	})
	if err != nil {
		uc.Log.Errorf("[MPPlanningUseCase.Update] " + err.Error())
		return nil, err
	}

	return &response.UpdateMPPlanningResponse{
		ID:                mpPlanningHeader.ID.String(),
		MPPPeriodID:       mpPlanningHeader.MPPPeriodID.String(),
		OrganizationID:    mpPlanningHeader.OrganizationID.String(),
		EmpOrganizationID: mpPlanningHeader.EmpOrganizationID.String(),
		DocumentNumber:    mpPlanningHeader.DocumentNumber,
		DocumentDate:      mpPlanningHeader.DocumentDate,
		Notes:             mpPlanningHeader.Notes,
		TotalRecruit:      mpPlanningHeader.TotalRecruit,
		TotalPromote:      mpPlanningHeader.TotalPromote,
		Status:            mpPlanningHeader.Status,
		RecommendedBy:     mpPlanningHeader.RecommendedBy,
		ApprovedBy:        mpPlanningHeader.ApprovedBy,
		RequestorID:       mpPlanningHeader.RequestorID.String(),
		NotesAttach:       mpPlanningHeader.NotesAttach,
		CreatedAt:         mpPlanningHeader.CreatedAt,
		UpdatedAt:         mpPlanningHeader.UpdatedAt,
	}, nil
}

func (uc *MPPlanningUseCase) Delete(req *request.DeleteHeaderMPPlanningRequest) error {
	err := uc.MPPlanningRepository.DeleteHeader(uuid.MustParse(req.ID))
	if err != nil {
		uc.Log.Errorf("[MPPlanningUseCase.Delete] " + err.Error())
		return err
	}

	return nil
}

func (uc *MPPlanningUseCase) FindAllLinesByHeaderIdPaginated(req *request.FindAllLinesByHeaderIdPaginatedMPPlanningLineRequest) (*response.FindAllLinesByHeaderIdPaginatedMPPlanningLineResponse, error) {
	mpPlanningLines, total, err := uc.MPPlanningRepository.FindAllLinesByHeaderIdPaginated(uuid.MustParse(req.HeaderID), req.Page, req.PageSize, req.Search)
	if err != nil {
		uc.Log.Errorf("[MPPlanningUseCase.FindAllLinesByHeaderIdPaginated] " + err.Error())
		return nil, err
	}

	return &response.FindAllLinesByHeaderIdPaginatedMPPlanningLineResponse{
		MPPlanningLines: mpPlanningLines,
		Total:           total,
	}, nil
}

func (uc *MPPlanningUseCase) FindLineById(req *request.FindLineByIdMPPlanningLineRequest) (*response.FindByIdMPPlanningLineResponse, error) {
	mpPlanningLine, err := uc.MPPlanningRepository.FindLineById(uuid.MustParse(req.ID))
	if err != nil {
		uc.Log.Errorf("[MPPlanningUseCase.FindLineById] " + err.Error())
		return nil, err
	}

	if mpPlanningLine == nil {
		uc.Log.Errorf("[MPPlanningUseCase.FindLineById] MP Planning Line not found")
		return nil, errors.New("MP Planning Line not found")
	}

	return &response.FindByIdMPPlanningLineResponse{
		MPPlanningLine: mpPlanningLine,
	}, nil
}

func (uc *MPPlanningUseCase) CreateLine(req *request.CreateLineMPPlanningLineRequest) (*response.CreateMPPlanningLineResponse, error) {
	mpPlanningLine, err := uc.MPPlanningRepository.CreateLine(&entity.MPPlanningLine{
		MPPlanningHeaderID:     req.MPPlanningHeaderID,
		OrganizationLocationID: &req.OrganizationLocationID,
		JobLevelID:             &req.JobLevelID,
		JobID:                  &req.JobID,
		Existing:               req.Existing,
		Recruit:                req.Recruit,
		SuggestedRecruit:       req.SuggestedRecruit,
		Promotion:              req.Promotion,
		Total:                  req.Total,
		RemainingBalance:       req.RemainingBalance,
		RecruitPH:              req.RecruitPH,
		RecruitMT:              req.RecruitMT,
	})
	if err != nil {
		uc.Log.Errorf("[MPPlanningUseCase.CreateLine] " + err.Error())
		return nil, err
	}

	return &response.CreateMPPlanningLineResponse{
		ID:                     mpPlanningLine.ID.String(),
		MPPlanningHeaderID:     mpPlanningLine.MPPlanningHeaderID.String(),
		OrganizationLocationID: mpPlanningLine.OrganizationLocationID.String(),
		JobLevelID:             mpPlanningLine.JobLevelID.String(),
		JobID:                  mpPlanningLine.JobID.String(),
		Existing:               mpPlanningLine.Existing,
		Recruit:                mpPlanningLine.Recruit,
		SuggestedRecruit:       mpPlanningLine.SuggestedRecruit,
		Promotion:              mpPlanningLine.Promotion,
		Total:                  mpPlanningLine.Total,
		RemainingBalance:       mpPlanningLine.RemainingBalance,
		RecruitPH:              mpPlanningLine.RecruitPH,
		RecruitMT:              mpPlanningLine.RecruitMT,
		CreatedAt:              mpPlanningLine.CreatedAt,
		UpdatedAt:              mpPlanningLine.UpdatedAt,
	}, nil
}

func (uc *MPPlanningUseCase) UpdateLine(req *request.UpdateLineMPPlanningLineRequest) (*response.UpdateMPPlanningLineResponse, error) {
	mpPlanningLine, err := uc.MPPlanningRepository.UpdateLine(&entity.MPPlanningLine{
		ID:                     req.ID,
		MPPlanningHeaderID:     req.MPPlanningHeaderID,
		OrganizationLocationID: &req.OrganizationLocationID,
		JobLevelID:             &req.JobLevelID,
		JobID:                  &req.JobID,
		Existing:               req.Existing,
		Recruit:                req.Recruit,
		SuggestedRecruit:       req.SuggestedRecruit,
		Promotion:              req.Promotion,
		Total:                  req.Total,
		RemainingBalance:       req.RemainingBalance,
		RecruitPH:              req.RecruitPH,
		RecruitMT:              req.RecruitMT,
	})
	if err != nil {
		uc.Log.Errorf("[MPPlanningUseCase.UpdateLine] " + err.Error())
		return nil, err
	}

	return &response.UpdateMPPlanningLineResponse{
		ID:                     mpPlanningLine.ID.String(),
		MPPlanningHeaderID:     mpPlanningLine.MPPlanningHeaderID.String(),
		OrganizationLocationID: mpPlanningLine.OrganizationLocationID.String(),
		JobLevelID:             mpPlanningLine.JobLevelID.String(),
		JobID:                  mpPlanningLine.JobID.String(),
		Existing:               mpPlanningLine.Existing,
		Recruit:                mpPlanningLine.Recruit,
		SuggestedRecruit:       mpPlanningLine.SuggestedRecruit,
		Promotion:              mpPlanningLine.Promotion,
		Total:                  mpPlanningLine.Total,
		RemainingBalance:       mpPlanningLine.RemainingBalance,
		RecruitPH:              mpPlanningLine.RecruitPH,
		RecruitMT:              mpPlanningLine.RecruitMT,
		CreatedAt:              mpPlanningLine.CreatedAt,
		UpdatedAt:              mpPlanningLine.UpdatedAt,
	}, nil
}

func (uc *MPPlanningUseCase) DeleteLine(req *request.DeleteLineMPPlanningLineRequest) error {
	err := uc.MPPlanningRepository.DeleteLine(uuid.MustParse(req.ID))
	if err != nil {
		uc.Log.Errorf("[MPPlanningUseCase.DeleteLine] " + err.Error())
		return err
	}

	return nil
}

func MPPlanningUseCaseFactory(log *logrus.Logger) IMPPlanningUseCase {
	repo := repository.MPPlanningRepositoryFactory(log)
	return NewMPPlanningUseCase(log, repo)
}
