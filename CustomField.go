package insightly

import (
	"fmt"
	"net/url"

	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
)

// CustomField stores CustomField from Service
//
type CustomField struct {
	FieldName     string                  `json:"FIELD_NAME"`
	FieldOrder    int64                   `json:"FIELD_ORDER"`
	FieldFor      string                  `json:"FIELD_FOR"`
	FieldLabel    string                  `json:"FIELD_LABEL"`
	FieldType     string                  `json:"FIELD_TYPE"`
	FieldHelpText *string                 `json:"FIELD_HELP_TEXT"`
	DefaultValue  *string                 `json:"DEFAULT_VALUE"`
	Editable      bool                    `json:"EDITABLE"`
	Visible       bool                    `json:"VISIBLE"`
	Options       []CustomFieldOption     `json:"CUSTOM_FIELD_OPTIONS"`
	Dependency    []CustomFieldDependency `json:"DEPENDENCY"`
	JoinObject    *string                 `json:"JOIN_OBJECT"`
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

type GetCustomFieldsConfig struct {
	ObjectName *string
	FieldName  *string
}

// GetCustomFields returns all customobjects
//
func (service *Service) GetCustomFields(config *GetCustomFieldsConfig) (*[]CustomField, *errortools.Error) {
	params := url.Values{}

	customFields := []CustomField{}
	isSearch := false

	objectName := "all"

	if config != nil {
		if config.ObjectName != nil {
			objectName = *config.ObjectName
		}
		if config.FieldName != nil {
			isSearch = true
			params.Set("field_name", *config.FieldName)
		}
	}

	endpoint := fmt.Sprintf("CustomFields/%s", objectName)
	if isSearch {
		endpoint += "/Search"
	}

	requestConfig := go_http.RequestConfig{
		URL:           service.url(fmt.Sprintf("%s?%s", endpoint, params.Encode())),
		ResponseModel: &customFields,
	}
	_, _, e := service.get(&requestConfig)
	if e != nil {
		return nil, e
	}

	return &customFields, nil
}
