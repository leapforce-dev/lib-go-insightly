package insightly

import (
	"fmt"
	"net/http"
	"net/url"

	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
)

// ProjectCategory stores ProjectCategory from Service
//
type ProjectCategory struct {
	CategoryID      int64  `json:"CATEGORY_ID"`
	CategoryName    string `json:"CATEGORY_NAME"`
	Active          bool   `json:"ACTIVE"`
	BackgroundColor string `json:"BACKGROUND_COLOR"`
}

type GetProjectCategoriesConfig struct {
	Skip       *uint64
	Top        *uint64
	CountTotal *bool
}

// GetProjectCategories returns all projectCategories
//
func (service *Service) GetProjectCategories(config *GetProjectCategoriesConfig) (*[]ProjectCategory, *errortools.Error) {
	params := url.Values{}

	endpoint := "ProjectCategories"
	projectCategories := []ProjectCategory{}
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

	for {
		params.Set("skip", fmt.Sprintf("%v", service.nextSkips[endpoint]))

		projectCategoriesBatch := []ProjectCategory{}

		requestConfig := go_http.RequestConfig{
			Method:        http.MethodGet,
			URL:           service.url(fmt.Sprintf("%s?%s", endpoint, params.Encode())),
			ResponseModel: &projectCategoriesBatch,
		}

		_, _, e := service.httpRequest(&requestConfig)
		if e != nil {
			return nil, e
		}

		projectCategories = append(projectCategories, projectCategoriesBatch...)

		if len(projectCategoriesBatch) < int(top) {
			delete(service.nextSkips, endpoint)
			break
		}

		service.nextSkips[endpoint] += top
		rowCount += top

		if rowCount >= service.maxRowCount {
			return &projectCategories, nil
		}
	}

	return &projectCategories, nil
}
