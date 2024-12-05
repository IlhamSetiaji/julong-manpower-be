package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MPPlanningHeaderAttachment struct {
	gorm.Model         `json:"-"`
	ID                 uuid.UUID `json:"id" gorm:"type:char(36);primaryKey;"`
	MPPlanningHeaderID uuid.UUID `json:"mp_planning_header_id" gorm:"type:char(36);"`
	FileName           string    `json:"file_name" gorm:"type:varchar(255);not null;"`
	FileType           string    `json:"file_type" gorm:"type:varchar(255);not null;"`
	FilePath           string    `json:"file_path" gorm:"type:text;not null;"`

	MPPlanningHeader MPPlanningHeader `json:"mp_planning_header" gorm:"foreignKey:MPPlanningHeaderID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (m *MPPlanningHeaderAttachment) BeforeCreate(tx *gorm.DB) (err error) {
	m.ID = uuid.New()
	m.CreatedAt = time.Now().Add(7 * time.Hour)
	m.UpdatedAt = time.Now().Add(7 * time.Hour)
	return nil
}

func (m *MPPlanningHeaderAttachment) BeforeUpdate(tx *gorm.DB) (err error) {
	m.UpdatedAt = time.Now().Add(7 * time.Hour)
	return nil
}

func (MPPlanningHeaderAttachment) TableName() string {
	return "mp_planning_header_attachments"
}
