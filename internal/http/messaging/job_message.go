package messaging

import (
	"errors"
	"fmt"
	"log"

	"github.com/IlhamSetiaji/julong-manpower-be/internal/http/request"
	"github.com/IlhamSetiaji/julong-manpower-be/internal/http/response"
	"github.com/IlhamSetiaji/julong-manpower-be/utils"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type IJobMessage interface {
	SendFindJobDataByIdMessage(request request.SendFindJobByIDMessageRequest) (*response.JobResponse, error)
	SendGetAllJobDataMessage() (*[]response.JobResponse, error)
	SendFindAllJobsIDsMessage(ids []string) (*[]response.JobResponse, error)
}

type JobMessage struct {
	Log *logrus.Logger
}

func NewJobMessage(log *logrus.Logger) IJobMessage {
	return &JobMessage{
		Log: log,
	}
}

func (m *JobMessage) SendFindJobDataByIdMessage(req request.SendFindJobByIDMessageRequest) (*response.JobResponse, error) {
	payload := map[string]interface{}{
		"job_id": req.ID,
	}

	docMsg := &request.RabbitMQRequest{
		ID:          uuid.New().String(),
		MessageType: "find_job_data_by_id",
		MessageData: payload,
		ReplyTo:     "julong_manpower",
	}

	log.Printf("INFO: document message: %v", docMsg)

	// create channel and add to rchans with uid
	rchan := make(chan response.RabbitMQResponse)
	utils.Rchans[docMsg.ID] = rchan

	// publish rabbit message
	msg := utils.RabbitMsg{
		QueueName: "julong_sso",
		Message:   *docMsg,
	}
	utils.Pchan <- msg

	// wait for reply
	resp, err := waitReply(docMsg.ID, rchan)
	if err != nil {
		return nil, err
	}

	log.Printf("INFO: response: %v", resp.MessageData["job"])

	if errMsg, ok := resp.MessageData["error"].(string); ok && errMsg != "" {
		return nil, errors.New("[SendFindJobDataByIdMessage] " + errMsg)
	}

	jobData := resp.MessageData["job"].(map[string]interface{})
	return convertInterfaceToJobResponse(jobData), nil
}

func (m *JobMessage) SendGetAllJobDataMessage() (*[]response.JobResponse, error) {
	docMsg := &request.RabbitMQRequest{
		ID:          uuid.New().String(),
		MessageType: "get_all_job_data",
		MessageData: nil,
		ReplyTo:     "julong_manpower",
	}

	log.Printf("INFO: document message: %v", docMsg)

	// create channel and add to rchans with uid
	rchan := make(chan response.RabbitMQResponse)
	utils.Rchans[docMsg.ID] = rchan

	// publish rabbit message
	msg := utils.RabbitMsg{
		QueueName: "julong_sso",
		Message:   *docMsg,
	}

	utils.Pchan <- msg

	// wait for reply
	resp, err := waitReply(docMsg.ID, rchan)
	if err != nil {
		return nil, err
	}

	log.Printf("INFO: response: %v", resp.MessageData["jobs"])

	if errMsg, ok := resp.MessageData["error"].(string); ok && errMsg != "" {
		return nil, errors.New("[SendGetAllJobDataMessage] " + errMsg)
	}

	jobsData := resp.MessageData["jobs"].([]interface{})
	var jobsResponse []response.JobResponse
	for _, job := range jobsData {
		if jobMap, ok := job.(map[string]interface{}); ok {
			jobsResponse = append(jobsResponse, *convertInterfaceToJobResponse(jobMap))
		}
	}

	return &jobsResponse, nil
}

func convertInterfaceToJobResponse(job map[string]interface{}) *response.JobResponse {
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

	// handle job level
	var jobLevelResponse *response.JobLevelResponse
	if jobLevel, ok := job["job_level"].(map[string]interface{}); ok {
		jobLevelID, _ := jobLevel["id"].(string)
		jobLevelName, _ := jobLevel["name"].(string)
		jobLevelLevel, _ := jobLevel["level"].(string)
		jobLevelResponse = &response.JobLevelResponse{
			ID:    uuid.MustParse(jobLevelID),
			Name:  jobLevelName,
			Level: jobLevelLevel,
		}
	}

	// Handle Children
	var childrenResponse []response.JobResponse
	if children, ok := job["children"].([]interface{}); ok {
		for _, child := range children {
			if childMap, ok := child.(map[string]interface{}); ok {
				childrenResponse = append(childrenResponse, *convertInterfaceToJobResponse(childMap))
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
		JobLevel: *jobLevelResponse,
	}
}

func (m *JobMessage) SendFindAllJobsIDsMessage(ids []string) (*[]response.JobResponse, error) {
	payload := map[string]interface{}{
		"included_ids": ids,
	}

	docMsg := &request.RabbitMQRequest{
		ID:          uuid.New().String(),
		MessageType: "find_all_jobs_by_ids",
		MessageData: payload,
		ReplyTo:     "julong_manpower",
	}

	log.Printf("INFO: document message: %v", docMsg)

	// create channel and add to rchans with uid
	rchan := make(chan response.RabbitMQResponse)
	utils.Rchans[docMsg.ID] = rchan

	// publish rabbit message
	msg := utils.RabbitMsg{
		QueueName: "julong_sso",
		Message:   *docMsg,
	}
	utils.Pchan <- msg

	// wait for reply
	resp, err := waitReply(docMsg.ID, rchan)
	if err != nil {
		return nil, err
	}

	log.Printf("INFO: response: %v", resp.MessageData["jobs"])

	if errMsg, ok := resp.MessageData["error"].(string); ok && errMsg != "" {
		return nil, errors.New("[SendFindAllJobsByHeaderIdMessage] " + errMsg)
	}

	jobsData := resp.MessageData["jobs"].([]interface{})
	var jobsResponse []response.JobResponse
	for _, job := range jobsData {
		if jobMap, ok := job.(map[string]interface{}); ok {
			jobsResponse = append(jobsResponse, *convertInterfaceToJobResponse(jobMap))
		}
	}

	return &jobsResponse, nil
}

func JobMessageFactory(log *logrus.Logger) IJobMessage {
	return NewJobMessage(log)
}
