package insightly

import (
	"fmt"
	"strconv"

	errortools "github.com/leapforce-libraries/go_errortools"
)

// LeadSource stores LeadSource from Insightly
//
type LeadSource struct {
	LeadSourceID int    `json:"LEAD_SOURCE_ID"`
	LeadSource   string `json:"LEAD_SOURCE"`
	DefaultValue bool   `json:"DEFAULT_VALUE"`
	FieldOrder   int    `json:"FIELD_ORDER"`
}

// GetLeadSources returns all leadSources
//
func (i *Insightly) GetLeadSources() ([]LeadSource, *errortools.Error) {
	endpointStr := "LeadSources?skip=%s&top=%s"
	skip := 0
	top := 500
	rowCount := top

	leadSources := []LeadSource{}

	for rowCount >= top {
		endpoint := fmt.Sprintf(endpointStr, strconv.Itoa(skip), strconv.Itoa(top))
		//fmt.Println(endpoint)

		ls := []LeadSource{}

		_, _, e := i.get(endpoint, nil, &ls)
		if e != nil {
			return nil, e
		}

		for _, l := range ls {
			leadSources = append(leadSources, l)
		}

		rowCount = len(ls)
		skip += top
	}

	if len(leadSources) == 0 {
		leadSources = nil
	}

	return leadSources, nil
}
