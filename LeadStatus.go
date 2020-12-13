package insightly

import (
	"fmt"
	"strconv"
	"strings"

	errortools "github.com/leapforce-libraries/go_errortools"
)

// LeadStatus stores LeadStatus from Insightly
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
func (i *Insightly) GetLeadStatuses(filter *GetLeadStatusesFilter) (*[]LeadStatus, *errortools.Error) {
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
		endpoint := fmt.Sprintf(endpointStr, searchString, strconv.Itoa(skip), strconv.Itoa(top))
		//fmt.Println(endpoint)

		cs := []LeadStatus{}

		_, _, e := i.get(endpoint, nil, &cs)
		if e != nil {
			return nil, e
		}

		leadStatuses = append(leadStatuses, cs...)

		rowCount = len(cs)
		//rowCount = 0
		skip += top
	}

	if len(leadStatuses) == 0 {
		leadStatuses = nil
	}

	return &leadStatuses, nil
}
