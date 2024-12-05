package response

import (
	"github.com/IlhamSetiaji/julong-manpower-be/internal/entity"
	"github.com/google/uuid"
)

type CheckJobExistMessageResponse struct {
	JobID uuid.UUID `json:"job_id"`
	Exist bool      `json:"exist"`
}

type FindAllPaginatedJobPlafonResponse struct {
	JobPlafons *[]entity.JobPlafon `json:"job_plafons"`
	Total      int64               `json:"total"`
}

type FindByIdJobPlafonResponse struct {
	JobPlafon *entity.JobPlafon `json:"job_plafon"`
}

type CreateJobPlafonResponse struct {
	JobPlafon *entity.JobPlafon `json:"job_plafon"`
}

type UpdateJobPlafonResponse struct {
	JobPlafon *entity.JobPlafon `json:"job_plafon"`
}

type FindByJobIdJobPlafonResponse struct {
	JobPlafon *entity.JobPlafon `json:"job_plafon"`
}
