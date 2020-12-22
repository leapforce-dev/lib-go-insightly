package insightly

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	errortools "github.com/leapforce-libraries/go_errortools"
)

type Contact struct {
	ContactID            int          `json:"CONTACT_ID"`
	Salutation           string       `json:"SALUTATION"`
	FirstName            string       `json:"FIRST_NAME"`
	LastName             string       `json:"LAST_NAME"`
	ImageURL             string       `json:"IMAGE_URL"`
	Background           string       `json:"BACKGROUND"`
	OwnerUserID          *int         `json:"OWNER_USER_ID"`
	DateCreatedUTC       DateUTC      `json:"DATE_CREATED_UTC"`
	DateUpdatedUTC       DateUTC      `json:"DATE_UPDATED_UTC"`
	SocialLinkedin       string       `json:"SOCIAL_LINKEDIN"`
	SocialFacebook       string       `json:"SOCIAL_FACEBOOK"`
	SocialTwitter        string       `json:"SOCIAL_TWITTER"`
	DateOfBirth          DateUTC      `json:"DATE_OF_BIRTH"`
	Phone                string       `json:"PHONE"`
	PhoneHome            string       `json:"PHONE_HOME"`
	PhoneMobile          string       `json:"PHONE_MOBILE"`
	PhoneOther           string       `json:"PHONE_OTHER"`
	PhoneAssistant       string       `json:"PHONE_ASSISTANT"`
	PhoneFax             string       `json:"PHONE_FAX"`
	EmailAddress         string       `json:"EMAIL_ADDRESS"`
	AssistantName        string       `json:"ASSISTANT_NAME"`
	AddressMailStreet    string       `json:"ADDRESS_MAIL_STREET"`
	AddressMailCity      string       `json:"ADDRESS_MAIL_CITY"`
	AddressMailState     string       `json:"ADDRESS_MAIL_STATE"`
	AddressMailPostcode  string       `json:"ADDRESS_MAIL_POSTCODE"`
	AddressMailCountry   string       `json:"ADDRESS_MAIL_COUNTRY"`
	AddressOtherStreet   string       `json:"ADDRESS_OTHER_STREET"`
	AddressOtherCity     string       `json:"ADDRESS_OTHER_CITY"`
	AddressOtherState    string       `json:"ADDRESS_OTHER_STATE"`
	AddressOtherPostcode string       `json:"ADDRESS_OTHER_POSTCODE"`
	LastActivityDateUTC  DateUTC      `json:"LAST_ACTIVITY_DATE_UTC"`
	NextActivityDateUTC  DateUTC      `json:"NEXT_ACTIVITY_DATE_UTC"`
	CreatedUserID        *int         `json:"CREATED_USER_ID"`
	OrganisationID       *int         `json:"ORGANISATION_ID"`
	Title                string       `json:"TITLE"`
	EmailOptedOut        bool         `json:"EMAIL_OPTED_OUT"`
	CustomFields         CustomFields `json:"CUSTOMFIELDS"`
	Tags                 []Tag        `json:"TAGS"`
	Dates                []Date       `json:"DATES"`
}

func (c *Contact) prepareMarshal() interface{} {
	if c == nil {
		return nil
	}

	return &struct {
		ContactID            int          `json:"CONTACT_ID"`
		Salutation           string       `json:"SALUTATION"`
		FirstName            string       `json:"FIRST_NAME"`
		LastName             string       `json:"LAST_NAME"`
		ImageURL             string       `json:"IMAGE_URL"`
		Background           string       `json:"BACKGROUND"`
		OwnerUserID          *int         `json:"OWNER_USER_ID"`
		SocialLinkedin       string       `json:"SOCIAL_LINKEDIN"`
		SocialFacebook       string       `json:"SOCIAL_FACEBOOK"`
		SocialTwitter        string       `json:"SOCIAL_TWITTER"`
		DateOfBirth          DateUTC      `json:"DATE_OF_BIRTH"`
		Phone                string       `json:"PHONE"`
		PhoneHome            string       `json:"PHONE_HOME"`
		PhoneMobile          string       `json:"PHONE_MOBILE"`
		PhoneOther           string       `json:"PHONE_OTHER"`
		PhoneAssistant       string       `json:"PHONE_ASSISTANT"`
		PhoneFax             string       `json:"PHONE_FAX"`
		EmailAddress         string       `json:"EMAIL_ADDRESS"`
		AssistantName        string       `json:"ASSISTANT_NAME"`
		AddressMailStreet    string       `json:"ADDRESS_MAIL_STREET"`
		AddressMailCity      string       `json:"ADDRESS_MAIL_CITY"`
		AddressMailState     string       `json:"ADDRESS_MAIL_STATE"`
		AddressMailPostcode  string       `json:"ADDRESS_MAIL_POSTCODE"`
		AddressMailCountry   string       `json:"ADDRESS_MAIL_COUNTRY"`
		AddressOtherStreet   string       `json:"ADDRESS_OTHER_STREET"`
		AddressOtherCity     string       `json:"ADDRESS_OTHER_CITY"`
		AddressOtherState    string       `json:"ADDRESS_OTHER_STATE"`
		AddressOtherPostcode string       `json:"ADDRESS_OTHER_POSTCODE"`
		OrganisationID       *int         `json:"ORGANISATION_ID"`
		Title                string       `json:"TITLE"`
		EmailOptedOut        bool         `json:"EMAIL_OPTED_OUT"`
		CustomFields         CustomFields `json:"CUSTOMFIELDS"`
	}{
		c.ContactID,
		c.Salutation,
		c.FirstName,
		c.LastName,
		c.ImageURL,
		c.Background,
		c.OwnerUserID,
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
		c.EmailOptedOut,
		c.CustomFields,
	}
}

// GetContact returns a specific contact
//
func (i *Insightly) GetContact(contactID int) (*Contact, *errortools.Error) {
	endpoint := fmt.Sprintf("Contacts/%v", contactID)

	contact := Contact{}

	_, _, e := i.get(endpoint, nil, &contact)
	if e != nil {
		return nil, e
	}

	return &contact, nil
}

type GetContactsFilter struct {
	UpdatedAfter *time.Time
	Field        *struct {
		FieldName  string
		FieldValue string
	}
}

// GetContacts returns all contacts
//
func (i *Insightly) GetContacts(filter *GetContactsFilter) (*[]Contact, *errortools.Error) {
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

	endpointStr := "Contacts%sskip=%s&top=%s"
	skip := 0
	top := 100
	rowCount := top

	contacts := []Contact{}

	for rowCount >= top {
		endpoint := fmt.Sprintf(endpointStr, searchString, strconv.Itoa(skip), strconv.Itoa(top))
		//fmt.Println(endpoint)

		cs := []Contact{}

		_, _, e := i.get(endpoint, nil, &cs)
		if e != nil {
			return nil, e
		}

		contacts = append(contacts, cs...)

		rowCount = len(cs)
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
func (i *Insightly) CreateContact(contact *Contact) (*Contact, *errortools.Error) {
	if contact == nil {
		return nil, nil
	}

	endpoint := "Contacts"

	contactNew := Contact{}

	_, _, e := i.post(endpoint, contact.prepareMarshal(), &contactNew)
	if e != nil {
		return nil, e
	}

	return &contactNew, nil
}

// UpdateContact updates an existing contract
//
func (i *Insightly) UpdateContact(contact *Contact) (*Contact, *errortools.Error) {
	if contact == nil {
		return nil, nil
	}

	endpoint := "Contacts"

	contactUpdated := Contact{}

	_, _, e := i.put(endpoint, contact.prepareMarshal(), &contactUpdated)
	if e != nil {
		return nil, e
	}

	return &contactUpdated, nil
}

// DeleteContact deletes a specific contact
//
func (i *Insightly) DeleteContact(contactID int) *errortools.Error {
	endpoint := fmt.Sprintf("Contacts/%v", contactID)

	_, _, e := i.delete(endpoint, nil, nil)
	if e != nil {
		return e
	}

	return nil
}

func (c *Contact) FullName() string {
	if c == nil {
		return ""
	}

	return fmt.Sprintf("%s %s", c.FirstName, c.LastName)
}

/*
func (c *Contact) ValidateEmail() *errortools.Error {
	// validate email
	if c.EmailAddress != "" {
		err := utilities.ValidateFormat(c.EmailAddress)
		if err != nil {
			return errortools.ErrorMessage(fmt.Sprintf("Invalid emailadress (between []): [%s] for contact: %s %s (%v)", c.EmailAddress, c.FirstName, c.LastName, c.ContactID))
		}
	}

	return nil
}*/
