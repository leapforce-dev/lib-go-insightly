package insightly

import (
	"fmt"
	"net/url"

	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
)

// Pipeline stores Pipeline from Service
//
type Pipeline struct {
	PipelineID       int64  `json:"PIPELINE_ID"`
	PipelineName     string `json:"PIPELINE_NAME"`
	ForOpportunities bool   `json:"FOR_OPPORTUNITIES"`
	ForProjects      bool   `json:"FOR_PROJECTS"`
	OwnerUserID      int64  `json:"OWNER_USER_ID"`
}

type GetPipelinesConfig struct {
	Skip       *uint64
	Top        *uint64
	CountTotal *bool
}

// GetPipelines returns all pipelines
//
func (service *Service) GetPipelines(config *GetPipelinesConfig) (*[]Pipeline, *errortools.Error) {
	params := url.Values{}

	endpoint := "Pipelines"
	pipelines := []Pipeline{}
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

		pipelinesBatch := []Pipeline{}

		requestConfig := go_http.RequestConfig{
			URL:           service.url(fmt.Sprintf("%s?%s", endpoint, params.Encode())),
			ResponseModel: &pipelinesBatch,
		}

		_, _, e := service.get(&requestConfig)
		if e != nil {
			return nil, e
		}

		pipelines = append(pipelines, pipelinesBatch...)

		if len(pipelinesBatch) < int(top) {
			delete(service.nextSkips, endpoint)
			break
		}

		service.nextSkips[endpoint] += top
		rowCount += top

		if rowCount >= service.maxRowCount {
			return &pipelines, nil
		}
	}

	return &pipelines, nil
}
