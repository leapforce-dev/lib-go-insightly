package insightly

import (
	"fmt"
	"net/url"

	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
)

// LeadSource stores LeadSource from Service
//
type LeadSource struct {
	LeadSourceID int64  `json:"LEAD_SOURCE_ID"`
	LeadSource   string `json:"LEAD_SOURCE"`
	DefaultValue bool   `json:"DEFAULT_VALUE"`
	FieldOrder   int64  `json:"FIELD_ORDER"`
}

type GetLeadSourcesConfig struct {
	Skip       *uint64
	Top        *uint64
	CountTotal *bool
}

// GetLeadSources returns all leadSources
//
func (service *Service) GetLeadSources(config *GetLeadSourcesConfig) (*[]LeadSource, *errortools.Error) {
	params := url.Values{}

	endpoint := "LeadSources"
	leadSources := []LeadSource{}
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

		leadSourcesBatch := []LeadSource{}

		requestConfig := go_http.RequestConfig{
			URL:           service.url(fmt.Sprintf("%s?%s", endpoint, params.Encode())),
			ResponseModel: &leadSourcesBatch,
		}

		_, _, e := service.get(&requestConfig)
		if e != nil {
			return nil, e
		}

		leadSources = append(leadSources, leadSourcesBatch...)

		if len(leadSourcesBatch) < int(top) {
			delete(service.nextSkips, endpoint)
			break
		}

		service.nextSkips[endpoint] += top
		rowCount += top

		if rowCount >= service.maxRowCount {
			return &leadSources, nil
		}
	}

	return &leadSources, nil
}
