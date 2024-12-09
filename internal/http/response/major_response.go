package response

import (
	"github.com/IlhamSetiaji/julong-manpower-be/internal/entity"
	"github.com/google/uuid"
)

type FindAllMajorResponse struct {
	ID             string                    `json:"id"`
	Major          string                    `json:"major"`
	EducationLevel entity.EducationLevelEnum `json:"education_level"`
}

type MajorResponse struct {
	ID             uuid.UUID                 `json:"id"`
	Major          string                    `json:"major"`
	EducationLevel entity.EducationLevelEnum `json:"education_level"`
}

type RequestMajorResponse struct {
	ID                string        `json:"id"`
	Major             MajorResponse `json:"major"`
	MPRequestHeaderID string        `json:"mp_request_header_id"`
}
