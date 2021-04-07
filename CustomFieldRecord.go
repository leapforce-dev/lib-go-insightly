package insightly

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"

	errortools "github.com/leapforce-libraries/go_errortools"
)

// CustomFieldRecord
//
type CustomFieldRecord struct {
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
func (cf *CustomFieldRecord) unmarshalValue() {
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
func (customFieldRecord *CustomFieldRecord) GetText() *string {
	if customFieldRecord == nil {
		return nil
	}

	customFieldRecord.unmarshalValue()
	return customFieldRecord.fieldValueText
}

func (customFieldRecord *CustomFieldRecord) GetNumeric() *float64 {
	if customFieldRecord == nil {
		return nil
	}

	customFieldRecord.unmarshalValue()
	return customFieldRecord.fieldValueNumeric
}

func (customFieldRecord *CustomFieldRecord) GetBit() *bool {
	if customFieldRecord == nil {
		return nil
	}

	customFieldRecord.unmarshalValue()
	return customFieldRecord.fieldValueBit
}

func (customFieldRecord *CustomFieldRecord) GetTime() *time.Time {
	text := customFieldRecord.GetText()
	if text == nil {
		return nil
	}

	time, err := time.Parse(dateTimeFormatCustomField, *text)
	if err != nil {
		return nil
	}

	return &time
}

func (customFieldRecord *CustomFieldRecord) Get() (*string, *float64, *bool) {
	if customFieldRecord == nil {
		return nil, nil, nil
	}

	customFieldRecord.unmarshalValue()
	return customFieldRecord.fieldValueText, customFieldRecord.fieldValueNumeric, customFieldRecord.fieldValueBit
}

func (customFieldRecord *CustomFieldRecord) SetText(value string) *errortools.Error {
	return customFieldRecord.set(&value, nil, nil)
}

func (customFieldRecord *CustomFieldRecord) SetNumeric(value float64) *errortools.Error {
	return customFieldRecord.set(nil, &value, nil)
}

func (customFieldRecord *CustomFieldRecord) SetNumericInt(value int) *errortools.Error {
	valueFloat := float64(value)
	return customFieldRecord.set(nil, &valueFloat, nil)
}

func (customFieldRecord *CustomFieldRecord) SetNumericInt32(value int32) *errortools.Error {
	valueFloat := float64(value)
	return customFieldRecord.set(nil, &valueFloat, nil)
}

func (customFieldRecord *CustomFieldRecord) SetNumericInt64(value int64) *errortools.Error {
	valueFloat := float64(value)
	return customFieldRecord.set(nil, &valueFloat, nil)
}

func (customFieldRecord *CustomFieldRecord) SetBit(value bool) *errortools.Error {
	return customFieldRecord.set(nil, nil, &value)
}

// set //
//
func (customFieldRecord *CustomFieldRecord) set(valueText *string, valueNumeric *float64, valueBit *bool) *errortools.Error {
	if customFieldRecord == nil {
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

	(*customFieldRecord).FieldValue = b
	(*customFieldRecord).fieldValueText = valueText
	(*customFieldRecord).fieldValueNumeric = valueNumeric
	(*customFieldRecord).fieldValueBit = valueBit

	return nil
}

// CustomFields
//
type CustomFields []CustomFieldRecord

// contains
func (customFields *CustomFields) Contains(fieldName string, fieldValue interface{}) bool {
	customFieldRecord := customFields.get(fieldName)

	if customFieldRecord == nil {
		return false
	}

	b := []byte{}
	b, _ = json.Marshal(fieldValue)

	if string(customFieldRecord.FieldValue) == string(b) {
		return true
	}

	return false
}

// get
//
func (customFields *CustomFields) get(fieldName string) *CustomFieldRecord {
	if customFields == nil {
		return nil
	}

	for _, customFieldRecord := range *customFields {
		if strings.ToLower(customFieldRecord.FieldName) == strings.ToLower(fieldName) {
			//customFieldRecord.unmarshalValue()
			return &customFieldRecord
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

	for i, customFieldRecord := range *customFields {
		if strings.ToLower(customFieldRecord.FieldName) == strings.ToLower(fieldName) {
			customFieldRecord.set(valueText, valueNumeric, valueBit)
			(*customFields)[i] = customFieldRecord
			return nil
		}
	}

	customFieldNew := CustomFieldRecord{
		FieldName: fieldName,
	}
	customFieldNew.set(valueText, valueNumeric, valueBit)

	*customFields = append(*customFields, customFieldNew)

	return nil
}
