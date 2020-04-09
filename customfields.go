package insightly

import (
	"encoding/json"
	"strconv"
	"strings"
)

// types
//
type CustomField struct {
	FIELD_NAME       string          `json:"FIELD_NAME"`
	FIELD_VALUE      json.RawMessage `json:"FIELD_VALUE"`
	CUSTOM_FIELD_ID  string          `json:"CUSTOM_FIELD_ID"`
	FieldValueString string
	FieldValueInt    *int
	FieldValueBool   *bool
}

// methods
//
func (i *Insightly) FindCustomField(cfs []CustomField, fieldName string) *CustomField {
	for _, cf := range cfs {
		if strings.ToLower(cf.FIELD_NAME) == strings.ToLower(fieldName) {
			return &cf
		}
	}

	return nil
}

func (i *Insightly) FindCustomFieldValue(cfs []CustomField, fieldName string) string {
	cf := i.FindCustomField(cfs, fieldName)

	if cf == nil {
		return ""
	} else {
		return cf.GetFieldValues()[0]
	}
}

func (i *Insightly) FindCustomFieldValueBool(cfs []CustomField, fieldName string) *bool {
	cf := i.FindCustomField(cfs, fieldName)

	if cf == nil {
		return nil
	} else {
		return cf.FieldValueBool
	}
}

func (cf *CustomField) GetFieldValues() []string {
	return strings.Split(string(cf.FieldValueString), ";")
}

func (cf *CustomField) UnmarshalValue() {
	j, _ := json.Marshal(&cf.FIELD_VALUE)
	// try unmarshalling to string
	err := json.Unmarshal(cf.FIELD_VALUE, &cf.FieldValueString)
	// try unmarshalling to int
	if err != nil {
		err = json.Unmarshal(cf.FIELD_VALUE, &cf.FieldValueInt)
	}
	// try unmarshalling to bool
	if err != nil {
		b, err1 := strconv.ParseBool(string(j))
		if err1 == nil {
			cf.FieldValueBool = &b
		} else {
			cf.FieldValueBool = nil
		}
	}
}
