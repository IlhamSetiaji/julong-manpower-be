package dto

import (
	"time"

	"github.com/IlhamSetiaji/julong-manpower-be/internal/http/response"
	"github.com/google/uuid"
)

func ConvertInterfaceToEmployeeResponse(data map[string]interface{}) *response.EmployeeResponse {
	// Extract values from the map
	id := data["id"].(string)
	organizationID := data["organization_id"].(string)
	name := data["name"].(string)
	endDate, _ := time.Parse("2006-01-02", data["end_date"].(string))
	retirementDate, _ := time.Parse("2006-01-02", data["retirement_date"].(string))
	email := data["email"].(string)
	mobilePhone := data["mobile_phone"].(string)

	return &response.EmployeeResponse{
		ID:             uuid.MustParse(id),
		OrganizationID: uuid.MustParse(organizationID),
		Name:           name,
		EndDate:        endDate,
		RetirementDate: retirementDate,
		Email:          email,
		MobilePhone:    mobilePhone,
	}
}
