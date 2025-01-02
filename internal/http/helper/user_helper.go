package helper

import (
	"errors"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type IUserHelper interface {
	CheckOrganizationLocation(user map[string]interface{}) (uuid.UUID, error)
	GetEmployeeId(user map[string]interface{}) (uuid.UUID, error)
	GetOrganizationStructureID(user map[string]interface{}) (uuid.UUID, error)
	GetOrganizationID(user map[string]interface{}) (uuid.UUID, error)
}

type UserHelper struct {
	Log *logrus.Logger
}

func NewUserHelper(log *logrus.Logger) IUserHelper {
	return &UserHelper{Log: log}
}

func UserHelperFactory(log *logrus.Logger) IUserHelper {
	return NewUserHelper(log)
}

func (h *UserHelper) CheckOrganizationLocation(user map[string]interface{}) (uuid.UUID, error) {
	// Check if the "user" key exists and is a map
	userData, ok := user["user"].(map[string]interface{})
	if !ok {
		h.Log.Errorf("User information is missing or invalid")
		return uuid.Nil, errors.New("User information is missing or invalid")
	}

	// Check if the "employee" key exists and is a map
	employee, ok := userData["employee"].(map[string]interface{})
	if !ok {
		h.Log.Errorf("Employee information is missing or invalid")
		return uuid.Nil, errors.New("Employee information is missing or invalid")
	}

	// Check if the "employee_job" key exists and is a map
	employeeJob, ok := employee["employee_job"].(map[string]interface{})
	if !ok {
		h.Log.Errorf("Employee job information is missing or invalid")
		return uuid.Nil, errors.New("Employee job information is missing or invalid")
	}

	// Check if the "OrganizationLocationID" key exists and is a string
	organizationLocationIDStr, ok := employeeJob["organization_location_id"].(string)
	if !ok {
		h.Log.Errorf("Organization location ID is missing or invalid")
		return uuid.Nil, errors.New("Organization location ID is missing or invalid")
	}

	// Parse the organization location ID to uuid.UUID
	organizationLocationID, err := uuid.Parse(organizationLocationIDStr)
	if err != nil {
		h.Log.Errorf("Invalid organization location ID format: %v", err)
		return uuid.Nil, errors.New("Invalid organization location ID format")
	}

	h.Log.Infof("Organization Location ID: %s", organizationLocationID)
	return organizationLocationID, nil
}

func (h *UserHelper) GetEmployeeId(user map[string]interface{}) (uuid.UUID, error) {
	// Check if the "user" key exists and is a map
	userData, ok := user["user"].(map[string]interface{})
	if !ok {
		h.Log.Errorf("User information is missing or invalid")
		return uuid.Nil, errors.New("User information is missing or invalid")
	}

	// Check if the "employee" key exists and is a map
	employee, ok := userData["employee"].(map[string]interface{})
	if !ok {
		h.Log.Errorf("Employee information is missing or invalid")
		return uuid.Nil, errors.New("Employee information is missing or invalid")
	}

	// Check if the "ID" key exists and is a string
	employeeIDStr, ok := employee["id"].(string)
	if !ok {
		h.Log.Errorf("Employee ID is missing or invalid")
		return uuid.Nil, errors.New("Employee ID is missing or invalid")
	}

	// Parse the employee ID to uuid.UUID
	employeeID, err := uuid.Parse(employeeIDStr)
	if err != nil {
		h.Log.Errorf("Invalid employee ID format: %v", err)
		return uuid.Nil, errors.New("Invalid employee ID format")
	}

	h.Log.Infof("Employee ID: %s", employeeID)
	return employeeID, nil
}

func (h *UserHelper) GetOrganizationStructureID(user map[string]interface{}) (uuid.UUID, error) {
	// Check if the "user" key exists and is a map
	userData, ok := user["user"].(map[string]interface{})
	if !ok {
		h.Log.Errorf("User information is missing or invalid")
		return uuid.Nil, errors.New("User information is missing or invalid")
	}

	// Check if the "employee" key exists and is a map
	employee, ok := userData["employee"].(map[string]interface{})
	if !ok {
		h.Log.Errorf("Employee information is missing or invalid")
		return uuid.Nil, errors.New("Employee information is missing or invalid")
	}

	// Check if the "employee_job" key exists and is a map
	employeeJob, ok := employee["employee_job"].(map[string]interface{})
	if !ok {
		h.Log.Errorf("Employee job information is missing or invalid")
		return uuid.Nil, errors.New("Employee job information is missing or invalid")
	}

	// check if the "organization_structure" key exists and is a map
	organizationStructure, ok := employeeJob["organization_structure"].(map[string]interface{})
	if !ok {
		h.Log.Errorf("Organization structure information is missing or invalid")
		return uuid.Nil, errors.New("Organization structure information is missing or invalid")
	}

	// Check if the "ID" key exists and is a string
	organizationStructureIDStr, ok := organizationStructure["id"].(string)
	if !ok {
		h.Log.Errorf("Organization structure ID is missing or invalid")
		return uuid.Nil, errors.New("Organization structure ID is missing or invalid")
	}

	// Parse the organization structure ID to uuid.UUID
	organizationStructureID, err := uuid.Parse(organizationStructureIDStr)
	if err != nil {
		h.Log.Errorf("Invalid organization structure ID format: %v", err)
		return uuid.Nil, errors.New("Invalid organization structure ID format")
	}

	h.Log.Infof("Organization Structure ID: %s", organizationStructureID)
	return organizationStructureID, nil
}

func (h *UserHelper) GetOrganizationID(user map[string]interface{}) (uuid.UUID, error) {
	// Check if the "user" key exists and is a map
	userData, ok := user["user"].(map[string]interface{})
	if !ok {
		h.Log.Errorf("User information is missing or invalid")
		return uuid.Nil, errors.New("User information is missing or invalid")
	}

	// Check if the "employee" key exists and is a map
	employee, ok := userData["employee"].(map[string]interface{})
	if !ok {
		h.Log.Errorf("Employee information is missing or invalid")
		return uuid.Nil, errors.New("Employee information is missing or invalid")
	}

	// Check if the "employee_job" key exists and is a map
	employeeJob, ok := employee["employee_job"].(map[string]interface{})
	if !ok {
		h.Log.Errorf("Employee job information is missing or invalid")
		return uuid.Nil, errors.New("Employee job information is missing or invalid")
	}

	// Check if the "organization" key exists and is a map
	organization, ok := employeeJob["organization"].(map[string]interface{})
	if !ok {
		h.Log.Errorf("Organization information is missing or invalid")
		return uuid.Nil, errors.New("Organization information is missing or invalid")
	}

	// Check if the "ID" key exists and is a string
	organizationIDStr, ok := organization["id"].(string)
	if !ok {
		h.Log.Errorf("Organization ID is missing or invalid")
		return uuid.Nil, errors.New("Organization ID is missing or invalid")
	}

	// Parse the organization ID to uuid.UUID
	organizationID, err := uuid.Parse(organizationIDStr)
	if err != nil {
		h.Log.Errorf("Invalid organization ID format: %v", err)
		return uuid.Nil, errors.New("Invalid organization ID format")
	}

	h.Log.Infof("Organization ID: %s", organizationID)
	return organizationID, nil
}
