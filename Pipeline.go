package insightly

import (
	"fmt"
	"strconv"
	"strings"

	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
)

// Pipeline stores Pipeline from Service
//
type Pipeline struct {
	PipelineID       int    `json:"PIPELINE_ID"`
	PipelineName     string `json:"PIPELINE_NAME"`
	ForOpportunities bool   `json:"FOR_OPPORTUNITIES"`
	ForProjects      bool   `json:"FOR_PROJECTS"`
	OwnerUserID      *int   `json:"OWNER_USER_ID"`
}

type GetPipelinesFilter struct {
}

// GetPipelines returns all pipelines
//
func (service *Service) GetPipelines(filter *GetPipelinesFilter) (*[]Pipeline, *errortools.Error) {
	searchString := "?"
	searchFilter := []string{}

	if filter != nil {
	}

	if len(searchFilter) > 0 {
		searchString = "/Search?" + strings.Join(searchFilter, "&")
	}

	endpointStr := "Pipelines%sskip=%s&top=%s"
	skip := 0
	top := 100
	rowCount := top

	pipelines := []Pipeline{}

	for rowCount >= top {
		_pipelines := []Pipeline{}

		requestConfig := go_http.RequestConfig{
			URL:           service.url(fmt.Sprintf(endpointStr, searchString, strconv.Itoa(skip), strconv.Itoa(top))),
			ResponseModel: &_pipelines,
		}
		_, _, e := service.get(&requestConfig)
		if e != nil {
			return nil, e
		}

		pipelines = append(pipelines, _pipelines...)

		rowCount = len(_pipelines)
		//rowCount = 0
		skip += top
	}

	if len(pipelines) == 0 {
		pipelines = nil
	}

	return &pipelines, nil
}
