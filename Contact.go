package insightly

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
	i_types "github.com/leapforce-libraries/go_insightly/types"
)

type Contact struct {
	ContactID            int64                   `json:"CONTACT_ID"`
	Salutation           *string                 `json:"SALUTATION,omitempty"`
	FirstName            *string                 `json:"FIRST_NAME,omitempty"`
	LastName             *string                 `json:"LAST_NAME,omitempty"`
	ImageUrl             *string                 `json:"IMAGE_URL,omitempty"`
	Background           *string                 `json:"BACKGROUND,omitempty"`
	OwnerUserID          *int64                  `json:"OWNER_USER_ID,omitempty"`
	DateCreatedUTC       *i_types.DateTimeString `json:"DATE_CREATED_UTC,omitempty"`
	DateUpdatedUTC       *i_types.DateTimeString `json:"DATE_UPDATED_UTC,omitempty"`
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
	AddressOtherCountry  *string                 `json:"ADDRESS_OTHER_COUNTRY,omitempty"`
	LastActivityDateUTC  *i_types.DateTimeString `json:"LAST_ACTIVITY_DATE_UTC,omitempty"`
	NextActivityDateUTC  *i_types.DateTimeString `json:"NEXT_ACTIVITY_DATE_UTC,omitempty"`
	CreatedUserID        *int64                  `json:"CREATED_USER_ID,omitempty"`
	OrganisationID       *int64                  `json:"ORGANISATION_ID,omitempty"`
	Title                *string                 `json:"TITLE,omitempty"`
	EmailOptedOut        *bool                   `json:"EMAIL_OPTED_OUT,omitempty"`
	CustomFields         *CustomFields           `json:"CUSTOMFIELDS,omitempty"`
	Tags                 *[]Tag                  `json:"TAGS,omitempty"`
	Dates                *[]Date                 `json:"DATES,omitempty"`
	Links                *[]Link                 `json:"LINKS,omitempty"`
}

// GetContact returns a specific contact
func (service *Service) GetContact(contactID int64) (*Contact, *errortools.Error) {
	contact := Contact{}

	requestConfig := go_http.RequestConfig{
		Method:        http.MethodGet,
		Url:           service.url(fmt.Sprintf("Contacts/%v", contactID)),
		ResponseModel: &contact,
	}
	_, _, e := service.httpRequest(&requestConfig)
	if e != nil {
		return nil, e
	}

	return &contact, nil
}

type GetContactsConfig struct {
	Skip         *uint64
	Top          *uint64
	Brief        *bool
	CountTotal   *bool
	UpdatedAfter *time.Time
	FieldFilter  *FieldFilter
}

// GetContacts returns all contacts
func (service *Service) GetContacts(config *GetContactsConfig) (*[]Contact, *errortools.Error) {
	params := url.Values{}

	endpoint := "Contacts"
	contacts := []Contact{}
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
		contactsBatch := []Contact{}

		requestConfig := go_http.RequestConfig{
			Method:        http.MethodGet,
			Url:           service.url(fmt.Sprintf("%s?%s", endpoint, params.Encode())),
			ResponseModel: &contactsBatch,
		}
		_, _, e := service.httpRequest(&requestConfig)
		if e != nil {
			return nil, e
		}

		contacts = append(contacts, contactsBatch...)

		if len(contactsBatch) < int(top) {
			delete(service.nextSkips, endpoint)
			break
		}

		service.nextSkips[endpoint] += top
		rowCount += top

		if rowCount >= service.maxRowCount {
			return &contacts, nil
		}
	}

	return &contacts, nil
}

// CreateContact creates a new contract
func (service *Service) CreateContact(contact *Contact) (*Contact, *errortools.Error) {
	if contact == nil {
		return nil, nil
	}

	contactNew := Contact{}

	requestConfig := go_http.RequestConfig{
		Method:        http.MethodPost,
		Url:           service.url("Contacts"),
		BodyModel:     contact,
		ResponseModel: &contactNew,
	}
	_, _, e := service.httpRequest(&requestConfig)
	if e != nil {
		return nil, e
	}

	return &contactNew, nil
}

// UpdateContact updates an existing contract
func (service *Service) UpdateContact(contact *Contact) (*Contact, *errortools.Error) {
	if contact == nil {
		return nil, nil
	}

	contactUpdated := Contact{}

	requestConfig := go_http.RequestConfig{
		Method:        http.MethodPut,
		Url:           service.url("Contacts"),
		BodyModel:     contact,
		ResponseModel: &contactUpdated,
	}
	_, _, e := service.httpRequest(&requestConfig)
	if e != nil {
		return nil, e
	}

	return &contactUpdated, nil
}

// DeleteContact deletes a specific contact
func (service *Service) DeleteContact(contactID int64) *errortools.Error {
	requestConfig := go_http.RequestConfig{
		Method: http.MethodDelete,
		Url:    service.url(fmt.Sprintf("Contacts/%v", contactID)),
	}
	_, _, e := service.httpRequest(&requestConfig)
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

// GetContactFileAttachments returns the file attachments of a specific email
func (service *Service) GetContactFileAttachments(id int64) (*[]FileAttachment, *errortools.Error) {
	var fileAttachments []FileAttachment

	requestConfig := go_http.RequestConfig{
		Method:        http.MethodGet,
		Url:           service.url(fmt.Sprintf("contacts/%v/fileattachments", id)),
		ResponseModel: &fileAttachments,
	}

	_, _, e := service.httpRequest(&requestConfig)
	if e != nil {
		return nil, e
	}

	return &fileAttachments, nil
}
