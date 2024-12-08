package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ManpowerAttachment struct {
	gorm.Model `json:"-"`
	ID         uuid.UUID `json:"id" gorm:"type:char(36);primaryKey;"`
	OwnerType  string    `json:"owner_type" gorm:"type:varchar(255);not null;"`
	OwnerID    uuid.UUID `json:"owner_id" gorm:"type:char(36);not null;"`
	FileName   string    `json:"file_name" gorm:"type:varchar(255);not null;"`
	FileType   string    `json:"file_type" gorm:"type:varchar(255);not null;"`
	FilePath   string    `json:"file_path" gorm:"type:text;not null;"`
}

func (m *ManpowerAttachment) BeforeCreate(tx *gorm.DB) (err error) {
	m.ID = uuid.New()
	m.CreatedAt = time.Now().Add(7 * time.Hour)
	m.UpdatedAt = time.Now().Add(7 * time.Hour)
	return nil
}

func (m *ManpowerAttachment) BeforeUpdate(tx *gorm.DB) (err error) {
	m.UpdatedAt = time.Now().Add(7 * time.Hour)
	return nil
}

func (ManpowerAttachment) TableName() string {
	return "manpower_attachments"
}
