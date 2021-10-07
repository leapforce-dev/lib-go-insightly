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

type User struct {
	UserID                 int64                  `json:"USER_ID"`
	ContactID              int64                  `json:"CONTACT_ID"`
	FirstName              string                 `json:"FIRST_NAME"`
	LastName               string                 `json:"LAST_NAME"`
	TimezoneID             string                 `json:"TIMEZONE_ID"`
	EmailAddress           string                 `json:"EMAIL_ADDRESS"`
	EmailDropboxIdentifier string                 `json:"EMAIL_DROPBOX_IDENTIFIER"`
	EmailDropboxAddress    string                 `json:"EMAIL_DROPBOX_ADDRESS"`
	Administrator          bool                   `json:"ADMINISTRATOR"`
	AccountOwner           bool                   `json:"ACCOUNT_OWNER"`
	Active                 bool                   `json:"ACTIVE"`
	DateCreatedUTC         i_types.DateTimeString `json:"DATE_CREATED_UTC"`
	DateUpdatedUTC         i_types.DateTimeString `json:"DATE_UPDATED_UTC"`
	UserCurrency           string                 `json:"USER_CURRENCY"`
	ContactDisplay         string                 `json:"CONTACT_DISPLAY"`
	ContactOrder           string                 `json:"CONTACT_ORDER"`
	TaskWeekStart          int64                  `json:"TASK_WEEK_START"`
	InstanceID             int64                  `json:"INSTANCE_ID"`
	ProfileID              *int64                 `json:"PROFILE_ID"`
	RoleID                 *int64                 `json:"ROLE_ID"`
}

// GetUser returns a specific user
//
func (service *Service) GetUser(userID int64) (*User, *errortools.Error) {
	user := User{}

	requestConfig := go_http.RequestConfig{
		Method:        http.MethodGet,
		URL:           service.url(fmt.Sprintf("Users/%v", userID)),
		ResponseModel: &user,
	}
	_, _, e := service.httpRequest(&requestConfig)
	if e != nil {
		return nil, e
	}

	return &user, nil
}

type GetUsersConfig struct {
	Skip         *uint64
	Top          *uint64
	Brief        *bool
	CountTotal   *bool
	UpdatedAfter *time.Time
	FieldFilter  *FieldFilter
}

// GetUsers returns all users
//
func (service *Service) GetUsers(config *GetUsersConfig) (*[]User, *errortools.Error) {
	params := url.Values{}

	endpoint := "Users"
	users := []User{}
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
		usersBatch := []User{}

		requestConfig := go_http.RequestConfig{
			Method:        http.MethodGet,
			URL:           service.url(fmt.Sprintf("%s?%s", endpoint, params.Encode())),
			ResponseModel: &usersBatch,
		}
		_, _, e := service.httpRequest(&requestConfig)
		if e != nil {
			return nil, e
		}

		users = append(users, usersBatch...)

		if len(usersBatch) < int(top) {
			delete(service.nextSkips, endpoint)
			break
		}

		service.nextSkips[endpoint] += top
		rowCount += top

		if rowCount >= service.maxRowCount {
			return &users, nil
		}
	}

	return &users, nil
}

func (u *User) FullName() string {
	if u == nil {
		return ""
	}

	return fmt.Sprintf("%s %s", u.FirstName, u.LastName)
}
