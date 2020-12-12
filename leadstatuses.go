package insightly

import (
	"fmt"
	"strconv"

	errortools "github.com/leapforce-libraries/go_errortools"
)

// LeadStatus stores LeadStatus from Insightly
//
type LeadStatus struct {
	LeadStatusID int    `json:"LEAD_STATUS_ID"`
	LeadSTatus    string `json:"LEAD_STATUS"`
	DefaultStatus bool   `json:"DEFAULT_STATUS"`
	StatusType    int    `json:"STATUS_TYPE"`
	FieldOrder    int    `json:"FIELD_ORDER"`
}

// GetLeadStatuses returns all leadStatuses
//
func (i *Insightly) GetLeadStatuses() ([]LeadStatus, *errortools.Error) {
	urlStr := "%sLeadStatuses?skip=%s&top=%s"
	skip := 0
	top := 500
	rowCount := top

	leadStatuses := []LeadStatus{}

	for rowCount >= top {
		url := fmt.Sprintf(urlStr, apiURL, strconv.Itoa(skip), strconv.Itoa(top))
		//fmt.Println(url)

		ls := []LeadStatus{}

		_, _, e := i.get(url, nil, &ls)
		if e != nil {
			return nil, e
		}

		for _, l := range ls {
			leadStatuses = append(leadStatuses, l)
		}

		rowCount = len(ls)
		skip += top
	}

	if len(leadStatuses) == 0 {
		leadStatuses = nil
	}

	return leadStatuses, nil
}
