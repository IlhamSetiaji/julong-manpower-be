package usecase

import (
	"errors"
	"time"

	"github.com/IlhamSetiaji/julong-manpower-be/internal/entity"
	"github.com/IlhamSetiaji/julong-manpower-be/internal/http/messaging"
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
	FindHeaderByMPPPeriodId(request *request.FindHeaderByMPPPeriodIdMPPlanningRequest) (*response.FindHeaderByMPPPeriodIdMPPlanningResponse, error)
	FindAllLinesByHeaderIdPaginated(request *request.FindAllLinesByHeaderIdPaginatedMPPlanningLineRequest) (*response.FindAllLinesByHeaderIdPaginatedMPPlanningLineResponse, error)
	FindLineById(request *request.FindLineByIdMPPlanningLineRequest) (*response.FindByIdMPPlanningLineResponse, error)
	CreateLine(request *request.CreateLineMPPlanningLineRequest) (*response.CreateMPPlanningLineResponse, error)
	UpdateLine(request *request.UpdateLineMPPlanningLineRequest) (*response.UpdateMPPlanningLineResponse, error)
	DeleteLine(request *request.DeleteLineMPPlanningLineRequest) error
	CreateOrUpdateBatchLineMPPlanningLines(request *request.CreateOrUpdateBatchLineMPPlanningLinesRequest) error
}

type MPPlanningUseCase struct {
	Log                  *logrus.Logger
	MPPlanningRepository repository.IMPPlanningRepository
	OrganizationMessage  messaging.IOrganizationMessage
	JobPlafonMessage     messaging.IJobPlafonMessage
	UserMessage          messaging.IUserMessage
}

func NewMPPlanningUseCase(log *logrus.Logger, repo repository.IMPPlanningRepository, message messaging.IOrganizationMessage, jpm messaging.IJobPlafonMessage, um messaging.IUserMessage) IMPPlanningUseCase {
	return &MPPlanningUseCase{
		Log:                  log,
		MPPlanningRepository: repo,
		OrganizationMessage:  message,
		JobPlafonMessage:     jpm,
		UserMessage:          um,
	}
}

func (uc *MPPlanningUseCase) FindAllHeadersPaginated(req *request.FindAllHeadersPaginatedMPPlanningRequest) (*response.FindAllHeadersPaginatedMPPlanningResponse, error) {
	mpPlanningHeaders, total, err := uc.MPPlanningRepository.FindAllHeadersPaginated(req.Page, req.PageSize, req.Search)
	if err != nil {
		uc.Log.Errorf("[MPPlanningUseCase.FindAllHeadersPaginated] " + err.Error())
		return nil, err
	}

	for i, header := range *mpPlanningHeaders {
		// Fetch organization names using RabbitMQ
		messageResponse, err := uc.OrganizationMessage.SendFindOrganizationByIDMessage(request.SendFindOrganizationByIDMessageRequest{
			ID: header.OrganizationID.String(),
		})
		if err != nil {
			uc.Log.Errorf("[MPPlanningUseCase.FindAllHeadersPaginated Message] " + err.Error())
			return nil, err
		}
		header.OrganizationName = messageResponse.Name

		// Fetch emp organization names using RabbitMQ
		message2Response, err := uc.OrganizationMessage.SendFindOrganizationByIDMessage(request.SendFindOrganizationByIDMessageRequest{
			ID: header.EmpOrganizationID.String(),
		})
		if err != nil {
			uc.Log.Errorf("[MPPlanningUseCase.FindAllHeadersPaginated Message] " + err.Error())
			return nil, err
		}
		header.EmpOrganizationName = message2Response.Name

		// Fetch job names using RabbitMQ
		messageJobResposne, err := uc.JobPlafonMessage.SendFindJobByIDMessage(request.SendFindJobByIDMessageRequest{
			ID: header.JobID.String(),
		})
		if err != nil {
			uc.Log.Errorf("[MPPlanningUseCase.FindAllHeadersPaginated Message] " + err.Error())
			return nil, err
		}
		header.JobName = messageJobResposne.Name

		// Fetch requestor names using RabbitMQ
		messageUserResponse, err := uc.UserMessage.SendFindUserByIDMessage(request.SendFindUserByIDMessageRequest{
			ID: header.RequestorID.String(),
		})
		if err != nil {
			uc.Log.Errorf("[MPPlanningUseCase.FindAllHeadersPaginated Message] " + err.Error())
			return nil, err
		}
		header.RequestorName = messageUserResponse.Name

		for i, line := range *&header.MPPlanningLines {
			// Fetch organization location names using RabbitMQ
			messageResponse, err := uc.OrganizationMessage.SendFindOrganizationLocationByIDMessage(request.SendFindOrganizationLocationByIDMessageRequest{
				ID: line.OrganizationLocationID.String(),
			})
			if err != nil {
				uc.Log.Errorf("[MPPlanningUseCase.FindAllLinesByHeaderIdPaginated Message] " + err.Error())
				return nil, err
			}
			line.OrganizationLocationName = messageResponse.Name

			// Fetch job level names using RabbitMQ
			message2Response, err := uc.JobPlafonMessage.SendFindJobLevelByIDMessage(request.SendFindJobLevelByIDMessageRequest{
				ID: line.JobLevelID.String(),
			})
			if err != nil {
				uc.Log.Errorf("[MPPlanningUseCase.FindAllLinesByHeaderIdPaginated Message] " + err.Error())
				return nil, err
			}
			line.JobLevelName = message2Response.Name

			// Fetch job names using RabbitMQ
			messageJobResposne, err := uc.JobPlafonMessage.SendFindJobByIDMessage(request.SendFindJobByIDMessageRequest{
				ID: line.JobID.String(),
			})
			if err != nil {
				uc.Log.Errorf("[MPPlanningUseCase.FindAllLinesByHeaderIdPaginated Message] " + err.Error())
				return nil, err
			}
			line.JobName = messageJobResposne.Name

			header.MPPlanningLines[i] = line
		}

		(*mpPlanningHeaders)[i] = header
	}

	return &response.FindAllHeadersPaginatedMPPlanningResponse{
		MPPlanningHeaders: func() []*response.MPPlanningHeaderResponse {
			var headers []*response.MPPlanningHeaderResponse
			for _, header := range *mpPlanningHeaders {
				headers = append(headers, &response.MPPlanningHeaderResponse{
					ID:                  header.ID,
					MPPPeriodID:         header.MPPPeriodID,
					OrganizationID:      header.OrganizationID,
					EmpOrganizationID:   header.EmpOrganizationID,
					JobID:               header.JobID,
					DocumentNumber:      header.DocumentNumber,
					DocumentDate:        header.DocumentDate,
					Notes:               header.Notes,
					TotalRecruit:        header.TotalRecruit,
					TotalPromote:        header.TotalPromote,
					Status:              header.Status,
					RecommendedBy:       header.RecommendedBy,
					ApprovedBy:          header.ApprovedBy,
					RequestorID:         header.RequestorID,
					NotesAttach:         header.NotesAttach,
					OrganizationName:    header.OrganizationName,
					EmpOrganizationName: header.EmpOrganizationName,
					JobName:             header.JobName,
					RequestorName:       header.RequestorName,
					CreatedAt:           header.CreatedAt,
					UpdatedAt:           header.UpdatedAt,
					MPPPeriod: &response.MPPeriodResponse{
						ID:        header.MPPPeriod.ID,
						Title:     header.MPPPeriod.Title,
						StartDate: header.MPPPeriod.StartDate.Format("2006-01-02"),
						EndDate:   header.MPPPeriod.EndDate.Format("2006-01-02"),
						CreatedAt: header.MPPPeriod.CreatedAt,
						UpdatedAt: header.MPPPeriod.UpdatedAt,
					},
					MPPlanningLines: func() []*response.MPPlanningLineResponse {
						var lines []*response.MPPlanningLineResponse
						for _, line := range header.MPPlanningLines {
							lines = append(lines, &response.MPPlanningLineResponse{
								ID:                       line.ID,
								MPPlanningHeaderID:       line.MPPlanningHeaderID,
								OrganizationLocationID:   *line.OrganizationLocationID,
								JobLevelID:               *line.JobLevelID,
								JobID:                    *line.JobID,
								Existing:                 line.Existing,
								Recruit:                  line.Recruit,
								SuggestedRecruit:         line.SuggestedRecruit,
								Promotion:                line.Promotion,
								Total:                    line.Total,
								RemainingBalancePH:       line.RemainingBalancePH,
								RemainingBalanceMT:       line.RemainingBalanceMT,
								RecruitPH:                line.RecruitPH,
								RecruitMT:                line.RecruitMT,
								OrganizationLocationName: line.OrganizationLocationName,
								JobLevelName:             line.JobLevelName,
								JobName:                  line.JobName,
							})
						}
						return lines
					}(),
				})
			}
			return headers
		}(),
		Total: total,
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

	// Fetch organization names using RabbitMQ
	messageResponse, err := uc.OrganizationMessage.SendFindOrganizationByIDMessage(request.SendFindOrganizationByIDMessageRequest{
		ID: mpPlanningHeader.OrganizationID.String(),
	})
	if err != nil {
		uc.Log.Errorf("[MPPlanningUseCase.FindById Message] " + err.Error())
		return nil, err
	}
	mpPlanningHeader.OrganizationName = messageResponse.Name

	// Fetch emp organization names using RabbitMQ
	message2Response, err := uc.OrganizationMessage.SendFindOrganizationByIDMessage(request.SendFindOrganizationByIDMessageRequest{
		ID: mpPlanningHeader.EmpOrganizationID.String(),
	})
	if err != nil {
		uc.Log.Errorf("[MPPlanningUseCase.FindById Message] " + err.Error())
		return nil, err
	}
	mpPlanningHeader.EmpOrganizationName = message2Response.Name

	// Fetch job names using RabbitMQ
	messageJobResposne, err := uc.JobPlafonMessage.SendFindJobByIDMessage(request.SendFindJobByIDMessageRequest{
		ID: mpPlanningHeader.JobID.String(),
	})
	if err != nil {
		uc.Log.Errorf("[MPPlanningUseCase.FindById Message] " + err.Error())
		return nil, err
	}
	mpPlanningHeader.JobName = messageJobResposne.Name

	// Fetch requestor names using RabbitMQ
	messageUserResponse, err := uc.UserMessage.SendFindUserByIDMessage(request.SendFindUserByIDMessageRequest{
		ID: mpPlanningHeader.RequestorID.String(),
	})
	if err != nil {
		uc.Log.Errorf("[MPPlanningUseCase.FindById Message] " + err.Error())
		return nil, err
	}
	mpPlanningHeader.RequestorName = messageUserResponse.Name

	for i, line := range *&mpPlanningHeader.MPPlanningLines {
		// Fetch organization location names using RabbitMQ
		messageResponse, err := uc.OrganizationMessage.SendFindOrganizationLocationByIDMessage(request.SendFindOrganizationLocationByIDMessageRequest{
			ID: line.OrganizationLocationID.String(),
		})
		if err != nil {
			uc.Log.Errorf("[MPPlanningUseCase.FindAllLinesByHeaderIdPaginated Message] " + err.Error())
			return nil, err
		}
		line.OrganizationLocationName = messageResponse.Name

		// Fetch job level names using RabbitMQ
		message2Response, err := uc.JobPlafonMessage.SendFindJobLevelByIDMessage(request.SendFindJobLevelByIDMessageRequest{
			ID: line.JobLevelID.String(),
		})
		if err != nil {
			uc.Log.Errorf("[MPPlanningUseCase.FindAllLinesByHeaderIdPaginated Message] " + err.Error())
			return nil, err
		}
		line.JobLevelName = message2Response.Name

		// Fetch job names using RabbitMQ
		messageJobResposne, err := uc.JobPlafonMessage.SendFindJobByIDMessage(request.SendFindJobByIDMessageRequest{
			ID: line.JobID.String(),
		})
		if err != nil {
			uc.Log.Errorf("[MPPlanningUseCase.FindAllLinesByHeaderIdPaginated Message] " + err.Error())
			return nil, err
		}
		line.JobName = messageJobResposne.Name

		mpPlanningHeader.MPPlanningLines[i] = line
	}

	return &response.FindByIdMPPlanningResponse{
		ID:                  mpPlanningHeader.ID,
		MPPPeriodID:         mpPlanningHeader.MPPPeriodID,
		OrganizationID:      mpPlanningHeader.OrganizationID,
		EmpOrganizationID:   mpPlanningHeader.EmpOrganizationID,
		JobID:               mpPlanningHeader.JobID,
		DocumentNumber:      mpPlanningHeader.DocumentNumber,
		DocumentDate:        mpPlanningHeader.DocumentDate,
		Notes:               mpPlanningHeader.Notes,
		TotalRecruit:        mpPlanningHeader.TotalRecruit,
		TotalPromote:        mpPlanningHeader.TotalPromote,
		Status:              mpPlanningHeader.Status,
		RecommendedBy:       mpPlanningHeader.RecommendedBy,
		ApprovedBy:          mpPlanningHeader.ApprovedBy,
		RequestorID:         mpPlanningHeader.RequestorID,
		NotesAttach:         mpPlanningHeader.NotesAttach,
		OrganizationName:    mpPlanningHeader.OrganizationName,
		EmpOrganizationName: mpPlanningHeader.EmpOrganizationName,
		JobName:             mpPlanningHeader.JobName,
		RequestorName:       mpPlanningHeader.RequestorName,
		CreatedAt:           mpPlanningHeader.CreatedAt,
		UpdatedAt:           mpPlanningHeader.UpdatedAt,
		MPPPeriod: &response.MPPeriodResponse{
			ID:        mpPlanningHeader.MPPPeriod.ID,
			Title:     mpPlanningHeader.MPPPeriod.Title,
			StartDate: mpPlanningHeader.MPPPeriod.StartDate.Format("2006-01-02"),
			EndDate:   mpPlanningHeader.MPPPeriod.EndDate.Format("2006-01-02"),
			CreatedAt: mpPlanningHeader.MPPPeriod.CreatedAt,
			UpdatedAt: mpPlanningHeader.MPPPeriod.UpdatedAt,
		},
		MPPlanningLines: func() []*response.MPPlanningLineResponse {
			var lines []*response.MPPlanningLineResponse
			for _, line := range mpPlanningHeader.MPPlanningLines {
				lines = append(lines, &response.MPPlanningLineResponse{
					ID:                       line.ID,
					MPPlanningHeaderID:       line.MPPlanningHeaderID,
					OrganizationLocationID:   *line.OrganizationLocationID,
					JobLevelID:               *line.JobLevelID,
					JobID:                    *line.JobID,
					Existing:                 line.Existing,
					Recruit:                  line.Recruit,
					SuggestedRecruit:         line.SuggestedRecruit,
					Promotion:                line.Promotion,
					Total:                    line.Total,
					RemainingBalancePH:       line.RemainingBalancePH,
					RemainingBalanceMT:       line.RemainingBalanceMT,
					RecruitPH:                line.RecruitPH,
					RecruitMT:                line.RecruitMT,
					OrganizationLocationName: line.OrganizationLocationName,
					JobLevelName:             line.JobLevelName,
					JobName:                  line.JobName,
				})
			}
			return lines
		}(),
	}, nil
}

func (uc *MPPlanningUseCase) FindHeaderByMPPPeriodId(req *request.FindHeaderByMPPPeriodIdMPPlanningRequest) (*response.FindHeaderByMPPPeriodIdMPPlanningResponse, error) {
	mpPlanningHeader, err := uc.MPPlanningRepository.FindHeaderByMPPPeriodId(uuid.MustParse(req.MPPPeriodID))
	if err != nil {
		uc.Log.Errorf("[MPPlanningUseCase.FindHeaderByMPPPeriodId] " + err.Error())
		return nil, err
	}

	if mpPlanningHeader == nil {
		uc.Log.Errorf("[MPPlanningUseCase.FindHeaderByMPPPeriodId] MP Planning Header not found")
		return nil, errors.New("MP Planning Header not found")
	}

	// Fetch organization names using RabbitMQ
	messageResponse, err := uc.OrganizationMessage.SendFindOrganizationByIDMessage(request.SendFindOrganizationByIDMessageRequest{
		ID: mpPlanningHeader.OrganizationID.String(),
	})
	if err != nil {
		uc.Log.Errorf("[MPPlanningUseCase.FindHeaderByMPPPeriodId Message] " + err.Error())
		return nil, err
	}
	mpPlanningHeader.OrganizationName = messageResponse.Name

	// Fetch emp organization names using RabbitMQ
	message2Response, err := uc.OrganizationMessage.SendFindOrganizationByIDMessage(request.SendFindOrganizationByIDMessageRequest{
		ID: mpPlanningHeader.EmpOrganizationID.String(),
	})
	if err != nil {
		uc.Log.Errorf("[MPPlanningUseCase.FindHeaderByMPPPeriodId Message] " + err.Error())
		return nil, err
	}
	mpPlanningHeader.EmpOrganizationName = message2Response.Name

	// Fetch job names using RabbitMQ
	messageJobResposne, err := uc.JobPlafonMessage.SendFindJobByIDMessage(request.SendFindJobByIDMessageRequest{
		ID: mpPlanningHeader.JobID.String(),
	})
	if err != nil {
		uc.Log.Errorf("[MPPlanningUseCase.FindHeaderByMPPPeriodId Message] " + err.Error())
		return nil, err
	}
	mpPlanningHeader.JobName = messageJobResposne.Name

	// Fetch requestor names using RabbitMQ
	messageUserResponse, err := uc.UserMessage.SendFindUserByIDMessage(request.SendFindUserByIDMessageRequest{
		ID: mpPlanningHeader.RequestorID.String(),
	})
	if err != nil {
		uc.Log.Errorf("[MPPlanningUseCase.FindHeaderByMPPPeriodId Message] " + err.Error())
		return nil, err
	}

	return &response.FindHeaderByMPPPeriodIdMPPlanningResponse{
		ID:                  mpPlanningHeader.ID,
		MPPPeriodID:         mpPlanningHeader.MPPPeriodID,
		OrganizationID:      mpPlanningHeader.OrganizationID,
		EmpOrganizationID:   mpPlanningHeader.EmpOrganizationID,
		JobID:               mpPlanningHeader.JobID,
		DocumentNumber:      mpPlanningHeader.DocumentNumber,
		DocumentDate:        mpPlanningHeader.DocumentDate,
		Notes:               mpPlanningHeader.Notes,
		TotalRecruit:        mpPlanningHeader.TotalRecruit,
		TotalPromote:        mpPlanningHeader.TotalPromote,
		Status:              mpPlanningHeader.Status,
		RecommendedBy:       mpPlanningHeader.RecommendedBy,
		ApprovedBy:          mpPlanningHeader.ApprovedBy,
		RequestorID:         mpPlanningHeader.RequestorID,
		NotesAttach:         mpPlanningHeader.NotesAttach,
		OrganizationName:    mpPlanningHeader.OrganizationName,
		EmpOrganizationName: mpPlanningHeader.EmpOrganizationName,
		JobName:             mpPlanningHeader.JobName,
		RequestorName:       messageUserResponse.Name,
		CreatedAt:           mpPlanningHeader.CreatedAt,
		UpdatedAt:           mpPlanningHeader.UpdatedAt,
	}, nil
}

func (uc *MPPlanningUseCase) Create(req *request.CreateHeaderMPPlanningRequest) (*response.CreateMPPlanningResponse, error) {
	// Check if organization exist
	orgExist, err := uc.OrganizationMessage.SendFindOrganizationByIDMessage(request.SendFindOrganizationByIDMessageRequest{
		ID: req.OrganizationID.String(),
	})

	if err != nil {
		uc.Log.Errorf("[MPPlanningUseCase.Create] " + err.Error())
		return nil, err
	}

	if orgExist == nil {
		uc.Log.Errorf("[MPPlanningUseCase.Create] Organization not found")
		return nil, errors.New("Organization not found")
	}

	// Check if emp organization exist
	empOrgExist, err := uc.OrganizationMessage.SendFindOrganizationByIDMessage(request.SendFindOrganizationByIDMessageRequest{
		ID: req.EmpOrganizationID.String(),
	})

	if err != nil {
		uc.Log.Errorf("[MPPlanningUseCase.Create] " + err.Error())
		return nil, err
	}

	if empOrgExist == nil {
		uc.Log.Errorf("[MPPlanningUseCase.Create] Emp Organization not found")
		return nil, errors.New("Emp Organization not found")
	}

	// Check if job exist
	jobExist, err := uc.JobPlafonMessage.SendFindJobByIDMessage(request.SendFindJobByIDMessageRequest{
		ID: req.JobID.String(),
	})

	if err != nil {
		uc.Log.Errorf("[MPPlanningUseCase.Create] " + err.Error())
		return nil, err
	}

	if jobExist == nil {
		uc.Log.Errorf("[MPPlanningUseCase.Create] Job not found")
		return nil, errors.New("Job not found")
	}

	// Check if requestor exist
	requestorExist, err := uc.UserMessage.SendFindUserByIDMessage(request.SendFindUserByIDMessageRequest{
		ID: req.RequestorID.String(),
	})

	if err != nil {
		uc.Log.Errorf("[MPPlanningUseCase.Create] " + err.Error())
		return nil, err
	}

	if requestorExist == nil {
		uc.Log.Errorf("[MPPlanningUseCase.Create] Requestor not found")
		return nil, errors.New("Requestor not found")
	}

	documentDate, err := time.Parse("2006-01-02", req.DocumentDate)
	if err != nil {
		uc.Log.Errorf("[MPPlanningUseCase.Create] " + err.Error())
		return nil, err
	}

	mpPlanningHeader, err := uc.MPPlanningRepository.CreateHeader(&entity.MPPlanningHeader{
		MPPPeriodID:       req.MPPPeriodID,
		OrganizationID:    &req.OrganizationID,
		EmpOrganizationID: &req.EmpOrganizationID,
		JobID:             &req.JobID,
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
		JobID:             mpPlanningHeader.JobID.String(),
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
	exist, err := uc.MPPlanningRepository.FindHeaderById(req.ID)
	if err != nil {
		uc.Log.Errorf("[MPPlanningUseCase.Update] " + err.Error())
		return nil, err
	}

	if exist == nil {
		uc.Log.Errorf("[MPPlanningUseCase.Update] MP Planning Header not found")
		return nil, errors.New("MP Planning Header not found")
	}

	documentDate, err := time.Parse("2006-01-02", req.DocumentDate)
	if err != nil {
		uc.Log.Errorf("[MPPlanningUseCase.Update] " + err.Error())
		return nil, err
	}

	mpPlanningHeader, err := uc.MPPlanningRepository.UpdateHeader(&entity.MPPlanningHeader{
		ID:                req.ID,
		MPPPeriodID:       req.MPPPeriodID,
		OrganizationID:    &req.OrganizationID,
		EmpOrganizationID: &req.EmpOrganizationID,
		JobID:             &req.JobID,
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
		JobID:             mpPlanningHeader.JobID.String(),
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
	exist, err := uc.MPPlanningRepository.FindHeaderById(uuid.MustParse(req.ID))
	if err != nil {
		uc.Log.Errorf("[MPPlanningUseCase.Update] " + err.Error())
		return err
	}

	if exist == nil {
		uc.Log.Errorf("[MPPlanningUseCase.Update] MP Planning Header not found")
		return errors.New("MP Planning Header not found")
	}

	err = uc.MPPlanningRepository.DeleteHeader(uuid.MustParse(req.ID))
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

	for i, line := range *mpPlanningLines {
		// Fetch organization location names using RabbitMQ
		messageResponse, err := uc.OrganizationMessage.SendFindOrganizationLocationByIDMessage(request.SendFindOrganizationLocationByIDMessageRequest{
			ID: line.OrganizationLocationID.String(),
		})
		if err != nil {
			uc.Log.Errorf("[MPPlanningUseCase.FindAllLinesByHeaderIdPaginated Message] " + err.Error())
			return nil, err
		}
		line.OrganizationLocationName = messageResponse.Name

		// Fetch job level names using RabbitMQ
		message2Response, err := uc.JobPlafonMessage.SendFindJobLevelByIDMessage(request.SendFindJobLevelByIDMessageRequest{
			ID: line.JobLevelID.String(),
		})
		if err != nil {
			uc.Log.Errorf("[MPPlanningUseCase.FindAllLinesByHeaderIdPaginated Message] " + err.Error())
			return nil, err
		}
		line.JobLevelName = message2Response.Name

		// Fetch job names using RabbitMQ
		messageJobResposne, err := uc.JobPlafonMessage.SendFindJobByIDMessage(request.SendFindJobByIDMessageRequest{
			ID: line.JobID.String(),
		})
		if err != nil {
			uc.Log.Errorf("[MPPlanningUseCase.FindAllLinesByHeaderIdPaginated Message] " + err.Error())
			return nil, err
		}
		line.JobName = messageJobResposne.Name

		(*mpPlanningLines)[i] = line
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

	// Fetch organization location names using RabbitMQ
	messageResponse, err := uc.OrganizationMessage.SendFindOrganizationLocationByIDMessage(request.SendFindOrganizationLocationByIDMessageRequest{
		ID: mpPlanningLine.OrganizationLocationID.String(),
	})
	if err != nil {
		uc.Log.Errorf("[MPPlanningUseCase.FindLineById Message] " + err.Error())
		return nil, err
	}
	mpPlanningLine.OrganizationLocationName = messageResponse.Name

	// Fetch job level names using RabbitMQ
	message2Response, err := uc.JobPlafonMessage.SendFindJobLevelByIDMessage(request.SendFindJobLevelByIDMessageRequest{
		ID: mpPlanningLine.JobLevelID.String(),
	})
	if err != nil {
		uc.Log.Errorf("[MPPlanningUseCase.FindLineById Message] " + err.Error())
		return nil, err
	}
	mpPlanningLine.JobLevelName = message2Response.Name

	// Fetch job names using RabbitMQ
	messageJobResposne, err := uc.JobPlafonMessage.SendFindJobByIDMessage(request.SendFindJobByIDMessageRequest{
		ID: mpPlanningLine.JobID.String(),
	})
	if err != nil {
		uc.Log.Errorf("[MPPlanningUseCase.FindLineById Message] " + err.Error())
		return nil, err
	}
	mpPlanningLine.JobName = messageJobResposne.Name

	return &response.FindByIdMPPlanningLineResponse{
		MPPlanningLine: mpPlanningLine,
	}, nil
}

func (uc *MPPlanningUseCase) CreateLine(req *request.CreateLineMPPlanningLineRequest) (*response.CreateMPPlanningLineResponse, error) {
	// Check if organization location exist
	orgLocExist, err := uc.OrganizationMessage.SendFindOrganizationLocationByIDMessage(request.SendFindOrganizationLocationByIDMessageRequest{
		ID: req.OrganizationLocationID.String(),
	})

	if err != nil {
		uc.Log.Errorf("[MPPlanningUseCase.CreateLine] " + err.Error())
		return nil, err
	}

	if orgLocExist == nil {
		uc.Log.Errorf("[MPPlanningUseCase.CreateLine] Organization Location not found")
		return nil, errors.New("Organization Location not found")
	}

	// Check if job level exist
	jobLevelExist, err := uc.JobPlafonMessage.SendFindJobLevelByIDMessage(request.SendFindJobLevelByIDMessageRequest{
		ID: req.JobLevelID.String(),
	})

	if err != nil {
		uc.Log.Errorf("[MPPlanningUseCase.CreateLine] " + err.Error())
		return nil, err
	}

	if jobLevelExist == nil {
		uc.Log.Errorf("[MPPlanningUseCase.CreateLine] Job Level not found")
		return nil, errors.New("Job Level not found")
	}

	// Check if job exist
	jobExist, err := uc.JobPlafonMessage.SendFindJobByIDMessage(request.SendFindJobByIDMessageRequest{
		ID: req.JobID.String(),
	})

	if err != nil {
		uc.Log.Errorf("[MPPlanningUseCase.CreateLine] " + err.Error())
		return nil, err
	}

	if jobExist == nil {
		uc.Log.Errorf("[MPPlanningUseCase.CreateLine] Job not found")
		return nil, errors.New("Job not found")
	}

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
		RemainingBalancePH:     req.RemainingBalancePH,
		RemainingBalanceMT:     req.RemainingBalanceMT,
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
		RemainingBalancePH:     mpPlanningLine.RemainingBalancePH,
		RemainingBalanceMT:     mpPlanningLine.RemainingBalanceMT,
		RecruitPH:              mpPlanningLine.RecruitPH,
		RecruitMT:              mpPlanningLine.RecruitMT,
		CreatedAt:              mpPlanningLine.CreatedAt,
		UpdatedAt:              mpPlanningLine.UpdatedAt,
	}, nil
}

func (uc *MPPlanningUseCase) UpdateLine(req *request.UpdateLineMPPlanningLineRequest) (*response.UpdateMPPlanningLineResponse, error) {
	// Check if organization location exist
	orgLocExist, err := uc.OrganizationMessage.SendFindOrganizationLocationByIDMessage(request.SendFindOrganizationLocationByIDMessageRequest{
		ID: req.OrganizationLocationID.String(),
	})

	if err != nil {
		uc.Log.Errorf("[MPPlanningUseCase.UpdateLine] " + err.Error())
		return nil, err
	}

	if orgLocExist == nil {
		uc.Log.Errorf("[MPPlanningUseCase.UpdateLine] Organization Location not found")
		return nil, errors.New("Organization Location not found")
	}

	// Check if job level exist
	jobLevelExist, err := uc.JobPlafonMessage.SendFindJobLevelByIDMessage(request.SendFindJobLevelByIDMessageRequest{
		ID: req.JobLevelID.String(),
	})

	if err != nil {
		uc.Log.Errorf("[MPPlanningUseCase.UpdateLine] " + err.Error())
		return nil, err
	}

	if jobLevelExist == nil {
		uc.Log.Errorf("[MPPlanningUseCase.UpdateLine] Job Level not found")
		return nil, errors.New("Job Level not found")
	}

	// Check if job exist
	jobExist, err := uc.JobPlafonMessage.SendFindJobByIDMessage(request.SendFindJobByIDMessageRequest{
		ID: req.JobID.String(),
	})

	if err != nil {
		uc.Log.Errorf("[MPPlanningUseCase.UpdateLine] " + err.Error())
		return nil, err
	}

	if jobExist == nil {
		uc.Log.Errorf("[MPPlanningUseCase.UpdateLine] Job not found")
		return nil, errors.New("Job not found")
	}

	exist, err := uc.MPPlanningRepository.FindLineById(req.ID)

	if err != nil {
		uc.Log.Errorf("[MPPlanningUseCase.UpdateLine] " + err.Error())
		return nil, err
	}

	if exist == nil {
		uc.Log.Errorf("[MPPlanningUseCase.UpdateLine] MP Planning Line not found")
		return nil, errors.New("MP Planning Line not found")
	}

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
		RemainingBalancePH:     req.RemainingBalancePH,
		RemainingBalanceMT:     req.RemainingBalanceMT,
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
		RemainingBalancePH:     mpPlanningLine.RemainingBalancePH,
		RemainingBalanceMT:     mpPlanningLine.RemainingBalanceMT,
		RecruitPH:              mpPlanningLine.RecruitPH,
		RecruitMT:              mpPlanningLine.RecruitMT,
		CreatedAt:              mpPlanningLine.CreatedAt,
		UpdatedAt:              mpPlanningLine.UpdatedAt,
	}, nil
}

func (uc *MPPlanningUseCase) DeleteLine(req *request.DeleteLineMPPlanningLineRequest) error {
	exist, err := uc.MPPlanningRepository.FindLineById(uuid.MustParse(req.ID))

	if err != nil {
		uc.Log.Errorf("[MPPlanningUseCase.UpdateLine] " + err.Error())
		return err
	}

	if exist == nil {
		uc.Log.Errorf("[MPPlanningUseCase.UpdateLine] MP Planning Line not found")
		return errors.New("MP Planning Line not found")
	}

	err = uc.MPPlanningRepository.DeleteLine(uuid.MustParse(req.ID))
	if err != nil {
		uc.Log.Errorf("[MPPlanningUseCase.DeleteLine] " + err.Error())
		return err
	}

	return nil
}

func (uc *MPPlanningUseCase) CreateOrUpdateBatchLineMPPlanningLines(req *request.CreateOrUpdateBatchLineMPPlanningLinesRequest) error {
	for _, line := range req.MPPlanningLines {
		// Check if organization location exist
		orgLocExist, err := uc.OrganizationMessage.SendFindOrganizationLocationByIDMessage(request.SendFindOrganizationLocationByIDMessageRequest{
			ID: line.OrganizationLocationID.String(),
		})

		if err != nil {
			uc.Log.Errorf("[MPPlanningUseCase.CreateOrUpdateBatchLineMPPlanningLines] " + err.Error())
			return err
		}

		if orgLocExist == nil {
			uc.Log.Errorf("[MPPlanningUseCase.CreateOrUpdateBatchLineMPPlanningLines] Organization Location not found")
			return errors.New("Organization Location not found")
		}

		// Check if job level exist
		jobLevelExist, err := uc.JobPlafonMessage.SendFindJobLevelByIDMessage(request.SendFindJobLevelByIDMessageRequest{
			ID: line.JobLevelID.String(),
		})

		if err != nil {
			uc.Log.Errorf("[MPPlanningUseCase.CreateOrUpdateBatchLineMPPlanningLines] " + err.Error())
			return err
		}

		if jobLevelExist == nil {
			uc.Log.Errorf("[MPPlanningUseCase.CreateOrUpdateBatchLineMPPlanningLines] Job Level not found")
			return errors.New("Job Level not found")
		}

		// Check if job exist
		jobExist, err := uc.JobPlafonMessage.SendFindJobByIDMessage(request.SendFindJobByIDMessageRequest{
			ID: line.JobID.String(),
		})

		if err != nil {
			uc.Log.Errorf("[MPPlanningUseCase.CreateOrUpdateBatchLineMPPlanningLines] " + err.Error())
			return err
		}

		if jobExist == nil {
			uc.Log.Errorf("[MPPlanningUseCase.CreateOrUpdateBatchLineMPPlanningLines] Job not found")
			return errors.New("Job not found")
		}

		exist, err := uc.MPPlanningRepository.FindLineById(line.ID)

		if err != nil {
			uc.Log.Errorf("[MPPlanningUseCase.CreateOrUpdateBatchLineMPPlanningLines] " + err.Error())
			return err
		}

		if exist == nil {
			_, err := uc.MPPlanningRepository.CreateLine(&entity.MPPlanningLine{
				ID:                     line.ID,
				MPPlanningHeaderID:     req.MPPlanningHeaderID,
				OrganizationLocationID: &line.OrganizationLocationID,
				JobLevelID:             &line.JobLevelID,
				JobID:                  &line.JobID,
				Existing:               line.Existing,
				Recruit:                line.Recruit,
				SuggestedRecruit:       line.SuggestedRecruit,
				Promotion:              line.Promotion,
				Total:                  line.Total,
				RemainingBalancePH:     line.RemainingBalancePH,
				RemainingBalanceMT:     line.RemainingBalanceMT,
				RecruitPH:              line.RecruitPH,
				RecruitMT:              line.RecruitMT,
			})
			if err != nil {
				uc.Log.Errorf("[MPPlanningUseCase.CreateOrUpdateBatchLineMPPlanningLines] " + err.Error())
				return err
			}
		} else {
			_, err := uc.MPPlanningRepository.UpdateLine(&entity.MPPlanningLine{
				ID:                     line.ID,
				MPPlanningHeaderID:     req.MPPlanningHeaderID,
				OrganizationLocationID: &line.OrganizationLocationID,
				JobLevelID:             &line.JobLevelID,
				JobID:                  &line.JobID,
				Existing:               line.Existing,
				Recruit:                line.Recruit,
				SuggestedRecruit:       line.SuggestedRecruit,
				Promotion:              line.Promotion,
				Total:                  line.Total,
				RemainingBalancePH:     line.RemainingBalancePH,
				RemainingBalanceMT:     line.RemainingBalanceMT,
				RecruitPH:              line.RecruitPH,
				RecruitMT:              line.RecruitMT,
			})
			if err != nil {
				uc.Log.Errorf("[MPPlanningUseCase.CreateOrUpdateBatchLineMPPlanningLines] " + err.Error())
				return err
			}
		}

	}
	return nil
}

func MPPlanningUseCaseFactory(log *logrus.Logger) IMPPlanningUseCase {
	repo := repository.MPPlanningRepositoryFactory(log)
	message := messaging.OrganizationMessageFactory(log)
	jpm := messaging.JobPlafonMessageFactory(log)
	um := messaging.UserMessageFactory(log)
	return NewMPPlanningUseCase(log, repo, message, jpm, um)
}
