package insightly

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

// Organisation stores Organisation from Insightly
//
type Organisation struct {
	ORGANISATION_ID          int           `json:"ORGANISATION_ID"`
	ORGANISATION_NAME        string        `json:"ORGANISATION_NAME"`
	BACKGROUND               string        `json:"BACKGROUND"`
	IMAGE_URL                string        `json:"IMAGE_URL"`
	OWNER_USER_ID            int           `json:"OWNER_USER_ID"`
	DATE_CREATED_UTC         string        `json:"DATE_CREATED_UTC"`
	DATE_UPDATED_UTC         string        `json:"DATE_UPDATED_UTC"`
	LAST_ACTIVITY_DATE_UTC   string        `json:"LAST_ACTIVITY_DATE_UTC"`
	NEXT_ACTIVITY_DATE_UTC   string        `json:"NEXT_ACTIVITY_DATE_UTC"`
	CREATED_USER_ID          int           `json:"CREATED_USER_ID"`
	PHONE                    string        `json:"PHONE"`
	PHONE_FAX                string        `json:"PHONE_FAX"`
	WEBSITE                  string        `json:"WEBSITE"`
	ADDRESS_BILLING_STREET   string        `json:"ADDRESS_BILLING_STREET"`
	ADDRESS_BILLING_CITY     string        `json:"ADDRESS_BILLING_CITY"`
	ADDRESS_BILLING_STATE    string        `json:"ADDRESS_BILLING_STATE"`
	ADDRESS_BILLING_COUNTRY  string        `json:"ADDRESS_BILLING_COUNTRY"`
	ADDRESS_BILLING_POSTCODE string        `json:"ADDRESS_BILLING_POSTCODE"`
	ADDRESS_SHIP_STREET      string        `json:"ADDRESS_SHIP_STREET"`
	ADDRESS_SHIP_CITY        string        `json:"ADDRESS_SHIP_CITY"`
	ADDRESS_SHIP_STATE       string        `json:"ADDRESS_SHIP_STATE"`
	ADDRESS_SHIP_COUNTRY     string        `json:"ADDRESS_SHIP_COUNTRY"`
	ADDRESS_SHIP_POSTCODE    string        `json:"ADDRESS_SHIP_POSTCODE"`
	SOCIAL_LINKEDIN          string        `json:"SOCIAL_LINKEDIN"`
	SOCIAL_FACEBOOK          string        `json:"SOCIAL_FACEBOOK"`
	SOCIAL_TWITTER           string        `json:"SOCIAL_TWITTER"`
	CUSTOMFIELDS             []CustomField `json:"CUSTOMFIELDS"`
	TAGS                     []Tag         `json:"TAGS"`
	DATES                    []Date        `json:"DATES"`
	EMAILDOMAINS             []EmailDomain `json:"EMAILDOMAINS"`
	DateCreated              *time.Time
	DateUpdated              *time.Time
	LastActivityDate         *time.Time
	NextActivityDate         *time.Time
}

func (i *Insightly) GetOrganisation(id int) (*Organisation, error) {
	urlStr := "%sOrganisations/%v"
	url := fmt.Sprintf(urlStr, i.apiURL, id)
	//fmt.Println(url)

	o := Organisation{}

	err := i.Get(url, &o)
	if err != nil {
		return nil, err
	}

	o.ParseDates()

	return &o, nil
}

// GetOrganisations returns all organisations
//
func (i *Insightly) GetOrganisations() ([]Organisation, error) {
	return i.GetOrganisationsInternal("")
}

// GetOrganisationsUpdatedAfter returns all organisations updated after certain date
//
func (i *Insightly) GetOrganisationsUpdatedAfter(updatedAfter time.Time) ([]Organisation, error) {
	from := updatedAfter.Format("2006-01-02")
	searchFilter := fmt.Sprintf("updated_after_utc=%s&", from)
	return i.GetOrganisationsInternal(searchFilter)
}

// GetOrganisationsFiltered returns all organisations fulfulling the specified filter
//
func (i *Insightly) GetOrganisationsFiltered(fieldname string, fieldvalue string) ([]Organisation, error) {
	searchFilter := fmt.Sprintf("field_name=%s&field_value=%s&", fieldname, fieldvalue)
	return i.GetOrganisationsInternal(searchFilter)
}

// GetOrganisationsInternal is the generic function retrieving organisations from Insightly
//
func (i *Insightly) GetOrganisationsInternal(searchFilter string) ([]Organisation, error) {
	searchString := ""

	if searchFilter != "" {
		searchString = "/Search?" + searchFilter
	} else {
		searchString = "?"
	}

	urlStr := "%sOrganisations%sskip=%s&top=%s"
	skip := 0
	top := 500
	rowCount := top

	organisations := []Organisation{}

	for rowCount >= top {
		url := fmt.Sprintf(urlStr, i.apiURL, searchString, strconv.Itoa(skip), strconv.Itoa(top))
		//fmt.Println(url)

		os := []Organisation{}

		err := i.Get(url, &os)
		if err != nil {
			return nil, err
		}

		for _, o := range os {
			o.ParseDates()
			organisations = append(organisations, o)
		}

		rowCount = len(os)
		skip += top
	}

	if len(organisations) == 0 {
		organisations = nil
	}

	return organisations, nil
}

func (o *Organisation) ParseDates() {
	// parse DATE_CREATED_UTC to time.Time
	if o.DATE_CREATED_UTC != "" {
		t, _ := time.Parse("2006-01-02 15:04:05 +0000 UTC", o.DATE_CREATED_UTC+" +0000 UTC")
		//errortools.Fatal(err)
		o.DateCreated = &t
	}

	// parse DATE_UPDATED_UTC to time.Time
	if o.DATE_UPDATED_UTC != "" {
		t, _ := time.Parse("2006-01-02 15:04:05 +0000 UTC", o.DATE_UPDATED_UTC+" +0000 UTC")
		//errortools.Fatal(err)
		o.DateUpdated = &t
	}

	// parse LAST_ACTIVITY_DATE_UTC to time.Time
	if o.LAST_ACTIVITY_DATE_UTC != "" {
		t, _ := time.Parse("2006-01-02 15:04:05 +0000 UTC", o.LAST_ACTIVITY_DATE_UTC+" +0000 UTC")
		//errortools.Fatal(err)
		o.LastActivityDate = &t
	}

	// parse NEXT_ACTIVITY_DATE_UTC to time.Time
	if o.NEXT_ACTIVITY_DATE_UTC != "" {
		t, _ := time.Parse("2006-01-02 15:04:05 +0000 UTC", o.NEXT_ACTIVITY_DATE_UTC+" +0000 UTC")
		//errortools.Fatal(err)
		o.NextActivityDate = &t
	}

	// parse dates in DATES
	for _, d := range o.DATES {
		d.ParseDates()
	}
}

func (i *Insightly) UpdateOrganisationRemoveCustomField(organisationID int, customFieldName string) error {
	urlStr := "%sOrganisations"
	url := fmt.Sprintf(urlStr, i.apiURL)

	type CustomFieldDelete struct {
		FIELD_NAME      string
		CUSTOM_FIELD_ID string
	}

	type OrganisationID struct {
		ORGANISATION_ID int
		CUSTOMFIELDS    []CustomFieldDelete
	}

	o1 := OrganisationID{}
	o1.ORGANISATION_ID = organisationID
	o1.CUSTOMFIELDS = make([]CustomFieldDelete, 1)
	o1.CUSTOMFIELDS[0] = CustomFieldDelete{customFieldName, customFieldName}

	b, err := json.Marshal(o1)
	if err != nil {
		fmt.Println("error:", err)
	}

	err = i.Put(url, b)
	if err != nil {
		fmt.Println("ERROR in UpdateOrganisationRemovePushToEO:", err)
		fmt.Println("url:", urlStr)
		return err
	}

	return nil
}
