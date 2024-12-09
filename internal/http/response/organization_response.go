package response

type SendFindOrganizationByIDMessageResponse struct {
	OrganizationID string `json:"organization_id"`
	Name           string `json:"name"`
}

type SendFindOrganizationLocationByIDMessageResponse struct {
	OrganizationLocationID string `json:"organization_location_id"`
	Name                   string `json:"name"`
}

type SendFindOrganizationStructureByIDMessageResponse struct {
	OrganizationStructureID string `json:"organization_structure_id"`
	Name                    string `json:"name"`
}
