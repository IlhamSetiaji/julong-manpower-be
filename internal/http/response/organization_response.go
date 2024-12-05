package response

type SendFindOrganizationByIDMessageResponse struct {
	OrganizationID string `json:"organization_id"`
	Name           string `json:"name"`
}

type SendFindOrganizationLocationByIDMessageResponse struct {
	OrganizationLocationID string `json:"organization_location_id"`
	Name                   string `json:"name"`
}
