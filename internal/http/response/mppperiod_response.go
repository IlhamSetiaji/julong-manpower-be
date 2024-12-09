package response

import (
	"time"

	"github.com/IlhamSetiaji/julong-manpower-be/internal/entity"
	"github.com/google/uuid"
)

type MPPeriodResponse struct {
	ID        uuid.UUID              `json:"id"`
	Title     string                 `json:"title"`
	StartDate string                 `json:"start_date"`
	EndDate   string                 `json:"end_date"`
	Status    entity.MPPPeriodStatus `json:"status"`
	CreatedAt time.Time              `json:"created_at"`
	UpdatedAt time.Time              `json:"updated_at"`
}

type FindAllPaginatedMPPPeriodResponse struct {
	MPPPeriods *[]entity.MPPPeriod `json:"mppperiods"`
	Total      int64               `json:"total"`
}

type FindByIdMPPPeriodResponse struct {
	MPPPeriod *entity.MPPPeriod `json:"mppperiod"`
}

type CreateMPPPeriodResponse struct {
	ID        uuid.UUID              `json:"id"`
	Title     string                 `json:"title"`
	StartDate string                 `json:"start_date"`
	EndDate   string                 `json:"end_date"`
	Status    entity.MPPPeriodStatus `json:"status"`
	CreatedAt time.Time              `json:"created_at"`
	UpdatedAt time.Time              `json:"updated_at"`
	DeletedAt *time.Time             `json:"deleted_at"`
}

type UpdateMPPPeriodResponse struct {
	ID        uuid.UUID              `json:"id"`
	Title     string                 `json:"title"`
	StartDate string                 `json:"start_date"`
	EndDate   string                 `json:"end_date"`
	Status    entity.MPPPeriodStatus `json:"status"`
	CreatedAt time.Time              `json:"created_at"`
	UpdatedAt time.Time              `json:"updated_at"`
	DeletedAt *time.Time             `json:"deleted_at"`
}

type FindByCurrentDateAndStatusMPPPeriodResponse struct {
	MPPPeriod *entity.MPPPeriod `json:"mppperiod"`
}
