package dto

import (
	"sort"
	"strconv"

	"github.com/IlhamSetiaji/julong-manpower-be/internal/entity"
	"github.com/IlhamSetiaji/julong-manpower-be/internal/http/messaging"
	"github.com/IlhamSetiaji/julong-manpower-be/internal/http/request"
	"github.com/IlhamSetiaji/julong-manpower-be/internal/http/response"
	"github.com/IlhamSetiaji/julong-manpower-be/internal/repository"
	"github.com/sirupsen/logrus"
)

type IBatchDTO interface {
	ConvertBatchHeaderEntityToResponse(batch *entity.BatchHeader) *response.BatchResponse
	ConvertToDocumentBatchResponse(batch *entity.BatchHeader, operatingUnit string) *response.DocumentBatchResponse
	ConvertToDocumentCalculationBatchResponse(mpPlanningLine entity.MPPlanningLine, isTotal bool, planningLines *[]entity.MPPlanningLine) *response.DocumentCalculationBatchResponse
	ConvertDocumentCalculationBatchResponses(mpPlanningLines []entity.MPPlanningLine) []response.DocumentCalculationBatchResponse
	ConvertRealDocumentBatchResponse(batch *entity.BatchHeader) *response.RealDocumentBatchResponse
}

type BatchDTO struct {
	Log              *logrus.Logger
	BatchLineDTO     IBatchLineDTO
	JobPlafonMessage messaging.IJobPlafonMessage
	OrgMessage       messaging.IOrganizationMessage
	mppPeriodRepo    repository.IMPPPeriodRepository
}

func NewBatchDTO(log *logrus.Logger, batchLineDTO IBatchLineDTO, jpm messaging.IJobPlafonMessage, orgMessage messaging.IOrganizationMessage, mppPeriodRepo repository.IMPPPeriodRepository) IBatchDTO {
	return &BatchDTO{
		Log:              log,
		BatchLineDTO:     batchLineDTO,
		JobPlafonMessage: jpm,
		OrgMessage:       orgMessage,
		mppPeriodRepo:    mppPeriodRepo,
	}
}

func (d *BatchDTO) ConvertRealDocumentBatchResponse(batch *entity.BatchHeader) *response.RealDocumentBatchResponse {
	return &response.RealDocumentBatchResponse{
		ApproverType: batch.ApproverType,
		Overall:      *d.ConvertToDocumentBatchResponse(batch, "Julong Group"),
		OrganizationOverall: func() []response.OrganizationOverallResponse {
			var organizationOverall []response.OrganizationOverallResponse
			// group batch lines by organization id
			groupedBatchLines := make(map[string][]entity.BatchLine)
			for _, bl := range batch.BatchLines {
				groupedBatchLines[bl.OrganizationID.String()] = append(groupedBatchLines[bl.OrganizationID.String()], bl)
			}
			for orgID, bls := range groupedBatchLines {
				// check org name
				messageResponse, err := d.OrgMessage.SendFindOrganizationByIDMessage(request.SendFindOrganizationByIDMessageRequest{
					ID: orgID,
				})
				if err != nil {
					d.Log.Errorf("[MPPlanningUseCase.FindAllLinesByHeaderIdPaginated Message] " + err.Error())
				}
				orgName := messageResponse.Name
				organizationOverall = append(organizationOverall, response.OrganizationOverallResponse{
					Overall: *d.ConvertToDocumentBatchResponse(&entity.BatchHeader{
						BatchLines: bls,
					}, orgName),
					LocationOverall: func() []response.DocumentBatchResponse {
						var locationOverall []response.DocumentBatchResponse
						// group batch lines by organization location id
						groupedBatchLines := make(map[string][]entity.BatchLine)
						for _, bl := range bls {
							groupedBatchLines[bl.OrganizationLocationID.String()] = append(groupedBatchLines[bl.OrganizationLocationID.String()], bl)
						}
						for locID, bls := range groupedBatchLines {
							// check location name
							messageResponse, err := d.OrgMessage.SendFindOrganizationLocationByIDMessage(request.SendFindOrganizationLocationByIDMessageRequest{
								ID: locID,
							})
							if err != nil {
								d.Log.Errorf("[MPPlanningUseCase.FindAllLinesByHeaderIdPaginated Message] " + err.Error())
							}
							orgLocationName := messageResponse.Name
							locationOverall = append(locationOverall, *d.ConvertToDocumentBatchResponse(&entity.BatchHeader{
								BatchLines: bls,
							}, orgLocationName+" ("+orgName+")"))
						}
						return locationOverall
					}(),
				})
			}
			return organizationOverall
		}(),
	}
}

func (d *BatchDTO) ConvertBatchHeaderEntityToResponse(batch *entity.BatchHeader) *response.BatchResponse {
	return &response.BatchResponse{
		ID:             batch.ID,
		DocumentNumber: batch.DocumentNumber,
		DocumentDate:   batch.DocumentDate,
		Status:         string(batch.Status),
		CreatedAt:      batch.CreatedAt,
		UpdatedAt:      batch.UpdatedAt,
		BatchLines:     d.BatchLineDTO.ConvertBatchLineEntitiesToResponse(&batch.BatchLines),
	}
}

func (d *BatchDTO) ConvertToDocumentBatchResponse(batch *entity.BatchHeader, operatingUnit string) *response.DocumentBatchResponse {
	currentMppPeriod, err := d.mppPeriodRepo.FindByStatus(entity.MPPeriodStatusOpen)
	if err != nil {
		d.Log.Errorf("[BatchDTO.ConvertToDocumentBatchResponse] " + err.Error())
	}
	var budgetYear string = "2024"
	var budgetRange string = "2025 Budget (Sep24-Aug25)"
	var existingDate string = "Sep24"
	if currentMppPeriod != nil {
		budgetYear = currentMppPeriod.BudgetStartDate.Format("2006") + "/" + currentMppPeriod.BudgetEndDate.Format("2006")
		budgetRange = currentMppPeriod.BudgetEndDate.Format("2006") + " Budget (" + currentMppPeriod.BudgetStartDate.Format("Jan06") + "-" + currentMppPeriod.BudgetEndDate.Format("Jan06") + ")"
		existingDate = currentMppPeriod.BudgetStartDate.Format("Jan06")
	}
	return &response.DocumentBatchResponse{
		OperatingUnit: operatingUnit,
		BudgetYear:    budgetYear,
		BudgetRange:   budgetRange,
		ExistingDate:  existingDate,
		Grade: func() response.GradeBatchResponse {
			var executivePromote int
			var totalExecutive int
			var totalNonExecutive int
			var totalExistingCount int
			var totalPromoteCount int
			var totalRecruitCount int
			return response.GradeBatchResponse{
				Executive: func() []response.DocumentCalculationBatchResponse {
					var executive []response.DocumentCalculationBatchResponse
					groupedByJobLevel := make(map[int]*response.DocumentCalculationBatchResponse)

					// First pass: Group data by job level and calculate existing, promote, and recruit
					for _, bl := range batch.BatchLines {
						for _, mpl := range bl.MPPlanningHeader.MPPlanningLines {
							// Check job level name
							message2Response, err := d.JobPlafonMessage.SendFindJobLevelByIDMessage(request.SendFindJobLevelByIDMessageRequest{
								ID: mpl.JobLevelID.String(),
							})
							if err != nil {
								d.Log.Errorf("[MPPlanningUseCase.FindAllLinesByHeaderIdPaginated Message] " + err.Error())
							}
							mpl.JobLevelName = message2Response.Name
							mpl.JobLevel = int(message2Response.Level)

							// Group mpl by job level (only for job levels > 3)
							if mpl.JobLevel > 3 {
								if _, exists := groupedByJobLevel[mpl.JobLevel]; !exists {
									groupedByJobLevel[mpl.JobLevel] = &response.DocumentCalculationBatchResponse{
										JobLevelName: strconv.Itoa(mpl.JobLevel),
										Existing:     0,
										Promote:      0,
										Recruit:      0,
										Total:        0,
									}
								}

								groupedByJobLevel[mpl.JobLevel].Existing += mpl.Existing
								groupedByJobLevel[mpl.JobLevel].Promote += mpl.Promotion
								groupedByJobLevel[mpl.JobLevel].Recruit += mpl.RecruitPH + mpl.RecruitMT
							}
						}
					}

					// Second pass: Calculate total by subtracting previous job level's promote
					// Sort job levels in descending order
					sortedJobLevels := make([]int, 0, len(groupedByJobLevel))
					for jobLevel := range groupedByJobLevel {
						sortedJobLevels = append(sortedJobLevels, jobLevel)
					}
					sort.Sort(sort.Reverse(sort.IntSlice(sortedJobLevels)))

					for i, jobLevel := range sortedJobLevels {
						if i > 0 {
							// Subtract the promote value of the previous job level
							previousJobLevel := sortedJobLevels[i-1]
							groupedByJobLevel[jobLevel].Total = groupedByJobLevel[jobLevel].Existing +
								groupedByJobLevel[jobLevel].Promote +
								groupedByJobLevel[jobLevel].Recruit -
								groupedByJobLevel[previousJobLevel].Promote

							if jobLevel == 4 {
								executivePromote = groupedByJobLevel[jobLevel].Promote
							}
						} else {
							// For the highest job level, no subtraction is needed
							groupedByJobLevel[jobLevel].Total = groupedByJobLevel[jobLevel].Existing +
								groupedByJobLevel[jobLevel].Promote +
								groupedByJobLevel[jobLevel].Recruit
						}

						totalExecutive += groupedByJobLevel[jobLevel].Total
					}

					// Convert the map to a slice
					for _, jobLevel := range sortedJobLevels {
						totalExistingCount += groupedByJobLevel[jobLevel].Existing
						totalPromoteCount += groupedByJobLevel[jobLevel].Promote
						totalRecruitCount += groupedByJobLevel[jobLevel].Recruit
						executive = append(executive, *groupedByJobLevel[jobLevel])
					}

					return executive
				}(),
				NonExecutive: func() []response.DocumentCalculationBatchResponse {
					var nonExecutive []response.DocumentCalculationBatchResponse
					groupedByJobLevel := make(map[int]*response.DocumentCalculationBatchResponse)

					// First pass: Group data by job level and calculate existing, promote, and recruit
					for _, bl := range batch.BatchLines {
						for _, mpl := range bl.MPPlanningHeader.MPPlanningLines {
							message2Response, err := d.JobPlafonMessage.SendFindJobLevelByIDMessage(request.SendFindJobLevelByIDMessageRequest{
								ID: mpl.JobLevelID.String(),
							})
							if err != nil {
								d.Log.Errorf("[MPPlanningUseCase.FindAllLinesByHeaderIdPaginated Message] " + err.Error())
								continue
							}
							mpl.JobLevelName = message2Response.Name
							mpl.JobLevel = int(message2Response.Level)

							// Group mpl by job level (only for job levels <= 3)
							if mpl.JobLevel <= 3 {
								if _, exists := groupedByJobLevel[mpl.JobLevel]; !exists {
									groupedByJobLevel[mpl.JobLevel] = &response.DocumentCalculationBatchResponse{
										JobLevelName: strconv.Itoa(mpl.JobLevel),
										Existing:     0,
										Promote:      0,
										Recruit:      0,
										Total:        0,
									}
								}

								groupedByJobLevel[mpl.JobLevel].Existing += mpl.Existing
								groupedByJobLevel[mpl.JobLevel].Promote += mpl.Promotion
								groupedByJobLevel[mpl.JobLevel].Recruit += mpl.RecruitPH + mpl.RecruitMT
							}
						}
					}

					// Second pass: Calculate total by subtracting previous job level's promote
					// Sort job levels in descending order
					sortedJobLevels := make([]int, 0, len(groupedByJobLevel))
					for jobLevel := range groupedByJobLevel {
						sortedJobLevels = append(sortedJobLevels, jobLevel)
					}
					sort.Sort(sort.Reverse(sort.IntSlice(sortedJobLevels)))

					for i, jobLevel := range sortedJobLevels {
						if sortedJobLevels[i] == 3 {
							groupedByJobLevel[jobLevel].Total = groupedByJobLevel[jobLevel].Existing +
								groupedByJobLevel[jobLevel].Promote +
								groupedByJobLevel[jobLevel].Recruit - executivePromote
						} else {
							if i > 0 {
								previousJobLevel := sortedJobLevels[i-1]
								// Check if the previous job level is exactly one level higher
								if previousJobLevel == jobLevel+1 {
									groupedByJobLevel[jobLevel].Total = groupedByJobLevel[jobLevel].Existing +
										groupedByJobLevel[jobLevel].Promote +
										groupedByJobLevel[jobLevel].Recruit -
										groupedByJobLevel[previousJobLevel].Promote
								} else {
									// If not, skip subtraction
									groupedByJobLevel[jobLevel].Total = groupedByJobLevel[jobLevel].Existing +
										groupedByJobLevel[jobLevel].Promote +
										groupedByJobLevel[jobLevel].Recruit
								}
							} else {
								// For the highest job level, no subtraction is needed
								groupedByJobLevel[jobLevel].Total = groupedByJobLevel[jobLevel].Existing +
									groupedByJobLevel[jobLevel].Promote +
									groupedByJobLevel[jobLevel].Recruit
							}
						}

						totalNonExecutive += groupedByJobLevel[jobLevel].Total
					}

					// Convert the map to a slice
					for _, jobLevel := range sortedJobLevels {
						totalExistingCount += groupedByJobLevel[jobLevel].Existing
						totalPromoteCount += groupedByJobLevel[jobLevel].Promote
						totalRecruitCount += groupedByJobLevel[jobLevel].Recruit
						nonExecutive = append(nonExecutive, *groupedByJobLevel[jobLevel])
					}

					return nonExecutive
				}(),
				Total: func() []response.DocumentCalculationBatchResponse {
					var total []response.DocumentCalculationBatchResponse
					var totalExisting, totalPromote, totalRecruit, totalOverall int
					for _, bl := range batch.BatchLines {
						for _, mpl := range bl.MPPlanningHeader.MPPlanningLines {
							totalExisting += mpl.Existing
							totalPromote += mpl.Promotion
							totalRecruit += mpl.RecruitPH + mpl.RecruitMT

							message2Response, err := d.JobPlafonMessage.SendFindJobLevelByIDMessage(request.SendFindJobLevelByIDMessageRequest{
								ID: mpl.JobLevelID.String(),
							})
							if err != nil {
								d.Log.Errorf("[MPPlanningUseCase.FindAllLinesByHeaderIdPaginated Message] " + err.Error())
							}
							mpl.JobLevelName = message2Response.Name
							mpl.JobLevel = int(message2Response.Level)

						}
					}
					totalOverall += totalExecutive + totalNonExecutive
					total = append(total, response.DocumentCalculationBatchResponse{
						JobLevelName: "Total",
						Existing:     totalExistingCount,
						Promote:      totalPromoteCount,
						Recruit:      totalRecruitCount,
						Total:        totalOverall,
						IsTotal:      true,
					})
					return total
				}(),
			}
		}(),
	}
}

func (d *BatchDTO) findPreviousMPPlanningLineByJobLevel(jobLevel int, mpPlanningLines []entity.MPPlanningLine) (*entity.MPPlanningLine, error) {
	for _, mpl := range mpPlanningLines {
		message2Response, err := d.JobPlafonMessage.SendFindJobLevelByIDMessage(request.SendFindJobLevelByIDMessageRequest{
			ID: mpl.JobLevelID.String(),
		})
		if err != nil {
			d.Log.Errorf("[MPPlanningUseCase.findPreviousMPPlanningLineByJobLevel Message] " + err.Error())
			return nil, err
		}
		mpl.JobLevelName = message2Response.Name
		mpl.JobLevel = int(message2Response.Level)
		if mpl.JobLevel == jobLevel {
			return &mpl, nil
		}
	}
	return nil, nil
}

func (d *BatchDTO) ConvertToDocumentCalculationBatchResponse(mpPlanningLine entity.MPPlanningLine, isTotal bool, planningLines *[]entity.MPPlanningLine) *response.DocumentCalculationBatchResponse {
	message2Response, err := d.JobPlafonMessage.SendFindJobLevelByIDMessage(request.SendFindJobLevelByIDMessageRequest{
		ID: mpPlanningLine.JobLevelID.String(),
	})
	if err != nil {
		d.Log.Errorf("[MPPlanningUseCase.FindAllLinesByHeaderIdPaginated Message] " + err.Error())
	}
	mpPlanningLine.JobLevelName = message2Response.Name
	mpPlanningLine.JobLevel = int(message2Response.Level)

	var totalOverall int = 0
	if planningLines != nil {
		previousMpPlanningLine, err := d.findPreviousMPPlanningLineByJobLevel(mpPlanningLine.JobLevel+1, *planningLines)
		if err != nil {
			d.Log.Errorf("[BatchDTO.ConvertToDocumentBatchResponse] " + err.Error())
		}
		if previousMpPlanningLine != nil {
			totalOverall += mpPlanningLine.Existing + mpPlanningLine.Promotion + mpPlanningLine.RecruitPH + mpPlanningLine.RecruitMT - previousMpPlanningLine.Promotion
		} else {
			totalOverall += mpPlanningLine.Existing + mpPlanningLine.Promotion + mpPlanningLine.RecruitPH + mpPlanningLine.RecruitMT
		}
	} else {
		totalOverall += mpPlanningLine.Existing + mpPlanningLine.Promotion + mpPlanningLine.RecruitPH + mpPlanningLine.RecruitMT
	}
	return &response.DocumentCalculationBatchResponse{
		JobLevelName: strconv.Itoa(mpPlanningLine.JobLevel),
		Existing:     mpPlanningLine.Existing,
		Promote:      mpPlanningLine.Promotion,
		Recruit:      mpPlanningLine.RecruitPH + mpPlanningLine.RecruitMT,
		Total:        totalOverall,
		IsTotal:      isTotal,
	}
}

func (d *BatchDTO) ConvertDocumentCalculationBatchResponses(mpPlanningLines []entity.MPPlanningLine) []response.DocumentCalculationBatchResponse {
	var responses []response.DocumentCalculationBatchResponse
	for _, mpPlanningLine := range mpPlanningLines {
		responses = append(responses, *d.ConvertToDocumentCalculationBatchResponse(mpPlanningLine, false, &mpPlanningLines))
	}
	return responses
}

func BatchDTOFactory(log *logrus.Logger) IBatchDTO {
	batchLineDTO := BatchLineDTOFactory(log)
	jpm := messaging.JobPlafonMessageFactory(log)
	orgMessage := messaging.OrganizationMessageFactory(log)
	mppPeriodRepo := repository.MPPPeriodRepositoryFactory(log)
	return NewBatchDTO(log, batchLineDTO, jpm, orgMessage, mppPeriodRepo)
}
