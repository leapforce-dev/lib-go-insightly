package insightly

import (
	"fmt"
	"strconv"

	errortools "github.com/leapforce-libraries/go_errortools"
)

// Pipeline stores Pipeline from Insightly
//
type Pipeline struct {
	PipelineID       int    `json:"PIPELINE_ID"`
	PipelineName     string `json:"PIPELINE_NAME"`
	ForOpportunities bool   `json:"FOR_OPPORTUNITIES"`
	ForProjects      bool   `json:"FOR_PROJECTS"`
	OwnerUserID      int    `json:"OWNER_USER_ID"`
}

// GetPipelines returns all pipelines
//
func (i *Insightly) GetPipelines() ([]Pipeline, *errortools.Error) {
	return i.GetPipelinesInternal()
}

// GetPipelinesInternal is the generic function retrieving pipelines from Insightly
//
func (i *Insightly) GetPipelinesInternal() ([]Pipeline, *errortools.Error) {
	endpointStr := "Pipelines?skip=%s&top=%s"
	skip := 0
	top := 500
	rowCount := top

	pipelines := []Pipeline{}

	for rowCount >= top {
		endpoint := fmt.Sprintf(endpointStr, strconv.Itoa(skip), strconv.Itoa(top))
		//fmt.Println(endpoint)

		ls := []Pipeline{}

		_, _, e := i.get(endpoint, nil, &ls)
		if e != nil {
			return nil, e
		}

		for _, l := range ls {
			pipelines = append(pipelines, l)
		}

		rowCount = len(ls)
		skip += top
	}

	if len(pipelines) == 0 {
		pipelines = nil
	}

	return pipelines, nil
}
