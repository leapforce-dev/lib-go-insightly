package insightly

import (
	"fmt"
	"net/url"

	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
)

// OpportunityStateReason stores OpportunityStateReason from Service
//
type OpportunityStateReason struct {
	StateReasonID       int64  `json:"STATE_REASON_ID"`
	StateReason         string `json:"STATE_REASON"`
	ForOpportunityState string `json:"FOR_OPPORTUNITY_STATE"`
}

type GetOpportunityStateReasonsConfig struct {
	Skip       *uint64
	Top        *uint64
	CountTotal *bool
}

// GetOpportunityStateReasons returns all opportunityStateReasons
//
func (service *Service) GetOpportunityStateReasons(config *GetOpportunityStateReasonsConfig) (*[]OpportunityStateReason, *errortools.Error) {
	params := url.Values{}

	endpoint := "OpportunityStateReasons"
	opportunityStateReasons := []OpportunityStateReason{}
	rowCount := uint64(0)
	top := defaultTop

	if config != nil {
		if config.Top != nil {
			top = *config.Top
		}
		if config.Skip != nil {
			service.nextSkips[endpoint] = *config.Skip
		}
		if config.CountTotal != nil {
			params.Set("count_total", fmt.Sprintf("%v", *config.CountTotal))
		}
	}

	params.Set("top", fmt.Sprintf("%v", top))

	for true {
		params.Set("skip", fmt.Sprintf("%v", service.nextSkips[endpoint]))

		opportunityStateReasonsBatch := []OpportunityStateReason{}

		requestConfig := go_http.RequestConfig{
			URL:           service.url(fmt.Sprintf("%s?%s", endpoint, params.Encode())),
			ResponseModel: &opportunityStateReasonsBatch,
		}

		_, _, e := service.get(&requestConfig)
		if e != nil {
			return nil, e
		}

		opportunityStateReasons = append(opportunityStateReasons, opportunityStateReasonsBatch...)

		if len(opportunityStateReasonsBatch) < int(top) {
			delete(service.nextSkips, endpoint)
			break
		}

		service.nextSkips[endpoint] += top
		rowCount += top

		if rowCount >= service.maxRowCount {
			return &opportunityStateReasons, nil
		}
	}

	return &opportunityStateReasons, nil
}
