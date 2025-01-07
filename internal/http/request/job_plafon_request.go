package request

type CheckJobExistMessageRequest struct {
	ID string `json:"id"`
}

type SendFindJobByIDMessageRequest struct {
	ID string `json:"id"`
}

type SendFindJobLevelByIDMessageRequest struct {
	ID string `json:"id"`
}

type FindAllPaginatedJobPlafonRequest struct {
	Page           int    `json:"page"`
	PageSize       int    `json:"page_size"`
	Search         string `json:"search"`
	RequestorID    string `json:"requestor_id"`
	OrganizationID string `json:"organization_id"`
}

type FindByIdJobPlafonRequest struct {
	ID string `json:"id"`
}

type CreateJobPlafonRequest struct {
	JobID  string `json:"job_id"`
	Plafon int    `json:"plafon"`
}

type UpdateJobPlafonRequest struct {
	ID     string `json:"id"`
	JobID  string `json:"job_id"`
	Plafon int    `json:"plafon"`
}

type DeleteJobPlafonRequest struct {
	ID string `json:"id"`
}

type FindByJobIdJobPlafonRequest struct {
	JobID string `json:"job_id"`
}

type CheckJobByJobLevelRequest struct {
	JobID      string `json:"job_id"`
	JobLevelID string `json:"job_level_id"`
}
