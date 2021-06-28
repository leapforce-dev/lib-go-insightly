package insightly

import (
	"fmt"
	"net/url"
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
	Links               *[]Link                 `json:"LINKS"`
}

// GetProject returns a specific project
//
func (service *Service) GetProject(projectID int64) (*Project, *errortools.Error) {
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
	Skip         *uint64
	Top          *uint64
	Brief        *bool
	CountTotal   *bool
	UpdatedAfter *time.Time
	FieldFilter  *FieldFilter
}

// GetProjects returns all projects
//
func (service *Service) GetProjects(config *GetProjectsConfig) (*[]Project, *errortools.Error) {
	params := url.Values{}

	endpoint := "Project"
	projects := []Project{}
	rowCount := uint64(0)
	top := defaultTop
	isSearch := false

	if config != nil {
		if config.Top != nil {
			top = *config.Top
		}
		if config.Skip != nil {
			service.nextSkips[endpoint] = *config.Skip
		}
		if config.Brief != nil {
			params.Set("brief", fmt.Sprintf("%v", *config.Brief))
		}
		if config.CountTotal != nil {
			params.Set("count_total", fmt.Sprintf("%v", *config.CountTotal))
		}
		if config.UpdatedAfter != nil {
			isSearch = true
			params.Set("updated_after_utc", fmt.Sprintf("%v", config.UpdatedAfter.Format(dateTimeFormat)))
		}
		if config.FieldFilter != nil {
			isSearch = true
			params.Set("field_name", config.FieldFilter.FieldName)
			params.Set("field_value", config.FieldFilter.FieldValue)
		}
	}

	if isSearch {
		endpoint += "/Search"
	}

	params.Set("top", fmt.Sprintf("%v", top))

	for {
		params.Set("skip", fmt.Sprintf("%v", service.nextSkips[endpoint]))
		projectsBatch := []Project{}

		requestConfig := go_http.RequestConfig{
			URL:           service.url(fmt.Sprintf("%s?%s", endpoint, params.Encode())),
			ResponseModel: &projectsBatch,
		}
		_, _, e := service.get(&requestConfig)
		if e != nil {
			return nil, e
		}

		projects = append(projects, projectsBatch...)

		if len(projectsBatch) < int(top) {
			delete(service.nextSkips, endpoint)
			break
		}

		service.nextSkips[endpoint] += top
		rowCount += top

		if rowCount >= service.maxRowCount {
			return &projects, nil
		}
	}

	return &projects, nil
}
