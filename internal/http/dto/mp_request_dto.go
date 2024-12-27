package dto

import (
	"time"

	"github.com/IlhamSetiaji/julong-manpower-be/internal/entity"
	"github.com/IlhamSetiaji/julong-manpower-be/internal/http/request"
	"github.com/IlhamSetiaji/julong-manpower-be/internal/http/response"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type IMPRequestDTO interface {
	ConvertToEntity(req *request.CreateMPRequestHeaderRequest) *entity.MPRequestHeader
	ConvertToResponse(ent *entity.MPRequestHeader) *response.MPRequestHeaderResponse
	ConvertEntityToRequest(ent *entity.MPRequestHeader) *request.CreateMPRequestHeaderRequest
	ConvertMPRequestApprovalHistoryToResponse(approvalHistories *entity.MPRequestApprovalHistory) *response.MPRequestApprovalHistoryResponse
	ConvertMPRequestApprovalHistoriesToResponse(approvalHistories *[]entity.MPRequestApprovalHistory) []*response.MPRequestApprovalHistoryResponse
}

type MPRequestDTO struct {
	Log           *logrus.Logger
	MPPlanningDTO IMPPlanningDTO
	Viper         *viper.Viper
}

func NewMPRequestDTO(log *logrus.Logger, mpPlanningDTO IMPPlanningDTO, viper *viper.Viper) IMPRequestDTO {
	return &MPRequestDTO{
		Log:           log,
		MPPlanningDTO: mpPlanningDTO,
		Viper:         viper,
	}
}

func MPRequestDTOFactory(log *logrus.Logger, viper *viper.Viper) IMPRequestDTO {
	mpPlanningDTO := MPPlanningDTOFactory(log)
	return NewMPRequestDTO(log, mpPlanningDTO, viper)
}

func (d *MPRequestDTO) ConvertToEntity(req *request.CreateMPRequestHeaderRequest) *entity.MPRequestHeader {
	expectedDate := parseDate(req.ExpectedDate)
	var parsedID uuid.NullUUID
	if req.ID != "" {
		parsedID = uuid.NullUUID{UUID: uuid.MustParse(req.ID), Valid: true}
	}
	return &entity.MPRequestHeader{
		ID:                         parsedID.UUID,
		OrganizationID:             &req.OrganizationID,
		OrganizationLocationID:     &req.OrganizationLocationID,
		ForOrganizationID:          &req.ForOrganizationID,
		ForOrganizationLocationID:  &req.ForOrganizationLocationID,
		ForOrganizationStructureID: &req.ForOrganizationStructureID,
		JobID:                      &req.JobID,
		RequestCategoryID:          req.RequestCategoryID,
		ExpectedDate:               &expectedDate,
		Experiences:                req.Experiences,
		DocumentDate:               parseDate(req.DocumentDate),
		DocumentNumber:             req.DocumentNumber,
		MaleNeeds:                  req.MaleNeeds,
		FemaleNeeds:                req.FemaleNeeds,
		MinimumAge:                 req.MinimumAge,
		MaximumAge:                 req.MaximumAge,
		MinimumExperience:          req.MinimumExperience,
		MaritalStatus:              req.MaritalStatus,
		MinimumEducation:           req.MinimumEducation,
		RequiredQualification:      req.RequiredQualification,
		Certificate:                req.Certificate,
		ComputerSkill:              req.ComputerSkill,
		LanguageSkill:              req.LanguageSkill,
		OtherSkill:                 req.OtherSkill,
		Jobdesc:                    req.Jobdesc,
		SalaryMin:                  req.SalaryMin,
		SalaryMax:                  req.SalaryMax,
		RequestorID:                req.RequestorID,
		DepartmentHead:             req.DepartmentHead,
		VpGmDirector:               req.VpGmDirector,
		CEO:                        req.CEO,
		HrdHoUnit:                  req.HrdHoUnit,
		MPPlanningHeaderID:         req.MPPlanningHeaderID,
		Status:                     req.Status,
		MPRequestType:              req.MPRequestType,
		MPPPeriodID:                *req.MPPPeriodID,
		EmpOrganizationID:          req.EmpOrganizationID,
		JobLevelID:                 req.JobLevelID,
		IsReplacement:              *req.IsReplacement,
		RecruitmentType:            req.RecruitmentType,
	}
}

func (d *MPRequestDTO) ConvertMPRequestApprovalHistoryToResponse(approvalHistories *entity.MPRequestApprovalHistory) *response.MPRequestApprovalHistoryResponse {
	return &response.MPRequestApprovalHistoryResponse{
		ID:                approvalHistories.ID,
		MPRequestHeaderID: approvalHistories.MPRequestHeaderID,
		ApproverID:        approvalHistories.ApproverID,
		ApproverName:      approvalHistories.ApproverName,
		Notes:             approvalHistories.Notes,
		Level:             approvalHistories.Level,
		Status:            approvalHistories.Status,
		CreatedAt:         approvalHistories.CreatedAt,
		UpdatedAt:         approvalHistories.UpdatedAt,
		Attachments:       ConvertManpowerAttachmentsToResponse(&approvalHistories.ManpowerAttachments, d.Viper),
	}
}

func (d *MPRequestDTO) ConvertMPRequestApprovalHistoriesToResponse(approvalHistories *[]entity.MPRequestApprovalHistory) []*response.MPRequestApprovalHistoryResponse {
	var response []*response.MPRequestApprovalHistoryResponse
	for _, approvalHistory := range *approvalHistories {
		response = append(response, d.ConvertMPRequestApprovalHistoryToResponse(&approvalHistory))
	}
	return response
}

func parseDate(dateStr string) time.Time {
	layout := "2006-01-02"
	date, err := time.Parse(layout, dateStr)
	if err != nil {
		return time.Time{}
	}
	return date
}

func (d *MPRequestDTO) ConvertToResponse(ent *entity.MPRequestHeader) *response.MPRequestHeaderResponse {
	var mpPlanningHeader response.MPPlanningHeaderResponse
	if ent.MPPlanningHeaderID != nil {
		mpPlanningHeader = *d.MPPlanningDTO.ConvertMPPlanningHeaderEntityToResponse(&ent.MPPlanningHeader)
	} else {
		mpPlanningHeader = response.MPPlanningHeaderResponse{}
	}
	var approvedByDeptHead, approvedByVpGmDirector, approvedByCEO, approvedByHrdHoUnit uuid.NullUUID
	if ent.DepartmentHead != nil {
		approvedByDeptHead = uuid.NullUUID{UUID: *ent.DepartmentHead, Valid: true}
	}
	if ent.VpGmDirector != nil {
		approvedByVpGmDirector = uuid.NullUUID{UUID: *ent.VpGmDirector, Valid: true}
	}
	if ent.CEO != nil {
		approvedByCEO = uuid.NullUUID{UUID: *ent.CEO, Valid: true}
	}
	if ent.HrdHoUnit != nil {
		approvedByHrdHoUnit = uuid.NullUUID{UUID: *ent.HrdHoUnit, Valid: true}
	}

	return &response.MPRequestHeaderResponse{
		ID:                         ent.ID,
		OrganizationID:             *ent.OrganizationID,
		OrganizationLocationID:     *ent.OrganizationLocationID,
		ForOrganizationID:          *ent.ForOrganizationID,
		ForOrganizationLocationID:  *ent.ForOrganizationLocationID,
		ForOrganizationStructureID: *ent.ForOrganizationStructureID,
		JobID:                      *ent.JobID,
		RequestCategoryID:          ent.RequestCategoryID,
		ExpectedDate:               *ent.ExpectedDate,
		Experiences:                ent.Experiences,
		DocumentNumber:             ent.DocumentNumber,
		DocumentDate:               ent.DocumentDate,
		MaleNeeds:                  ent.MaleNeeds,
		FemaleNeeds:                ent.FemaleNeeds,
		MinimumAge:                 ent.MinimumAge,
		MaximumAge:                 ent.MaximumAge,
		MinimumExperience:          ent.MinimumExperience,
		MaritalStatus:              ent.MaritalStatus,
		MinimumEducation:           ent.MinimumEducation,
		RequiredQualification:      ent.RequiredQualification,
		Certificate:                ent.Certificate,
		ComputerSkill:              ent.ComputerSkill,
		LanguageSkill:              ent.LanguageSkill,
		OtherSkill:                 ent.OtherSkill,
		Jobdesc:                    ent.Jobdesc,
		SalaryMin:                  ent.SalaryMin,
		SalaryMax:                  ent.SalaryMax,
		RequestorID:                ent.RequestorID,
		DepartmentHead:             ent.DepartmentHead,
		VpGmDirector:               ent.VpGmDirector,
		CEO:                        ent.CEO,
		HrdHoUnit:                  ent.HrdHoUnit,
		MPPlanningHeaderID:         ent.MPPlanningHeaderID,
		Status:                     ent.Status,
		MPRequestType:              ent.MPRequestType,
		RecruitmentType:            ent.RecruitmentType,
		MPPPeriodID:                &ent.MPPPeriodID,
		IsReplacement:              ent.IsReplacement,
		EmpOrganizationID:          ent.EmpOrganizationID,
		JobLevelID:                 ent.JobLevelID,
		CreatedAt:                  ent.CreatedAt,
		UpdatedAt:                  ent.UpdatedAt,

		RequestCategory: map[string]interface{}{
			"ID":            ent.RequestCategory.ID,
			"Name":          ent.RequestCategory.Name,
			"IsReplacement": ent.RequestCategory.IsReplacement,
		},
		RequestMajors: func() []map[string]interface{} {
			var majors []map[string]interface{}
			for _, major := range ent.RequestMajors {
				majors = append(majors, map[string]interface{}{
					"ID": major.ID.String(),
					"Major": map[string]interface{}{
						"ID":    major.MajorID,
						"Major": major.Major.Major,
					},
				})
			}
			return majors
		}(),
		MPPlanningHeader: &mpPlanningHeader,
		// MPPlanningHeader: map[string]interface{}{
		// 	"ID":             ent.MPPlanningHeader.ID,
		// 	"DocumentNumber": ent.MPPlanningHeader.DocumentNumber,
		// 	"DocumentDate":   ent.MPPlanningHeader.DocumentDate,
		// },
		MPPPeriod: response.MPPeriodResponse{
			ID:              ent.MPPPeriod.ID,
			Title:           ent.MPPPeriod.Title,
			StartDate:       ent.MPPPeriod.StartDate.Format("2006-01-02"),
			EndDate:         ent.MPPPeriod.EndDate.Format("2006-01-02"),
			BudgetStartDate: ent.MPPPeriod.BudgetStartDate.Format("2006-01-02"),
			BudgetEndDate:   ent.MPPPeriod.BudgetEndDate.Format("2006-01-02"),
			Status:          ent.MPPPeriod.Status,
			CreatedAt:       ent.MPPPeriod.CreatedAt,
			UpdatedAt:       ent.MPPPeriod.UpdatedAt,
		},

		OrganizationName:         ent.OrganizationName,
		OrganizationCategory:     ent.OrganizationCategory,
		OrganizationLocationName: ent.OrganizationLocationName,
		ForOrganizationName:      ent.ForOrganizationName,
		ForOrganizationLocation:  ent.ForOrganizationLocation,
		ForOrganizationStructure: ent.ForOrganizationStructure,
		JobName:                  ent.JobName,
		RequestorName:            ent.RequestorName,
		DepartmentHeadName:       ent.DepartmentHeadName,
		HrdHoUnitName:            ent.HrdHoUnitName,
		VpGmDirectorName:         ent.VpGmDirectorName,
		CeoName:                  ent.CeoName,
		EmpOrganizationName:      ent.EmpOrganizationName,
		JobLevelName:             ent.JobLevelName,
		JobLevel:                 ent.JobLevel,
		ApprovedByDepartmentHead: approvedByDeptHead.Valid,
		ApprovedByVpGmDirector:   approvedByVpGmDirector.Valid,
		ApprovedByCEO:            approvedByCEO.Valid,
		ApprovedByHrdHoUnit:      approvedByHrdHoUnit.Valid,
		RequestorEmployeeJob:     ent.RequestorEmployeeJob,
	}
}

func (d *MPRequestDTO) ConvertEntityToRequest(ent *entity.MPRequestHeader) *request.CreateMPRequestHeaderRequest {
	return &request.CreateMPRequestHeaderRequest{
		OrganizationID:             *ent.OrganizationID,
		OrganizationLocationID:     *ent.OrganizationLocationID,
		ForOrganizationID:          *ent.ForOrganizationID,
		ForOrganizationLocationID:  *ent.ForOrganizationLocationID,
		ForOrganizationStructureID: *ent.ForOrganizationStructureID,
		JobID:                      *ent.JobID,
		RequestCategoryID:          ent.RequestCategoryID,
		ExpectedDate:               ent.ExpectedDate.Format("2006-01-02"),
		Experiences:                ent.Experiences,
		DocumentNumber:             ent.DocumentNumber,
		DocumentDate:               ent.DocumentDate.Format("2006-01-02"),
		MaleNeeds:                  ent.MaleNeeds,
		FemaleNeeds:                ent.FemaleNeeds,
		MinimumAge:                 ent.MinimumAge,
		MaximumAge:                 ent.MaximumAge,
		MinimumExperience:          ent.MinimumExperience,
		MaritalStatus:              ent.MaritalStatus,
		MinimumEducation:           ent.MinimumEducation,
		RequiredQualification:      ent.RequiredQualification,
		Certificate:                ent.Certificate,
		ComputerSkill:              ent.ComputerSkill,
		LanguageSkill:              ent.LanguageSkill,
		OtherSkill:                 ent.OtherSkill,
		Jobdesc:                    ent.Jobdesc,
		SalaryMin:                  ent.SalaryMin,
		SalaryMax:                  ent.SalaryMax,
		RequestorID:                ent.RequestorID,
		DepartmentHead:             ent.DepartmentHead,
		VpGmDirector:               ent.VpGmDirector,
		CEO:                        ent.CEO,
		HrdHoUnit:                  ent.HrdHoUnit,
		MPPlanningHeaderID:         ent.MPPlanningHeaderID,
		Status:                     ent.Status,
		MPRequestType:              ent.MPRequestType,
		RecruitmentType:            ent.RecruitmentType,
		MPPPeriodID:                &ent.MPPPeriodID,
		EmpOrganizationID:          ent.EmpOrganizationID,
		JobLevelID:                 ent.JobLevelID,
		IsReplacement:              &ent.IsReplacement,
	}
}
