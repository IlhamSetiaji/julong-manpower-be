package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EducationLevelEnum string

const (
	EducationLevelEnumDoctoral EducationLevelEnum = "1 - Doctoral / Professor"
	EducationLevelEnumMaster   EducationLevelEnum = "2 - Master Degree"
	EducationLevelEnumBachelor EducationLevelEnum = "3 - Bachelor"
	EducationLevelEnumD1       EducationLevelEnum = "4 - Diploma 1"
	EducationLevelEnumD2       EducationLevelEnum = "5 - Diploma 2"
	EducationLevelEnumD3       EducationLevelEnum = "6 - Diploma 3"
	EducationLevelEnumD4       EducationLevelEnum = "7 - Diploma 4"
	EducationLevelEnumSD       EducationLevelEnum = "8 - Elementary School"
	EducationLevelEnumSMA      EducationLevelEnum = "9 - Senior High School"
	EducationLevelEnumSMP      EducationLevelEnum = "10 - Junior High School"
	EducationLevelEnumTK       EducationLevelEnum = "11 - Unschooled"
)

type Major struct {
	gorm.Model     `json:"-"`
	ID             uuid.UUID          `json:"id" gorm:"type:char(36);primaryKey;"`
	Major          string             `json:"major" gorm:"type:varchar(255);not null;"`
	EducationLevel EducationLevelEnum `json:"education_level" gorm:"type:varchar(255);not null;"`

	RequestMajors []RequestMajor `json:"request_majors" gorm:"foreignKey:MajorID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (m *Major) BeforeCreate(tx *gorm.DB) (err error) {
	m.ID = uuid.New()
	m.CreatedAt = time.Now().Add(7 * time.Hour)
	m.UpdatedAt = time.Now().Add(7 * time.Hour)
	return nil
}

func (m *Major) BeforeUpdate(tx *gorm.DB) (err error) {
	m.UpdatedAt = time.Now().Add(7 * time.Hour)
	return nil
}

func (Major) TableName() string {
	return "majors"
}
