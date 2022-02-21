package insightly

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
	i_types "github.com/leapforce-libraries/go_insightly/types"
)

// Opportunity stores Opportunity from Service
//
type Opportunity struct {
	OpportunityID       int64                   `json:"OPPORTUNITY_ID"`
	OpportunityName     *string                 `json:"OPPORTUNITY_NAME,omitempty"`
	OpportunityDetails  *string                 `json:"OPPORTUNITY_DETAILS,omitempty"`
	OpportunityState    *string                 `json:"OPPORTUNITY_STATE,omitempty"`
	ResponsibleUserID   *int64                  `json:"RESPONSIBLE_USER_ID,omitempty"`
	CategoryID          *int64                  `json:"CATEGORY_ID,omitempty"`
	ImageUrl            *string                 `json:"IMAGE_Url,omitempty"`
	BidCurrency         *string                 `json:"BID_CURRENCY,omitempty"`
	BidAmount           *float64                `json:"BID_AMOUNT,omitempty"`
	BidType             *string                 `json:"BID_TYPE,omitempty"`
	BidDuration         *int64                  `json:"BID_DURATION,omitempty"`
	ActualCloseDate     *i_types.DateTimeString `json:"ACTUAL_CLOSE_DATE,omitempty"`
	DateCreatedUTC      *i_types.DateTimeString `json:"DATE_CREATED_UTC,omitempty"`
	DateUpdatedUTC      *i_types.DateTimeString `json:"DATE_UPDATED_UTC,omitempty"`
	OpportunityValue    *float64                `json:"OPPORTUNITY_VALUE,omitempty"`
	Probability         *int64                  `json:"PROBABILITY,omitempty"`
	ForecastCloseDate   *i_types.DateTimeString `json:"FORECAST_CLOSE_DATE,omitempty"`
	OwnerUserID         *int64                  `json:"OWNER_USER_ID,omitempty"`
	LastActivityDateUTC *i_types.DateTimeString `json:"LAST_ACTIVITY_DATE_UTC,omitempty"`
	NextActivityDateUTC *i_types.DateTimeString `json:"NEXT_ACTIVITY_DATE_UTC,omitempty"`
	PipelineID          *int64                  `json:"PIPELINE_ID,omitempty"`
	StageID             *int64                  `json:"STAGE_ID,omitempty"`
	CreatedUserID       *int64                  `json:"CREATED_USER_ID,omitempty"`
	OrganisationID      *int64                  `json:"ORGANISATION_ID,omitempty"`
	CustomFields        *CustomFields           `json:"CUSTOMFIELDS,omitempty"`
	Tags                *[]Tag                  `json:"TAGS,omitempty"`
	Links               *[]Link                 `json:"LINKS,omitempty"`
}

// GetOpportunity returns a specific opportunity
//
func (service *Service) GetOpportunity(opportunityID int64) (*Opportunity, *errortools.Error) {
	opportunity := Opportunity{}

	requestConfig := go_http.RequestConfig{
		Method:        http.MethodGet,
		Url:           service.url(fmt.Sprintf("Opportunities/%v", opportunityID)),
		ResponseModel: &opportunity,
	}
	_, _, e := service.httpRequest(&requestConfig)
	if e != nil {
		return nil, e
	}

	return &opportunity, nil
}

type GetOpportunitiesConfig struct {
	Skip         *uint64
	Top          *uint64
	Brief        *bool
	CountTotal   *bool
	UpdatedAfter *time.Time
	FieldFilter  *FieldFilter
}

// GetOpportunities returns all opportunities
//
func (service *Service) GetOpportunities(config *GetOpportunitiesConfig) (*[]Opportunity, *errortools.Error) {
	params := url.Values{}

	endpoint := "Opportunities"
	opportunities := []Opportunity{}
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
		opportunitiesBatch := []Opportunity{}

		requestConfig := go_http.RequestConfig{
			Method:        http.MethodGet,
			Url:           service.url(fmt.Sprintf("%s?%s", endpoint, params.Encode())),
			ResponseModel: &opportunitiesBatch,
		}
		_, _, e := service.httpRequest(&requestConfig)
		if e != nil {
			return nil, e
		}

		opportunities = append(opportunities, opportunitiesBatch...)

		if len(opportunitiesBatch) < int(top) {
			delete(service.nextSkips, endpoint)
			break
		}

		service.nextSkips[endpoint] += top
		rowCount += top

		if rowCount >= service.maxRowCount {
			return &opportunities, nil
		}
	}

	return &opportunities, nil
}

// CreateOpportunity creates a new contract
//
func (service *Service) CreateOpportunity(opportunity *Opportunity) (*Opportunity, *errortools.Error) {
	if opportunity == nil {
		return nil, nil
	}

	opportunityNew := Opportunity{}

	requestConfig := go_http.RequestConfig{
		Method:        http.MethodPost,
		Url:           service.url("Opportunities"),
		BodyModel:     opportunity,
		ResponseModel: &opportunityNew,
	}
	_, _, e := service.httpRequest(&requestConfig)
	if e != nil {
		return nil, e
	}

	return &opportunityNew, nil
}

// UpdateOpportunity updates an existing contract
//
func (service *Service) UpdateOpportunity(opportunity *Opportunity) (*Opportunity, *errortools.Error) {
	if opportunity == nil {
		return nil, nil
	}

	opportunityUpdated := Opportunity{}

	requestConfig := go_http.RequestConfig{
		Method:        http.MethodPut,
		Url:           service.url("Opportunities"),
		BodyModel:     opportunity,
		ResponseModel: &opportunityUpdated,
	}
	_, _, e := service.httpRequest(&requestConfig)
	if e != nil {
		return nil, e
	}

	return &opportunityUpdated, nil
}

// DeleteOpportunity deletes a specific opportunity
//
func (service *Service) DeleteOpportunity(opportunityID int64) *errortools.Error {
	requestConfig := go_http.RequestConfig{
		Method: http.MethodDelete,
		Url:    service.url(fmt.Sprintf("Opportunities/%v", opportunityID)),
	}
	_, _, e := service.httpRequest(&requestConfig)
	if e != nil {
		return e
	}

	return nil
}

// GetOpportunityLinks returns links for a specific opportunity
//
func (service *Service) GetOpportunityLinks(opportunityID int64) (*[]Link, *errortools.Error) {
	links := []Link{}

	requestConfig := go_http.RequestConfig{
		Method:        http.MethodGet,
		Url:           service.url(fmt.Sprintf("Opportunity/%v/Links", opportunityID)),
		ResponseModel: &links,
	}
	_, _, e := service.httpRequest(&requestConfig)
	if e != nil {
		return nil, e
	}

	return &links, nil
}
