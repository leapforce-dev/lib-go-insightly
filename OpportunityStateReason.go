package insightly

import (
	"fmt"
	"strconv"
	"strings"

	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
)

// OpportunityStateReason stores OpportunityStateReason from Service
//
type OpportunityStateReason struct {
	StateReasonID       int    `json:"STATE_REASON_ID"`
	StateReason         string `json:"STATE_REASON"`
	ForOpportunityState string `json:"FOR_OPPORTUNITY_STATE"`
}

type GetOpportunityStateReasonsFilter struct {
}

// GetOpportunityStateReasons returns all opportunityStateReasons
//
func (service *Service) GetOpportunityStateReasons(filter *GetOpportunityStateReasonsFilter) (*[]OpportunityStateReason, *errortools.Error) {
	searchString := "?"
	searchFilter := []string{}

	if filter != nil {
	}

	if len(searchFilter) > 0 {
		searchString = "/Search?" + strings.Join(searchFilter, "&")
	}

	endpointStr := "OpportunityStateReasons%sskip=%s&top=%s"
	skip := 0
	top := 100
	rowCount := top

	opportunityStateReasons := []OpportunityStateReason{}

	for rowCount >= top {
		_opportunityStateReasons := []OpportunityStateReason{}

		requestConfig := go_http.RequestConfig{
			URL:           service.url(fmt.Sprintf(endpointStr, searchString, strconv.Itoa(skip), strconv.Itoa(top))),
			ResponseModel: &_opportunityStateReasons,
		}
		_, _, e := service.get(&requestConfig)
		if e != nil {
			return nil, e
		}

		opportunityStateReasons = append(opportunityStateReasons, _opportunityStateReasons...)

		rowCount = len(_opportunityStateReasons)
		//rowCount = 0
		skip += top
	}

	if len(opportunityStateReasons) == 0 {
		opportunityStateReasons = nil
	}

	return &opportunityStateReasons, nil
}
