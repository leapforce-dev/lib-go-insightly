package insightly

import (
	"fmt"
	"net/http"
	"net/url"

	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
)

// PipelineStage stores PipelineStage from Service
//
type PipelineStage struct {
	StageID       int64  `json:"STAGE_ID"`
	PipelineID    int64  `json:"PIPELINE_ID"`
	StageName     string `json:"STAGE_NAME"`
	StageOrder    int64  `json:"STAGE_ORDER"`
	ActivitySetID *int64 `json:"ACTIVITYSET_ID"`
	OwnerUserID   int64  `json:"OWNER_USER_ID"`
}

type GetPipelineStagesConfig struct {
	Skip       *uint64
	Top        *uint64
	CountTotal *bool
}

// GetPipelineStages returns all pipelineStages
//
func (service *Service) GetPipelineStages(config *GetPipelineStagesConfig) (*[]PipelineStage, *errortools.Error) {
	params := url.Values{}

	endpoint := "PipelineStages"
	pipelineStages := []PipelineStage{}
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

		pipelineStagesBatch := []PipelineStage{}

		requestConfig := go_http.RequestConfig{
			Method:        http.MethodGet,
			Url:           service.url(fmt.Sprintf("%s?%s", endpoint, params.Encode())),
			ResponseModel: &pipelineStagesBatch,
		}

		_, _, e := service.httpRequest(&requestConfig)
		if e != nil {
			return nil, e
		}

		pipelineStages = append(pipelineStages, pipelineStagesBatch...)

		if len(pipelineStagesBatch) < int(top) {
			delete(service.nextSkips, endpoint)
			break
		}

		service.nextSkips[endpoint] += top
		rowCount += top

		if rowCount >= service.maxRowCount {
			return &pipelineStages, nil
		}
	}

	return &pipelineStages, nil
}
