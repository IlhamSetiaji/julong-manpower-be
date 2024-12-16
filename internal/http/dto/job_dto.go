package dto

import (
	"fmt"

	"github.com/IlhamSetiaji/julong-manpower-be/internal/http/response"
	"github.com/google/uuid"
)

func ConvertInterfaceToJobResponse(job map[string]interface{}) *response.JobResponse {
	// Extract values from the map
	id, _ := job["id"].(string)
	name, _ := job["name"].(string)
	organizationStructureID, _ := job["organization_structure_id"].(string)
	fmt.Println("ini cek", organizationStructureID)
	organizationStructureName, _ := job["organization_structure_name"].(string)
	organizationID, _ := job["organization_id"].(string)
	organizationName, _ := job["organization_name"].(string)
	level, _ := job["level"].(int)
	parentID, _ := job["parent_id"].(string)
	path, _ := job["path"].(string)
	existing, _ := job["existing"].(int)

	// Handle Parent
	var parentResponse *response.ParentJobResponse
	if parent, ok := job["parent"].(map[string]interface{}); ok {
		parentIDStr, _ := parent["id"].(string)
		parentName, _ := parent["name"].(string)
		parentID, _ := uuid.Parse(parentIDStr)
		parentResponse = &response.ParentJobResponse{ID: parentID, Name: parentName}
	}

	// Handle Children
	var childrenResponse []response.JobResponse
	if children, ok := job["children"].([]interface{}); ok {
		for _, child := range children {
			if childMap, ok := child.(map[string]interface{}); ok {
				childrenResponse = append(childrenResponse, *ConvertInterfaceToJobResponse(childMap))
			}
		}
	}

	return &response.JobResponse{
		ID:                        uuid.MustParse(id),
		Name:                      name,
		OrganizationStructureID:   uuid.MustParse(organizationStructureID),
		OrganizationStructureName: organizationStructureName,
		OrganizationID:            uuid.MustParse(organizationID),
		OrganizationName:          organizationName,
		Level:                     level,
		ParentID: func(id string) *uuid.UUID {
			parsedID, err := uuid.Parse(id)
			if err != nil {
				return nil
			}
			return &parsedID
		}(parentID),
		Path:     path,
		Existing: existing,
		Parent:   parentResponse,
		Children: childrenResponse,
	}
}
