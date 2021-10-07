package insightly

import (
	"fmt"
	"net/http"
	"net/url"
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
	DateCreatedUTC  i_types.DateTimeString  `json:"DATE_CREATED_UTC"`
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
		Method:        http.MethodGet,
		URL:           service.url(fmt.Sprintf("Events/%v", eventID)),
		ResponseModel: &event,
	}
	_, _, e := service.httpRequest(&requestConfig)
	if e != nil {
		return nil, e
	}

	return &event, nil
}

type GetEventsConfig struct {
	Skip         *uint64
	Top          *uint64
	Brief        *bool
	CountTotal   *bool
	UpdatedAfter *time.Time
	FieldFilter  *FieldFilter
}

// GetEvents returns all events
//
func (service *Service) GetEvents(config *GetEventsConfig) (*[]Event, *errortools.Error) {
	params := url.Values{}

	endpoint := "Events"
	events := []Event{}
	rowCount := uint64(0)
	top := defaultTop
	isSearch := false

	if config != nil {
		if config.Top != nil {
			top = *config.Top
		}
		if config.Skip != nil {
			service.nextSkips[endpoint] = *config.Skip
		}
		if config.Brief != nil {
			params.Set("brief", fmt.Sprintf("%v", *config.Brief))
		}
		if config.CountTotal != nil {
			params.Set("count_total", fmt.Sprintf("%v", *config.CountTotal))
		}
		if config.UpdatedAfter != nil {
			isSearch = true
			params.Set("updated_after_utc", fmt.Sprintf("%v", config.UpdatedAfter.Format(dateTimeFormat)))
		}
		if config.FieldFilter != nil {
			isSearch = true
			params.Set("field_name", config.FieldFilter.FieldName)
			params.Set("field_value", config.FieldFilter.FieldValue)
		}
	}

	if isSearch {
		endpoint += "/Search"
	}

	params.Set("top", fmt.Sprintf("%v", top))

	for {
		params.Set("skip", fmt.Sprintf("%v", service.nextSkips[endpoint]))
		eventsBatch := []Event{}

		requestConfig := go_http.RequestConfig{
			Method:        http.MethodGet,
			URL:           service.url(fmt.Sprintf("%s?%s", endpoint, params.Encode())),
			ResponseModel: &eventsBatch,
		}
		_, _, e := service.httpRequest(&requestConfig)
		if e != nil {
			return nil, e
		}

		events = append(events, eventsBatch...)

		if len(eventsBatch) < int(top) {
			delete(service.nextSkips, endpoint)
			break
		}

		service.nextSkips[endpoint] += top
		rowCount += top

		if rowCount >= service.maxRowCount {
			return &events, nil
		}
	}

	return &events, nil
}
