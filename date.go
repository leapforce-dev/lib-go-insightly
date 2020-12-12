package insightly

import (
	"time"
)

// types
//
type Date struct {
	DateID           int        `json:"DATE_ID"`
	OccasionName     string     `json:"OCCASION_NAME"`
	OccasionDate     string     `json:"OCCASION_DATE"`
	RepeatYearly     bool       `json:"REPEAT_YEARLY"`
	CreateTaskYearly bool       `json:"CREATE_TASK_YEARLY"`
	OccasionDateT    *time.Time `json:"OccasionDate"`
}

func (d *Date) parseDates() {
	// parse OCCASION_DATE to time.Time
	if d.OccasionDate != "" {
		t, _ := time.Parse("2006-01-02 15:04:05 +0000 UTC", d.OccasionDate+" +0000 UTC")
		//errortools.Fatal(err)
		d.OccasionDateT = &t
	}
}
