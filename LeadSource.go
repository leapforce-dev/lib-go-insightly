package insightly

import (
	"fmt"
	"strconv"
	"strings"

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
}

// GetLeadSources returns all leadSources
//
func (service *Service) GetLeadSources(config *GetLeadSourcesConfig) (*[]LeadSource, *errortools.Error) {
	searchString := "?"
	searchFilter := []string{}

	if config != nil {
	}

	if len(searchFilter) > 0 {
		searchString = "/Search?" + strings.Join(searchFilter, "&")
	}

	endpointStr := "LeadSources%sskip=%s&top=%s"
	skip := 0
	top := 100
	rowCount := top

	leadSources := []LeadSource{}

	for rowCount >= top {
		_leadSources := []LeadSource{}

		requestConfig := go_http.RequestConfig{
			URL:           service.url(fmt.Sprintf(endpointStr, searchString, strconv.Itoa(skip), strconv.Itoa(top))),
			ResponseModel: &_leadSources,
		}
		_, _, e := service.get(&requestConfig)
		if e != nil {
			return nil, e
		}

		leadSources = append(leadSources, _leadSources...)

		rowCount = len(_leadSources)
		//rowCount = 0
		skip += top
	}

	if len(leadSources) == 0 {
		leadSources = nil
	}

	return &leadSources, nil
}
