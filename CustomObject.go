package insightly

import (
	"net/http"

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

	requestConfig := go_http.RequestConfig{
		Method:        http.MethodGet,
		Url:           service.url("CustomObjects"),
		ResponseModel: &customObjects,
	}
	_, _, e := service.httpRequest(&requestConfig)
	if e != nil {
		return nil, e
	}

	return &customObjects, nil
}
