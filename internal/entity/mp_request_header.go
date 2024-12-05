package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MPRequestStatus string

const (
	MPRequestStatusDraft     MPRequestStatus = "DRAFT"
	MPRequestStatusSubmitted MPRequestStatus = "SUBMITTED"
	MPRequestStatusRejected  MPRequestStatus = "REJECTED"
	MPRequestStatusApproved  MPRequestStatus = "APPROVED"
)

type MaritalStatusEnum string

const (
	MaritalStatusEnumSingle   MaritalStatusEnum = "single"
	MaritalStatusEnumMarried  MaritalStatusEnum = "married"
	MaritalStatusEnumDivorced MaritalStatusEnum = "divorced"
	MaritalStatusEnumWidowed  MaritalStatusEnum = "widowed"
)

type EducationEnum string

const (
	EducationEnumSD  EducationEnum = "sd"
	EducationEnumSMP EducationEnum = "smp"
	EducationEnumSMA EducationEnum = "sma"
	EducationEnumD3  EducationEnum = "d3"
	EducationEnumS1  EducationEnum = "s1"
	EducationEnumS2  EducationEnum = "s2"
	EducationEnumS3  EducationEnum = "s3"
)

type MPRequestTypeEnum string

const (
	MPRequestTypeEnumOnBudget  MPRequestTypeEnum = "ON_BUDGET"
	MPRequestTypeEnumOffBudget MPRequestTypeEnum = "OFF_BUDGET"
)

type MPRequestHeader struct {
	gorm.Model             `json:"-"`
	ID                     uuid.UUID         `json:"id" gorm:"type:char(36);primaryKey;"`
	OrganizationID         *uuid.UUID        `json:"organization_id" gorm:"type:char(36);not null;"`
	OrganizationLocationID *uuid.UUID        `json:"organization_location_id" gorm:"type:char(36);not null"`
	JobID                  *uuid.UUID        `json:"job_id" gorm:"type:char(36);not null"`
	RequestCategoryID      *uuid.UUID        `json:"request_category_id" gorm:"type:char(36);not null"`
	ExpectedDate           *time.Time        `json:"expected_date" gorm:"type:date;null;"`
	Experiences            string            `json:"experiences" gorm:"type:text;default:null"` // Experiences in years
	DocumentNumber         string            `json:"document_number" gorm:"type:varchar(255);not null;unique;"`
	DocumentDate           time.Time         `json:"document_date" gorm:"type:date;not null;"`
	MaleNeeds              int               `json:"male_needs" gorm:"type:int;default:0"`
	FemaleNeeds            int               `json:"female_needs" gorm:"type:int;default:0"`
	MinimumAge             int               `json:"minimum_age" gorm:"type:int;default:0"`
	MaximumAge             int               `json:"maximum_age" gorm:"type:int;default:0"`
	MinimumExperience      int               `json:"minimum_experience" gorm:"type:int;default:0"`
	MaritalStatus          MaritalStatusEnum `json:"marital_status" gorm:"default:'single'not null"`
	MinimumEducation       EducationEnum     `json:"minimum_education" gorm:"default:'s1';not null"`
	RequiredQualification  string            `json:"required_qualification" gorm:"type:text;default:null"`
	Certificate            string            `json:"certificate" gorm:"type:text;default:null"`
	ComputerSkill          string            `json:"computer_skill" gorm:"type:text;default:null"`
	LanguageSkill          string            `json:"language_skill" gorm:"type:text;default:null"`
	OtherSkill             string            `json:"other_skill" gorm:"type:text;default:null"`
	Jobdesc                string            `json:"jobdesc" gorm:"type:text;not null"`
	SalaryRange            string            `json:"salary_range" gorm:"type:text;not null"`
	RequestorID            *uuid.UUID        `json:"requestor_id" gorm:"type:char(36);not null"`
	DepartmentHead         *uuid.UUID        `json:"department_head" gorm:"type:char(36);not null"`
	VpGmDirector           string            `json:"vp_gm_director" gorm:"type:text;default:null"`
	CEO                    string            `json:"ceo" gorm:"type:text;default:null"`
	HrdHoUnit              *uuid.UUID        `json:"hrd_ho_unit" gorm:"type:char(36);not null"`
	MPPlanningHeaderID     *uuid.UUID        `json:"mp_planning_header_id" gorm:"type:char(36);not null"`
	Status                 MPRequestStatus   `json:"status" gorm:"default:'DRAFT'"`
	MPRequestType          MPRequestTypeEnum `json:"mp_request_type" gorm:"default:'ON_BUDGET'"`

	RequestCategory  RequestCategory  `json:"request_category" gorm:"foreignKey:RequestCategoryID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	RequestMajors    []RequestMajor   `json:"request_majors" gorm:"foreignKey:MPRequestHeaderID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	MPPlanningHeader MPPlanningHeader `json:"mp_planning_header" gorm:"foreignKey:MPPlanningHeaderID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
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

func (MPRequestHeader) TableName() string {
	return "mp_request_headers"
}
