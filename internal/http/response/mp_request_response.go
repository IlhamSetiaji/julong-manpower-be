package response

import (
	"time"

	"github.com/IlhamSetiaji/julong-manpower-be/internal/entity"
	"github.com/google/uuid"
)

type MPRequestPaginatedResponse struct {
	MPRequestHeader []MPRequestHeaderResponse `json:"mp_request_header"`
	Total           int64                     `json:"total"`
}

type MPRequestHeaderResponse struct {
	ID                         uuid.UUID                  `json:"id"`
	OrganizationID             uuid.UUID                  `json:"organization_id"`
	OrganizationLocationID     uuid.UUID                  `json:"organization_location_id"`
	ForOrganizationID          uuid.UUID                  `json:"for_organization_id"`
	ForOrganizationLocationID  uuid.UUID                  `json:"for_organization_location_id"`
	ForOrganizationStructureID uuid.UUID                  `json:"for_organization_structure_id"`
	JobID                      uuid.UUID                  `json:"job_id"`
	RequestCategoryID          uuid.UUID                  `json:"request_category_id"`
	GradeID                    *uuid.UUID                 `json:"grade_id"`
	ExpectedDate               time.Time                  `json:"expected_date"`
	Experiences                string                     `json:"experiences"`
	DocumentNumber             string                     `json:"document_number"`
	DocumentDate               time.Time                  `json:"document_date"`
	MaleNeeds                  int                        `json:"male_needs"`
	FemaleNeeds                int                        `json:"female_needs"`
	AnyGender                  int                        `json:"any_gender"`
	TotalNeeds                 int                        `json:"total_needs"`
	MinimumAge                 int                        `json:"minimum_age"`
	MaximumAge                 int                        `json:"maximum_age"`
	MinimumExperience          int                        `json:"minimum_experience"`
	MaritalStatus              entity.MaritalStatusEnum   `json:"marital_status"`
	MinimumEducation           entity.EducationLevelEnum  `json:"minimum_education"`
	RequiredQualification      string                     `json:"required_qualification"`
	Certificate                string                     `json:"certificate"`
	ComputerSkill              string                     `json:"computer_skill"`
	LanguageSkill              string                     `json:"language_skill"`
	OtherSkill                 string                     `json:"other_skill"`
	Jobdesc                    string                     `json:"jobdesc"`
	SalaryMin                  string                     `json:"salary_min"`
	SalaryMax                  string                     `json:"salary_max"`
	RequestorID                *uuid.UUID                 `json:"requestor_id"`
	DepartmentHead             *uuid.UUID                 `json:"department_head"`
	VpGmDirector               *uuid.UUID                 `json:"vp_gm_director"`
	CEO                        *uuid.UUID                 `json:"ceo"`
	HrdHoUnit                  *uuid.UUID                 `json:"hrd_ho_unit"`
	MPPlanningHeaderID         *uuid.UUID                 `json:"mp_planning_header_id"`
	Status                     entity.MPRequestStatus     `json:"status"`
	MPRequestType              entity.MPRequestTypeEnum   `json:"mp_request_type"`
	RecruitmentType            entity.RecruitmentTypeEnum `json:"recruitment_type"`
	MPPPeriodID                *uuid.UUID                 `json:"mpp_period_id"`
	EmpOrganizationID          *uuid.UUID                 `json:"emp_organization_id"`
	JobLevelID                 *uuid.UUID                 `json:"job_level_id"`
	IsReplacement              bool                       `json:"is_replacement"`
	Revised                    int                        `json:"revised"`
	CreatedAt                  time.Time                  `json:"created_at"`
	UpdatedAt                  time.Time                  `json:"updated_at"`

	RequestCategory map[string]interface{}   `json:"request_category" gorm:"foreignKey:RequestCategoryID"`
	RequestMajors   []map[string]interface{} `json:"request_majors" gorm:"foreignKey:MPRequestHeaderID"`
	// MPPlanningHeader map[string]interface{}   `json:"mp_planning_header" gorm:"foreignKey:MPPlanningHeaderID"`
	MPPlanningHeader *MPPlanningHeaderResponse `json:"mp_planning_header" gorm:"foreignKey:MPPlanningHeaderID"`
	MPPPeriod        MPPeriodResponse          `json:"mpp_period" gorm:"foreignKey:MPPPeriodID"`

	GradeName                 string                 `json:"grade_name" gorm:"-"`
	OrganizationName          string                 `json:"organization_name" gorm:"-"`
	OrganizationCategory      string                 `json:"organization_category" gorm:"-"`
	OrganizationLocationName  string                 `json:"organization_location_name" gorm:"-"`
	ForOrganizationName       string                 `json:"for_organization_name" gorm:"-"`
	ForOrganizationLocation   string                 `json:"for_organization_location" gorm:"-"`
	ForOrganizationStructure  string                 `json:"for_organization_structure" gorm:"-"`
	JobName                   string                 `json:"job_name" gorm:"-"`
	RequestorName             string                 `json:"requestor_name" gorm:"-"`
	DepartmentHeadName        string                 `json:"department_head_name" gorm:"-"`
	HrdHoUnitName             string                 `json:"hrd_ho_unit_name" gorm:"-"`
	VpGmDirectorName          string                 `json:"vp_gm_director_name" gorm:"-"`
	CeoName                   string                 `json:"ceo_name" gorm:"-"`
	EmpOrganizationName       string                 `json:"emp_organization_name" gorm:"-"`
	JobLevelName              string                 `json:"job_level_name" gorm:"-"`
	JobLevel                  int                    `json:"job_level" gorm:"-"`
	ApprovedByDepartmentHead  bool                   `json:"approved_by_department_head" gorm:"-"`
	ApprovedByVpGmDirector    bool                   `json:"approved_by_vp_gm_director" gorm:"-"`
	ApprovedByCEO             bool                   `json:"approved_by_ceo" gorm:"-"`
	ApprovedByHrdHoUnit       bool                   `json:"approved_by_hrd_ho_unit" gorm:"-"`
	RequestorEmployeeJob      map[string]interface{} `json:"requestor_employee_job" gorm:"-"`
	DepartmentHeadEmployeeJob map[string]interface{} `json:"department_head_employee_job" gorm:"-"`
	VpGmDirectorEmployeeJob   map[string]interface{} `json:"vp_gm_director_employee_job" gorm:"-"`
	CeoEmployeeJob            map[string]interface{} `json:"ceo_employee_job" gorm:"-"`
}

type CheckPortalDataMPRequestResponse struct {
	GradeName                    string                 `json:"grade_name"`
	OrganizationName             string                 `json:"organization_name"`
	OrganizationCategory         string                 `json:"organization_category"`
	OrganizationLocationName     string                 `json:"organization_location_name"`
	ForOrganizationName          string                 `json:"for_organization_name"`
	ForOrganizationLocationName  string                 `json:"for_organization_location_name"`
	ForOrganizationStructureName string                 `json:"for_organization_structure_name"`
	JobName                      string                 `json:"job_name"`
	RequestorName                string                 `json:"requestor_name"`
	DepartmentHeadName           string                 `json:"department_head_name"`
	VpGmDirectorName             string                 `json:"vp_gm_director_name"`
	CeoName                      string                 `json:"ceo_name"`
	HrdHoUnitName                string                 `json:"hrd_ho_unit_name"`
	EmpOrganizationName          string                 `json:"emp_organization_name"`
	JobLevelName                 string                 `json:"job_level_name"`
	JobLevel                     int                    `json:"job_level"`
	RequestorEmployeeJob         map[string]interface{} `json:"requestor_employee_job"`
	DepartmentHeadEmployeeJob    map[string]interface{} `json:"department_head_employee_job"`
	VpGmDirectorEmployeeJob      map[string]interface{} `json:"vp_gm_director_employee_job"`
	CeoEmployeeJob               map[string]interface{} `json:"ceo_employee_job"`
}

type MPRequestApprovalHistoryResponse struct {
	ID                uuid.UUID                             `json:"id"`
	MPRequestHeaderID uuid.UUID                             `json:"mp_request_header_id"`
	ApproverID        uuid.UUID                             `json:"approver_id"`
	ApproverName      string                                `json:"approver_name"`
	Notes             string                                `json:"notes"`
	Level             string                                `json:"level"`
	Status            entity.MPRequestApprovalHistoryStatus `json:"status"`
	CreatedAt         time.Time                             `json:"created_at"`
	UpdatedAt         time.Time                             `json:"updated_at"`
	Attachments       []*ManpowerAttachmentResponse         `json:"attachments"`
}
