package insightly

import (
	"fmt"
	"strconv"

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

// GetPipelineStages returns all pipelinestages
//
func (i *Insightly) GetPipelineStages() ([]PipelineStage, *errortools.Error) {
	return i.GetPipelineStagesInternal()
}

// GetPipelineStagesInternal is the generic function retrieving pipelinestages from Insightly
//
func (i *Insightly) GetPipelineStagesInternal() ([]PipelineStage, *errortools.Error) {
	endpointStr := "PipelineStages?skip=%s&top=%s"
	skip := 0
	top := 500
	rowCount := top

	pipelineStages := []PipelineStage{}

	for rowCount >= top {
		endpoint := fmt.Sprintf(endpointStr, strconv.Itoa(skip), strconv.Itoa(top))
		//fmt.Println(endpoint)

		oc := []PipelineStage{}

		_, _, e := i.get(endpoint, nil, &oc)
		if e != nil {
			return nil, e
		}

		for _, o := range oc {
			pipelineStages = append(pipelineStages, o)
		}

		rowCount = len(oc)
		skip += top
	}

	if len(pipelineStages) == 0 {
		pipelineStages = nil
	}

	return pipelineStages, nil
}
