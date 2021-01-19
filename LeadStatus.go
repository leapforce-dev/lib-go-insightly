package insightly

import (
	"fmt"
	"strconv"
	"strings"

	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
)

// LeadStatus stores LeadStatus from Service
//
type LeadStatus struct {
	LeadStatusID  int    `json:"LEAD_STATUS_ID"`
	LeadSTatus    string `json:"LEAD_STATUS"`
	DefaultStatus bool   `json:"DEFAULT_STATUS"`
	StatusType    int    `json:"STATUS_TYPE"`
	FieldOrder    int    `json:"FIELD_ORDER"`
}

type GetLeadStatusesFilter struct {
}

// GetLeadStatuses returns all leadStatuses
//
func (service *Service) GetLeadStatuses(filter *GetLeadStatusesFilter) (*[]LeadStatus, *errortools.Error) {
	searchString := "?"
	searchFilter := []string{}

	if filter != nil {
	}

	if len(searchFilter) > 0 {
		searchString = "/Search?" + strings.Join(searchFilter, "&")
	}

	endpointStr := "LeadStatuses%sskip=%s&top=%s"
	skip := 0
	top := 100
	rowCount := top

	leadStatuses := []LeadStatus{}

	for rowCount >= top {
		_leadStatuses := []LeadStatus{}

		requestConfig := go_http.RequestConfig{
			URL:           service.url(fmt.Sprintf(endpointStr, searchString, strconv.Itoa(skip), strconv.Itoa(top))),
			ResponseModel: &_leadStatuses,
		}
		_, _, e := service.get(&requestConfig)
		if e != nil {
			return nil, e
		}

		leadStatuses = append(leadStatuses, _leadStatuses...)

		rowCount = len(_leadStatuses)
		//rowCount = 0
		skip += top
	}

	if len(leadStatuses) == 0 {
		leadStatuses = nil
	}

	return &leadStatuses, nil
}
