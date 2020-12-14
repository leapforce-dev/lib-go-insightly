package insightly

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	errortools "github.com/leapforce-libraries/go_errortools"
)

type User struct {
	UserID               int           `json:"CONTACT_ID"`
	Salutation           string        `json:"SALUTATION"`
	FirstName            string        `json:"FIRST_NAME"`
	LastName             string        `json:"LAST_NAME"`
	ImageURL             string        `json:"IMAGE_URL"`
	Background           string        `json:"BACKGROUND"`
	OwnerUserID          int           `json:"OWNER_USER_ID"`
	DateCreatedUTC       DateUTC       `json:"DATE_CREATED_UTC"`
	DateUpdatedUTC       DateUTC       `json:"DATE_UPDATED_UTC"`
	SocialLinkedin       string        `json:"SOCIAL_LINKEDIN"`
	SocialFacebook       string        `json:"SOCIAL_FACEBOOK"`
	SocialTwitter        string        `json:"SOCIAL_TWITTER"`
	DateOfBirth          DateUTC       `json:"DATE_OF_BIRTH"`
	Phone                string        `json:"PHONE"`
	PhoneHome            string        `json:"PHONE_HOME"`
	PhoneMobile          string        `json:"PHONE_MOBILE"`
	PhoneOther           string        `json:"PHONE_OTHER"`
	PhoneAssistant       string        `json:"PHONE_ASSISTANT"`
	PhoneFax             string        `json:"PHONE_FAX"`
	EmailAddress         string        `json:"EMAIL_ADDRESS"`
	AssistantName        string        `json:"ASSISTANT_NAME"`
	AddressMailStreet    string        `json:"ADDRESS_MAIL_STREET"`
	AddressMailCity      string        `json:"ADDRESS_MAIL_CITY"`
	AddressMailState     string        `json:"ADDRESS_MAIL_STATE"`
	AddressMailPostcode  string        `json:"ADDRESS_MAIL_POSTCODE"`
	AddressMailCountry   string        `json:"ADDRESS_MAIL_COUNTRY"`
	AddressOtherStreet   string        `json:"ADDRESS_OTHER_STREET"`
	AddressOtherCity     string        `json:"ADDRESS_OTHER_CITY"`
	AddressOtherState    string        `json:"ADDRESS_OTHER_STATE"`
	AddressOtherPostcode string        `json:"ADDRESS_OTHER_POSTCODE"`
	LastActivityDateUTC  DateUTC       `json:"LAST_ACTIVITY_DATE_UTC"`
	NextActivityDateUTC  DateUTC       `json:"NEXT_ACTIVITY_DATE_UTC"`
	CreatedUserID        int           `json:"CREATED_USER_ID"`
	OrganisationID       int           `json:"ORGANISATION_ID"`
	Title                string        `json:"TITLE"`
	EmailOptedOut        bool          `json:"EMAIL_OPTED_OUT"`
	CustomFields         []CustomField `json:"CUSTOMFIELDS"`
	Tags                 []Tag         `json:"TAGS"`
	Dates                []Date        `json:"DATES"`
}

func (u *User) prepareMarshal() interface{} {
	if u == nil {
		return nil
	}

	return &struct {
		UserID               int           `json:"CONTACT_ID"`
		Salutation           string        `json:"SALUTATION"`
		FirstName            string        `json:"FIRST_NAME"`
		LastName             string        `json:"LAST_NAME"`
		ImageURL             string        `json:"IMAGE_URL"`
		Background           string        `json:"BACKGROUND"`
		OwnerUserID          int           `json:"OWNER_USER_ID"`
		SocialLinkedin       string        `json:"SOCIAL_LINKEDIN"`
		SocialFacebook       string        `json:"SOCIAL_FACEBOOK"`
		SocialTwitter        string        `json:"SOCIAL_TWITTER"`
		DateOfBirth          DateUTC       `json:"DATE_OF_BIRTH"`
		Phone                string        `json:"PHONE"`
		PhoneHome            string        `json:"PHONE_HOME"`
		PhoneMobile          string        `json:"PHONE_MOBILE"`
		PhoneOther           string        `json:"PHONE_OTHER"`
		PhoneAssistant       string        `json:"PHONE_ASSISTANT"`
		PhoneFax             string        `json:"PHONE_FAX"`
		EmailAddress         string        `json:"EMAIL_ADDRESS"`
		AssistantName        string        `json:"ASSISTANT_NAME"`
		AddressMailStreet    string        `json:"ADDRESS_MAIL_STREET"`
		AddressMailCity      string        `json:"ADDRESS_MAIL_CITY"`
		AddressMailState     string        `json:"ADDRESS_MAIL_STATE"`
		AddressMailPostcode  string        `json:"ADDRESS_MAIL_POSTCODE"`
		AddressMailCountry   string        `json:"ADDRESS_MAIL_COUNTRY"`
		AddressOtherStreet   string        `json:"ADDRESS_OTHER_STREET"`
		AddressOtherCity     string        `json:"ADDRESS_OTHER_CITY"`
		AddressOtherState    string        `json:"ADDRESS_OTHER_STATE"`
		AddressOtherPostcode string        `json:"ADDRESS_OTHER_POSTCODE"`
		LastActivityDateUTC  DateUTC       `json:"LAST_ACTIVITY_DATE_UTC"`
		NextActivityDateUTC  DateUTC       `json:"NEXT_ACTIVITY_DATE_UTC"`
		OrganisationID       int           `json:"ORGANISATION_ID"`
		Title                string        `json:"TITLE"`
		EmailOptedOut        bool          `json:"EMAIL_OPTED_OUT"`
		CustomFields         []CustomField `json:"CUSTOMFIELDS"`
	}{
		u.UserID,
		u.Salutation,
		u.FirstName,
		u.LastName,
		u.ImageURL,
		u.Background,
		u.OwnerUserID,
		u.SocialLinkedin,
		u.SocialFacebook,
		u.SocialTwitter,
		u.DateOfBirth,
		u.Phone,
		u.PhoneHome,
		u.PhoneMobile,
		u.PhoneOther,
		u.PhoneAssistant,
		u.PhoneFax,
		u.EmailAddress,
		u.AssistantName,
		u.AddressMailStreet,
		u.AddressMailCity,
		u.AddressMailState,
		u.AddressMailPostcode,
		u.AddressMailCountry,
		u.AddressOtherStreet,
		u.AddressOtherCity,
		u.AddressOtherState,
		u.AddressOtherPostcode,
		u.LastActivityDateUTC,
		u.NextActivityDateUTC,
		u.OrganisationID,
		u.Title,
		u.EmailOptedOut,
		u.CustomFields,
	}
}

// GetUser returns a specific user
//
func (i *Insightly) GetUser(userID int) (*User, *errortools.Error) {
	endpoint := fmt.Sprintf("Users/%v", userID)

	user := User{}

	_, _, e := i.get(endpoint, nil, &user)
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
func (i *Insightly) GetUsers(filter *GetUsersFilter) (*[]User, *errortools.Error) {
	searchString := "?"
	searchFilter := []string{}

	if filter != nil {
		if filter.UpdatedAfter != nil {
			from := filter.UpdatedAfter.Format(ISO8601Format)
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
		endpoint := fmt.Sprintf(endpointStr, searchString, strconv.Itoa(skip), strconv.Itoa(top))
		//fmt.Println(endpoint)

		cs := []User{}

		_, _, e := i.get(endpoint, nil, &cs)
		if e != nil {
			return nil, e
		}

		users = append(users, cs...)

		rowCount = len(cs)
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
func (i *Insightly) CreateUser(user *User) (*User, *errortools.Error) {
	if user == nil {
		return nil, nil
	}

	endpoint := "Users"

	userNew := User{}

	_, _, e := i.post(endpoint, user.prepareMarshal(), &userNew)
	if e != nil {
		return nil, e
	}

	return &userNew, nil
}

// UpdateUser updates an existing contract
//
func (i *Insightly) UpdateUser(user *User) (*User, *errortools.Error) {
	if user == nil {
		return nil, nil
	}

	endpoint := "Users"

	userUpdated := User{}

	_, _, e := i.put(endpoint, user.prepareMarshal(), &userUpdated)
	if e != nil {
		return nil, e
	}

	return &userUpdated, nil
}

// DeleteUser deletes a specific user
//
func (i *Insightly) DeleteUser(userID int) *errortools.Error {
	endpoint := fmt.Sprintf("Users/%v", userID)

	_, _, e := i.delete(endpoint, nil, nil)
	if e != nil {
		return e
	}

	return nil
}
