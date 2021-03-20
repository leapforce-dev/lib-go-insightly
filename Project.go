package insightly

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
	i_types "github.com/leapforce-libraries/go_insightly/types"
)

type Project struct {
	ProjectID           int64                   `json:"PROJECT_ID"`
	ProjectName         string                  `json:"PROJECT_NAME"`
	Status              string                  `json:"STATUS"`
	ProjectDetails      *string                 `json:"PROJECT_DETAILS"`
	StartedDate         *i_types.DateTimeString `json:"STARTED_DATE"`
	CompletedDate       *i_types.DateTimeString `json:"COMPLETED_DATE"`
	OpportunityID       *int64                  `json:"OPPORTUNITY_ID"`
	CategoryID          int64                   `json:"CATEGORY_ID"`
	PipelineID          int64                   `json:"PIPELINE_ID"`
	StageID             int64                   `json:"STAGE_ID"`
	ImageURL            *string                 `json:"IMAGE_URL"`
	OwnerUserID         int64                   `json:"OWNER_USER_ID"`
	DateCreatedUTC      i_types.DateTimeString  `json:"DATE_CREATED_UTC"`
	DateUpdatedUTC      i_types.DateTimeString  `json:"DATE_UPDATED_UTC"`
	LastActivityDateUTC *i_types.DateTimeString `json:"LAST_ACTIVITY_DATE_UTC"`
	NextActivityDateUTC *i_types.DateTimeString `json:"NEXT_ACTIVITY_DATE_UTC"`
	CreatedUserID       int64                   `json:"CREATED_USER_ID"`
	ResponsibleUserID   *int64                  `json:"RESPONSIBLE_USER_ID"`
	CustomFields        *CustomFields           `json:"CUSTOMFIELDS"`
	Tags                *[]Tag                  `json:"TAGS"`
}

// GetProject returns a specific project
//
func (service *Service) GetProject(projectID int) (*Project, *errortools.Error) {
	project := Project{}

	requestConfig := go_http.RequestConfig{
		URL:           service.url(fmt.Sprintf("Project/%v", projectID)),
		ResponseModel: &project,
	}
	_, _, e := service.get(&requestConfig)
	if e != nil {
		return nil, e
	}

	return &project, nil
}

type GetProjectsConfig struct {
	UpdatedAfter *time.Time
	FieldFilter  *FieldFilter
}

// GetProjects returns all projects
//
func (service *Service) GetProjects(config *GetProjectsConfig) (*[]Project, *errortools.Error) {
	searchString := "?"
	searchFilter := []string{}

	if config != nil {
		if config.UpdatedAfter != nil {
			from := config.UpdatedAfter.Format(DateTimeFormat)
			searchFilter = append(searchFilter, fmt.Sprintf("updated_after_utc=%s&", from))
		}

		if config.FieldFilter != nil {
			searchFilter = append(searchFilter, fmt.Sprintf("field_name=%s&field_value=%s&", config.FieldFilter.FieldName, config.FieldFilter.FieldValue))
		}
	}

	if len(searchFilter) > 0 {
		searchString = "/Search?" + strings.Join(searchFilter, "&")
	}

	endpointStr := "Project%sskip=%s&top=%s"
	skip := 0
	top := 100
	rowCount := top

	projects := []Project{}

	for rowCount >= top {
		_projects := []Project{}

		requestConfig := go_http.RequestConfig{
			URL:           service.url(fmt.Sprintf(endpointStr, searchString, strconv.Itoa(skip), strconv.Itoa(top))),
			ResponseModel: &_projects,
		}
		_, _, e := service.get(&requestConfig)
		if e != nil {
			return nil, e
		}

		projects = append(projects, _projects...)

		rowCount = len(_projects)
		//rowCount = 0
		skip += top
	}

	if len(projects) == 0 {
		projects = nil
	}

	return &projects, nil
}
