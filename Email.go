package insightly

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
	i_types "github.com/leapforce-libraries/go_insightly/types"
	go_types "github.com/leapforce-libraries/go_types"
)

// Email stores Email from Service
//
type Email struct {
	EmailID           int64                   `json:"EMAIL_ID"`
	EmailFrom         string                  `json:"EMAIL_FROM"`
	EmailTo           string                  `json:"EMAIL_TO"`
	EmailCC           *go_types.String        `json:"EMAIL_CC"`
	Subject           string                  `json:"SUBJECT"`
	Body              string                  `json:"BODY"`
	EmailDateUTC      i_types.DateTimeString  `json:"EMAIL_DATE_UTC"`
	Format            string                  `json:"FORMAT"`
	Size              int64                   `json:"SIZE"`
	OwnerUserID       int64                   `json:"OWNER_USER_ID"`
	DateCreatedUTC    i_types.DateTimeString  `json:"DATE_CREATED_UTC"`
	QueuedSendDateUTC *i_types.DateTimeString `json:"QUEUED_SEND_DATE_UTC"`
	CreatedUserID     int64                   `json:"CREATED_USER_ID"`
	Tags              *[]Tag                  `json:"TAGS"`
	Links             *[]Link                 `json:"LINKS"`
}

type GetEmailsConfig struct {
	Skip         *uint64
	Top          *uint64
	Brief        *bool
	CountTotal   *bool
	UpdatedAfter *time.Time
	FieldFilter  *FieldFilter
}

// GetEmails returns all emails
//

func (service *Service) GetEmails(config *GetEmailsConfig) (*[]Email, *errortools.Error) {
	params := url.Values{}

	endpoint := "Emails"
	emails := []Email{}
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

		emailsBatch := []Email{}

		requestConfig := go_http.RequestConfig{
			Method:        http.MethodGet,
			Url:           service.url(fmt.Sprintf("%s?%s", endpoint, params.Encode())),
			ResponseModel: &emailsBatch,
		}

		_, _, e := service.httpRequest(&requestConfig)
		if e != nil {
			return nil, e
		}

		emails = append(emails, emailsBatch...)

		if len(emailsBatch) < int(top) {
			delete(service.nextSkips, endpoint)
			break
		}

		service.nextSkips[endpoint] += top
		rowCount += top

		if rowCount >= service.maxRowCount {
			return &emails, nil
		}
	}

	return &emails, nil
}
