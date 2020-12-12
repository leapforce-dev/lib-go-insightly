package insightly

import (
	"fmt"
	"strconv"
	"time"

	errortools "github.com/leapforce-libraries/go_errortools"
	utilities "github.com/leapforce-libraries/go_utilities"
)

type Contact struct {
	ContactID            int           `json:"CONTACT_ID"`
	Salutation           string        `json:"SALUTATION"`
	FirstName            string        `json:"FIRST_NAME"`
	LastName             string        `json:"LAST_NAME"`
	ImageURL             string        `json:"IMAGE_URL"`
	Background           string        `json:"BACKGROUND"`
	OwnerUserID          int           `json:"OWNER_USER_ID"`
	DateCreatedUTC       string        `json:"DATE_CREATED_UTC"`
	DateUpdateUTC        string        `json:"DATE_UPDATED_UTC"`
	SocialLinkedin       string        `json:"SOCIAL_LINKEDIN"`
	SocialFacebook       string        `json:"SOCIAL_FACEBOOK"`
	SocialTwitter        string        `json:"SOCIAL_TWITTER"`
	DateOfBirth          string        `json:"DATE_OF_BIRTH"`
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
	LastActivityDateUTC  string        `json:"LAST_ACTIVITY_DATE_UTC"`
	NextActivityDateUTC  string        `json:"NEXT_ACTIVITY_DATE_UTC"`
	CreatedUserID        int           `json:"CREATED_USER_ID"`
	OrganisationID       int           `json:"ORGANISATION_ID"`
	Title                string        `json:"TITLE"`
	EmailOptedOut        bool          `json:"EMAIL_OPTED_OUT"`
	CustomFields         []CustomField `json:"CUSTOMFIELDS"`
	Tags                 []Tag         `json:"TAGS"`
	Dates                []Date        `json:"DATES"`
	DateCreatedT         *time.Time
	DateUpdatedT         *time.Time
	DateOfBirthT         *time.Time
	LastActivityDateT    *time.Time
	NextActivityDateT    *time.Time
}

// GetContacts returns all contacts
//
func (i *Insightly) GetContacts() ([]Contact, *errortools.Error) {
	return i.GetContactsInternal("")
}

// GetContactsUpdatedAfter returns all contacts updated after certain date
//
func (i *Insightly) GetContactsUpdatedAfter(updatedAfter time.Time) ([]Contact, *errortools.Error) {
	from := updatedAfter.Format("2006-01-02")
	searchFilter := fmt.Sprintf("updated_after_utc=%s&", from)
	return i.GetContactsInternal(searchFilter)
}

// GetContactsFiltered returns all contacts fulfulling the specified filter
//
func (i *Insightly) GetContactsFiltered(fieldname string, fieldvalue string) ([]Contact, *errortools.Error) {
	searchFilter := fmt.Sprintf("field_name=%s&field_value=%s&", fieldname, fieldvalue)
	return i.GetContactsInternal(searchFilter)
}

// GetContactsInternal is the generic function retrieving Contacts from Insightly
//
func (i *Insightly) GetContactsInternal(searchFilter string) ([]Contact, *errortools.Error) {
	searchString := ""

	if searchFilter != "" {
		searchString = "/Search?" + searchFilter
	} else {
		searchString = "?"
	}

	urlStr := "%sContacts%sskip=%s&top=%s"
	skip := 0
	top := 500
	rowCount := top

	contacts := []Contact{}

	for rowCount >= top {
		url := fmt.Sprintf(urlStr, apiURL, searchString, strconv.Itoa(skip), strconv.Itoa(top))
		//fmt.Printf(url)

		cs := []Contact{}

		_, _, err := i.get(url, nil, &cs)
		if err != nil {
			return nil, err
		}

		for _, c := range cs {
			c.parseDates()
			contacts = append(contacts, c)
		}

		rowCount = len(cs)
		//rowCount = 0
		skip += top
	}

	if len(contacts) == 0 {
		contacts = nil
	}

	return contacts, nil
}

func (c *Contact) ValidateEmail() *errortools.Error {
	// validate email
	if c.EmailAddress != "" {
		err := utilities.ValidateFormat(c.EmailAddress)
		if err != nil {
			return errortools.ErrorMessage(fmt.Sprintf("Invalid emailadress (between []): [%s] for contact: %s %s (%v)", c.EmailAddress, c.FirstName, c.LastName, c.ContactID))
		}
	}

	return nil
}

func (c *Contact) parseDates() {
	// parse DATE_CREATED_UTC to time.Time
	if c.DateCreatedUTC != "" {
		t, _ := time.Parse("2006-01-02 15:04:05 +0000 UTC", c.DateCreatedUTC+" +0000 UTC")
		//errortools.Fatal(err)
		c.DateCreatedT = &t
	}

	// parse DATE_UPDATED_UTC to time.Time
	if c.DateUpdateUTC != "" {
		t, _ := time.Parse("2006-01-02 15:04:05 +0000 UTC", c.DateUpdateUTC+" +0000 UTC")
		//errortools.Fatal(err)
		c.DateUpdatedT = &t
	}

	// parse DATE_OF_BIRTH to time.Time
	if c.DateOfBirth != "" {
		t, _ := time.Parse("2006-01-02 15:04:05 +0000 UTC", c.DateOfBirth+" +0000 UTC")
		//errortools.Fatal(err)
		c.DateOfBirthT = &t
	}

	// parse LAST_ACTIVITY_DATE_UTC to time.Time
	if c.LastActivityDateUTC != "" {
		t, _ := time.Parse("2006-01-02 15:04:05 +0000 UTC", c.LastActivityDateUTC+" +0000 UTC")
		//errortools.Fatal(err)
		c.LastActivityDateT = &t
	}

	// parse NEXT_ACTIVITY_DATE_UTC to time.Time
	if c.NextActivityDateUTC != "" {
		t, _ := time.Parse("2006-01-02 15:04:05 +0000 UTC", c.NextActivityDateUTC+" +0000 UTC")
		//errortools.Fatal(err)
		c.NextActivityDateT = &t
	}

	// parse dates in DATES
	for _, d := range c.Dates {
		d.parseDates()
	}
}
