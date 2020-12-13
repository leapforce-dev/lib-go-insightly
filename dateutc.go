package insightly

import (
	"fmt"
	"regexp"
	"time"
)

// DateUTC allows for unmarshalling the date objects returned by Exact.
type DateUTC struct {
	time.Time
}

// IsSet returns a boolean if the Date is actually set.
func (d *DateUTC) IsSet() bool {
	return !d.IsZero()
}

// UnmarshalJSON unmarshals the date format returned from the
// Exact Online API.
func (d *DateUTC) UnmarshalJSON(b []byte) error {
	s := string(b)

	if s == "" {
		d.Time = time.Time{}
		return nil
	}

	re := regexp.MustCompile(`[0-9][0-9][0-9][0-9]-[0-9][0-9]-[0-9][0-9] [0-9][0-9]:[0-9][0-9]:[0-9][0-9]`)
	s1 := re.FindString(s)
	if s1 == "" {
		d.Time = time.Time{}
		return nil
	}

	t, err := time.Parse("2006-01-02 15:04:05", s1)
	if err != nil {
		return fmt.Errorf("DateUTC.UnmarshalJSON() error: %v", err)
	}

	d.Time = t
	return nil
}

// MarshalJSON marshals the date to a format expected by the
// Exact Online API.
func (d *DateUTC) MarshalJSON() ([]byte, error) {
	if d.IsZero() {
		return []byte("null"), nil
	}

	return d.Time.MarshalJSON()
}
