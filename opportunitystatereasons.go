package insightly

import (
	"fmt"
	"strconv"

	errortools "github.com/leapforce-libraries/go_errortools"
)

// OpportunityStateReason stores OpportunityStateReason from Insightly
//
type OpportunityStateReason struct {
	StateReasonID       int    `json:"STATE_REASON_ID"`
	StateReason         string `json:"STATE_REASON"`
	ForOpportunityState string `json:"FOR_OPPORTUNITY_STATE"`
}

// GetOpportunityStateReasons returns all opportunityStateReasons
//
func (i *Insightly) GetOpportunityStateReasons() ([]OpportunityStateReason, *errortools.Error) {
	endpointStr := "OpportunityStateReasons?skip=%s&top=%s"
	skip := 0
	top := 500
	rowCount := top

	opportunityStateReasons := []OpportunityStateReason{}

	for rowCount >= top {
		endpoint := fmt.Sprintf(endpointStr, strconv.Itoa(skip), strconv.Itoa(top))
		//fmt.Println(endpoint)

		ls := []OpportunityStateReason{}

		_, _, e := i.get(endpoint, nil, &ls)
		if e != nil {
			return nil, e
		}

		for _, l := range ls {
			opportunityStateReasons = append(opportunityStateReasons, l)
		}

		rowCount = len(ls)
		skip += top
	}

	if len(opportunityStateReasons) == 0 {
		opportunityStateReasons = nil
	}

	return opportunityStateReasons, nil
}
