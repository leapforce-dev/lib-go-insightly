package insightly

import "fmt"

type FieldFilter struct {
	FieldName  string
	FieldValue string
}

func (f *FieldFilter) Search() string {
	if f == nil {
		return "?"
	}
	return fmt.Sprintf("/Search?field_name=%s&field_value=%s&", f.FieldName, f.FieldValue)
}
