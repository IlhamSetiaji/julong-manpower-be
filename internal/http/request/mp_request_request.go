package request

import (
	"github.com/IlhamSetiaji/julong-manpower-be/internal/entity"
	"github.com/google/uuid"
)

type CreateMPRequestHeaderRequest struct {
	ID                         string                     `json:"id" validate:"omitempty,uuid"` // ID
	OrganizationID             uuid.UUID                  `json:"organization_id" validate:"required,uuid"`
	OrganizationLocationID     uuid.UUID                  `json:"organization_location_id" validate:"required,uuid"`
	ForOrganizationID          uuid.UUID                  `json:"for_organization_id" validate:"required,uuid"`
	ForOrganizationLocationID  uuid.UUID                  `json:"for_organization_location_id" validate:"required,uuid"`
	ForOrganizationStructureID uuid.UUID                  `json:"for_organization_structure_id" validate:"required,uuid"`
	JobID                      uuid.UUID                  `json:"job_id" validate:"required,uuid"`
	RequestCategoryID          uuid.UUID                  `json:"request_category_id" validate:"required,uuid"`
	ExpectedDate               string                     `json:"expected_date" validate:"required"`
	Experiences                string                     `json:"experiences" validate:"required"`
	DocumentNumber             string                     `json:"document_number" validate:"required"`
	DocumentDate               string                     `json:"document_date" validate:"required"`
	MaleNeeds                  int                        `json:"male_needs" validate:"omitempty"`
	FemaleNeeds                int                        `json:"female_needs" validate:"omitempty"`
	MinimumAge                 int                        `json:"minimum_age" validate:"required"`
	MaximumAge                 int                        `json:"maximum_age" validate:"required"`
	MinimumExperience          int                        `json:"minimum_experience" validate:"omitempty"`
	MaritalStatus              entity.MaritalStatusEnum   `json:"marital_status" validate:"required,MaritalStatusValidation"`
	MinimumEducation           entity.EducationLevelEnum  `json:"minimum_education" validate:"required,MinimumEducationValidation"`
	RequiredQualification      string                     `json:"required_qualification" validate:"required"`
	Certificate                string                     `json:"certificate" validate:"omitempty"`
	ComputerSkill              string                     `json:"computer_skill" validate:"omitempty"`
	LanguageSkill              string                     `json:"language_skill" validate:"omitempty"`
	OtherSkill                 string                     `json:"other_skill" validate:"omitempty"`
	Jobdesc                    string                     `json:"jobdesc" validate:"required"`
	SalaryMin                  string                     `json:"salary_min" validate:"required"`
	SalaryMax                  string                     `json:"salary_max" validate:"required"`
	RequestorID                *uuid.UUID                 `json:"requestor_id" validate:"required,uuid"`
	DepartmentHead             *uuid.UUID                 `json:"department_head" validate:"omitempty,uuid"`
	VpGmDirector               *uuid.UUID                 `json:"vp_gm_director" validate:"omitempty"`
	CEO                        *uuid.UUID                 `json:"ceo" validate:"omitempty"`
	HrdHoUnit                  *uuid.UUID                 `json:"hrd_ho_unit" validate:"omitempty,uuid"`
	MPPlanningHeaderID         *uuid.UUID                 `json:"mp_planning_header_id" validate:"omitempty,uuid"`
	Status                     entity.MPRequestStatus     `json:"status" validate:"required,MPRequestStatusValidation"`
	MPRequestType              entity.MPRequestTypeEnum   `json:"mp_request_type" validate:"required,MPRequestTypeEnumValidation"`
	RecruitmentType            entity.RecruitmentTypeEnum `json:"recruitment_type" validate:"required,RecruitmentTypeEnumValidation"`
	MajorIDs                   []uuid.UUID                `json:"major_ids" validate:"omitempty,dive,uuid"`
	MPPPeriodID                *uuid.UUID                 `json:"mpp_period_id" validate:"omitempty,uuid"`
	EmpOrganizationID          *uuid.UUID                 `json:"emp_organization_id" validate:"required,uuid"`
	JobLevelID                 *uuid.UUID                 `json:"job_level_id" validate:"required,uuid"`
	IsReplacement              bool                       `json:"is_replacement" validate:"required"`
}

type UpdateMPRequestHeaderRequest struct {
	ID          string                               `json:"id" validate:"required"`
	Status      entity.MPRequestStatus               `json:"status" validate:"required,MPRequestStatusValidation"`
	Notes       string                               `json:"notes" validate:"omitempty"`
	Level       entity.MPRequestApprovalHistoryLevel `json:"level" validate:"required,MPRequestApprovalHistoryLevelValidation"`
	Attachments []ManpowerAttachmentRequest          `json:"attachments" validate:"omitempty,dive"`
	ApprovedBy  string                               `json:"approved_by" validate:"required"`
	ApproverID  uuid.UUID                            `json:"approver_id" validate:"required"`
}
