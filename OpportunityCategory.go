package insightly

import (
	"fmt"
	"strconv"
	"strings"

	errortools "github.com/leapforce-libraries/go_errortools"
)

// OpportunityCategory stores OpportunityCategory from Insightly
//
type OpportunityCategory struct {
	CategoryID      int    `json:"CATEGORY_ID"`
	CategoryName    string `json:"CATEGORY_NAME"`
	Active          bool   `json:"ACTIVE"`
	BackgroundColor string `json:"BACKGROUND_COLOR"`
}

type GetOpportunityCategoriesFilter struct {
}

// GetOpportunityCategories returns all opportunityCategorys
//
func (i *Insightly) GetOpportunityCategories(filter *GetOpportunityCategoriesFilter) (*[]OpportunityCategory, *errortools.Error) {
	searchString := "?"
	searchFilter := []string{}

	if filter != nil {
	}

	if len(searchFilter) > 0 {
		searchString = "/Search?" + strings.Join(searchFilter, "&")
	}

	endpointStr := "OpportunityCategories%sskip=%s&top=%s"
	skip := 0
	top := 100
	rowCount := top

	opportunityCategorys := []OpportunityCategory{}

	for rowCount >= top {
		endpoint := fmt.Sprintf(endpointStr, searchString, strconv.Itoa(skip), strconv.Itoa(top))
		//fmt.Println(endpoint)

		cs := []OpportunityCategory{}

		_, _, e := i.get(endpoint, nil, &cs)
		if e != nil {
			return nil, e
		}

		opportunityCategorys = append(opportunityCategorys, cs...)

		rowCount = len(cs)
		//rowCount = 0
		skip += top
	}

	if len(opportunityCategorys) == 0 {
		opportunityCategorys = nil
	}

	return &opportunityCategorys, nil
}
