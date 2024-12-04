package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MPRequestHeader struct {
	gorm.Model             `json:"-"`
	ID                     uuid.UUID  `json:"id" gorm:"type:char(32);primaryKey;"`
	OrganizationID         *uuid.UUID `json:"organization_id" gorm:"type:char(32);not null;"`
	OrganizationLocationID *uuid.UUID `json:"organization_location_id" gorm:"type:char(32);not null"`
	JobID                  *uuid.UUID `json:"job_id" gorm:"type:char(32);not null"`
	RequestTypeID          *uuid.UUID `json:"request_type_id" gorm:"type:char(32);not null"`
	DocumentNumber         string     `json:"document_number" gorm:"type:varchar(255);not null;unique;"`
	DocumentDate           time.Time  `json:"document_date" gorm:"type:date;not null;"`
	MaleNeeds              int        `json:"male_needs" gorm:"type:int;default:0"`
	FemaleNeeds            int        `json:"female_needs" gorm:"type:int;default:0"`
	MinimumAge             int        `json:"minimum_age" gorm:"type:int;default:0"`
	MaximumAge             int        `json:"maximum_age" gorm:"type:int;default:0"`
	MinimumExperience      int        `json:"minimum_experience" gorm:"type:int;default:0"`
	MaritalStatus          string     `json:"marital_status" gorm:"type:enum('single', 'married', 'divorced', 'widowed');default:'single'"`
	MinimumEducation       string     `json:"minimum_education" gorm:"type:enum('sd', 'smp', 'sma', 'd3', 's1', 's2', 's3');default:'s1'"`
	JobMajorID             uuid.UUID  `json:"job_major_id" gorm:"type:char(32);"`
	RequiredQualification  string     `json:"required_qualification" gorm:"type:text;"`
	Certificate            string     `json:"certificate" gorm:"type:text;default:null"`
	ComputerSkill          string     `json:"computer_skill" gorm:"type:text;default:null"`
	LanguageSkill          string     `json:"language_skill" gorm:"type:text;default:null"`
	OtherSkill             string     `json:"other_skill" gorm:"type:text;default:null"`
	Jobdesc                string     `json:"jobdesc" gorm:"type:text;default:null"`
	Salary                 string     `json:"salary" gorm:"type:text;not null"`
	CreatedBy              *uuid.UUID `json:"created_by" gorm:"type:char(32);not null"`

	RequestType RequestType `json:"request_type" gorm:"foreignKey:RequestTypeID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	JobMajor    JobMajor    `json:"job_major" gorm:"foreignKey:JobMajorID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
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
