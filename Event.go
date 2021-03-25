package insightly

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
	i_types "github.com/leapforce-libraries/go_insightly/types"
)

// Event stores Event from Service
//
type Event struct {
	EventID         int64                   `json:"EVENT_ID"`
	Title           string                  `json:"TITLE"`
	Location        string                  `json:"LOCATION"`
	LastName        string                  `json:"LAST_NAME"`
	StartDateUTC    i_types.DateTimeString  `json:"START_DATE_UTC"`
	EndDateUTC      i_types.DateTimeString  `json:"END_DATE_UTC"`
	AllDay          bool                    `json:"ALL_DAY"`
	Details         string                  `json:"DETAILS"`
	DateCreateUTC   i_types.DateTimeString  `json:"DATE_CREATED_UTC"`
	DateUpdatedUTC  i_types.DateTimeString  `json:"DATE_UPDATED_UTC"`
	ReminderDateUTC *i_types.DateTimeString `json:"REMINDER_DATE_UTC"`
	ReminderSent    bool                    `json:"REMINDER_SENT"`
	OwnerUserID     int64                   `json:"OWNER_USER_ID"`
	CustomFields    *CustomFields           `json:"CUSTOMFIELDS"`
	Links           *[]Link                 `json:"LINKS"`
}

// GetEvent returns a specific event
//
func (service *Service) GetEvent(eventID int64) (*Event, *errortools.Error) {
	event := Event{}

	requestConfig := go_http.RequestConfig{
		URL:           service.url(fmt.Sprintf("Events/%v", eventID)),
		ResponseModel: &event,
	}
	_, _, e := service.get(&requestConfig)
	if e != nil {
		return nil, e
	}

	return &event, nil
}

type GetEventsConfig struct {
	UpdatedAfter *time.Time
	FieldFilter  *FieldFilter
}

// GetEvents returns all events
//
func (service *Service) GetEvents(config *GetEventsConfig) (*[]Event, *errortools.Error) {
	searchString := "?"
	searchFilter := []string{}

	if config != nil {
		if config.UpdatedAfter != nil {
			from := config.UpdatedAfter.Format(DateTimeFormat)
			searchFilter = append(searchFilter, fmt.Sprintf("updated_after_utc=%s&", from))
		}

		if config.FieldFilter != nil {
			searchFilter = append(searchFilter, fmt.Sprintf("field_name=%s&field_value=%s&", config.FieldFilter.FieldName, config.FieldFilter.FieldValue))
		}
	}

	if len(searchFilter) > 0 {
		searchString = "/Search?" + strings.Join(searchFilter, "&")
	}

	endpointStr := "Events%sskip=%s&top=%s"
	skip := 0
	top := 100
	rowCount := top

	events := []Event{}

	for rowCount >= top {
		_events := []Event{}

		requestConfig := go_http.RequestConfig{
			URL:           service.url(fmt.Sprintf(endpointStr, searchString, strconv.Itoa(skip), strconv.Itoa(top))),
			ResponseModel: &_events,
		}
		_, _, e := service.get(&requestConfig)
		if e != nil {
			return nil, e
		}

		events = append(events, _events...)

		rowCount = len(_events)
		//rowCount = 0
		skip += top
	}

	if len(events) == 0 {
		events = nil
	}

	return &events, nil
}
