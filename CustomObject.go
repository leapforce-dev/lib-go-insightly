package insightly

import (
	errortools "github.com/leapforce-libraries/go_errortools"
)

// CustomObject stores CustomObject from Service
//
type CustomObject struct {
	ObjectName              string  `json:"OBJECT_NAME"`
	SingularLabel           string  `json:"SINGULAR_LABEL"`
	PluralLabel             string  `json:"PLURAL_LABEL"`
	Description             string  `json:"DESCRIPTION"`
	RecordNameLabel         string  `json:"RECORD_NAME_LABEL"`
	RecordNameType          string  `json:"RECORD_NAME_TYPE"`
	RecordNameDisplayFormat string  `json:"RECORD_NAME_DISPLAY_FORMAT"`
	EnableNavbar            bool    `json:"ENABLE_NAVBAR"`
	EnableWorkflows         bool    `json:"ENABLE_WORKFLOWS"`
	CreatedUserID           *int    `json:"CREATED_USER_ID"`
	DateCreatedUTC          DateUTC `json:"DATE_CREATED_UTC"`
	DateUpdatedUTC          DateUTC `json:"DATE_UPDATED_UTC"`
}

// GetCustomObjects returns all customobjects
//
func (service *Service) GetCustomObjects() ([]CustomObject, *errortools.Error) {
	customobjects := []CustomObject{}

	endpoint := "CustomObjects"

	os := []CustomObject{}

	_, _, err := service.get(endpoint, nil, &os)
	if err != nil {
		return nil, err
	}

	for _, o := range os {
		customobjects = append(customobjects, o)
	}

	if len(customobjects) == 0 {
		customobjects = nil
	}

	return customobjects, nil
}
