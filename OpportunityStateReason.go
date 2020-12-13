package insightly

import (
	"fmt"
	"strconv"
	"strings"

	errortools "github.com/leapforce-libraries/go_errortools"
)

// OpportunityStateReason stores OpportunityStateReason from Insightly
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
func (i *Insightly) GetOpportunityStateReasons(filter *GetOpportunityStateReasonsFilter) (*[]OpportunityStateReason, *errortools.Error) {
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
		endpoint := fmt.Sprintf(endpointStr, searchString, strconv.Itoa(skip), strconv.Itoa(top))
		//fmt.Println(endpoint)

		cs := []OpportunityStateReason{}

		_, _, e := i.get(endpoint, nil, &cs)
		if e != nil {
			return nil, e
		}

		opportunityStateReasons = append(opportunityStateReasons, cs...)

		rowCount = len(cs)
		//rowCount = 0
		skip += top
	}

	if len(opportunityStateReasons) == 0 {
		opportunityStateReasons = nil
	}

	return &opportunityStateReasons, nil
}
