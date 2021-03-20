package insightly

import (
	"fmt"

	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
)

// CustomField stores CustomField from Service
//
type CustomField struct {
	FieldName          string                  `json:"FIELD_NAME"`
	FieldOrder         int64                   `json:"FIELD_ORDER"`
	FieldFor           string                  `json:"FIELD_FOR"`
	FieldLabel         string                  `json:"FIELD_LABEL"`
	FieldType          string                  `json:"FIELD_TYPE"`
	FieldHelpText      *string                 `json:"FIELD_HELP_TEXT"`
	DefaultValue       *string                 `json:"DEFAULT_VALUE"`
	Editable           bool                    `json:"EDITABLE"`
	Visible            bool                    `json:"VISIBLE"`
	CustomFieldOptions []CustomFieldOption     `json:"CUSTOM_FIELD_OPTIONS"`
	Dependency         []CustomFieldDependency `json:"DEPENDENCY"`
	JoinObject         *string                 `json:"JOIN_OBJECT"`
}

type CustomFieldOption struct {
	OptionID      int64  `json:"OPTION_ID"`
	OptionValue   string `json:"OPTION_VALUE"`
	OptionDefault bool   `json:"OPTION_DEFAULT"`
}

type CustomFieldDependency struct {
	ControllingFieldID string                     `json:"CONTROLLING_FIELD_ID"`
	OptionsFilters     []CustomFieldOptionsFilter `json:"OPTIONS_FILTERS"`
}

type CustomFieldOptionsFilter struct {
	ControllingValue string  `json:"CONTROLLING_VALUE"`
	OptionIDs        []int64 `json:"OPTION_IDS"`
}

// GetCustomFields returns all customobjects
//
func (service *Service) GetCustomFields(objectName string) (*[]CustomField, *errortools.Error) {
	customFields := []CustomField{}

	_customFields := []CustomField{}

	requestConfig := go_http.RequestConfig{
		URL:           service.url(fmt.Sprintf("CustomFields/%s", objectName)),
		ResponseModel: &_customFields,
	}
	_, _, e := service.get(&requestConfig)
	if e != nil {
		return nil, e
	}

	for _, customField := range _customFields {
		customFields = append(customFields, customField)
	}

	if len(customFields) == 0 {
		customFields = nil
	}

	return &customFields, nil
}
