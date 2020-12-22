package insightly

import (
	"fmt"
	"strconv"
	"strings"

	errortools "github.com/leapforce-libraries/go_errortools"
)

// Pipeline stores Pipeline from Insightly
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
func (i *Insightly) GetPipelines(filter *GetPipelinesFilter) (*[]Pipeline, *errortools.Error) {
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
		endpoint := fmt.Sprintf(endpointStr, searchString, strconv.Itoa(skip), strconv.Itoa(top))
		//fmt.Println(endpoint)

		cs := []Pipeline{}

		_, _, e := i.get(endpoint, nil, &cs)
		if e != nil {
			return nil, e
		}

		pipelines = append(pipelines, cs...)

		rowCount = len(cs)
		//rowCount = 0
		skip += top
	}

	if len(pipelines) == 0 {
		pipelines = nil
	}

	return &pipelines, nil
}
