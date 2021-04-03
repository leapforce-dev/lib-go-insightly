package insightly

import (
	"fmt"
	"net/url"

	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
)

// Tag stores Tag from Service
//
type Tag struct {
	TagName string `json:"TAG_NAME"`
}

type GetTagsConfig struct {
	Skip       *uint64
	Top        *uint64
	CountTotal *bool
	RecordType RecordType
}

// GetTags returns all tags
//
func (service *Service) GetTags(config *GetTagsConfig) (*[]Tag, *errortools.Error) {
	if config == nil {
		return nil, nil
	}

	params := url.Values{}

	endpoint := "Tags"
	tags := []Tag{}
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
		params.Set("record_type", string(config.RecordType))
	}

	params.Set("top", fmt.Sprintf("%v", top))

	for true {
		params.Set("skip", fmt.Sprintf("%v", service.nextSkips[endpoint]))

		tagsBatch := []Tag{}

		requestConfig := go_http.RequestConfig{
			URL:           service.url(fmt.Sprintf("%s?%s", endpoint, params.Encode())),
			ResponseModel: &tagsBatch,
		}

		_, _, e := service.get(&requestConfig)
		if e != nil {
			return nil, e
		}

		tags = append(tags, tagsBatch...)

		if len(tagsBatch) < int(top) {
			delete(service.nextSkips, endpoint)
			break
		}

		service.nextSkips[endpoint] += top
		rowCount += top

		if rowCount >= service.maxRowCount {
			return &tags, nil
		}
	}

	return &tags, nil
}
