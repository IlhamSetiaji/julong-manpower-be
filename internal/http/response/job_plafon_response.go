package response

import "github.com/google/uuid"

type CheckJobExistMessageResponse struct {
	JobID uuid.UUID `json:"job_id"`
	Exist bool      `json:"exist"`
}
