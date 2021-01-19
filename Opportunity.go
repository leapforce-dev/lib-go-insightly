package insightly

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
)

// Opportunity stores Opportunity from Service
//
type Opportunity struct {
	OpportunityID       int           `json:"OPPORTUNITY_ID"`
	OpportunityName     *string       `json:"OPPORTUNITY_NAME"`
	OpportunityDetails  *string       `json:"OPPORTUNITY_DETAILS"`
	OpportunityState    *string       `json:"OPPORTUNITY_STATE"`
	ResponsibleUserID   *int          `json:"RESPONSIBLE_USER_ID"`
	CategoryID          *int          `json:"CATEGORY_ID"`
	ImageURL            *string       `json:"IMAGE_URL"`
	BidCurrency         *string       `json:"BID_CURRENCY"`
	BidAmount           *float32      `json:"BID_AMOUNT"`
	BidType             *string       `json:"BID_TYPE"`
	BidDuration         *int          `json:"BID_DURATION"`
	ActualCloseDate     *DateUTC      `json:"ACTUAL_CLOSE_DATE"`
	DateCreatedUTC      *DateUTC      `json:"DATE_CREATED_UTC"`
	DateUpdatedUTC      *DateUTC      `json:"DATE_UPDATED_UTC"`
	OpportunityValue    *float32      `json:"OPPORTUNITY_VALUE"`
	Probability         *int          `json:"PROBABILITY"`
	ForecastCloseDate   *DateUTC      `json:"FORECAST_CLOSE_DATE"`
	OwnerUserID         *int          `json:"OWNER_USER_ID"`
	LastActivityDateUTC *DateUTC      `json:"LAST_ACTIVITY_DATE_UTC"`
	NextActivityDateUTC *DateUTC      `json:"NEXT_ACTIVITY_DATE_UTC"`
	PipelineID          *int          `json:"PIPELINE_ID"`
	StageID             *int          `json:"STAGE_ID"`
	CreatedUserID       *int          `json:"CREATED_USER_ID"`
	OrganisationID      *int          `json:"ORGANISATION_ID"`
	CustomFields        *CustomFields `json:"CUSTOMFIELDS"`
	Tags                *[]Tag        `json:"TAGS"`
}

func (o *Opportunity) prepareMarshal() interface{} {
	if o == nil {
		return nil
	}

	return &struct {
		OpportunityID      int           `json:"OPPORTUNITY_ID"`
		OpportunityName    *string       `json:"OPPORTUNITY_NAME,omitempty"`
		OpportunityDetails *string       `json:"OPPORTUNITY_DETAILS,omitempty"`
		OpportunityState   *string       `json:"OPPORTUNITY_STATE,omitempty"`
		ResponsibleUserID  *int          `json:"RESPONSIBLE_USER_ID,omitempty"`
		CategoryID         *int          `json:"CATEGORY_ID,omitempty"`
		ImageURL           *string       `json:"IMAGE_URL,omitempty"`
		BidCurrency        *string       `json:"BID_CURRENCY,omitempty"`
		BidAmount          *float32      `json:"BID_AMOUNT,omitempty"`
		BidType            *string       `json:"BID_TYPE,omitempty"`
		BidDuration        *int          `json:"BID_DURATION,omitempty"`
		ActualCloseDate    *DateUTC      `json:"ACTUAL_CLOSE_DATE,omitempty"`
		OpportunityValue   *float32      `json:"OPPORTUNITY_VALUE,omitempty"`
		Probability        *int          `json:"PROBABILITY,omitempty"`
		ForecastCloseDate  *DateUTC      `json:"FORECAST_CLOSE_DATE,omitempty"`
		OwnerUserID        *int          `json:"OWNER_USER_ID,omitempty"`
		PipelineID         *int          `json:"PIPELINE_ID,omitempty"`
		StageID            *int          `json:"STAGE_ID,omitempty"`
		OrganisationID     *int          `json:"ORGANISATION_ID,omitempty"`
		CustomFields       *CustomFields `json:"CUSTOMFIELDS,omitempty"`
	}{
		o.OpportunityID,
		o.OpportunityName,
		o.OpportunityDetails,
		o.OpportunityState,
		o.ResponsibleUserID,
		o.CategoryID,
		o.ImageURL,
		o.BidCurrency,
		o.BidAmount,
		o.BidType,
		o.BidDuration,
		o.ActualCloseDate,
		o.OpportunityValue,
		o.Probability,
		o.ForecastCloseDate,
		o.OwnerUserID,
		o.PipelineID,
		o.StageID,
		o.OrganisationID,
		o.CustomFields,
	}
}

// GetOpportunity returns a specific opportunity
//
func (service *Service) GetOpportunity(opportunityID int) (*Opportunity, *errortools.Error) {
	opportunity := Opportunity{}

	requestConfig := go_http.RequestConfig{
		URL:           service.url(fmt.Sprintf("Opportunities/%v", opportunityID)),
		ResponseModel: &opportunity,
	}
	_, _, e := service.get(&requestConfig)
	if e != nil {
		return nil, e
	}

	return &opportunity, nil
}

type GetOpportunitiesFilter struct {
	UpdatedAfter *time.Time
	Field        *struct {
		FieldName  string
		FieldValue string
	}
}

// GetOpportunities returns all opportunities
//
func (service *Service) GetOpportunities(filter *GetOpportunitiesFilter) (*[]Opportunity, *errortools.Error) {
	searchString := "?"
	searchFilter := []string{}

	if filter != nil {
		if filter.UpdatedAfter != nil {
			from := filter.UpdatedAfter.Format(DateFormat)
			searchFilter = append(searchFilter, fmt.Sprintf("updated_after_utc=%s&", from))
		}

		if filter.Field != nil {
			searchFilter = append(searchFilter, fmt.Sprintf("field_name=%s&field_value=%s&", filter.Field.FieldName, filter.Field.FieldValue))
		}
	}

	if len(searchFilter) > 0 {
		searchString = "/Search?" + strings.Join(searchFilter, "&")
	}

	endpointStr := "Opportunities%sskip=%s&top=%s"
	skip := 0
	top := 100
	rowCount := top

	opportunities := []Opportunity{}

	for rowCount >= top {
		_opportunities := []Opportunity{}

		requestConfig := go_http.RequestConfig{
			URL:           service.url(fmt.Sprintf(endpointStr, searchString, strconv.Itoa(skip), strconv.Itoa(top))),
			ResponseModel: &_opportunities,
		}
		_, _, e := service.get(&requestConfig)
		if e != nil {
			return nil, e
		}

		opportunities = append(opportunities, _opportunities...)

		rowCount = len(_opportunities)
		//rowCount = 0
		skip += top
	}

	if len(opportunities) == 0 {
		opportunities = nil
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
		URL:           service.url("Opportunities"),
		BodyModel:     opportunity.prepareMarshal(),
		ResponseModel: &opportunityNew,
	}
	_, _, e := service.post(&requestConfig)
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
		URL:           service.url("Opportunities"),
		BodyModel:     opportunity.prepareMarshal(),
		ResponseModel: &opportunityUpdated,
	}
	_, _, e := service.put(&requestConfig)
	if e != nil {
		return nil, e
	}

	return &opportunityUpdated, nil
}

// DeleteOpportunity deletes a specific opportunity
//
func (service *Service) DeleteOpportunity(opportunityID int) *errortools.Error {
	requestConfig := go_http.RequestConfig{
		URL: service.url(fmt.Sprintf("Opportunities/%v", opportunityID)),
	}
	_, _, e := service.delete(&requestConfig)
	if e != nil {
		return e
	}

	return nil
}

// GetOpportunityLinks returns links for a specific opportunity
//
func (service *Service) GetOpportunityLinks(opportunityID int) (*[]Link, *errortools.Error) {
	links := []Link{}

	requestConfig := go_http.RequestConfig{
		URL:           service.url(fmt.Sprintf("Opportunity/%v/Links", opportunityID)),
		ResponseModel: &links,
	}
	_, _, e := service.get(&requestConfig)
	if e != nil {
		return nil, e
	}

	return &links, nil
}
