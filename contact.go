package insightly

import (
	"fmt"
	"strconv"
	"time"

	types "github.com/Leapforce-nl/go_types"
)

type Contact struct {
	CONTACT_ID             int           `json:"CONTACT_ID"`
	SALUTATION             string        `json:"SALUTATION"`
	FIRST_NAME             string        `json:"FIRST_NAME"`
	LAST_NAME              string        `json:"LAST_NAME"`
	IMAGE_URL              string        `json:"IMAGE_URL"`
	BACKGROUND             string        `json:"BACKGROUND"`
	OWNER_USER_ID          int           `json:"OWNER_USER_ID"`
	DATE_CREATED_UTC       string        `json:"DATE_CREATED_UTC"`
	DATE_UPDATED_UTC       string        `json:"DATE_UPDATED_UTC"`
	SOCIAL_LINKEDIN        string        `json:"SOCIAL_LINKEDIN"`
	SOCIAL_FACEBOOK        string        `json:"SOCIAL_FACEBOOK"`
	SOCIAL_TWITTER         string        `json:"SOCIAL_TWITTER"`
	DATE_OF_BIRTH          string        `json:"DATE_OF_BIRTH"`
	PHONE                  string        `json:"PHONE"`
	PHONE_HOME             string        `json:"PHONE_HOME"`
	PHONE_MOBILE           string        `json:"PHONE_MOBILE"`
	PHONE_OTHER            string        `json:"PHONE_OTHER"`
	PHONE_ASSISTANT        string        `json:"PHONE_ASSISTANT"`
	PHONE_FAX              string        `json:"PHONE_FAX"`
	EMAIL_ADDRESS          string        `json:"EMAIL_ADDRESS"`
	ASSISTANT_NAME         string        `json:"ASSISTANT_NAME"`
	ADDRESS_MAIL_STREET    string        `json:"ADDRESS_MAIL_STREET"`
	ADDRESS_MAIL_CITY      string        `json:"ADDRESS_MAIL_CITY"`
	ADDRESS_MAIL_STATE     string        `json:"ADDRESS_MAIL_STATE"`
	ADDRESS_MAIL_POSTCODE  string        `json:"ADDRESS_MAIL_POSTCODE"`
	ADDRESS_MAIL_COUNTRY   string        `json:"ADDRESS_MAIL_COUNTRY"`
	ADDRESS_OTHER_STREET   string        `json:"ADDRESS_OTHER_STREET"`
	ADDRESS_OTHER_CITY     string        `json:"ADDRESS_OTHER_CITY"`
	ADDRESS_OTHER_STATE    string        `json:"ADDRESS_OTHER_STATE"`
	ADDRESS_OTHER_POSTCODE string        `json:"ADDRESS_OTHER_POSTCODE"`
	LAST_ACTIVITY_DATE_UTC string        `json:"LAST_ACTIVITY_DATE_UTC"`
	NEXT_ACTIVITY_DATE_UTC string        `json:"NEXT_ACTIVITY_DATE_UTC"`
	CREATED_USER_ID        int           `json:"CREATED_USER_ID"`
	ORGANISATION_ID        int           `json:"ORGANISATION_ID"`
	TITLE                  string        `json:"TITLE"`
	EMAIL_OPTED_OUT        bool          `json:"EMAIL_OPTED_OUT"`
	CUSTOMFIELDS           []CustomField `json:"CUSTOMFIELDS"`
	TAGS                   []Tag         `json:"TAGS"`
	DATES                  []Date        `json:"DATES"`
	DateCreated            *time.Time
	DateUpdated            *time.Time
	DateOfBirth            *time.Time
	LastActivityDate       *time.Time
	NextActivityDate       *time.Time
}

// GetContacts returns all contacts
//
func (i *Insightly) GetContacts() ([]Contact, error) {
	return i.GetContactsInternal("")
}

// GetContactsUpdatedAfter returns all contacts updated after certain date
//
func (i *Insightly) GetContactsUpdatedAfter(updatedAfter time.Time) ([]Contact, error) {
	from := updatedAfter.Format("2006-01-02")
	searchFilter := fmt.Sprintf("updated_after_utc=%s&", from)
	return i.GetContactsInternal(searchFilter)
}

// GetContactsFiltered returns all contacts fulfulling the specified filter
//
func (i *Insightly) GetContactsFiltered(fieldname string, fieldvalue string) ([]Contact, error) {
	searchFilter := fmt.Sprintf("field_name=%s&field_value=%s&", fieldname, fieldvalue)
	return i.GetContactsInternal(searchFilter)
}

// GetContactsInternal is the generic function retrieving Contacts from Insightly
//
func (i *Insightly) GetContactsInternal(searchFilter string) ([]Contact, error) {
	searchString := ""

	if searchFilter != "" {
		searchString = "/Search?" + searchFilter
	} else {
		searchString = "?"
	}

	urlStr := "%sContacts%sskip=%s&top=%s"
	skip := 0
	top := 500
	rowCount := 1

	contacts := []Contact{}

	for rowCount > 0 {
		url := fmt.Sprintf(urlStr, i.apiURL, searchString, strconv.Itoa(skip), strconv.Itoa(top))
		//fmt.Printf(url)

		cs := []Contact{}

		err := i.Get(url, &cs)
		if err != nil {
			return nil, err
		}

		for _, c := range cs {
			c.ParseDates()
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

func (c *Contact) ValidateEmail() error {
	// validate email
	if c.EMAIL_ADDRESS != "" {
		err := ValidateFormat(c.EMAIL_ADDRESS)
		if err != nil {
			return &types.ErrorString{fmt.Sprintf("invalid emailadress in Insightly: %s", c.EMAIL_ADDRESS)}
		}
	}

	return nil
}

func (c *Contact) ParseDates() {
	// parse DATE_CREATED_UTC to time.Time
	if c.DATE_CREATED_UTC != "" {
		t, _ := time.Parse("2006-01-02 15:04:05 +0000 UTC", c.DATE_CREATED_UTC+" +0000 UTC")
		//errortools.Fatal(err)
		c.DateCreated = &t
	}

	// parse DATE_UPDATED_UTC to time.Time
	if c.DATE_UPDATED_UTC != "" {
		t, _ := time.Parse("2006-01-02 15:04:05 +0000 UTC", c.DATE_UPDATED_UTC+" +0000 UTC")
		//errortools.Fatal(err)
		c.DateUpdated = &t
	}

	// parse DATE_OF_BIRTH to time.Time
	if c.DATE_OF_BIRTH != "" {
		t, _ := time.Parse("2006-01-02 15:04:05 +0000 UTC", c.DATE_OF_BIRTH+" +0000 UTC")
		//errortools.Fatal(err)
		c.DateOfBirth = &t
	}

	// parse LAST_ACTIVITY_DATE_UTC to time.Time
	if c.LAST_ACTIVITY_DATE_UTC != "" {
		t, _ := time.Parse("2006-01-02 15:04:05 +0000 UTC", c.LAST_ACTIVITY_DATE_UTC+" +0000 UTC")
		//errortools.Fatal(err)
		c.LastActivityDate = &t
	}

	// parse NEXT_ACTIVITY_DATE_UTC to time.Time
	if c.NEXT_ACTIVITY_DATE_UTC != "" {
		t, _ := time.Parse("2006-01-02 15:04:05 +0000 UTC", c.NEXT_ACTIVITY_DATE_UTC+" +0000 UTC")
		//errortools.Fatal(err)
		c.NextActivityDate = &t
	}

	// parse dates in DATES
	for _, d := range c.DATES {
		d.ParseDates()
	}
}
