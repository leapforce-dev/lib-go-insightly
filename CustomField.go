package insightly

import (
	"encoding/json"
	"strconv"
	"strings"

	errortools "github.com/leapforce-libraries/go_errortools"
)

// types
//
type CustomField struct {
	FieldName         string          `json:"FIELD_NAME"`
	FieldValue        json.RawMessage `json:"FIELD_VALUE"`
	CustomFieldID     string          `json:"CUSTOM_FIELD_ID"`
	fieldValueText    *string
	fieldValueNumeric *float64
	fieldValueBit     *bool
}

type CustomFields []CustomField

// get //
//
func (customFields *CustomFields) get(fieldName string) *CustomField {
	if customFields == nil {
		return nil
	}

	for _, customField := range *customFields {
		if strings.ToLower(customField.FieldName) == strings.ToLower(fieldName) {
			customField.unmarshalValue()
			return &customField
		}
	}

	return nil
}

func (customFields *CustomFields) GetText(fieldName string) *string {
	if customFields == nil {
		return nil
	}

	cf := customFields.get(fieldName)

	if cf == nil {
		return nil
	} else {
		return cf.fieldValueText
	}
}

func (customFields *CustomFields) GetNumeric(fieldName string) *float64 {
	if customFields == nil {
		return nil
	}

	cf := customFields.get(fieldName)

	if cf == nil {
		return nil
	} else {
		return cf.fieldValueNumeric
	}
}

func (customFields *CustomFields) GetBit(fieldName string) *bool {
	if customFields == nil {
		return nil
	}

	cf := customFields.get(fieldName)

	if cf == nil {
		return nil
	} else {
		return cf.fieldValueBit
	}
}

// unmarshalValue //
//
func (cf *CustomField) unmarshalValue() {
	j, _ := json.Marshal(&cf.FieldValue)
	// try unmarshalling to string
	err := json.Unmarshal(cf.FieldValue, &cf.fieldValueText)
	// try unmarshalling to float64
	if err != nil {
		cf.fieldValueNumeric = nil
		err = json.Unmarshal(cf.FieldValue, &cf.fieldValueNumeric)
	}
	// try unmarshalling to bool
	if err != nil {
		cf.fieldValueNumeric = nil
		b, err1 := strconv.ParseBool(string(j))
		if err1 == nil {
			cf.fieldValueBit = &b
		} else {
			cf.fieldValueBit = nil
		}
	}
}

func (customFields *CustomFields) SetText(fieldName string, value string) *errortools.Error {
	return customFields.set(fieldName, &value, nil, nil)
}

func (customFields *CustomFields) SetNumeric(fieldName string, value float64) *errortools.Error {
	return customFields.set(fieldName, nil, &value, nil)
}

func (customFields *CustomFields) SetBit(fieldName string, value bool) *errortools.Error {
	return customFields.set(fieldName, nil, nil, &value)
}

func (customFields *CustomFields) Delete(fieldName string) *errortools.Error {
	return customFields.set(fieldName, nil, nil, nil)
}

// set //
//
func (customFields *CustomFields) set(fieldName string, valueText *string, valueNumeric *float64, valueBit *bool) *errortools.Error {
	if customFields == nil {
		return nil
	}

	b := []byte{}
	if valueText != nil {
		b, _ = json.Marshal(*valueText)
	} else if valueNumeric != nil {
		b, _ = json.Marshal(*valueNumeric)
	} else if valueBit != nil {
		b, _ = json.Marshal(*valueBit)
	} else {
		b = nil
	}

	for i, customField := range *customFields {
		if strings.ToLower(customField.FieldName) == strings.ToLower(fieldName) {
			customField.FieldValue = b
			customField.fieldValueText = valueText
			customField.fieldValueNumeric = valueNumeric
			customField.fieldValueBit = valueBit

			(*customFields)[i] = customField

			return nil
		}
	}

	customFieldNew := CustomField{
		FieldName:         fieldName,
		FieldValue:        b,
		fieldValueText:    valueText,
		fieldValueNumeric: valueNumeric,
		fieldValueBit:     valueBit,
	}

	*customFields = append(*customFields, customFieldNew)

	return nil
}
