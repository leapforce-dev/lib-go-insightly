package insightly

import (
	"fmt"
	"strconv"
	"strings"

	errortools "github.com/leapforce-libraries/go_errortools"
)

// PipelineStage stores PipelineStage from Insightly
//
type PipelineStage struct {
	StageID       int    `json:"STAGE_ID"`
	PipelineID    int    `json:"PIPELINE_ID"`
	StageName     string `json:"STAGE_NAME"`
	StageOrder    int    `json:"STAGE_ORDER"`
	ActivitySetID int    `json:"ACTIVITYSET_ID,omitempty"`
	OwnerUserID   int    `json:"OWNER_USER_ID"`
}

type GetPipelineStagesFilter struct {
}

// GetPipelineStages returns all pipelineStages
//
func (i *Insightly) GetPipelineStages(filter *GetPipelineStagesFilter) (*[]PipelineStage, *errortools.Error) {
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
		endpoint := fmt.Sprintf(endpointStr, searchString, strconv.Itoa(skip), strconv.Itoa(top))
		//fmt.Println(endpoint)

		cs := []PipelineStage{}

		_, _, e := i.get(endpoint, nil, &cs)
		if e != nil {
			return nil, e
		}

		pipelineStages = append(pipelineStages, cs...)

		rowCount = len(cs)
		//rowCount = 0
		skip += top
	}

	if len(pipelineStages) == 0 {
		pipelineStages = nil
	}

	return &pipelineStages, nil
}
