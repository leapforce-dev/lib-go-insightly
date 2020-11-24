package insightly

import (
	"fmt"
	"strconv"

	errortools "github.com/leapforce-libraries/go_errortools"
)

// Pipeline stores Pipeline from Insightly
//
type Pipeline struct {
	PIPELINE_ID       int    `json:"PIPELINE_ID"`
	PIPELINE_NAME     string `json:"PIPELINE_NAME"`
	FOR_OPPORTUNITIES bool   `json:"FOR_OPPORTUNITIES"`
	FOR_PROJECTS      bool   `json:"FOR_PROJECTS"`
	OWNER_USER_ID     int    `json:"OWNER_USER_ID"`
}

// GetPipelines returns all pipelines
//
func (i *Insightly) GetPipelines() ([]Pipeline, *errortools.Error) {
	return i.GetPipelinesInternal()
}

// GetPipelinesInternal is the generic function retrieving pipelines from Insightly
//
func (i *Insightly) GetPipelinesInternal() ([]Pipeline, *errortools.Error) {
	urlStr := "%sPipelines?skip=%s&top=%s"
	skip := 0
	top := 500
	rowCount := top

	pipelines := []Pipeline{}

	for rowCount >= top {
		url := fmt.Sprintf(urlStr, i.apiURL, strconv.Itoa(skip), strconv.Itoa(top))
		//fmt.Println(url)

		ls := []Pipeline{}

		e := i.Get(url, &ls)
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
