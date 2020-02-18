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
	FieldValueBool   bool
}

type CustomFieldBQ struct {
	FIELD_NAME         string `json:"FIELD_NAME"`
	FIELD_VALUE_BOOL   bool   `json:"FIELD_VALUE_BOOL"`
	FIELD_VALUE_STRING string `json:"FIELD_VALUE_STRING"`
	CUSTOM_FIELD_ID    string `json:"CUSTOM_FIELD_ID"`
}

func (c *CustomField) ToBQ() CustomFieldBQ {
	c.UnmarshalValue()

	return CustomFieldBQ{
		c.FIELD_NAME,
		c.FieldValueBool,
		c.FieldValueString,
		c.CUSTOM_FIELD_ID,
	}
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

func (i *Insightly) FindCustomFieldValueBool(cfs []CustomField, fieldName string) bool {
	cf := i.FindCustomField(cfs, fieldName)

	if cf == nil {
		return false
	} else {
		return cf.FieldValueBool
	}
}

func (cf *CustomField) GetFieldValues() []string {
	return strings.Split(string(cf.FieldValueString), ";")
}

func (cf *CustomField) UnmarshalValue() {
	j, _ := json.Marshal(&cf.FIELD_VALUE)
	json.Unmarshal(cf.FIELD_VALUE, &cf.FieldValueString)
	b, _ := strconv.ParseBool(string(j))
	cf.FieldValueBool = b
}
