package insightly

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"

	errortools "github.com/leapforce-libraries/go_errortools"
)

// CustomField
//
type CustomField struct {
	FieldName         string          `json:"FIELD_NAME"`
	FieldValue        json.RawMessage `json:"FIELD_VALUE"`
	CustomFieldID     string          `json:"CUSTOM_FIELD_ID"`
	fieldValueText    *string
	fieldValueNumeric *float64
	fieldValueBit     *bool
	unmarshalled      bool
}

// unmarshalValue
//
func (cf *CustomField) unmarshalValue() {
	if cf == nil {
		return
	}

	if cf.unmarshalled {
		return
	}

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

	cf.unmarshalled = true
}

// get
//
func (customField *CustomField) GetText() *string {
	if customField == nil {
		return nil
	}

	customField.unmarshalValue()
	return customField.fieldValueText
}

func (customField *CustomField) GetNumeric() *float64 {
	if customField == nil {
		return nil
	}

	customField.unmarshalValue()
	return customField.fieldValueNumeric
}

func (customField *CustomField) GetBit() *bool {
	if customField == nil {
		return nil
	}

	customField.unmarshalValue()
	return customField.fieldValueBit
}

func (customField *CustomField) GetTime() *time.Time {
	text := customField.GetText()
	if text == nil {
		return nil
	}

	time, err := time.Parse(DateFormat, *text)
	if err != nil {
		return nil
	}

	return &time
}

func (customField *CustomField) SetText(value string) *errortools.Error {
	return customField.set(&value, nil, nil)
}

func (customField *CustomField) SetNumeric(value float64) *errortools.Error {
	return customField.set(nil, &value, nil)
}

func (customField *CustomField) SetNumericInt(value int) *errortools.Error {
	valueFloat := float64(value)
	return customField.set(nil, &valueFloat, nil)
}

func (customField *CustomField) SetNumericInt32(value int32) *errortools.Error {
	valueFloat := float64(value)
	return customField.set(nil, &valueFloat, nil)
}

func (customField *CustomField) SetNumericInt64(value int64) *errortools.Error {
	valueFloat := float64(value)
	return customField.set(nil, &valueFloat, nil)
}

func (customField *CustomField) SetBit(value bool) *errortools.Error {
	return customField.set(nil, nil, &value)
}

// set //
//
func (customField *CustomField) set(valueText *string, valueNumeric *float64, valueBit *bool) *errortools.Error {
	if customField == nil {
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

	(*customField).FieldValue = b
	(*customField).fieldValueText = valueText
	(*customField).fieldValueNumeric = valueNumeric
	(*customField).fieldValueBit = valueBit

	return nil
}

// CustomFields
//
type CustomFields []CustomField

// contains
func (customFields *CustomFields) Contains(fieldName string, fieldValue interface{}) bool {
	customField := customFields.get(fieldName)

	if customField == nil {
		return false
	}

	b := []byte{}
	b, _ = json.Marshal(fieldValue)

	if string(customField.FieldValue) == string(b) {
		return true
	}

	return false
}

// get
//
func (customFields *CustomFields) get(fieldName string) *CustomField {
	if customFields == nil {
		return nil
	}

	for _, customField := range *customFields {
		if strings.ToLower(customField.FieldName) == strings.ToLower(fieldName) {
			//customField.unmarshalValue()
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
		return cf.GetText()
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
		return cf.GetNumeric()
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
		return cf.GetBit()
	}
}

func (customFields *CustomFields) GetTime(fieldName string) *time.Time {
	if customFields == nil {
		return nil
	}

	cf := customFields.get(fieldName)

	if cf == nil {
		return nil
	} else {
		return cf.GetTime()
	}
}

func (customFields *CustomFields) SetText(fieldName string, value string) *errortools.Error {
	return customFields.set(fieldName, &value, nil, nil)
}

func (customFields *CustomFields) SetNumeric(fieldName string, value float64) *errortools.Error {
	return customFields.set(fieldName, nil, &value, nil)
}

func (customFields *CustomFields) SetNumericInt(fieldName string, value int) *errortools.Error {
	valueFloat := float64(value)
	return customFields.set(fieldName, nil, &valueFloat, nil)
}

func (customFields *CustomFields) SetNumericInt32(fieldName string, value int32) *errortools.Error {
	valueFloat := float64(value)
	return customFields.set(fieldName, nil, &valueFloat, nil)
}

func (customFields *CustomFields) SetNumericInt64(fieldName string, value int64) *errortools.Error {
	valueFloat := float64(value)
	return customFields.set(fieldName, nil, &valueFloat, nil)
}

func (customFields *CustomFields) SetBit(fieldName string, value bool) *errortools.Error {
	return customFields.set(fieldName, nil, nil, &value)
}

func (customFields *CustomFields) Delete(fieldName string) *errortools.Error {
	return customFields.set(fieldName, nil, nil, nil)
}

// set
//
func (customFields *CustomFields) set(fieldName string, valueText *string, valueNumeric *float64, valueBit *bool) *errortools.Error {
	if customFields == nil {
		return nil
	}

	for i, customField := range *customFields {
		if strings.ToLower(customField.FieldName) == strings.ToLower(fieldName) {
			customField.set(valueText, valueNumeric, valueBit)
			(*customFields)[i] = customField
			return nil
		}
	}

	customFieldNew := CustomField{
		FieldName: fieldName,
	}
	customFieldNew.set(valueText, valueNumeric, valueBit)

	*customFields = append(*customFields, customFieldNew)

	return nil
}
