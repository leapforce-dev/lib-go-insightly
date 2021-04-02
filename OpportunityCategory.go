package insightly

import (
	"fmt"
	"net/url"

	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
)

// OpportunityCategory stores OpportunityCategory from Service
//
type OpportunityCategory struct {
	CategoryID      int64  `json:"CATEGORY_ID"`
	CategoryName    string `json:"CATEGORY_NAME"`
	Active          bool   `json:"ACTIVE"`
	BackgroundColor string `json:"BACKGROUND_COLOR"`
}

type GetOpportunityCategoriesConfig struct {
	Skip       *uint64
	Top        *uint64
	CountTotal *bool
}

// GetOpportunityCategories returns all opportunityCategories
//
func (service *Service) GetOpportunityCategories(config *GetOpportunityCategoriesConfig) (*[]OpportunityCategory, *errortools.Error) {
	params := url.Values{}

	endpoint := "OpportunityCategories"
	opportunityCategories := []OpportunityCategory{}
	rowCount := uint64(0)
	top := defaultTop

	if config != nil {
		if config.Top != nil {
			top = *config.Top
		}
		if config.Skip != nil {
			service.nextSkips[endpoint] = *config.Skip
		}
		if config.CountTotal != nil {
			params.Set("count_total", fmt.Sprintf("%v", *config.CountTotal))
		}
	}

	params.Set("top", fmt.Sprintf("%v", top))

	for true {
		params.Set("skip", fmt.Sprintf("%v", service.nextSkips[endpoint]))

		opportunityCategoriesBatch := []OpportunityCategory{}

		requestConfig := go_http.RequestConfig{
			URL:           service.url(fmt.Sprintf("%s?%s", endpoint, params.Encode())),
			ResponseModel: &opportunityCategoriesBatch,
		}

		_, _, e := service.get(&requestConfig)
		if e != nil {
			return nil, e
		}

		opportunityCategories = append(opportunityCategories, opportunityCategoriesBatch...)

		if len(opportunityCategoriesBatch) < int(top) {
			delete(service.nextSkips, endpoint)
			break
		}

		service.nextSkips[endpoint] += top
		rowCount += top

		if rowCount >= service.maxRowCount {
			return &opportunityCategories, nil
		}
	}

	return &opportunityCategories, nil
}
