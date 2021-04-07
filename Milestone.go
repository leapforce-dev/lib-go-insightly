package insightly

import (
	"fmt"
	"net/url"
	"time"

	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
	i_types "github.com/leapforce-libraries/go_insightly/types"
)

// Milestone stores Milestone from Service
//
type Milestone struct {
	MilestoneID       int64                   `json:"MILESTONE_ID"`
	Title             string                  `json:"TITLE"`
	Completed         bool                    `json:"COMPLETED"`
	DueDate           i_types.DateTimeString  `json:"DUE_DATE"`
	OwnerUserID       int64                   `json:"OWNER_USER_ID"`
	DateCreatedUTC    i_types.DateTimeString  `json:"DATE_CREATED_UTC"`
	DateUpdatedUTC    i_types.DateTimeString  `json:"DATE_UPDATED_UTC"`
	CompletedDateUTC  *i_types.DateTimeString `json:"COMPLETED_DATE_UTC"`
	ProjectID         int64                   `json:"PROJECT_ID"`
	ResponsibleUserID int64                   `json:"RESPONSIBLE_USER"`
}

// GetMilestone returns a specific milestone
//
func (service *Service) GetMilestone(milestoneID int64) (*Milestone, *errortools.Error) {
	milestone := Milestone{}

	requestConfig := go_http.RequestConfig{
		URL:           service.url(fmt.Sprintf("Milestones/%v", milestoneID)),
		ResponseModel: &milestone,
	}
	_, _, e := service.get(&requestConfig)
	if e != nil {
		return nil, e
	}

	return &milestone, nil
}

type GetMilestonesConfig struct {
	Skip         *uint64
	Top          *uint64
	Brief        *bool
	CountTotal   *bool
	UpdatedAfter *time.Time
	FieldFilter  *FieldFilter
}

// GetMilestones returns all milestones
//
func (service *Service) GetMilestones(config *GetMilestonesConfig) (*[]Milestone, *errortools.Error) {
	params := url.Values{}

	endpoint := "Milestones"
	milestones := []Milestone{}
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

	for true {
		params.Set("skip", fmt.Sprintf("%v", service.nextSkips[endpoint]))
		milestonesBatch := []Milestone{}

		requestConfig := go_http.RequestConfig{
			URL:           service.url(fmt.Sprintf("%s?%s", endpoint, params.Encode())),
			ResponseModel: &milestonesBatch,
		}
		_, _, e := service.get(&requestConfig)
		if e != nil {
			return nil, e
		}

		milestones = append(milestones, milestonesBatch...)

		if len(milestonesBatch) < int(top) {
			delete(service.nextSkips, endpoint)
			break
		}

		service.nextSkips[endpoint] += top
		rowCount += top

		if rowCount >= service.maxRowCount {
			return &milestones, nil
		}
	}

	return &milestones, nil
}
