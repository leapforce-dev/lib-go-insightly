package insightly

import (
	"fmt"
	"strconv"
	"time"

	errortools "github.com/leapforce-libraries/go_errortools"
)

// Lead stores Lead from Insightly
//
type Lead struct {
	LeadID                  int           `json:"LEAD_ID"`
	Salutation              string        `json:"SALUTATION"`
	FirstName               string        `json:"FIRST_NAME"`
	LastName                string        `json:"LAST_NAME"`
	LeadSourceID            int           `json:"LEAD_SOURCE_ID"`
	LeadStatusID            int           `json:"LEAD_STATUS_ID"`
	Title                   string        `json:"TITLE"`
	Converted               bool          `json:"CONVERTED"`
	ConvertedContactID      int           `json:"CONVERTED_CONTACT_ID"`
	ConvertedDateUTC        string        `json:"CONVERTED_DATE_UTC"`
	ConvertedOpportunityID  int           `json:"CONVERTED_OPPORTUNITY_ID"`
	ConvertedOrganisationID int           `json:"CONVERTED_ORGANISATION_ID"`
	DateCreateUTC           string        `json:"DATE_CREATED_UTC"`
	DateUpdatedUTC          string        `json:"DATE_UPDATED_UTC"`
	Email                   string        `json:"EMAIL"`
	EmployeeCount           int           `json:"EMPLOYEE_COUNT"`
	Fax                     string        `json:"FAX"`
	Industry                string        `json:"INDUSTRY"`
	LeadDescription         string        `json:"LEAD_DESCRIPTION"`
	LeadRating              int           `json:"LEAD_RATING"`
	Mobile                  string        `json:"MOBILE"`
	OwnerUserID             int           `json:"OWNER_USER_ID"`
	Phone                   string        `json:"PHONE"`
	ResponsibleUserID       int           `json:"RESPONSIBLE_USER_ID"`
	Website                 string        `json:"WEBSITE"`
	AddressStreet           string        `json:"ADDRESS_STREET"`
	AddressCity             string        `json:"ADDRESS_CITY"`
	AddressState            string        `json:"ADDRESS_STATE"`
	AddressPostcode         string        `json:"ADDRESS_POSTCODE"`
	AddressCountry          string        `json:"ADDRESS_COUNTRY"`
	LastActivityDateUTC     string        `json:"LAST_ACTIVITY_DATE_UTC"`
	NextActivityDateUTC     string        `json:"NEXT_ACTIVITY_DATE_UTC"`
	OrganisationName        string        `json:"ORGANISATION_NAME"`
	CreatedUserID           int           `json:"CREATED_USER_ID"`
	ImageURL                string        `json:"IMAGE_URL"`
	EmailOptedOut           bool          `json:"EMAIL_OPTED_OUT"`
	CustomFields            []CustomField `json:"CUSTOMFIELDS"`
	Tags                    []Tag         `json:"TAGS"`
	ConvertedDateT          *time.Time
	DateCreatedT            *time.Time
	DateUpdatedT            *time.Time
	LastActivityDateT       *time.Time
	NextActivityDateT       *time.Time
}

// GetLeads returns all leads
//
func (i *Insightly) GetLeads() ([]Lead, *errortools.Error) {
	return i.GetLeadsInternal("")
}

// GetLeadsUpdatedAfter returns all leads updated after certain date
//
func (i *Insightly) GetLeadsUpdatedAfter(updatedAfter time.Time) ([]Lead, *errortools.Error) {
	from := updatedAfter.Format("2006-01-02")
	searchFilter := fmt.Sprintf("updated_after_utc=%s&", from)
	return i.GetLeadsInternal(searchFilter)
}

// GetLeadsFiltered returns all leads fulfulling the specified filter
//
func (i *Insightly) GetLeadsFiltered(fieldname string, fieldvalue string) ([]Lead, *errortools.Error) {
	searchFilter := fmt.Sprintf("field_name=%s&field_value=%s&", fieldname, fieldvalue)
	return i.GetLeadsInternal(searchFilter)
}

// GetLeadsInternal is the generic function retrieving leads from Insightly
//
func (i *Insightly) GetLeadsInternal(searchFilter string) ([]Lead, *errortools.Error) {
	searchString := ""

	if searchFilter != "" {
		searchString = "/Search?" + searchFilter
	} else {
		searchString = "?"
	}

	urlStr := "%sLeads%sskip=%s&top=%s"
	skip := 0
	top := 500
	rowCount := top

	leads := []Lead{}

	for rowCount >= top {
		url := fmt.Sprintf(urlStr, apiURL, searchString, strconv.Itoa(skip), strconv.Itoa(top))
		//fmt.Println(url)

		ls := []Lead{}

		_, _, e := i.get(url, nil, &ls)
		if e != nil {
			return nil, e
		}

		for _, l := range ls {
			l.parseDates()
			leads = append(leads, l)
		}

		rowCount = len(ls)
		skip += top
	}

	if len(leads) == 0 {
		leads = nil
	}

	return leads, nil
}

func (l *Lead) parseDates() {
	// parse CONVERTED_DATE_UTC to time.Time
	if l.ConvertedDateUTC != "" {
		t, _ := time.Parse("2006-01-02 15:04:05 +0000 UTC", l.ConvertedDateUTC+" +0000 UTC")
		//errortools.Fatal(err)
		l.ConvertedDateT = &t
	}

	// parse DATE_CREATED_UTC to time.Time
	if l.DateCreateUTC != "" {
		t, _ := time.Parse("2006-01-02 15:04:05 +0000 UTC", l.DateCreateUTC+" +0000 UTC")
		//errortools.Fatal(err)
		l.DateCreatedT = &t
	}

	// parse DATE_UPDATED_UTC to time.Time
	if l.DateUpdatedUTC != "" {
		t, _ := time.Parse("2006-01-02 15:04:05 +0000 UTC", l.DateUpdatedUTC+" +0000 UTC")
		//errortools.Fatal(err)
		l.DateUpdatedT = &t
	}

	// parse LAST_ACTIVITY_DATE_UTC to time.Time
	if l.LastActivityDateUTC != "" {
		t, _ := time.Parse("2006-01-02 15:04:05 +0000 UTC", l.LastActivityDateUTC+" +0000 UTC")
		//errortools.Fatal(err)
		l.LastActivityDateT = &t
	}

	// parse NEXT_ACTIVITY_DATE_UTC to time.Time
	if l.NextActivityDateUTC != "" {
		t, _ := time.Parse("2006-01-02 15:04:05 +0000 UTC", l.NextActivityDateUTC+" +0000 UTC")
		//errortools.Fatal(err)
		l.NextActivityDateT = &t
	}
}
