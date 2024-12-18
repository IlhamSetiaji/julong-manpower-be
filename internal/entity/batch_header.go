package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BatchHeaderApprovalStatus string

const (
	BatchHeaderApprovalStatusApproved     BatchHeaderApprovalStatus = "APPROVED"
	BatchHeaderApprovalStatusRejected     BatchHeaderApprovalStatus = "REJECTED"
	BatchHeaderApprovalStatusNeedApproval BatchHeaderApprovalStatus = "NEED_APPROVAL"
)

type BatchHeader struct {
	gorm.Model     `json:"-"`
	ID             uuid.UUID                 `json:"id" gorm:"type:char(36);primaryKey;"`
	DocumentNumber string                    `json:"document_number" gorm:"type:varchar(255);not null;"`
	DocumentDate   time.Time                 `json:"document_date" gorm:"default:null;"`
	ApproverID     *uuid.UUID                `json:"approver_id" gorm:"type:char(36);default:null;"`
	ApproverName   string                    `json:"approver_name" gorm:"type:varchar(255);default:null;"`
	Status         BatchHeaderApprovalStatus `json:"status" gorm:"default:null"`

	BatchLines []BatchLine `json:"batch_lines" gorm:"foreignKey:BatchHeaderID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (m *BatchHeader) BeforeCreate(tx *gorm.DB) (err error) {
	m.ID = uuid.New()
	m.CreatedAt = time.Now().Add(time.Hour * 7)
	m.UpdatedAt = time.Now().Add(time.Hour * 7)
	return nil
}

func (m *BatchHeader) BeforeUpdate(tx *gorm.DB) (err error) {
	m.UpdatedAt = time.Now().Add(time.Hour * 7)
	return nil
}

func (BatchHeader) TableName() string {
	return "batch_headers"
}
