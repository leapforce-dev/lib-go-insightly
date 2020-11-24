package insightly

import (
	"fmt"
	"strconv"

	errortools "github.com/leapforce-libraries/go_errortools"
)

// PipelineStage stores PipelineStage from Insightly
//
type PipelineStage struct {
	STAGE_ID       int    `json:"STAGE_ID"`
	PIPELINE_ID    int    `json:"PIPELINE_ID"`
	STAGE_NAME     string `json:"STAGE_NAME"`
	STAGE_ORDER    int    `json:"STAGE_ORDER"`
	ACTIVITYSET_ID int    `json:"ACTIVITYSET_ID,omitempty"`
	OWNER_USER_ID  int    `json:"OWNER_USER_ID"`
}

// GetPipelineStages returns all pipelinestages
//
func (i *Insightly) GetPipelineStages() ([]PipelineStage, *errortools.Error) {
	return i.GetPipelineStagesInternal()
}

// GetPipelineStagesInternal is the generic function retrieving pipelinestages from Insightly
//
func (i *Insightly) GetPipelineStagesInternal() ([]PipelineStage, *errortools.Error) {
	urlStr := "%sPipelineStages?skip=%s&top=%s"
	skip := 0
	top := 500
	rowCount := top

	pipelineStages := []PipelineStage{}

	for rowCount >= top {
		url := fmt.Sprintf(urlStr, i.apiURL, strconv.Itoa(skip), strconv.Itoa(top))
		//fmt.Println(url)

		oc := []PipelineStage{}

		e := i.Get(url, &oc)
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
