package insightly

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	errortools "github.com/leapforce-libraries/go_errortools"
)

const (
	timeFormat string = "3:04PM"
)

type TimeString time.Time

func (d *TimeString) UnmarshalJSON(b []byte) error {
	var returnError = func() error {
		errortools.CaptureError(fmt.Sprintf("Cannot parse '%s' to TimeString", string(b)))
		return nil
	}

	var s string

	err := json.Unmarshal(b, &s)
	if err != nil {
		fmt.Println("TimeString", string(b))
		return returnError()
	}

	s = strings.ReplaceAll(s, " ", "")

	if s == "" || s == "0000-00-00 00:00:00" {
		d = nil
		return nil
	}

	_t, err := time.Parse(timeFormat, s)
	if err != nil {
		return returnError()
	}

	*d = TimeString(_t)
	return nil
}

func (d *TimeString) ValuePtr() *time.Time {
	if d == nil {
		return nil
	}

	_d := time.Time(*d)
	return &_d
}

func (d TimeString) Value() time.Time {
	return time.Time(d)
}
