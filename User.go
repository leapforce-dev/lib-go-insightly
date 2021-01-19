package insightly

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
)

type User struct {
	UserID                 int     `json:"USER_ID"`
	ContactID              *int    `json:"CONTACT_ID"`
	FirstName              string  `json:"FIRST_NAME"`
	LastName               string  `json:"LAST_NAME"`
	TimezoneID             string  `json:"TIMEZONE_ID"`
	EmailAddress           string  `json:"EMAIL_ADDRESS"`
	EmailDropboxIdentifier string  `json:"EMAIL_DROPBOX_IDENTIFIER"`
	EmailDropboxAddress    string  `json:"EMAIL_DROPBOX_ADDRESS"`
	Administrator          bool    `json:"ADMINISTRATOR"`
	AccountOwner           bool    `json:"ACCOUNT_OWNER"`
	Active                 bool    `json:"ACTIVE"`
	DateCreatedUTC         DateUTC `json:"DATE_CREATED_UTC"`
	DateUpdatedUTC         DateUTC `json:"DATE_UPDATED_UTC"`
	UserCurrency           string  `json:"USER_CURRENCY"`
	ContactDisplay         string  `json:"CONTACT_DISPLAY"`
	ContactOrder           string  `json:"CONTACT_ORDER"`
	TaskWeekStart          int     `json:"TASK_WEEK_START"`
	InstanceID             *int    `json:"INSTANCE_ID"`
	ProfileID              *int    `json:"PROFILE_ID"`
	RoleID                 *int    `json:"ROLE_ID"`
}

func (u *User) prepareMarshal() interface{} {
	if u == nil {
		return nil
	}

	return &struct {
		UserID                 int    `json:"USER_ID"`
		ContactID              *int   `json:"CONTACT_ID"`
		FirstName              string `json:"FIRST_NAME"`
		LastName               string `json:"LAST_NAME"`
		TimezoneID             string `json:"TIMEZONE_ID"`
		EmailAddress           string `json:"EMAIL_ADDRESS"`
		EmailDropboxIdentifier string `json:"EMAIL_DROPBOX_IDENTIFIER"`
		EmailDropboxAddress    string `json:"EMAIL_DROPBOX_ADDRESS"`
		Administrator          bool   `json:"ADMINISTRATOR"`
		AccountOwner           bool   `json:"ACCOUNT_OWNER"`
		Active                 bool   `json:"ACTIVE"`
		UserCurrency           string `json:"USER_CURRENCY"`
		ContactDisplay         string `json:"CONTACT_DISPLAY"`
		ContactOrder           string `json:"CONTACT_ORDER"`
		TaskWeekStart          int    `json:"TASK_WEEK_START"`
		InstanceID             *int   `json:"INSTANCE_ID"`
		ProfileID              *int   `json:"PROFILE_ID"`
		RoleID                 *int   `json:"ROLE_ID"`
	}{
		u.UserID,
		u.ContactID,
		u.FirstName,
		u.LastName,
		u.TimezoneID,
		u.EmailAddress,
		u.EmailDropboxIdentifier,
		u.EmailDropboxAddress,
		u.Administrator,
		u.AccountOwner,
		u.Active,
		u.UserCurrency,
		u.ContactDisplay,
		u.ContactOrder,
		u.TaskWeekStart,
		u.InstanceID,
		u.ProfileID,
		u.RoleID,
	}
}

// GetUser returns a specific user
//
func (service *Service) GetUser(userID int) (*User, *errortools.Error) {
	user := User{}

	requestConfig := go_http.RequestConfig{
		URL:           service.url(fmt.Sprintf("Users/%v", userID)),
		ResponseModel: &user,
	}
	_, _, e := service.get(&requestConfig)
	if e != nil {
		return nil, e
	}

	return &user, nil
}

type GetUsersFilter struct {
	UpdatedAfter *time.Time
	Field        *struct {
		FieldName  string
		FieldValue string
	}
}

// GetUsers returns all users
//
func (service *Service) GetUsers(filter *GetUsersFilter) (*[]User, *errortools.Error) {
	searchString := "?"
	searchFilter := []string{}

	if filter != nil {
		if filter.UpdatedAfter != nil {
			from := filter.UpdatedAfter.Format(time.RFC3339)
			searchFilter = append(searchFilter, fmt.Sprintf("updated_after_utc=%s&", from))
		}

		if filter.Field != nil {
			searchFilter = append(searchFilter, fmt.Sprintf("field_name=%s&field_value=%s&", filter.Field.FieldName, filter.Field.FieldValue))
		}
	}

	if len(searchFilter) > 0 {
		searchString = "/Search?" + strings.Join(searchFilter, "&")
	}

	endpointStr := "Users%sskip=%s&top=%s"
	skip := 0
	top := 100
	rowCount := top

	users := []User{}

	for rowCount >= top {
		_users := []User{}

		requestConfig := go_http.RequestConfig{
			URL:           service.url(fmt.Sprintf(endpointStr, searchString, strconv.Itoa(skip), strconv.Itoa(top))),
			ResponseModel: &_users,
		}
		_, _, e := service.get(&requestConfig)
		if e != nil {
			return nil, e
		}

		users = append(users, _users...)

		rowCount = len(_users)
		//rowCount = 0
		skip += top
	}

	if len(users) == 0 {
		users = nil
	}

	return &users, nil
}

// CreateUser creates a new contract
//
func (service *Service) CreateUser(user *User) (*User, *errortools.Error) {
	if user == nil {
		return nil, nil
	}

	userNew := User{}

	requestConfig := go_http.RequestConfig{
		URL:           service.url("Users"),
		BodyModel:     user.prepareMarshal(),
		ResponseModel: &userNew,
	}
	_, _, e := service.post(&requestConfig)
	if e != nil {
		return nil, e
	}

	return &userNew, nil
}

// UpdateUser updates an existing contract
//
func (service *Service) UpdateUser(user *User) (*User, *errortools.Error) {
	if user == nil {
		return nil, nil
	}

	userUpdated := User{}

	requestConfig := go_http.RequestConfig{
		URL:           service.url("Users"),
		BodyModel:     user.prepareMarshal(),
		ResponseModel: &userUpdated,
	}
	_, _, e := service.put(&requestConfig)
	if e != nil {
		return nil, e
	}

	return &userUpdated, nil
}

// DeleteUser deletes a specific user
//
func (service *Service) DeleteUser(userID int) *errortools.Error {
	requestConfig := go_http.RequestConfig{
		URL: service.url(fmt.Sprintf("Users/%v", userID)),
	}
	_, _, e := service.delete(&requestConfig)
	if e != nil {
		return e
	}

	return nil
}

func (u *User) FullName() string {
	if u == nil {
		return ""
	}

	return fmt.Sprintf("%s %s", u.FirstName, u.LastName)
}
