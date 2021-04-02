package insightly

import (
	"fmt"
	"net/url"

	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
)

// FileCategory stores FileCategory from Service
//
type FileCategory struct {
	CategoryID      int64  `json:"CATEGORY_ID"`
	CategoryName    string `json:"CATEGORY_NAME"`
	Active          bool   `json:"ACTIVE"`
	BackgroundColor string `json:"BACKGROUND_COLOR"`
}

type GetFileCategoriesConfig struct {
	Skip       *uint64
	Top        *uint64
	CountTotal *bool
}

// GetFileCategories returns all fileCategories
//
func (service *Service) GetFileCategories(config *GetFileCategoriesConfig) (*[]FileCategory, *errortools.Error) {
	params := url.Values{}

	endpoint := "FileCategories"
	fileCategories := []FileCategory{}
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

		fileCategoriesBatch := []FileCategory{}

		requestConfig := go_http.RequestConfig{
			URL:           service.url(fmt.Sprintf("%s?%s", endpoint, params.Encode())),
			ResponseModel: &fileCategoriesBatch,
		}

		_, _, e := service.get(&requestConfig)
		if e != nil {
			return nil, e
		}

		fileCategories = append(fileCategories, fileCategoriesBatch...)

		if len(fileCategoriesBatch) < int(top) {
			delete(service.nextSkips, endpoint)
			break
		}

		service.nextSkips[endpoint] += top
		rowCount += top

		if rowCount >= service.maxRowCount {
			return &fileCategories, nil
		}
	}

	return &fileCategories, nil
}
