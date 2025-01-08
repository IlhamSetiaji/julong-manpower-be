package usecase

import (
	"errors"
	"fmt"
	"time"

	"github.com/IlhamSetiaji/julong-manpower-be/internal/entity"
	"github.com/IlhamSetiaji/julong-manpower-be/internal/http/dto"
	"github.com/IlhamSetiaji/julong-manpower-be/internal/http/messaging"
	"github.com/IlhamSetiaji/julong-manpower-be/internal/http/request"
	"github.com/IlhamSetiaji/julong-manpower-be/internal/http/response"
	"github.com/IlhamSetiaji/julong-manpower-be/internal/repository"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type IBatchUsecase interface {
	CreateBatchHeaderAndLines(req *request.CreateBatchHeaderAndLinesRequest) (*response.BatchResponse, error)
	FindByStatus(status entity.BatchHeaderApprovalStatus, approverType string, orgID string) (*response.BatchResponse, error)
	FindById(id string) (*response.BatchResponse, error)
	GetOrganizationsForBatchApproval(id string) (*[]response.OrganizationResponse, error)
	FindDocumentByID(id string) (*response.RealDocumentBatchResponse, error)
	FindByNeedApproval(approverType string, orgID string) (*response.RealDocumentBatchResponse, error)
	FindByCurrentDocumentDateAndStatus(status entity.BatchHeaderApprovalStatus) (*response.BatchResponse, error)
	UpdateStatusBatchHeader(req *request.UpdateStatusBatchHeaderRequest) (*response.BatchResponse, error)
	GetCompletedBatchHeader() (*[]response.CompletedBatchResponse, error)
}

type BatchUsecase struct {
	Viper          *viper.Viper
	Log            *logrus.Logger
	Repo           repository.IBatchRepository
	OrgMessage     messaging.IOrganizationMessage
	EmpMessage     messaging.IEmployeeMessage
	batchDTO       dto.IBatchDTO
	mpPlanningRepo repository.IMPPlanningRepository
}

func NewBatchUsecase(viper *viper.Viper, log *logrus.Logger, repo repository.IBatchRepository, orgMessage messaging.IOrganizationMessage, empMessage messaging.IEmployeeMessage, batchDTO dto.IBatchDTO, mpPlanningRepo repository.IMPPlanningRepository) IBatchUsecase {
	return &BatchUsecase{
		Viper:          viper,
		Log:            log,
		Repo:           repo,
		OrgMessage:     orgMessage,
		EmpMessage:     empMessage,
		batchDTO:       batchDTO,
		mpPlanningRepo: mpPlanningRepo,
	}
}

func (uc *BatchUsecase) GetCompletedBatchHeader() (*[]response.CompletedBatchResponse, error) {
	batchHeaders, err := uc.Repo.GetBatchHeadersByStatus(entity.BatchHeaderApprovalStatusCompleted)
	if err != nil {
		uc.Log.Errorf("[BatchUsecase.GetCompletedBatchHeader] " + err.Error())
		return nil, err
	}

	if len(batchHeaders) == 0 {
		return nil, nil
	}

	completedBatchResponses := make([]response.CompletedBatchResponse, len(batchHeaders))
	// get one mp planning header
	mpPlanningHeader, err := uc.mpPlanningRepo.FindHeaderById(batchHeaders[0].BatchLines[0].MPPlanningHeaderID)
	if err != nil {
		uc.Log.Errorf("[BatchUsecase.GetCompletedBatchHeader] " + err.Error())
		return nil, err
	}

	// embed batch headers to completed batch responses
	for i, bh := range batchHeaders {
		completedBatchResponses[i] = response.CompletedBatchResponse{
			ID:             bh.ID,
			DocumentNumber: bh.DocumentNumber,
			DocumentDate:   bh.DocumentDate,
			CreatedAt:      bh.CreatedAt,
			UpdatedAt:      bh.UpdatedAt,
			MPPPeriod: response.MPPeriodResponse{
				ID:              mpPlanningHeader.MPPPeriod.ID,
				Title:           mpPlanningHeader.MPPPeriod.Title,
				StartDate:       mpPlanningHeader.MPPPeriod.StartDate.Format("2006-01-02"),
				EndDate:         mpPlanningHeader.MPPPeriod.EndDate.Format("2006-01-02"),
				BudgetStartDate: mpPlanningHeader.MPPPeriod.BudgetStartDate.Format("2006-01-02"),
				BudgetEndDate:   mpPlanningHeader.MPPPeriod.BudgetEndDate.Format("2006-01-02"),
				Status:          mpPlanningHeader.MPPPeriod.Status,
				CreatedAt:       mpPlanningHeader.MPPPeriod.CreatedAt,
				UpdatedAt:       mpPlanningHeader.MPPPeriod.UpdatedAt,
			},
		}
	}

	return &completedBatchResponses, nil
}

func (uc *BatchUsecase) GetOrganizationsForBatchApproval(id string) (*[]response.OrganizationResponse, error) {
	batchHeader, err := uc.Repo.FindById(id)
	if err != nil {
		uc.Log.Errorf("[BatchUsecase.GetOrganizationsForBatchApproval] " + err.Error())
		return nil, err
	}

	if batchHeader == nil {
		return nil, errors.New("Batch not found")
	}

	orgIds := make([]string, len(batchHeader.BatchLines))
	for i, bl := range batchHeader.BatchLines {
		orgIds[i] = bl.OrganizationID.String()
	}

	ogrs, err := uc.OrgMessage.SendFindAllOrganizationMessage(orgIds)

	if err != nil {
		uc.Log.Errorf("[BatchUsecase.GetOrganizationsForBatchApproval] " + err.Error())
		return nil, err
	}

	if ogrs == nil {
		return nil, errors.New("Organization not found")
	}

	organizationResponses := make([]response.OrganizationResponse, len(*ogrs))
	for i, ogr := range *ogrs {
		organizationResponses[i] = response.OrganizationResponse{
			ID:                 ogr.ID,
			OrganizationTypeID: ogr.OrganizationTypeID,
			Name:               ogr.Name,
		}
	}

	return &organizationResponses, nil
}

func (uc *BatchUsecase) CreateBatchHeaderAndLines(req *request.CreateBatchHeaderAndLinesRequest) (*response.BatchResponse, error) {
	dateNow := time.Now()
	documentNumber := ""
	approverType := entity.BatchHeaderApproverTypeCEO

	foundBatchHeader, err := uc.Repo.GetHeadersByDocumentDate(dateNow.Format("2006-01-02"))
	if err != nil {
		uc.Log.Errorf("[MPPlanningUseCase.GenerateDocumentNumber] " + err.Error())
		return nil, err
	}

	if req.ApproverType != "" {
		if req.ApproverType == entity.BatchHeaderApproverTypeCEO {
			if foundBatchHeader == nil {
				documentNumber = "MPP/BATCH/" + dateNow.Format("20060102") + "/001"
			} else {
				documentNumber = "MPP/BATCH/" + dateNow.Format("20060102") + "/" + fmt.Sprintf("%03d", len(*&foundBatchHeader)+1)
			}
		} else {
			if foundBatchHeader == nil {
				documentNumber = "MPP/BATCH/DIR/" + dateNow.Format("20060102") + "/001"
			} else {
				documentNumber = "MPP/BATCH/DIR/" + dateNow.Format("20060102") + "/" + fmt.Sprintf("%03d", len(*&foundBatchHeader)+1)
			}
			approverType = entity.BatchHeaderApproverTypeDirector
		}
	} else {
		if foundBatchHeader == nil {
			documentNumber = "MPP/BATCH/" + dateNow.Format("20060102") + "/001"
		} else {
			documentNumber = "MPP/BATCH/" + dateNow.Format("20060102") + "/" + fmt.Sprintf("%03d", len(*&foundBatchHeader)+1)
		}
	}

	var batchHeader *entity.BatchHeader

	var orgID uuid.UUID
	if req.OrganizationID != "" {
		orgID = uuid.MustParse(req.OrganizationID)
	}
	if req.DocumentNumber != "" {
		batchHeader = &entity.BatchHeader{
			DocumentNumber: req.DocumentNumber,
			DocumentDate:   dateNow,
			Status:         req.Status,
			ApproverType:   approverType,
			OrganizationID: &orgID,
		}
	} else {
		batchHeader = &entity.BatchHeader{
			DocumentNumber: documentNumber,
			DocumentDate:   dateNow,
			Status:         entity.BatchHeaderApprovalStatusNeedApproval,
			ApproverType:   approverType,
			OrganizationID: &orgID,
		}
	}

	// batchLines := make([]entity.BatchLine, len(req.BatchLines))
	var batchLines []entity.BatchLine
	for _, bl := range req.BatchLines {
		mpPlanningHeaderExist, err := uc.mpPlanningRepo.FindHeaderById(uuid.MustParse(bl.MPPlanningHeaderID))
		if err != nil {
			uc.Log.Errorf("[MPPlanningUseCase.CreateOrUpdateBatchLineMPPlanningLines] " + err.Error())
			return nil, err
		}

		if mpPlanningHeaderExist == nil {
			uc.Log.Errorf("[MPPlanningUseCase.CreateOrUpdateBatchLineMPPlanningLines] MP Planning Header not found")
			return nil, errors.New("MP Planning Header not found: " + bl.MPPlanningHeaderID)
		}

		if approverType == entity.BatchHeaderApproverTypeDirector {
			if mpPlanningHeaderExist.OrganizationID.String() != req.OrganizationID {
				uc.Log.Errorf("[MPPlanningUseCase.CreateOrUpdateBatchLineMPPlanningLines] Organization ID not match")
				return nil, errors.New("Organization ID not match")
			}

			if mpPlanningHeaderExist.Status != entity.MPPlaningStatusNeedApproval && mpPlanningHeaderExist.RecommendedBy != "" {
				uc.Log.Errorf("[MPPlanningUseCase.CreateOrUpdateBatchLineMPPlanningLines] MP Planning Header not in Need Approval status")
				return nil, errors.New("MP Planning Header not in Need Approval status")
			}
		}
		mpLine, err := uc.mpPlanningRepo.FindLineByHeaderID(uuid.MustParse(bl.MPPlanningHeaderID))
		if err != nil {
			uc.Log.Errorf("[MPPlanningUseCase.CreateOrUpdateBatchLineMPPlanningLines] " + err.Error())
			return nil, err
		}

		if mpLine == nil {
			uc.Log.Errorf("[MPPlanningUseCase.CreateOrUpdateBatchLineMPPlanningLines] MP Planning Line not found")
			return nil, errors.New("MP Planning Line not found")
		}

		if mpLine.OrganizationLocationID != nil {
			// Check if organization location exist
			orgLocExist, err := uc.OrgMessage.SendFindOrganizationLocationByIDMessage(request.SendFindOrganizationLocationByIDMessageRequest{
				ID: mpLine.OrganizationLocationID.String(),
			})

			if err != nil {
				uc.Log.Errorf("[MPPlanningUseCase.CreateOrUpdateBatchLineMPPlanningLines] " + err.Error())
				return nil, err
			}

			if orgLocExist == nil {
				uc.Log.Errorf("[MPPlanningUseCase.CreateOrUpdateBatchLineMPPlanningLines] Organization Location not found")
				return nil, errors.New("Organization Location not found")
			}

		}

		if mpLine.MPPlanningHeader.OrganizationID != nil {
			// Check if organization exist
			orgExist, err := uc.OrgMessage.SendFindOrganizationByIDMessage(request.SendFindOrganizationByIDMessageRequest{
				ID: mpLine.MPPlanningHeader.OrganizationID.String(),
			})

			if err != nil {
				uc.Log.Errorf("[MPPlanningUseCase.CreateOrUpdateBatchLineMPPlanningLines] " + err.Error())
				return nil, err
			}

			if orgExist == nil {
				uc.Log.Errorf("[MPPlanningUseCase.CreateOrUpdateBatchLineMPPlanningLines] Organization not found")
				return nil, errors.New("Organization not found")
			}

		}

		mpHeaderByStatus, err := uc.mpPlanningRepo.FindHeaderById(uuid.MustParse(bl.MPPlanningHeaderID))

		if err != nil {
			uc.Log.Errorf("[BatchUsecase.CreateBatchHeaderAndLines] " + err.Error())
			return nil, err
		}

		if approverType == entity.BatchHeaderApproverTypeCEO {
			if mpHeaderByStatus.Status != entity.MPPlaningStatusApproved {
				uc.Log.Errorf("[BatchUsecase.CreateBatchHeaderAndLines] MP Planning Header not in Approved status")
				continue
			}
		} else {
			if mpHeaderByStatus.Status != entity.MPPlaningStatusNeedApproval {
				uc.Log.Errorf("[BatchUsecase.CreateBatchHeaderAndLines] MP Planning Header not in Need Approval status")
				continue
			}
		}

		batchLines = append(batchLines, entity.BatchLine{
			BatchHeaderID:          batchHeader.ID,
			MPPlanningHeaderID:     uuid.MustParse(bl.MPPlanningHeaderID),
			OrganizationID:         *mpLine.MPPlanningHeader.OrganizationID,
			OrganizationLocationID: *mpLine.OrganizationLocationID,
		})
	}

	// uc.Log.Infof("batchLines hahahahaha: %+v", batchLines[0].OrganizationID)
	var batchHeaderExists = &entity.BatchHeader{}
	if approverType == entity.BatchHeaderApproverTypeCEO {
		batchHeaderExists, err = uc.Repo.FindByStatus(entity.BatchHeaderApprovalStatusNeedApproval, string(approverType), "")
		if err != nil {
			uc.Log.Errorf("[BatchUsecase.CreateBatchHeaderAndLines] " + err.Error())
			return nil, err
		}
	} else {
		batchHeaderExists, err = uc.Repo.FindByStatus(entity.BatchHeaderApprovalStatusNeedApproval, string(approverType), req.OrganizationID)
		if err != nil {
			uc.Log.Errorf("[BatchUsecase.CreateBatchHeaderAndLines] " + err.Error())
			return nil, err
		}
	}

	if batchHeaderExists != nil {
		err = uc.Repo.InsertLinesByBatchHeaderID(batchHeaderExists.ID.String(), batchLines)
		if err != nil {
			uc.Log.Errorf("[BatchUsecase.CreateBatchHeaderAndLines] " + err.Error())
			return nil, err
		}

		err = uc.Repo.DeleteLinesNotInBatchLines(batchHeaderExists.ID.String(), batchLines)
		if err != nil {
			uc.Log.Errorf("[BatchUsecase.CreateBatchHeaderAndLines] " + err.Error())
			return nil, err
		}

		findBatchHeader, err := uc.Repo.FindById(batchHeaderExists.ID.String())
		if err != nil {
			uc.Log.Errorf("[BatchUsecase.CreateBatchHeaderAndLines] " + err.Error())
			return nil, err
		}

		if approverType == entity.BatchHeaderApproverTypeCEO {
			err = uc.updateMpPlanningHeaderStatus(batchLines, uuid.MustParse(req.ApproverID), req.ApproverName)

			if err != nil {
				uc.Log.Errorf("[BatchUsecase.CreateBatchHeaderAndLines] " + err.Error())
				return nil, err
			}
		} else {
			err = uc.updateMpPlanningHeaderStatusDirector(batchLines, uuid.MustParse(req.ApproverID), req.ApproverName)

			if err != nil {
				uc.Log.Errorf("[BatchUsecase.CreateBatchHeaderAndLines] " + err.Error())
				return nil, err
			}
		}

		return uc.batchDTO.ConvertBatchHeaderEntityToResponse(findBatchHeader), nil
	}

	resp, err := uc.Repo.CreateBatchHeaderAndLines(batchHeader, batchLines)
	if err != nil {
		return nil, err
	}

	if approverType == entity.BatchHeaderApproverTypeCEO {
		err = uc.updateMpPlanningHeaderStatus(batchLines, uuid.MustParse(req.ApproverID), req.ApproverName)

		if err != nil {
			uc.Log.Errorf("[BatchUsecase.CreateBatchHeaderAndLines] " + err.Error())
			return nil, err
		}
	} else {
		err = uc.updateMpPlanningHeaderStatusDirector(batchLines, uuid.MustParse(req.ApproverID), req.ApproverName)

		if err != nil {
			uc.Log.Errorf("[BatchUsecase.CreateBatchHeaderAndLines] " + err.Error())
			return nil, err
		}
	}

	return uc.batchDTO.ConvertBatchHeaderEntityToResponse(resp), nil
}

func (uc *BatchUsecase) updateMpPlanningHeaderStatusDirector(batchLines []entity.BatchLine, approverID uuid.UUID, approverName string) error {
	for _, bl := range batchLines {
		approvalHistory := &entity.MPPlanningApprovalHistory{
			MPPlanningHeaderID: bl.MPPlanningHeaderID,
			Notes:              "",
			ApproverID:         approverID,
			ApproverName:       approverName,
			Level:              string(entity.MPPlanningApprovalHistoryLevelDirekturUnit),
			Status:             entity.MPPlanningApprovalHistoryStatusNeedApproval,
		}

		uc.Log.Infof("approver id direktur: %s", approverID.String())

		err := uc.mpPlanningRepo.UpdateStatusHeader(bl.MPPlanningHeaderID, string(entity.MPPlanningApprovalHistoryStatusNeedApproval), "", approvalHistory)
		if err != nil {
			uc.Log.Errorf("[MPPlanningUseCase.UpdateStatusMPPlanningHeader] " + err.Error())
			return err
		}
	}

	return nil
}

func (uc *BatchUsecase) updateMpPlanningHeaderStatus(batchLines []entity.BatchLine, approverID uuid.UUID, approverName string) error {
	for _, bl := range batchLines {
		approvalHistory := &entity.MPPlanningApprovalHistory{
			MPPlanningHeaderID: bl.MPPlanningHeaderID,
			ApproverID:         approverID,
			ApproverName:       approverName,
			Notes:              "",
			Level:              string(entity.MPPlanningApprovalHistoryLevelRecruitment),
			Status:             entity.MPPlanningApprovalHistoryStatusNeedApproval,
		}

		uc.Log.Infof("approver id: %s", approverID.String())

		err := uc.mpPlanningRepo.UpdateStatusHeader(bl.MPPlanningHeaderID, string(entity.MPPlanningApprovalHistoryStatusNeedApproval), approverID.String(), approvalHistory)
		if err != nil {
			uc.Log.Errorf("[MPPlanningUseCase.UpdateStatusMPPlanningHeader] " + err.Error())
			return err
		}
	}

	return nil
}

func (uc *BatchUsecase) FindByStatus(status entity.BatchHeaderApprovalStatus, approverType string, orgID string) (*response.BatchResponse, error) {
	resp, err := uc.Repo.FindByStatus(status, approverType, orgID)
	if err != nil {
		return nil, err
	}

	if resp == nil {
		return nil, errors.New("Batch not found")
	}

	return uc.batchDTO.ConvertBatchHeaderEntityToResponse(resp), nil
}

func (uc *BatchUsecase) FindById(id string) (*response.BatchResponse, error) {
	resp, err := uc.Repo.FindById(id)
	if err != nil {
		return nil, err
	}

	if resp == nil {
		return nil, errors.New("Batch not found")
	}

	return uc.batchDTO.ConvertBatchHeaderEntityToResponse(resp), nil
}

func (uc *BatchUsecase) UpdateStatusBatchHeader(req *request.UpdateStatusBatchHeaderRequest) (*response.BatchResponse, error) {
	batchHeader, err := uc.Repo.FindById(req.ID)
	if err != nil {
		uc.Log.Errorf("[BatchUsecase.UpdateStatusBatchHeader] " + err.Error())
		return nil, err
	}

	if batchHeader == nil {
		return nil, errors.New("Batch not found")
	}

	if req.ApproverType == "" || req.ApproverType == entity.BatchHeaderApproverTypeCEO {
		err = uc.Repo.UpdateStatusBatchHeader(batchHeader, req.Status, req.ApprovedBy, req.ApproverName)
		if err != nil {
			uc.Log.Errorf("[BatchUsecase.UpdateStatusBatchHeader] " + err.Error())
			return nil, err
		}
	} else {
		err = uc.Repo.UpdateStatusBatchHeaderDirector(batchHeader, req.Status, req.ApprovedBy, req.ApproverName)
		if err != nil {
			uc.Log.Errorf("[BatchUsecase.UpdateStatusBatchHeader] " + err.Error())
			return nil, err
		}
	}

	return uc.batchDTO.ConvertBatchHeaderEntityToResponse(batchHeader), nil
}

func (uc *BatchUsecase) FindDocumentByID(id string) (*response.RealDocumentBatchResponse, error) {
	resp, err := uc.Repo.FindById(id)
	if err != nil {
		return nil, err
	}

	if resp == nil {
		return nil, errors.New("Batch not found")
	}

	return uc.batchDTO.ConvertRealDocumentBatchResponse(resp), nil
}

func (uc *BatchUsecase) FindByNeedApproval(approverType string, orgID string) (*response.RealDocumentBatchResponse, error) {
	resp, err := uc.Repo.FindByNeedApproval(approverType, orgID)
	if err != nil {
		return nil, err
	}

	if resp == nil {
		return nil, errors.New("Batch not found")
	}

	return uc.batchDTO.ConvertRealDocumentBatchResponse(resp), nil
}

func (uc *BatchUsecase) FindByCurrentDocumentDateAndStatus(status entity.BatchHeaderApprovalStatus) (*response.BatchResponse, error) {
	resp, err := uc.Repo.FindByCurrentDocumentDateAndStatus(status)
	if err != nil {
		return nil, err
	}

	return uc.batchDTO.ConvertBatchHeaderEntityToResponse(resp), nil
}

func BatchUsecaseFactory(viper *viper.Viper, log *logrus.Logger) IBatchUsecase {
	repo := repository.BatchRepositoryFactory(log)
	orgMessage := messaging.OrganizationMessageFactory(log)
	empMessage := messaging.EmployeeMessageFactory(log)
	batchDTO := dto.BatchDTOFactory(log)
	mpPlanningRepo := repository.MPPlanningRepositoryFactory(log)
	return NewBatchUsecase(viper, log, repo, orgMessage, empMessage, batchDTO, mpPlanningRepo)
}
