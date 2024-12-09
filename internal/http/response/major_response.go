package response

import "github.com/IlhamSetiaji/julong-manpower-be/internal/entity"

type FindAllMajorResponse struct {
	ID             string                    `json:"id"`
	Major          string                    `json:"major"`
	EducationLevel entity.EducationLevelEnum `json:"education_level"`
}
