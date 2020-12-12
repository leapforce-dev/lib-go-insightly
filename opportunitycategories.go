package insightly

import (
	"fmt"
	"strconv"

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

// GetOpportunityCategories returns all opportunitycategories
//
func (i *Insightly) GetOpportunityCategories() ([]OpportunityCategory, *errortools.Error) {
	return i.GetOpportunityCategoriesInternal()
}

// GetOpportunityCategoriesInternal is the generic function retrieving opportunitycategories from Insightly
//
func (i *Insightly) GetOpportunityCategoriesInternal() ([]OpportunityCategory, *errortools.Error) {
	urlStr := "%sOpportunityCategories?skip=%s&top=%s"
	skip := 0
	top := 500
	rowCount := top

	opportunityCategories := []OpportunityCategory{}

	for rowCount >= top {
		url := fmt.Sprintf(urlStr, apiURL, strconv.Itoa(skip), strconv.Itoa(top))
		//fmt.Println(url)

		oc := []OpportunityCategory{}

		_, _, e := i.get(url, nil, &oc)
		if e != nil {
			return nil, e
		}

		for _, o := range oc {
			opportunityCategories = append(opportunityCategories, o)
		}

		rowCount = len(oc)
		skip += top
	}

	if len(opportunityCategories) == 0 {
		opportunityCategories = nil
	}

	return opportunityCategories, nil
}
