package helper

import (
	"errors"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type IUserHelper interface {
	CheckOrganizationLocation(user map[string]interface{}) (uuid.UUID, error)
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
