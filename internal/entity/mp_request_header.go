package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MPRequestStatus string

const (
	MPRequestStatusDraft        MPRequestStatus = "DRAFT"
	MPRequestStatusSubmitted    MPRequestStatus = "SUBMITTED"
	MPRequestStatusRejected     MPRequestStatus = "REJECTED"
	MPRequestStatusApproved     MPRequestStatus = "APPROVED"
	MPRequestStatusNeedApproval MPRequestStatus = "NEED APPROVAL"
	MPRequestStatusCompleted    MPRequestStatus = "COMPLETED"
	MPRequestStatusInProgress   MPRequestStatus = "IN PROGRESS"
)

type MaritalStatusEnum string

const (
	MaritalStatusEnumSingle   MaritalStatusEnum = "single"
	MaritalStatusEnumMarried  MaritalStatusEnum = "married"
	MaritalStatusEnumDivorced MaritalStatusEnum = "divorced"
	MaritalStatusEnumWidowed  MaritalStatusEnum = "widowed"
	MaritalStatusEnumAny      MaritalStatusEnum = "any"
)

type MPRequestTypeEnum string

const (
	MPRequestTypeEnumOnBudget  MPRequestTypeEnum = "ON_BUDGET"
	MPRequestTypeEnumOffBudget MPRequestTypeEnum = "OFF_BUDGET"
)

type RecruitmentTypeEnum string

const (
	RecruitmentTypeEnumMT RecruitmentTypeEnum = "MT_Management Trainee"
	RecruitmentTypeEnumPH RecruitmentTypeEnum = "PH_Professional Hire"
	RecruitmentTypeEnumNS RecruitmentTypeEnum = "NS_Non Staff to Staff"
)

type MPRequestHeader struct {
	gorm.Model                 `json:"-"`
	ID                         uuid.UUID           `json:"id" gorm:"type:char(36);primaryKey;"`
	OrganizationID             *uuid.UUID          `json:"organization_id" gorm:"type:char(36);not null;"`
	OrganizationLocationID     *uuid.UUID          `json:"organization_location_id" gorm:"type:char(36);not null"`
	ForOrganizationID          *uuid.UUID          `json:"for_organization_id" gorm:"type:char(36);not null"`           // For Organization ID
	ForOrganizationLocationID  *uuid.UUID          `json:"for_organization_location_id" gorm:"type:char(36);not null"`  // For Organization Location ID
	ForOrganizationStructureID *uuid.UUID          `json:"for_organization_structure_id" gorm:"type:char(36);not null"` // For Organization Structure ID
	JobID                      *uuid.UUID          `json:"job_id" gorm:"type:char(36);not null"`
	RequestCategoryID          uuid.UUID           `json:"request_category_id" gorm:"type:char(36);not null"`
	ExpectedDate               *time.Time          `json:"expected_date" gorm:"type:date;null;"`
	Experiences                string              `json:"experiences" gorm:"type:text;default:null"` // Experiences in years
	DocumentNumber             string              `json:"document_number" gorm:"type:varchar(255);not null;unique;"`
	DocumentDate               time.Time           `json:"document_date" gorm:"type:date;not null;"`
	MaleNeeds                  int                 `json:"male_needs" gorm:"type:int;default:0"`
	FemaleNeeds                int                 `json:"female_needs" gorm:"type:int;default:0"`
	AnyGender                  int                 `json:"any_gender" gorm:"type:int;default:0"`
	MinimumAge                 int                 `json:"minimum_age" gorm:"type:int;default:0"`
	MaximumAge                 int                 `json:"maximum_age" gorm:"type:int;default:0"`
	MinimumExperience          int                 `json:"minimum_experience" gorm:"type:int;default:0"`
	MaritalStatus              MaritalStatusEnum   `json:"marital_status" gorm:"default:'single'not null"`
	MinimumEducation           EducationLevelEnum  `json:"minimum_education" gorm:"default:'s1';not null"`
	RequiredQualification      string              `json:"required_qualification" gorm:"type:text;default:null"`
	Certificate                string              `json:"certificate" gorm:"type:text;default:null"`
	ComputerSkill              string              `json:"computer_skill" gorm:"type:text;default:null"`
	LanguageSkill              string              `json:"language_skill" gorm:"type:text;default:null"`
	OtherSkill                 string              `json:"other_skill" gorm:"type:text;default:null"`
	Jobdesc                    string              `json:"jobdesc" gorm:"type:text;not null"`
	SalaryMin                  string              `json:"salary_min" gorm:"type:varchar(255);not null"`
	SalaryMax                  string              `json:"salary_max" gorm:"type:varchar(255);not null"`
	RequestorID                *uuid.UUID          `json:"requestor_id" gorm:"type:char(36);not null"`
	DepartmentHead             *uuid.UUID          `json:"department_head" gorm:"type:char(36);null"`
	VpGmDirector               *uuid.UUID          `json:"vp_gm_director" gorm:"type:text;default:null"`
	CEO                        *uuid.UUID          `json:"ceo" gorm:"type:text;default:null"`
	HrdHoUnit                  *uuid.UUID          `json:"hrd_ho_unit" gorm:"type:char(36);null"` // verificator tim rekrutmen
	MPPlanningHeaderID         *uuid.UUID          `json:"mp_planning_header_id" gorm:"type:char(36);null"`
	Status                     MPRequestStatus     `json:"status" gorm:"default:'DRAFT'"`
	MPRequestType              MPRequestTypeEnum   `json:"mp_request_type" gorm:"default:'ON_BUDGET'"`
	RecruitmentType            RecruitmentTypeEnum `json:"recruitment_type" gorm:"type:text;default:not null"`
	MPPPeriodID                uuid.UUID           `json:"mpp_period_id" gorm:"type:char(36);null"`
	NotesDepartmentHead        string              `json:"notes_department_head" gorm:"type:text;default:null"`
	NotesVpGmDirector          string              `json:"notes_vp_gm_director" gorm:"type:text;default:null"`
	NotesCEO                   string              `json:"notes_ceo" gorm:"type:text;default:null"`
	NotesHrdHo                 string              `json:"notes_hrd_ho" gorm:"type:text;default:null"`
	TotalNeeds                 int                 `json:"total_needs" gorm:"type:int;default:0"`
	EmpOrganizationID          *uuid.UUID          `json:"emp_organization_id" gorm:"type:char(36);null"`
	JobLevelID                 *uuid.UUID          `json:"job_level_id" gorm:"type:char(36);null"`
	IsReplacement              bool                `json:"is_replacement" gorm:"default:false"`
	GradeID                    *uuid.UUID          `json:"grade_id" gorm:"type:char(36);null"`

	RequestCategory            RequestCategory            `json:"request_category" gorm:"foreignKey:RequestCategoryID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	RequestMajors              []RequestMajor             `json:"request_majors" gorm:"foreignKey:MPRequestHeaderID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	MPPlanningHeader           MPPlanningHeader           `json:"mp_planning_header" gorm:"foreignKey:MPPlanningHeaderID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	MPRequestApprovalHistories []MPRequestApprovalHistory `json:"mp_request_approval_histories" gorm:"foreignKey:MPRequestHeaderID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	MPPPeriod                  MPPPeriod                  `json:"mpp_period" gorm:"foreignKey:MPPPeriodID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`

	GradeName                 string                 `json:"grade_name" gorm:"-"`
	OrganizationName          string                 `json:"organization_name" gorm:"-"`
	OrganizationCategory      string                 `json:"organization_category" gorm:"-"`
	OrganizationLocationName  string                 `json:"organization_location_name" gorm:"-"`
	ForOrganizationName       string                 `json:"for_organization_name" gorm:"-"`
	ForOrganizationLocation   string                 `json:"for_organization_location" gorm:"-"`
	ForOrganizationStructure  string                 `json:"for_organization_structure" gorm:"-"`
	JobName                   string                 `json:"job_name" gorm:"-"`
	RequestorName             string                 `json:"requestor_name" gorm:"-"`
	DepartmentHeadName        string                 `json:"department_head_name" gorm:"-"`
	VpGmDirectorName          string                 `json:"vp_gm_director_name" gorm:"-"`
	CeoName                   string                 `json:"ceo_name" gorm:"-"`
	HrdHoUnitName             string                 `json:"hrd_ho_unit_name" gorm:"-"`
	EmpOrganizationName       string                 `json:"emp_organization_name" gorm:"-"`
	JobLevelName              string                 `json:"job_level_name" gorm:"-"`
	JobLevel                  int                    `json:"job_level" gorm:"-"`
	RequestorEmployeeJob      map[string]interface{} `json:"requestor_employee_job" gorm:"-"`
	DepartmentHeadEmployeeJob map[string]interface{} `json:"department_head_employee_job" gorm:"-"`
	VpGmDirectorEmployeeJob   map[string]interface{} `json:"vp_gm_director_employee_job" gorm:"-"`
	CeoEmployeeJob            map[string]interface{} `json:"ceo_employee_job" gorm:"-"`
}

func (m *MPRequestHeader) BeforeCreate(tx *gorm.DB) (err error) {
	m.ID = uuid.New()
	m.CreatedAt = time.Now().Add(7 * time.Hour)
	m.UpdatedAt = time.Now().Add(7 * time.Hour)
	return nil
}

func (m *MPRequestHeader) BeforeUpdate(tx *gorm.DB) (err error) {
	m.UpdatedAt = time.Now().Add(7 * time.Hour)
	return nil
}

func (m *MPRequestHeader) BeforeDelete(tx *gorm.DB) (err error) {
	if m.DeletedAt.Valid {
		return nil
	}

	randomString := uuid.New().String()

	m.DocumentNumber = m.DocumentNumber + "_deleted" + randomString

	if err := tx.Model(&m).Where("id = ?", m.ID).Updates((map[string]interface{}{
		"document_number": m.DocumentNumber,
	})).Error; err != nil {
		return err
	}

	return nil
}

func (MPRequestHeader) TableName() string {
	return "mp_request_headers"
}
