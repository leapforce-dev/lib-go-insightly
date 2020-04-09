package insightly

import (
	"fmt"
	"strconv"
	"time"
)

// Lead stores Lead from Insightly
//
type Lead struct {
	LEAD_ID                   int           `json:"LEAD_ID"`
	SALUTATION                string        `json:"SALUTATION"`
	FIRST_NAME                string        `json:"FIRST_NAME"`
	LAST_NAME                 string        `json:"LAST_NAME"`
	LEAD_SOURCE_ID            int           `json:"LEAD_SOURCE_ID"`
	LEAD_STATUS_ID            int           `json:"LEAD_STATUS_ID"`
	TITLE                     string        `json:"TITLE"`
	CONVERTED                 bool          `json:"CONVERTED"`
	CONVERTED_CONTACT_ID      int           `json:"CONVERTED_CONTACT_ID"`
	CONVERTED_DATE_UTC        string        `json:"CONVERTED_DATE_UTC"`
	CONVERTED_OPPORTUNITY_ID  int           `json:"CONVERTED_OPPORTUNITY_ID"`
	CONVERTED_ORGANISATION_ID int           `json:"CONVERTED_ORGANISATION_ID"`
	DATE_CREATED_UTC          string        `json:"DATE_CREATED_UTC"`
	DATE_UPDATED_UTC          string        `json:"DATE_UPDATED_UTC"`
	EMAIL                     string        `json:"EMAIL"`
	EMPLOYEE_COUNT            int           `json:"EMPLOYEE_COUNT"`
	FAX                       string        `json:"FAX"`
	INDUSTRY                  string        `json:"INDUSTRY"`
	LEAD_DESCRIPTION          string        `json:"LEAD_DESCRIPTION"`
	LEAD_RATING               int           `json:"LEAD_RATING"`
	MOBILE                    string        `json:"MOBILE"`
	OWNER_USER_ID             int           `json:"OWNER_USER_ID"`
	PHONE                     string        `json:"PHONE"`
	RESPONSIBLE_USER_ID       int           `json:"RESPONSIBLE_USER_ID"`
	WEBSITE                   string        `json:"WEBSITE"`
	ADDRESS_STREET            string        `json:"ADDRESS_STREET"`
	ADDRESS_CITY              string        `json:"ADDRESS_CITY"`
	ADDRESS_STATE             string        `json:"ADDRESS_STATE"`
	ADDRESS_POSTCODE          string        `json:"ADDRESS_POSTCODE"`
	ADDRESS_COUNTRY           string        `json:"ADDRESS_COUNTRY"`
	LAST_ACTIVITY_DATE_UTC    string        `json:"LAST_ACTIVITY_DATE_UTC"`
	NEXT_ACTIVITY_DATE_UTC    string        `json:"NEXT_ACTIVITY_DATE_UTC"`
	ORGANISATION_NAME         string        `json:"ORGANISATION_NAME"`
	CREATED_USER_ID           int           `json:"CREATED_USER_ID"`
	IMAGE_URL                 string        `json:"IMAGE_URL"`
	EMAIL_OPTED_OUT           bool          `json:"EMAIL_OPTED_OUT"`
	CUSTOMFIELDS              []CustomField `json:"CUSTOMFIELDS"`
	TAGS                      []Tag         `json:"TAGS"`
	ConvertedDate             *time.Time
	DateCreated               *time.Time
	DateUpdated               *time.Time
	LastActivityDate          *time.Time
	NextActivityDate          *time.Time
}

// GetLeads returns all leads
//
func (i *Insightly) GetLeads() ([]Lead, error) {
	return i.GetLeadsInternal("")
}

// GetLeadsUpdatedAfter returns all leads updated after certain date
//
func (i *Insightly) GetLeadsUpdatedAfter(updatedAfter time.Time) ([]Lead, error) {
	from := updatedAfter.Format("2006-01-02")
	searchFilter := fmt.Sprintf("updated_after_utc=%s&", from)
	return i.GetLeadsInternal(searchFilter)
}

// GetLeadsFiltered returns all leads fulfulling the specified filter
//
func (i *Insightly) GetLeadsFiltered(fieldname string, fieldvalue string) ([]Lead, error) {
	searchFilter := fmt.Sprintf("field_name=%s&field_value=%s&", fieldname, fieldvalue)
	return i.GetLeadsInternal(searchFilter)
}

// GetLeadsInternal is the generic function retrieving leads from Insightly
//
func (i *Insightly) GetLeadsInternal(searchFilter string) ([]Lead, error) {
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
		url := fmt.Sprintf(urlStr, i.apiURL, searchString, strconv.Itoa(skip), strconv.Itoa(top))
		//fmt.Println(url)

		ls := []Lead{}

		err := i.Get(url, &ls)
		if err != nil {
			return nil, err
		}

		for _, l := range ls {
			l.ParseDates()
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

func (l *Lead) ParseDates() {
	// parse CONVERTED_DATE_UTC to time.Time
	if l.CONVERTED_DATE_UTC != "" {
		t, _ := time.Parse("2006-01-02 15:04:05 +0000 UTC", l.CONVERTED_DATE_UTC+" +0000 UTC")
		//errortools.Fatal(err)
		l.ConvertedDate = &t
	}

	// parse DATE_CREATED_UTC to time.Time
	if l.DATE_CREATED_UTC != "" {
		t, _ := time.Parse("2006-01-02 15:04:05 +0000 UTC", l.DATE_CREATED_UTC+" +0000 UTC")
		//errortools.Fatal(err)
		l.DateCreated = &t
	}

	// parse DATE_UPDATED_UTC to time.Time
	if l.DATE_UPDATED_UTC != "" {
		t, _ := time.Parse("2006-01-02 15:04:05 +0000 UTC", l.DATE_UPDATED_UTC+" +0000 UTC")
		//errortools.Fatal(err)
		l.DateUpdated = &t
	}

	// parse LAST_ACTIVITY_DATE_UTC to time.Time
	if l.LAST_ACTIVITY_DATE_UTC != "" {
		t, _ := time.Parse("2006-01-02 15:04:05 +0000 UTC", l.LAST_ACTIVITY_DATE_UTC+" +0000 UTC")
		//errortools.Fatal(err)
		l.LastActivityDate = &t
	}

	// parse NEXT_ACTIVITY_DATE_UTC to time.Time
	if l.NEXT_ACTIVITY_DATE_UTC != "" {
		t, _ := time.Parse("2006-01-02 15:04:05 +0000 UTC", l.NEXT_ACTIVITY_DATE_UTC+" +0000 UTC")
		//errortools.Fatal(err)
		l.NextActivityDate = &t
	}
}
