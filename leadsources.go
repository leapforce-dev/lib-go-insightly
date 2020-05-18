package insightly

import (
	"fmt"
	"strconv"
)

// LeadSource stores LeadSource from Insightly
//
type LeadSource struct {
	LEAD_SOURCE_ID int    `json:"LEAD_SOURCE_ID"`
	LEAD_SOURCE    string `json:"LEAD_SOURCE"`
	DEFAULT_VALUE  bool   `json:"DEFAULT_VALUE"`
	FIELD_ORDER    int    `json:"FIELD_ORDER"`
}

// GetLeadSources returns all leadSources
//
func (i *Insightly) GetLeadSources() ([]LeadSource, error) {
	urlStr := "%sLeadSources?skip=%s&top=%s"
	skip := 0
	top := 500
	rowCount := top

	leadSources := []LeadSource{}

	for rowCount >= top {
		url := fmt.Sprintf(urlStr, i.apiURL, strconv.Itoa(skip), strconv.Itoa(top))
		//fmt.Println(url)

		ls := []LeadSource{}

		err := i.Get(url, &ls)
		if err != nil {
			return nil, err
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
