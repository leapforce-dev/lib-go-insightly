package insightly

import (
	"fmt"
	"net/http"
	"net/url"

	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
)

// TaskCategory stores TaskCategory from Service
//
type TaskCategory struct {
	CategoryID      int64  `json:"CATEGORY_ID"`
	CategoryName    string `json:"CATEGORY_NAME"`
	Active          bool   `json:"ACTIVE"`
	BackgroundColor string `json:"BACKGROUND_COLOR"`
}

type GetTaskCategoriesConfig struct {
	Skip       *uint64
	Top        *uint64
	CountTotal *bool
}

// GetTaskCategories returns all taskCategories
//
func (service *Service) GetTaskCategories(config *GetTaskCategoriesConfig) (*[]TaskCategory, *errortools.Error) {
	params := url.Values{}

	endpoint := "TaskCategories"
	taskCategories := []TaskCategory{}
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

		taskCategoriesBatch := []TaskCategory{}

		requestConfig := go_http.RequestConfig{
			Method:        http.MethodGet,
			Url:           service.url(fmt.Sprintf("%s?%s", endpoint, params.Encode())),
			ResponseModel: &taskCategoriesBatch,
		}

		_, _, e := service.httpRequest(&requestConfig)
		if e != nil {
			return nil, e
		}

		taskCategories = append(taskCategories, taskCategoriesBatch...)

		if len(taskCategoriesBatch) < int(top) {
			delete(service.nextSkips, endpoint)
			break
		}

		service.nextSkips[endpoint] += top
		rowCount += top

		if rowCount >= service.maxRowCount {
			return &taskCategories, nil
		}
	}

	return &taskCategories, nil
}
