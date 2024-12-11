package dto

import (
	"time"

	"github.com/IlhamSetiaji/julong-manpower-be/internal/entity"
	"github.com/IlhamSetiaji/julong-manpower-be/internal/http/request"
	"github.com/IlhamSetiaji/julong-manpower-be/internal/http/response"
)

func ConvertToEntity(req *request.CreateMPRequestHeaderRequest) *entity.MPRequestHeader {
	expectedDate := parseDate(req.ExpectedDate)
	return &entity.MPRequestHeader{
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
	}
}

func parseDate(dateStr string) time.Time {
	layout := "2006-01-02"
	date, err := time.Parse(layout, dateStr)
	if err != nil {
		return time.Time{}
	}
	return date
}

func ConvertToResponse(ent *entity.MPRequestHeader) *response.MPRequestHeaderResponse {
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

		RequestCategory: map[string]interface{}{
			"ID":   ent.RequestCategory.ID,
			"Name": ent.RequestCategory.Name,
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
		MPPlanningHeader: map[string]interface{}{
			"ID":             ent.MPPlanningHeader.ID,
			"DocumentNumber": ent.MPPlanningHeader.DocumentNumber,
			"DocumentDate":   ent.MPPlanningHeader.DocumentDate,
		},

		OrganizationName:         ent.OrganizationName,
		OrganizationLocationName: ent.OrganizationLocationName,
		ForOrganizationName:      ent.ForOrganizationName,
		ForOrganizationLocation:  ent.ForOrganizationLocation,
		ForOrganizationStructure: ent.ForOrganizationStructure,
		JobName:                  ent.JobName,
		RequestorName:            ent.RequestorName,
		DepartmentHeadName:       ent.DepartmentHeadName,
		HrdHoUnitName:            ent.HrdHoUnitName,
	}
}
