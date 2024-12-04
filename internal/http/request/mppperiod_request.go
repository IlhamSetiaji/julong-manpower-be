package request

import (
	"time"

	"github.com/IlhamSetiaji/julong-manpower-be/internal/entity"
	"github.com/google/uuid"
)

type FindAllPaginatedMPPPeriodRequest struct {
	Page     int    `json:"page"`
	PageSize int    `json:"page_size"`
	Search   string `json:"search"`
}

type FindByIdMPPPeriodRequest struct {
	ID uuid.UUID `json:"id"`
}

type CreateMPPPeriodRequest struct {
	Title     string                 `json:"title" validate:"required"`
	StartDate time.Time              `json:"start_date" validate:"required"`
	EndDate   time.Time              `json:"end_date" validate:"required"`
	Status    entity.MPPPeriodStatus `json:"status" validate:"required,MPPPeriodStatusValidation"`
}

type UpdateMPPPeriodRequest struct {
	ID        uuid.UUID              `json:"id" validate:"required"`
	Title     string                 `json:"title" validate:"required"`
	StartDate time.Time              `json:"start_date" validate:"required"`
	EndDate   time.Time              `json:"end_date" validate:"required"`
	Status    entity.MPPPeriodStatus `json:"status" validate:"required,MPPPeriodStatusValidation"`
}

type DeleteMPPPeriodRequest struct {
	ID uuid.UUID `json:"id"`
}
