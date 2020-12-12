package insightly

import (
	"fmt"
	"strconv"
	"time"

	errortools "github.com/leapforce-libraries/go_errortools"
)

// Organisation stores Organisation from Insightly
//
type Organisation struct {
	OrganisationID         int           `json:"ORGANISATION_ID"`
	OrganisationName       string        `json:"ORGANISATION_NAME"`
	Background             string        `json:"BACKGROUND"`
	ImageURL               string        `json:"IMAGE_URL"`
	OwnerUserID            int           `json:"OWNER_USER_ID"`
	DateCreatedUTC         string        `json:"DATE_CREATED_UTC"`
	DateUpdatedUTC         string        `json:"DATE_UPDATED_UTC"`
	LastActivityDateUTC    string        `json:"LAST_ACTIVITY_DATE_UTC"`
	NextActivityDateUTC    string        `json:"NEXT_ACTIVITY_DATE_UTC"`
	CreatedUserID          int           `json:"CREATED_USER_ID"`
	Phone                  string        `json:"PHONE"`
	PhoneFax               string        `json:"PHONE_FAX"`
	Website                string        `json:"WEBSITE"`
	AddressBillingStreet   string        `json:"ADDRESS_BILLING_STREET"`
	AddressBillingCity     string        `json:"ADDRESS_BILLING_CITY"`
	AddressBillingState    string        `json:"ADDRESS_BILLING_STATE"`
	AddressBillingCountry  string        `json:"ADDRESS_BILLING_COUNTRY"`
	AddressBillingPostcode string        `json:"ADDRESS_BILLING_POSTCODE"`
	AddressShipStreet      string        `json:"ADDRESS_SHIP_STREET"`
	AddressShipCity        string        `json:"ADDRESS_SHIP_CITY"`
	AddressShipState       string        `json:"ADDRESS_SHIP_STATE"`
	AddressShipCountry     string        `json:"ADDRESS_SHIP_COUNTRY"`
	AddressShipPostcode    string        `json:"ADDRESS_SHIP_POSTCODE"`
	SocialLinkedin         string        `json:"SOCIAL_LINKEDIN"`
	SocialFacebook         string        `json:"SOCIAL_FACEBOOK"`
	SocialTwitter          string        `json:"SOCIAL_TWITTER"`
	CustomFields           []CustomField `json:"CUSTOMFIELDS"`
	Tags                   []Tag         `json:"TAGS"`
	Dates                  []Date        `json:"DATES"`
	EmailDomains           []EmailDomain `json:"EMAILDOMAINS"`
	DateCreatedT           *time.Time
	DateUpdatedT           *time.Time
	LastActivityDateT      *time.Time
	NextActivityDateT      *time.Time
}

func (i *Insightly) GetOrganisation(id int) (*Organisation, *errortools.Error) {
	endpointStr := "%sOrganisations/%v"
	endpoint := fmt.Sprintf(endpointStr, apiURL, id)
	//fmt.Println(endpoint)

	o := Organisation{}

	_, _, e := i.get(endpoint, nil, &o)
	if e != nil {
		return nil, e
	}

	o.parseDates()

	return &o, nil
}

// GetOrganisations returns all organisations
//
func (i *Insightly) GetOrganisations() ([]Organisation, *errortools.Error) {
	return i.GetOrganisationsInternal("")
}

// GetOrganisationsUpdatedAfter returns all organisations updated after certain date
//
func (i *Insightly) GetOrganisationsUpdatedAfter(updatedAfter time.Time) ([]Organisation, *errortools.Error) {
	from := updatedAfter.Format("2006-01-02")
	searchFilter := fmt.Sprintf("updated_after_utc=%s&", from)
	return i.GetOrganisationsInternal(searchFilter)
}

// GetOrganisationsFiltered returns all organisations fulfulling the specified filter
//
func (i *Insightly) GetOrganisationsFiltered(fieldname string, fieldvalue string) ([]Organisation, *errortools.Error) {
	searchFilter := fmt.Sprintf("field_name=%s&field_value=%s&", fieldname, fieldvalue)
	return i.GetOrganisationsInternal(searchFilter)
}

// GetOrganisationsInternal is the generic function retrieving organisations from Insightly
//
func (i *Insightly) GetOrganisationsInternal(searchFilter string) ([]Organisation, *errortools.Error) {
	searchString := ""

	if searchFilter != "" {
		searchString = "/Search?" + searchFilter
	} else {
		searchString = "?"
	}

	endpointStr := "Organisations%sskip=%s&top=%s"
	skip := 0
	top := 500
	rowCount := top

	organisations := []Organisation{}

	for rowCount >= top {
		endpoint := fmt.Sprintf(endpointStr, searchString, strconv.Itoa(skip), strconv.Itoa(top))
		//fmt.Println(endpoint)

		os := []Organisation{}

		_, _, e := i.get(endpoint, nil, &os)
		if e != nil {
			return nil, e
		}

		for _, o := range os {
			o.parseDates()
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

func (o *Organisation) parseDates() {
	// parse DATE_CREATED_UTC to time.Time
	if o.DateCreatedUTC != "" {
		t, _ := time.Parse("2006-01-02 15:04:05 +0000 UTC", o.DateCreatedUTC+" +0000 UTC")
		//errortools.Fatal(err)
		o.DateCreatedT = &t
	}

	// parse DATE_UPDATED_UTC to time.Time
	if o.DateUpdatedUTC != "" {
		t, _ := time.Parse("2006-01-02 15:04:05 +0000 UTC", o.DateUpdatedUTC+" +0000 UTC")
		//errortools.Fatal(err)
		o.DateUpdatedT = &t
	}

	// parse LAST_ACTIVITY_DATE_UTC to time.Time
	if o.LastActivityDateUTC != "" {
		t, _ := time.Parse("2006-01-02 15:04:05 +0000 UTC", o.LastActivityDateUTC+" +0000 UTC")
		//errortools.Fatal(err)
		o.LastActivityDateT = &t
	}

	// parse NEXT_ACTIVITY_DATE_UTC to time.Time
	if o.NextActivityDateUTC != "" {
		t, _ := time.Parse("2006-01-02 15:04:05 +0000 UTC", o.NextActivityDateUTC+" +0000 UTC")
		//errortools.Fatal(err)
		o.NextActivityDateT = &t
	}

	// parse dates in DATES
	for _, d := range o.Dates {
		d.parseDates()
	}
}

func (i *Insightly) UpdateOrganisationRemoveCustomField(organisationID int, customFieldName string) *errortools.Error {
	endpoint := "Organisations"

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

	_, _, e := i.put(endpoint, o1, nil)
	if e != nil {
		return e
	}

	return nil
}
