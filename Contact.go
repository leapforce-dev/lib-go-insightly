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

type Contact struct {
	ContactID            int64                   `json:"CONTACT_ID"`
	Salutation           *string                 `json:"SALUTATION"`
	FirstName            *string                 `json:"FIRST_NAME"`
	LastName             *string                 `json:"LAST_NAME"`
	ImageURL             *string                 `json:"IMAGE_URL"`
	Background           *string                 `json:"BACKGROUND"`
	OwnerUserID          int64                   `json:"OWNER_USER_ID"`
	DateCreatedUTC       i_types.DateTimeString  `json:"DATE_CREATED_UTC"`
	DateUpdatedUTC       i_types.DateTimeString  `json:"DATE_UPDATED_UTC"`
	SocialLinkedin       *string                 `json:"SOCIAL_LINKEDIN"`
	SocialFacebook       *string                 `json:"SOCIAL_FACEBOOK"`
	SocialTwitter        *string                 `json:"SOCIAL_TWITTER"`
	DateOfBirth          *i_types.DateTimeString `json:"DATE_OF_BIRTH"`
	Phone                *string                 `json:"PHONE"`
	PhoneHome            *string                 `json:"PHONE_HOME"`
	PhoneMobile          *string                 `json:"PHONE_MOBILE"`
	PhoneOther           *string                 `json:"PHONE_OTHER"`
	PhoneAssistant       *string                 `json:"PHONE_ASSISTANT"`
	PhoneFax             *string                 `json:"PHONE_FAX"`
	EmailAddress         *string                 `json:"EMAIL_ADDRESS"`
	AssistantName        *string                 `json:"ASSISTANT_NAME"`
	AddressMailStreet    *string                 `json:"ADDRESS_MAIL_STREET"`
	AddressMailCity      *string                 `json:"ADDRESS_MAIL_CITY"`
	AddressMailState     *string                 `json:"ADDRESS_MAIL_STATE"`
	AddressMailPostcode  *string                 `json:"ADDRESS_MAIL_POSTCODE"`
	AddressMailCountry   *string                 `json:"ADDRESS_MAIL_COUNTRY"`
	AddressOtherStreet   *string                 `json:"ADDRESS_OTHER_STREET"`
	AddressOtherCity     *string                 `json:"ADDRESS_OTHER_CITY"`
	AddressOtherState    *string                 `json:"ADDRESS_OTHER_STATE"`
	AddressOtherPostcode *string                 `json:"ADDRESS_OTHER_POSTCODE"`
	LastActivityDateUTC  *i_types.DateTimeString `json:"LAST_ACTIVITY_DATE_UTC"`
	NextActivityDateUTC  *i_types.DateTimeString `json:"NEXT_ACTIVITY_DATE_UTC"`
	CreatedUserID        int64                   `json:"CREATED_USER_ID"`
	OrganisationID       *int64                  `json:"ORGANISATION_ID"`
	Title                *string                 `json:"TITLE"`
	EmailOptedOut        bool                    `json:"EMAIL_OPTED_OUT"`
	CustomFields         *CustomFields           `json:"CUSTOMFIELDS"`
	Tags                 *[]Tag                  `json:"TAGS"`
	Dates                *[]Date                 `json:"DATES"`
}

func (c *Contact) prepareMarshal() interface{} {
	if c == nil {
		return nil
	}

	return &struct {
		ContactID            *int64                  `json:"CONTACT_ID,omitempty"`
		Salutation           *string                 `json:"SALUTATION,omitempty"`
		FirstName            *string                 `json:"FIRST_NAME,omitempty"`
		LastName             *string                 `json:"LAST_NAME,omitempty"`
		ImageURL             *string                 `json:"IMAGE_URL,omitempty"`
		Background           *string                 `json:"BACKGROUND,omitempty"`
		OwnerUserID          *int64                  `json:"OWNER_USER_ID,omitempty"`
		SocialLinkedin       *string                 `json:"SOCIAL_LINKEDIN,omitempty"`
		SocialFacebook       *string                 `json:"SOCIAL_FACEBOOK,omitempty"`
		SocialTwitter        *string                 `json:"SOCIAL_TWITTER,omitempty"`
		DateOfBirth          *i_types.DateTimeString `json:"DATE_OF_BIRTH,omitempty"`
		Phone                *string                 `json:"PHONE,omitempty"`
		PhoneHome            *string                 `json:"PHONE_HOME,omitempty"`
		PhoneMobile          *string                 `json:"PHONE_MOBILE,omitempty"`
		PhoneOther           *string                 `json:"PHONE_OTHER,omitempty"`
		PhoneAssistant       *string                 `json:"PHONE_ASSISTANT,omitempty"`
		PhoneFax             *string                 `json:"PHONE_FAX,omitempty"`
		EmailAddress         *string                 `json:"EMAIL_ADDRESS,omitempty"`
		AssistantName        *string                 `json:"ASSISTANT_NAME,omitempty"`
		AddressMailStreet    *string                 `json:"ADDRESS_MAIL_STREET,omitempty"`
		AddressMailCity      *string                 `json:"ADDRESS_MAIL_CITY,omitempty"`
		AddressMailState     *string                 `json:"ADDRESS_MAIL_STATE,omitempty"`
		AddressMailPostcode  *string                 `json:"ADDRESS_MAIL_POSTCODE,omitempty"`
		AddressMailCountry   *string                 `json:"ADDRESS_MAIL_COUNTRY,omitempty"`
		AddressOtherStreet   *string                 `json:"ADDRESS_OTHER_STREET,omitempty"`
		AddressOtherCity     *string                 `json:"ADDRESS_OTHER_CITY,omitempty"`
		AddressOtherState    *string                 `json:"ADDRESS_OTHER_STATE,omitempty"`
		AddressOtherPostcode *string                 `json:"ADDRESS_OTHER_POSTCODE,omitempty"`
		OrganisationID       *int64                  `json:"ORGANISATION_ID,omitempty"`
		Title                *string                 `json:"TITLE,omitempty"`
		EmailOptedOut        *bool                   `json:"EMAIL_OPTED_OUT,omitempty"`
		CustomFields         *CustomFields           `json:"CUSTOMFIELDS,omitempty"`
	}{
		&c.ContactID,
		c.Salutation,
		c.FirstName,
		c.LastName,
		c.ImageURL,
		c.Background,
		&c.OwnerUserID,
		c.SocialLinkedin,
		c.SocialFacebook,
		c.SocialTwitter,
		c.DateOfBirth,
		c.Phone,
		c.PhoneHome,
		c.PhoneMobile,
		c.PhoneOther,
		c.PhoneAssistant,
		c.PhoneFax,
		c.EmailAddress,
		c.AssistantName,
		c.AddressMailStreet,
		c.AddressMailCity,
		c.AddressMailState,
		c.AddressMailPostcode,
		c.AddressMailCountry,
		c.AddressOtherStreet,
		c.AddressOtherCity,
		c.AddressOtherState,
		c.AddressOtherPostcode,
		c.OrganisationID,
		c.Title,
		&c.EmailOptedOut,
		c.CustomFields,
	}
}

// GetContact returns a specific contact
//
func (service *Service) GetContact(contactID int) (*Contact, *errortools.Error) {
	contact := Contact{}

	requestConfig := go_http.RequestConfig{
		URL:           service.url(fmt.Sprintf("Contacts/%v", contactID)),
		ResponseModel: &contact,
	}
	_, _, e := service.get(&requestConfig)
	if e != nil {
		return nil, e
	}

	return &contact, nil
}

type GetContactsConfig struct {
	UpdatedAfter *time.Time
	Field        *struct {
		FieldName  string
		FieldValue string
	}
}

// GetContacts returns all contacts
//
func (service *Service) GetContacts(config *GetContactsConfig) (*[]Contact, *errortools.Error) {
	searchString := "?"
	searchFilter := []string{}

	if config != nil {
		if config.UpdatedAfter != nil {
			from := config.UpdatedAfter.Format(DateTimeFormat)
			searchFilter = append(searchFilter, fmt.Sprintf("updated_after_utc=%s&", from))
		}

		if config.Field != nil {
			searchFilter = append(searchFilter, fmt.Sprintf("field_name=%s&field_value=%s&", config.Field.FieldName, config.Field.FieldValue))
		}
	}

	if len(searchFilter) > 0 {
		searchString = "/Search?" + strings.Join(searchFilter, "&")
	}

	endpointStr := "Contacts%sskip=%s&top=%s"
	skip := 0
	top := 100
	rowCount := top

	contacts := []Contact{}

	for rowCount >= top {
		_contacts := []Contact{}

		requestConfig := go_http.RequestConfig{
			URL:           service.url(fmt.Sprintf(endpointStr, searchString, strconv.Itoa(skip), strconv.Itoa(top))),
			ResponseModel: &_contacts,
		}
		_, _, e := service.get(&requestConfig)
		if e != nil {
			return nil, e
		}

		contacts = append(contacts, _contacts...)

		rowCount = len(_contacts)
		//rowCount = 0
		skip += top
	}

	if len(contacts) == 0 {
		contacts = nil
	}

	return &contacts, nil
}

// CreateContact creates a new contract
//
func (service *Service) CreateContact(contact *Contact) (*Contact, *errortools.Error) {
	if contact == nil {
		return nil, nil
	}

	contactNew := Contact{}

	requestConfig := go_http.RequestConfig{
		URL:           service.url("Contacts"),
		BodyModel:     contact.prepareMarshal(),
		ResponseModel: &contactNew,
	}
	_, _, e := service.post(&requestConfig)
	if e != nil {
		return nil, e
	}

	return &contactNew, nil
}

// UpdateContact updates an existing contract
//
func (service *Service) UpdateContact(contact *Contact) (*Contact, *errortools.Error) {
	if contact == nil {
		return nil, nil
	}

	contactUpdated := Contact{}

	requestConfig := go_http.RequestConfig{
		URL:           service.url("Contacts"),
		BodyModel:     contact.prepareMarshal(),
		ResponseModel: &contactUpdated,
	}
	_, _, e := service.put(&requestConfig)
	if e != nil {
		return nil, e
	}

	return &contactUpdated, nil
}

// DeleteContact deletes a specific contact
//
func (service *Service) DeleteContact(contactID int) *errortools.Error {
	requestConfig := go_http.RequestConfig{
		URL: service.url(fmt.Sprintf("Contacts/%v", contactID)),
	}
	_, _, e := service.delete(&requestConfig)
	if e != nil {
		return e
	}

	return nil
}

func (c *Contact) FullName() string {
	if c == nil {
		return ""
	}

	name := ""
	if c.LastName != nil {
		name = *c.LastName
	}
	if c.FirstName != nil {
		name = fmt.Sprintf("%s %s", *c.FirstName, name)
	}

	return strings.Trim(name, " ")
}
