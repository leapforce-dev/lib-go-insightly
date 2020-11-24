package insightly

import (
	"fmt"
	"strconv"

	errortools "github.com/leapforce-libraries/go_errortools"
)

// OpportunityStateReason stores OpportunityStateReason from Insightly
//
type OpportunityStateReason struct {
	STATE_REASON_ID       int    `json:"STATE_REASON_ID"`
	STATE_REASON          string `json:"STATE_REASON"`
	FOR_OPPORTUNITY_STATE string `json:"FOR_OPPORTUNITY_STATE"`
}

// GetOpportunityStateReasons returns all opportunityStateReasons
//
func (i *Insightly) GetOpportunityStateReasons() ([]OpportunityStateReason, *errortools.Error) {
	urlStr := "%sOpportunityStateReasons?skip=%s&top=%s"
	skip := 0
	top := 500
	rowCount := top

	opportunityStateReasons := []OpportunityStateReason{}

	for rowCount >= top {
		url := fmt.Sprintf(urlStr, i.apiURL, strconv.Itoa(skip), strconv.Itoa(top))
		//fmt.Println(url)

		ls := []OpportunityStateReason{}

		e := i.Get(url, &ls)
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
