package insightly

import (
	"fmt"
	"strconv"
	"strings"

	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
)

// OpportunityCategory stores OpportunityCategory from Service
//
type OpportunityCategory struct {
	CategoryID      int    `json:"CATEGORY_ID"`
	CategoryName    string `json:"CATEGORY_NAME"`
	Active          bool   `json:"ACTIVE"`
	BackgroundColor string `json:"BACKGROUND_COLOR"`
}

type GetOpportunityCategoriesFilter struct {
}

// GetOpportunityCategories returns all opportunityCategories
//
func (service *Service) GetOpportunityCategories(filter *GetOpportunityCategoriesFilter) (*[]OpportunityCategory, *errortools.Error) {
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

	opportunityCategories := []OpportunityCategory{}

	for rowCount >= top {
		_opportunityCategories := []OpportunityCategory{}

		requestConfig := go_http.RequestConfig{
			URL:           service.url(fmt.Sprintf(endpointStr, searchString, strconv.Itoa(skip), strconv.Itoa(top))),
			ResponseModel: &_opportunityCategories,
		}
		_, _, e := service.get(&requestConfig)
		if e != nil {
			return nil, e
		}

		opportunityCategories = append(opportunityCategories, _opportunityCategories...)

		rowCount = len(_opportunityCategories)
		//rowCount = 0
		skip += top
	}

	if len(opportunityCategories) == 0 {
		opportunityCategories = nil
	}

	return &opportunityCategories, nil
}
