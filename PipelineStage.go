package insightly

import (
	"fmt"
	"strconv"
	"strings"

	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
)

// PipelineStage stores PipelineStage from Service
//
type PipelineStage struct {
	StageID       int    `json:"STAGE_ID"`
	PipelineID    int    `json:"PIPELINE_ID"`
	StageName     string `json:"STAGE_NAME"`
	StageOrder    int    `json:"STAGE_ORDER"`
	ActivitySetID *int   `json:"ACTIVITYSET_ID"`
	OwnerUserID   *int   `json:"OWNER_USER_ID"`
}

type GetPipelineStagesFilter struct {
}

// GetPipelineStages returns all pipelineStages
//
func (service *Service) GetPipelineStages(filter *GetPipelineStagesFilter) (*[]PipelineStage, *errortools.Error) {
	searchString := "?"
	searchFilter := []string{}

	if filter != nil {
	}

	if len(searchFilter) > 0 {
		searchString = "/Search?" + strings.Join(searchFilter, "&")
	}

	endpointStr := "PipelineStages%sskip=%s&top=%s"
	skip := 0
	top := 100
	rowCount := top

	pipelineStages := []PipelineStage{}

	for rowCount >= top {
		_pipelineStages := []PipelineStage{}

		requestConfig := go_http.RequestConfig{
			URL:           service.url(fmt.Sprintf(endpointStr, searchString, strconv.Itoa(skip), strconv.Itoa(top))),
			ResponseModel: &_pipelineStages,
		}
		_, _, e := service.get(&requestConfig)
		if e != nil {
			return nil, e
		}

		pipelineStages = append(pipelineStages, _pipelineStages...)

		rowCount = len(_pipelineStages)
		//rowCount = 0
		skip += top
	}

	if len(pipelineStages) == 0 {
		pipelineStages = nil
	}

	return &pipelineStages, nil
}
