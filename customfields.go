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
	FieldValueString string          `json:"-"`
	FieldValueInt    *int            `json:"-"`
	FieldValueBool   *bool           `json:"-"`
}

// methods //
//
func (i *Insightly) FindCustomField(cfs []CustomField, fieldName string) *CustomField {
	for _, cf := range cfs {
		if strings.ToLower(cf.FIELD_NAME) == strings.ToLower(fieldName) {
			cf.UnmarshalValue()
			return &cf
		}
	}

	return nil
}

// FindCustomFieldValue //
//
func (i *Insightly) FindCustomFieldValue(cfs []CustomField, fieldName string) string {
	cf := i.FindCustomField(cfs, fieldName)

	if cf == nil {
		return ""
	} else {
		return cf.GetFieldValues()[0]
	}
}

// FindCustomFieldValueBool //
//
func (i *Insightly) FindCustomFieldValueBool(cfs []CustomField, fieldName string) *bool {
	cf := i.FindCustomField(cfs, fieldName)

	if cf == nil {
		return nil
	} else {
		return cf.FieldValueBool
	}
}

// GetFieldValues //
//
func (cf *CustomField) GetFieldValues() []string {
	return strings.Split(string(cf.FieldValueString), ";")
}

// UnmarshalValue //
//
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

// SetCustomField //
//
func (i *Insightly) SetCustomField(cfs []CustomField, fieldName string, valueString string, valueInt *int, valueBool *bool) error {
	for index, cf := range cfs {
		if strings.ToLower(cf.FIELD_NAME) == strings.ToLower(fieldName) {
			b := []byte{}
			cfs[index].FieldValueBool = valueBool
			if valueInt != nil {
				cfs[index].FieldValueString = ""
				cfs[index].FieldValueInt = valueInt
				cfs[index].FieldValueBool = nil
				_b, err := json.Marshal(valueInt)
				if err != nil {
					return err
				}
				b = _b
			} else if valueBool != nil {
				cfs[index].FieldValueString = ""
				cfs[index].FieldValueInt = nil
				cfs[index].FieldValueBool = valueBool
				_b, err := json.Marshal(valueBool)
				if err != nil {
					return err
				}
				b = _b
			} else {
				cfs[index].FieldValueString = valueString
				cfs[index].FieldValueInt = nil
				cfs[index].FieldValueBool = nil
				_b, err := json.Marshal(valueString)
				if err != nil {
					return err
				}
				b = _b
			}

			cfs[index].FIELD_VALUE = b

			//fmt.Println(cfs[index])
		}
	}

	return nil
}
