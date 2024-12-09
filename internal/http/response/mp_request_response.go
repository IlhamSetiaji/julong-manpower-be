package response

import (
	"time"

	"github.com/IlhamSetiaji/julong-manpower-be/internal/entity"
	"github.com/google/uuid"
)

type MPRequestHeaderResponse struct {
	ID                         uuid.UUID                  `json:"id"`
	OrganizationID             uuid.UUID                  `json:"organization_id"`
	OrganizationLocationID     uuid.UUID                  `json:"organization_location_id"`
	ForOrganizationID          uuid.UUID                  `json:"for_organization_id"`
	ForOrganizationLocationID  uuid.UUID                  `json:"for_organization_location_id"`
	ForOrganizationStructureID uuid.UUID                  `json:"for_organization_structure_id"`
	JobID                      uuid.UUID                  `json:"job_id"`
	RequestCategoryID          uuid.UUID                  `json:"request_category_id"`
	ExpectedDate               time.Time                  `json:"expected_date"`
	Experiences                string                     `json:"experiences"`
	DocumentNumber             string                     `json:"document_number"`
	DocumentDate               time.Time                  `json:"document_date"`
	MaleNeeds                  int                        `json:"male_needs"`
	FemaleNeeds                int                        `json:"female_needs"`
	MinimumAge                 int                        `json:"minimum_age"`
	MaximumAge                 int                        `json:"maximum_age"`
	MinimumExperience          int                        `json:"minimum_experience"`
	MaritalStatus              entity.MaritalStatusEnum   `json:"marital_status"`
	MinimumEducation           entity.EducationEnum       `json:"minimum_education"`
	RequiredQualification      string                     `json:"required_qualification"`
	Certificate                string                     `json:"certificate"`
	ComputerSkill              string                     `json:"computer_skill"`
	LanguageSkill              string                     `json:"language_skill"`
	OtherSkill                 string                     `json:"other_skill"`
	Jobdesc                    string                     `json:"jobdesc"`
	SalaryMin                  string                     `json:"salary_min"`
	SalaryMax                  string                     `json:"salary_max"`
	RequestorID                *uuid.UUID                 `json:"requestor_id"`
	DepartmentHead             *uuid.UUID                 `json:"department_head"`
	VpGmDirector               string                     `json:"vp_gm_director"`
	CEO                        string                     `json:"ceo"`
	HrdHoUnit                  *uuid.UUID                 `json:"hrd_ho_unit"`
	MPPlanningHeaderID         *uuid.UUID                 `json:"mp_planning_header_id"`
	Status                     entity.MPRequestStatus     `json:"status"`
	MPRequestType              entity.MPRequestTypeEnum   `json:"mp_request_type"`
	RecruitmentType            entity.RecruitmentTypeEnum `json:"recruitment_type"`

	RequestCategory  RequestCategoryResponse  `json:"request_category" gorm:"foreignKey:RequestCategoryID"`
	RequestMajors    []RequestMajorResponse   `json:"request_majors" gorm:"foreignKey:MPRequestHeaderID"`
	MPPlanningHeader MPPlanningHeaderResponse `json:"mp_planning_header" gorm:"foreignKey:MPPlanningHeaderID"`

	OrganizationName         string `json:"organization_name" gorm:"-"`
	OrganizationLocationName string `json:"organization_location_name" gorm:"-"`
	ForOrganizationName      string `json:"for_organization_name" gorm:"-"`
	ForOrganizationLocation  string `json:"for_organization_location" gorm:"-"`
	ForOrganizationStructure string `json:"for_organization_structure" gorm:"-"`
	JobName                  string `json:"job_name" gorm:"-"`
	RequestorName            string `json:"requestor_name" gorm:"-"`
	DepartmentHeadName       string `json:"department_head_name" gorm:"-"`
	HrdHoUnitName            string `json:"hrd_ho_unit_name" gorm:"-"`
}
