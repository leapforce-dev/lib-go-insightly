package insightly

import (
	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
	i_types "github.com/leapforce-libraries/go_insightly/types"
)

// CustomObject stores CustomObject from Service
//
type CustomObject struct {
	ObjectName              string                 `json:"OBJECT_NAME"`
	SingularLabel           string                 `json:"SINGULAR_LABEL"`
	PluralLabel             string                 `json:"PLURAL_LABEL"`
	Description             *string                `json:"DESCRIPTION"`
	RecordNameLabel         string                 `json:"RECORD_NAME_LABEL"`
	RecordNameType          string                 `json:"RECORD_NAME_TYPE"`
	RecordNameDisplayFormat *string                `json:"RECORD_NAME_DISPLAY_FORMAT"`
	EnableNavbar            bool                   `json:"ENABLE_NAVBAR"`
	EnableWorkflows         bool                   `json:"ENABLE_WORKFLOWS"`
	CreatedUserID           int64                  `json:"CREATED_USER_ID"`
	DateCreatedUTC          i_types.DateTimeString `json:"DATE_CREATED_UTC"`
	DateUpdatedUTC          i_types.DateTimeString `json:"DATE_UPDATED_UTC"`
}

// GetCustomObjects returns all customobjects
//
func (service *Service) GetCustomObjects() (*[]CustomObject, *errortools.Error) {
	customObjects := []CustomObject{}

	_customObjects := []CustomObject{}

	requestConfig := go_http.RequestConfig{
		URL:           service.url("CustomObjects"),
		ResponseModel: &_customObjects,
	}
	_, _, e := service.get(&requestConfig)
	if e != nil {
		return nil, e
	}

	for _, customObject := range _customObjects {
		customObjects = append(customObjects, customObject)
	}

	if len(customObjects) == 0 {
		customObjects = nil
	}

	return &customObjects, nil
}
