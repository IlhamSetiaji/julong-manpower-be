package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MPPlanningApprovalHistoryStatus string

const (
	MPPlanningApprovalHistoryStatusApproved MPPlanningApprovalHistoryStatus = "Approved"
	MPPlanningApprovalHistoryStatusRejected MPPlanningApprovalHistoryStatus = "Rejected"
)

type MPPlanningApprovalHistory struct {
	gorm.Model         `json:"-"`
	ID                 uuid.UUID                       `json:"id" gorm:"type:char(36);primaryKey;"`
	MPPlanningHeaderID uuid.UUID                       `json:"mp_planning_header_id" gorm:"type:char(36);"`
	ApproverID         uuid.UUID                       `json:"approver_id" gorm:"type:char(36);"`
	ApproverName       string                          `json:"approver_name" gorm:"type:varchar(255);"`
	Notes              string                          `json:"notes" gorm:"type:text;"`
	Status             MPPlanningApprovalHistoryStatus `json:"status" gorm:"not null"`

	MPPlanningHeader    MPPlanningHeader     `gorm:"foreignKey:MPPlanningHeaderID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	ManpowerAttachments []ManpowerAttachment `json:"manpower_attachments" gorm:"polymorphicType:OwnerType;polymorphicId:OwnerID;polymorphicValue:mp_planning_approval_histories" constraint:"OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (m *MPPlanningApprovalHistory) BeforeCreate(tx *gorm.DB) (err error) {
	m.ID = uuid.New()
	m.CreatedAt = time.Now().Add(7 * time.Hour)
	m.UpdatedAt = time.Now().Add(7 * time.Hour)
	return nil
}

func (m *MPPlanningApprovalHistory) BeforeUpdate(tx *gorm.DB) (err error) {
	m.UpdatedAt = time.Now().Add(7 * time.Hour)
	return nil
}

func (MPPlanningApprovalHistory) TableName() string {
	return "mp_planning_approval_histories"
}
