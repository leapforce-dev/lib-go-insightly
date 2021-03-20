package insightly

import (
	i_types "github.com/leapforce-libraries/go_insightly/types"
)

type Date struct {
	DateID           int                    `json:"DATE_ID"`
	OccasionName     string                 `json:"OCCASION_NAME"`
	OccasionDate     i_types.DateTimeString `json:"OCCASION_DATE"`
	RepeatYearly     bool                   `json:"REPEAT_YEARLY"`
	CreateTaskYearly bool                   `json:"CREATE_TASK_YEARLY"`
}
