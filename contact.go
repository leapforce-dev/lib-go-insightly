package insightly

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/getsentry/sentry-go"

	errortools "github.com/Leapforce-nl/go_errortools"
)

type Contact struct {
	CONTACT_ID       int           `json:"CONTACT_ID"`
	FIRST_NAME       string        `json:"FIRST_NAME"`
	LAST_NAME        string        `json:"LAST_NAME"`
	ORGANISATION_ID  int           `json:"ORGANISATION_ID"`
	EMAIL_ADDRESS    string        `json:"EMAIL_ADDRESS"`
	DATE_UPDATED_UTC string        `json:"DATE_UPDATED_UTC"`
	CUSTOMFIELDS     []CustomField `json:"CUSTOMFIELDS"`
	DateUpdated      time.Time
	Initials         string
	Gender           string
	Title            string
	Email            string
	IsMainContact    bool
	PushToEO         bool
}

/*
type iContacts struct {
	Contacts []Contact
}*/

func (c *Contact) GenderToGender(gender string) string {
	if strings.ToLower(gender) == "man" {
		return "Mannelijk"
	} else if strings.ToLower(gender) == "vrouw" {
		return "Vrouwelijk"
	}

	return "Onbekend"
}

func (c *Contact) GenderToTitle(gender string) string {
	if strings.ToLower(gender) == "man" {
		return "DHR"
	} else if strings.ToLower(gender) == "vrouw" {
		return "MEVR"
	}

	return ""
}

func (c *Contact) Updated(i *Insightly) bool {
	return c.DateUpdated.After(i.FromTimestamp)
}

// GetContacts returns all contacts updated after FromTimestamp date
//
func (i *Insightly) GetContacts() error {
	co, err := i.GetContactsInternal(true, "", "")

	i.Contacts = co

	//for _, c := range i.Contacts {
	//	fmt.Println(c)
	//}

	return err
}

// GetContactsFiltered returns all contacts fulfulling the specified filter
//
func (i *Insightly) GetContactsFiltered(fieldname string, fieldvalue string) ([]Contact, error) {
	return i.GetContactsInternal(false, fieldname, fieldvalue)
}

// GetContactsInternal is the generic function retrieving Contacts from Insightly
//
func (i *Insightly) GetContactsInternal(updatedAfterUTC bool, fieldName, fieldValue string) ([]Contact, error) {
	search := false
	updated := ""
	fieldname := ""
	fieldvalue := ""
	searchstring := ""
	if updatedAfterUTC {
		search = true
		from := i.FromTimestamp.Format("2006-01-02")
		updated = fmt.Sprintf("updated_after_utc=%s&", from)
	}
	if fieldName != "" && fieldValue != "" {
		search = true
		updated = fmt.Sprintf("field_name=%s&field_value=%s&", fieldName, fieldValue)
	}

	if search {
		searchstring = "/Search?" + updated + fieldname + fieldvalue
	} else {
		searchstring = "?"
	}

	urlStr := "%sContacts%sskip=%s&top=%s"
	skip := 0
	top := 500
	rowCount := 1
	isMainContactCount := 1
	pushToEOCount := 1

	contacts := []Contact{}

	for rowCount > 0 {
		url := fmt.Sprintf(urlStr, i.ApiUrl, searchstring, strconv.Itoa(skip), strconv.Itoa(top))
		//fmt.Printf(url)

		cs := []Contact{}

		err := i.Get(url, &cs)
		if err != nil {
			return nil, err
		}

		for _, c := range cs {
			// unmarshal custom fields
			for ii := range c.CUSTOMFIELDS {
				c.CUSTOMFIELDS[ii].UnmarshalValue()
			}

			// get Initials from custom field
			c.Initials = i.FindCustomFieldValue(c.CUSTOMFIELDS, customFieldNameInitials)

			// get Gender from custom field
			c.Gender = c.GenderToGender(i.FindCustomFieldValue(c.CUSTOMFIELDS, customFieldNameGender))

			// get Title from custom field
			c.Title = c.GenderToTitle(i.FindCustomFieldValue(c.CUSTOMFIELDS, customFieldNameGender))

			// parse DATE_UPDATED_UTC to time.Time
			t, err := time.Parse("2006-01-02 15:04:05 +0000 UTC", c.DATE_UPDATED_UTC+" +0000 UTC")
			errortools.Fatal(err)
			c.DateUpdated = t

			//fmt.Println("o.DATE_UPDATED_UTC", c.DATE_UPDATED_UTC, "o.DateUpdated", c.DateUpdated, "Now", time.Now(), "Diff", time.Now().Sub(c.DateUpdated))

			// validate email
			if c.EMAIL_ADDRESS != "" {
				err := ValidateFormat(c.EMAIL_ADDRESS)
				if err != nil {
					message := fmt.Sprintf("invalid emailadress in Insightly: %s", c.EMAIL_ADDRESS)
					fmt.Println(message)
					if i.IsLive {
						sentry.CaptureMessage(message)
					}
				} else {
					c.Email = c.EMAIL_ADDRESS
				}
			}

			c.IsMainContact = i.FindCustomFieldValueBool(c.CUSTOMFIELDS, customFieldNameMainContactPerson)
			if c.IsMainContact {
				isMainContactCount++
			}

			c.PushToEO = i.FindCustomFieldValueBool(c.CUSTOMFIELDS, customFieldNamePushToEO)
			if c.PushToEO {
				pushToEOCount++
			}

			contacts = append(contacts, c)
		}

		rowCount = len(cs)
		skip += top
	}

	//fmt.Println("isMainContactCount:", isMainContactCount)
	//fmt.Println("pushToEOCount (Contact):", pushToEOCount)

	if len(contacts) == 0 {
		contacts = nil
	}

	return contacts, nil
}
