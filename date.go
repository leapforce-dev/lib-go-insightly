package insightly

import (
	"time"
)

// types
//
type Date struct {
	DATE_ID            int        `json:"DATE_ID"`
	OCCASION_NAME      string     `json:"OCCASION_NAME"`
	OCCASION_DATE      string     `json:"OCCASION_DATE"`
	REPEAT_YEARLY      bool       `json:"REPEAT_YEARLY"`
	CREATE_TASK_YEARLY bool       `json:"CREATE_TASK_YEARLY"`
	OccasionDate       *time.Time `json:"OccasionDate"`
}

func (d *Date) ParseDates() {
	// parse OCCASION_DATE to time.Time
	if d.OCCASION_DATE != "" {
		t, _ := time.Parse("2006-01-02 15:04:05 +0000 UTC", d.OCCASION_DATE+" +0000 UTC")
		//errortools.Fatal(err)
		d.OccasionDate = &t
	}
}
