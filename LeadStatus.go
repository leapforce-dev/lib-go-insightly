package insightly

import (
	"fmt"
	"net/url"

	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
)

// LeadStatus stores LeadStatus from Service
//
type LeadStatus struct {
	LeadStatusID  int64  `json:"LEAD_STATUS_ID"`
	LeadStatus    string `json:"LEAD_STATUS"`
	DefaultStatus bool   `json:"DEFAULT_STATUS"`
	StatusType    int64  `json:"STATUS_TYPE"`
	FieldOrder    int64  `json:"FIELD_ORDER"`
}

type GetLeadStatusesConfig struct {
	Skip             *uint64
	Top              *uint64
	CountTotal       *bool
	IncludeConverted *bool
}

// GetLeadStatuses returns all leadStatuses
//
func (service *Service) GetLeadStatuses(config *GetLeadStatusesConfig) (*[]LeadStatus, *errortools.Error) {
	params := url.Values{}

	endpoint := "LeadStatuses"
	leadStatuses := []LeadStatus{}
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
		if config.IncludeConverted != nil {
			params.Set("include_converted", fmt.Sprintf("%v", *config.IncludeConverted))
		}
	}

	params.Set("top", fmt.Sprintf("%v", top))

	for true {
		params.Set("skip", fmt.Sprintf("%v", service.nextSkips[endpoint]))

		leadStatusesBatch := []LeadStatus{}

		requestConfig := go_http.RequestConfig{
			URL:           service.url(fmt.Sprintf("%s?%s", endpoint, params.Encode())),
			ResponseModel: &leadStatusesBatch,
		}

		_, _, e := service.get(&requestConfig)
		if e != nil {
			return nil, e
		}

		leadStatuses = append(leadStatuses, leadStatusesBatch...)

		if len(leadStatusesBatch) < int(top) {
			delete(service.nextSkips, endpoint)
			break
		}

		service.nextSkips[endpoint] += top
		rowCount += top

		if rowCount >= service.maxRowCount {
			return &leadStatuses, nil
		}
	}

	return &leadStatuses, nil
}
