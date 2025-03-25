package response

type GradeResponse struct {
	ID           string `json:"id"`
	JobLevelID   string `json:"job_level_id"`
	Name         string `json:"name"`
	JobLevelName string `json:"job_level_name"`
}
