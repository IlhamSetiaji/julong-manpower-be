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
	GetBatchedMPPlanningHeaders(approverType string, orgID string) (*[]response.MPPlanningHeaderResponse, error)
	FindByStatus(status entity.BatchHeaderApprovalStatus, approverType string, orgID string) (*response.BatchResponse, error)
	FindById(id string) (*response.BatchResponse, error)
	GetOrganizationsForBatchApproval(id string) (*[]response.OrganizationResponse, error)
	FindDocumentByID(id string) (*response.RealDocumentBatchResponse, error)
	FindByNeedApproval(approverType string, orgID string) (*response.RealDocumentBatchResponse, error)
	FindByCurrentDocumentDateAndStatus(status entity.BatchHeaderApprovalStatus) (*response.BatchResponse, error)
	UpdateStatusBatchHeader(req *request.UpdateStatusBatchHeaderRequest) (*response.BatchResponse, error)
	GetCompletedBatchHeader(page, pageSize int, search string, sort map[string]interface{}, employeeID uuid.UUID) (*[]response.CompletedBatchResponse, int64, error)
	GetBatchHeadersByStatus(status entity.BatchHeaderApprovalStatus, approverType string, orgID string) (*[]response.CompletedBatchResponse, error)
	GetBatchHeadersByStatusPaginated(status entity.BatchHeaderApprovalStatus, approverType string, orgID string, page, pageSize int, search string, sort map[string]interface{}, employeeID uuid.UUID) (*[]response.CompletedBatchResponse, int64, error)
	TriggerCreate(approverType string, orgID string) (bool, error)
	MPPlanningDetailsByBatchHeader(batchHeaderID uuid.UUID, page, pageSize int, search string, sort map[string]interface{}) (*[]response.MPPlanningHeaderResponse, int64, error)
}

type BatchUsecase struct {
	Viper            *viper.Viper
	Log              *logrus.Logger
	Repo             repository.IBatchRepository
	OrgMessage       messaging.IOrganizationMessage
	EmpMessage       messaging.IEmployeeMessage
	batchDTO         dto.IBatchDTO
	mpPlanningRepo   repository.IMPPlanningRepository
	JobPlafonMessage messaging.IJobPlafonMessage
	MPPlanningDTO    dto.IMPPlanningDTO
}

func NewBatchUsecase(
	viper *viper.Viper,
	log *logrus.Logger,
	repo repository.IBatchRepository,
	orgMessage messaging.IOrganizationMessage,
	empMessage messaging.IEmployeeMessage,
	batchDTO dto.IBatchDTO,
	mpPlanningRepo repository.IMPPlanningRepository,
	jpMessage messaging.IJobPlafonMessage,
	mpPlanningDTO dto.IMPPlanningDTO,
) IBatchUsecase {
	return &BatchUsecase{
		Viper:            viper,
		Log:              log,
		Repo:             repo,
		OrgMessage:       orgMessage,
		EmpMessage:       empMessage,
		batchDTO:         batchDTO,
		mpPlanningRepo:   mpPlanningRepo,
		JobPlafonMessage: jpMessage,
		MPPlanningDTO:    mpPlanningDTO,
	}
}

func (uc *BatchUsecase) GetCompletedBatchHeader(page, pageSize int, search string, sort map[string]interface{}, employeeID uuid.UUID) (*[]response.CompletedBatchResponse, int64, error) {
	batchHeaders, total, err := uc.Repo.GetBatchHeadersByStatusPaginated(entity.BatchHeaderApprovalStatusCompleted, entity.BatchHeaderApproverTypeCEO, "", page, pageSize, search, sort, employeeID)
	if err != nil {
		uc.Log.Errorf("[BatchUsecase.GetCompletedBatchHeader] " + err.Error())
		return nil, 0, err
	}

	if len(batchHeaders) == 0 {
		return nil, 0, nil
	}

	// var mpPlanningHeaderID uuid.UUID
	// for _, bh := range batchHeaders {
	// 	for _, bl := range bh.BatchLines {
	// 		if &bl.MPPlanningHeaderID != nil && bl.MPPlanningHeaderID != uuid.Nil {
	// 			mpPlanningHeaderID = bl.MPPlanningHeaderID
	// 			break
	// 		}
	// 	}
	// }

	// uc.Log.Infof("mpPlanningHeaderID: %s", mpPlanningHeaderID.String())

	// // get one mp planning header
	// mpPlanningHeader, err := uc.mpPlanningRepo.FindHeaderById(mpPlanningHeaderID)
	// if err != nil {
	// 	uc.Log.Errorf("[BatchUsecase.GetCompletedBatchHeader] " + err.Error())
	// 	return nil, 0, err
	// }

	// embed batch headers to completed batch responses
	completedBatchResponses := make([]response.CompletedBatchResponse, len(batchHeaders))

	var mpPlanningHeaderID uuid.UUID
	for i, bh := range batchHeaders {
		mpPlanningHeaderID = uuid.Nil
		for _, bl := range bh.BatchLines {
			if &bl.MPPlanningHeaderID != nil && bl.MPPlanningHeaderID != uuid.Nil {
				mpPlanningHeaderID = bl.MPPlanningHeaderID
				break
			}
		}
		uc.Log.Infof("mpPlanningHeaderID: %s", mpPlanningHeaderID.String())

		// get one mp planning header
		mpPlanningHeader, err := uc.mpPlanningRepo.FindHeaderById(mpPlanningHeaderID)
		if err != nil {
			uc.Log.Errorf("[BatchUsecase.GetCompletedBatchHeader] " + err.Error())
			return nil, 0, err
		}
		var mpPeriodResponse *response.MPPeriodResponse
		if mpPlanningHeader != nil {
			mpPeriodResponse = &response.MPPeriodResponse{
				ID:              mpPlanningHeader.MPPPeriod.ID,
				Title:           mpPlanningHeader.MPPPeriod.Title,
				StartDate:       mpPlanningHeader.MPPPeriod.StartDate.Format("2006-01-02"),
				EndDate:         mpPlanningHeader.MPPPeriod.EndDate.Format("2006-01-02"),
				BudgetStartDate: mpPlanningHeader.MPPPeriod.BudgetStartDate.Format("2006-01-02"),
				BudgetEndDate:   mpPlanningHeader.MPPPeriod.BudgetEndDate.Format("2006-01-02"),
				Status:          mpPlanningHeader.MPPPeriod.Status,
				CreatedAt:       mpPlanningHeader.MPPPeriod.CreatedAt,
				UpdatedAt:       mpPlanningHeader.MPPPeriod.UpdatedAt,
			}
		}
		completedBatchResponses[i] = response.CompletedBatchResponse{
			ID:             bh.ID,
			DocumentNumber: bh.DocumentNumber,
			DocumentDate:   bh.DocumentDate,
			Status:         bh.Status,
			CreatedAt:      bh.CreatedAt,
			UpdatedAt:      bh.UpdatedAt,
			MPPPeriod: func() response.MPPeriodResponse {
				if mpPeriodResponse != nil {
					return *mpPeriodResponse
				}
				return response.MPPeriodResponse{}
			}(),
		}
	}

	return &completedBatchResponses, total, nil
}

func (uc *BatchUsecase) GetBatchHeadersByStatus(status entity.BatchHeaderApprovalStatus, approverType string, orgID string) (*[]response.CompletedBatchResponse, error) {
	batchHeaders, err := uc.Repo.GetBatchHeadersByStatus(status, entity.BatchHeaderApproverType(approverType), orgID)
	if err != nil {
		uc.Log.Errorf("[BatchUsecase.GetBatchHeadersByStatus] " + err.Error())
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
			Status:         bh.Status,
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

func (uc *BatchUsecase) GetBatchHeadersByStatusPaginated(status entity.BatchHeaderApprovalStatus, approverType string, orgID string, page, pageSize int, search string, sort map[string]interface{}, employeeID uuid.UUID) (*[]response.CompletedBatchResponse, int64, error) {
	batchHeaders, total, err := uc.Repo.GetBatchHeadersByStatusPaginated(status, entity.BatchHeaderApproverType(approverType), orgID, page, pageSize, search, sort, employeeID)
	if err != nil {
		uc.Log.Errorf("[BatchUsecase.GetBatchHeadersByStatus] " + err.Error())
		return nil, 0, err
	}

	if len(batchHeaders) == 0 {
		return nil, 0, nil
	}

	completedBatchResponses := make([]response.CompletedBatchResponse, len(batchHeaders))

	// embed batch headers to completed batch responses
	for i, bh := range batchHeaders {
		if len(bh.BatchLines) == 0 {
			uc.Log.Warnf("[BatchUsecase.GetCompletedBatchHeader] No batch lines found for batch header ID: %s", bh.ID)
			continue
		}

		mpPlanningHeaderID := bh.BatchLines[0].MPPlanningHeaderID
		if mpPlanningHeaderID == uuid.Nil {
			uc.Log.Warnf("[BatchUsecase.GetCompletedBatchHeader] MPPlanningHeaderID is nil for batch header ID: %s", bh.ID)
			continue
		}

		mpPlanningHeader, err := uc.mpPlanningRepo.FindHeaderById(mpPlanningHeaderID)
		if err != nil {
			uc.Log.Errorf("[BatchUsecase.GetCompletedBatchHeader] " + err.Error())
			continue
		}

		if mpPlanningHeader == nil {
			uc.Log.Warnf("[BatchUsecase.GetCompletedBatchHeader] MPPlanningHeader not found for ID: %s", mpPlanningHeaderID)
			continue
		}

		mpPeriodResponse := &response.MPPeriodResponse{
			ID:              mpPlanningHeader.MPPPeriod.ID,
			Title:           mpPlanningHeader.MPPPeriod.Title,
			StartDate:       mpPlanningHeader.MPPPeriod.StartDate.Format("2006-01-02"),
			EndDate:         mpPlanningHeader.MPPPeriod.EndDate.Format("2006-01-02"),
			BudgetStartDate: mpPlanningHeader.MPPPeriod.BudgetStartDate.Format("2006-01-02"),
			BudgetEndDate:   mpPlanningHeader.MPPPeriod.BudgetEndDate.Format("2006-01-02"),
			Status:          mpPlanningHeader.MPPPeriod.Status,
			CreatedAt:       mpPlanningHeader.MPPPeriod.CreatedAt,
			UpdatedAt:       mpPlanningHeader.MPPPeriod.UpdatedAt,
		}

		completedBatchResponses[i] = response.CompletedBatchResponse{
			ID:             bh.ID,
			DocumentNumber: bh.DocumentNumber,
			DocumentDate:   bh.DocumentDate,
			Status:         bh.Status,
			CreatedAt:      bh.CreatedAt,
			UpdatedAt:      bh.UpdatedAt,
			MPPPeriod: func() response.MPPeriodResponse {
				if mpPeriodResponse != nil {
					return *mpPeriodResponse
				}
				return response.MPPeriodResponse{}
			}(),
		}
	}

	return &completedBatchResponses, total, nil
}

func (uc *BatchUsecase) GetBatchedMPPlanningHeaders(approverType string, orgID string) (*[]response.MPPlanningHeaderResponse, error) {
	mpPlanningHeaders := make([]response.MPPlanningHeaderResponse, 0)
	entMpPlanningHeaders := make([]entity.MPPlanningHeader, 0)
	if approverType == "" || approverType == string(entity.BatchHeaderApproverTypeCEO) {
		headers, err := uc.mpPlanningRepo.GetAllHeadersGroupedApproverByOrg("", entity.MPPlaningStatusApproved, "direktur", "")
		if err != nil {
			uc.Log.Errorf("[BatchUsecase.GetBatchedMPPlanningHeaders] " + err.Error())
			return nil, err
		}

		for _, header := range *headers {
			entMpPlanningHeaders = append(entMpPlanningHeaders, entity.MPPlanningHeader{
				ID:             header.ID,
				DocumentNumber: header.DocumentNumber,
				DocumentDate:   header.DocumentDate,
				Status:         header.Status,
			})
		}
	} else {
		headers, err := uc.mpPlanningRepo.GetAllHeadersGroupedApproverByOrg(orgID, entity.MPPlanningStatusInProgress, "", "")
		if err != nil {
			uc.Log.Errorf("[BatchUsecase.GetBatchedMPPlanningHeaders] " + err.Error())
			return nil, err
		}

		for _, header := range *headers {
			entMpPlanningHeaders = append(entMpPlanningHeaders, entity.MPPlanningHeader{
				ID:             header.ID,
				DocumentNumber: header.DocumentNumber,
				DocumentDate:   header.DocumentDate,
				Status:         header.Status,
			})
		}
	}

	for _, header := range entMpPlanningHeaders {
		mpPlanningHeaders = append(mpPlanningHeaders, response.MPPlanningHeaderResponse{
			ID:             header.ID,
			DocumentNumber: header.DocumentNumber,
			DocumentDate:   header.DocumentDate,
			Status:         header.Status,
		})
	}

	return &mpPlanningHeaders, nil
}

func (uc *BatchUsecase) TriggerCreate(approverType string, orgID string) (bool, error) {
	mpPlanningExists := &entity.MPPlanningHeader{}
	var err error
	if approverType == "" || approverType == string(entity.BatchHeaderApproverTypeCEO) {
		mpPlanningExists, err = uc.mpPlanningRepo.FindAllHeadersGroupedApproverByOrg("", entity.MPPlaningStatusApproved, "direktur", "")
		if err != nil {
			uc.Log.Errorf("[BatchUsecase.TriggerCreate] " + err.Error())
			return false, err
		}
	} else {
		mpPlanningExists, err = uc.mpPlanningRepo.FindAllHeadersGroupedApproverByOrg(orgID, entity.MPPlanningStatusInProgress, "", "")
		if err != nil {
			uc.Log.Errorf("[BatchUsecase.TriggerCreate] " + err.Error())
			return false, err
		}
	}

	if mpPlanningExists == nil {
		return false, nil
	}

	return true, nil
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
				uc.Log.Errorf("[BatchUsecase.CreateBatchHeaderAndLines] MP Planning Header not in APPROVED status")
				continue
			}
		} else {
			if mpHeaderByStatus.Status != entity.MPPlanningStatusInProgress {
				uc.Log.Errorf("[BatchUsecase.CreateBatchHeaderAndLines] MP Planning Header not in IN PROGRESS status")
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
			Level:              string(entity.MPPlanningApprovalHistoryLevelHRDUnit),
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
		err = uc.Repo.UpdateStatusBatchHeaderForDirector(batchHeader, req.Status, req.ApprovedBy, req.ApproverName)
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

	// // get job level for mp planning lines
	// for i, bl := range resp.BatchLines {
	// 	for j, mpl := range bl.MPPlanningHeader.MPPlanningLines {
	// 		if mpl.JobLevelName == "" {
	// 			jobLevel, err := uc.JobPlafonMessage.SendFindJobLevelByIDMessage(request.SendFindJobLevelByIDMessageRequest{
	// 				ID: mpl.JobLevelID.String(),
	// 			})
	// 			if err != nil {
	// 				uc.Log.Errorf("[BatchUsecase.FindDocumentByID] " + err.Error())
	// 			}
	// 			mpl.JobLevelName = jobLevel.Name
	// 			mpl.JobLevel = int(jobLevel.Level)

	// 			resp.BatchLines[i].MPPlanningHeader.MPPlanningLines[j] = mpl
	// 		}
	// 	}
	// }

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

func (uc *BatchUsecase) MPPlanningDetailsByBatchHeader(batchHeaderID uuid.UUID, page, pageSize int, search string, sort map[string]interface{}) (*[]response.MPPlanningHeaderResponse, int64, error) {
	batch, err := uc.Repo.FindById(batchHeaderID.String())
	if err != nil {
		uc.Log.Errorf("[BatchUsecase.MPPlanningDetailsByBatchHeader] " + err.Error())
		return nil, 0, err
	}

	if batch == nil {
		return nil, 0, errors.New("Batch not found")
	}

	var filteredBatchLines []entity.BatchLine
	for _, bl := range batch.BatchLines {
		mpPlanningHeader, err := uc.mpPlanningRepo.FindHeaderById(bl.MPPlanningHeaderID)
		if err != nil {
			uc.Log.Errorf("[BatchUsecase.MPPlanningDetailsByBatchHeader] " + err.Error())
			return nil, 0, err
		}

		if search == "" || (search != "" && mpPlanningHeader.DocumentNumber == search) {
			filteredBatchLines = append(filteredBatchLines, bl)
		}
	}

	total := int64(len(filteredBatchLines))

	// Apply pagination
	start := (page - 1) * pageSize
	end := start + pageSize
	if start > len(filteredBatchLines) {
		start = len(filteredBatchLines)
	}
	if end > len(filteredBatchLines) {
		end = len(filteredBatchLines)
	}
	paginatedBatchLines := filteredBatchLines[start:end]

	mpPlanningHeaders := make([]response.MPPlanningHeaderResponse, len(paginatedBatchLines))
	for i, bl := range paginatedBatchLines {
		mpPlanningHeader, err := uc.mpPlanningRepo.FindHeaderById(bl.MPPlanningHeaderID)
		if err != nil {
			uc.Log.Errorf("[BatchUsecase.MPPlanningDetailsByBatchHeader] " + err.Error())
			return nil, 0, err
		}

		mpPlanningHeaders[i] = *uc.MPPlanningDTO.ConvertMPPlanningHeaderEntityToResponse(mpPlanningHeader)
	}

	return &mpPlanningHeaders, total, nil
}

func BatchUsecaseFactory(viper *viper.Viper, log *logrus.Logger) IBatchUsecase {
	repo := repository.BatchRepositoryFactory(log)
	orgMessage := messaging.OrganizationMessageFactory(log)
	empMessage := messaging.EmployeeMessageFactory(log)
	batchDTO := dto.BatchDTOFactory(log)
	mpPlanningRepo := repository.MPPlanningRepositoryFactory(log)
	jpMessage := messaging.JobPlafonMessageFactory(log)
	mpPlanningDTO := dto.MPPlanningDTOFactory(log)
	return NewBatchUsecase(viper, log, repo, orgMessage, empMessage, batchDTO, mpPlanningRepo, jpMessage, mpPlanningDTO)
}
