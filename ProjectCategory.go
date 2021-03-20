package insightly

import (
	"fmt"
	"strconv"
	"strings"

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
}

// GetProjectCategories returns all projectCategories
//
func (service *Service) GetProjectCategories(config *GetProjectCategoriesConfig) (*[]ProjectCategory, *errortools.Error) {
	searchString := "?"
	searchFilter := []string{}

	if config != nil {
	}

	if len(searchFilter) > 0 {
		searchString = "/Search?" + strings.Join(searchFilter, "&")
	}

	endpointStr := "ProjectCategories%sskip=%s&top=%s"
	skip := 0
	top := 100
	rowCount := top

	projectCategories := []ProjectCategory{}

	for rowCount >= top {
		_projectCategories := []ProjectCategory{}

		requestConfig := go_http.RequestConfig{
			URL:           service.url(fmt.Sprintf(endpointStr, searchString, strconv.Itoa(skip), strconv.Itoa(top))),
			ResponseModel: &_projectCategories,
		}
		_, _, e := service.get(&requestConfig)
		if e != nil {
			return nil, e
		}

		projectCategories = append(projectCategories, _projectCategories...)

		rowCount = len(_projectCategories)
		//rowCount = 0
		skip += top
	}

	if len(projectCategories) == 0 {
		projectCategories = nil
	}

	return &projectCategories, nil
}
