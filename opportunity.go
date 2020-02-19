package insightly

import (
	"fmt"
	"strconv"
	"time"
)

// Opportunity stores Opportunity from Insightly
//
type Opportunity struct {
	OPPORTUNITY_ID         int           `json:"OPPORTUNITY_ID"`
	OPPORTUNITY_NAME       string        `json:"OPPORTUNITY_NAME"`
	OPPORTUNITY_DETAILS    string        `json:"OPPORTUNITY_DETAILS"`
	OPPORTUNITY_STATE      string        `json:"OPPORTUNITY_STATE"`
	RESPONSIBLE_USER_ID    int           `json:"RESPONSIBLE_USER_ID"`
	CATEGORY_ID            int           `json:"CATEGORY_ID"`
	IMAGE_URL              string        `json:"IMAGE_URL"`
	BID_CURRENCY           string        `json:"BID_CURRENCY"`
	BID_AMOUNT             string        `json:"BID_AMOUNT"`
	BID_TYPE               string        `json:"BID_TYPE"`
	BID_DURATION           int           `json:"BID_DURATION"`
	ACTUAL_CLOSE_DATE      string        `json:"ACTUAL_CLOSE_DATE"`
	DATE_CREATED_UTC       string        `json:"DATE_CREATED_UTC"`
	DATE_UPDATED_UTC       string        `json:"DATE_UPDATED_UTC"`
	OPPORTUNITY_VALUE      string        `json:"OPPORTUNITY_VALUE"`
	PROBABILITY            int           `json:"PROBABILITY"`
	FORECAST_CLOSE_DATE    string        `json:"FORECAST_CLOSE_DATE"`
	OWNER_USER_ID          int           `json:"OWNER_USER_ID"`
	LAST_ACTIVITY_DATE_UTC string        `json:"LAST_ACTIVITY_DATE_UTC"`
	NEXT_ACTIVITY_DATE_UTC string        `json:"NEXT_ACTIVITY_DATE_UTC"`
	PIPELINE_ID            int           `json:"PIPELINE_ID"`
	STAGE_ID               int           `json:"STAGE_ID"`
	CREATED_USER_ID        int           `json:"CREATED_USER_ID"`
	ORGANISATION_ID        int           `json:"ORGANISATION_ID"`
	CUSTOMFIELDS           []CustomField `json:"CUSTOMFIELDS"`
	TAGS                   []Tag         `json:"TAGS"`
	ActualCloseDate        *time.Time
	DateCreated            *time.Time
	DateUpdated            *time.Time
	ForecastCloseDate      *time.Time
	LastActivityDate       *time.Time
	NextActivityDate       *time.Time
}

func (i *Insightly) GetOpportunity(id int) (*Opportunity, error) {
	urlStr := "%sOpportunities/%v"
	url := fmt.Sprintf(urlStr, i.apiURL, id)
	//fmt.Println(url)

	o := Opportunity{}

	err := i.Get(url, &o)
	if err != nil {
		return nil, err
	}

	o.ParseDates()

	return &o, nil
}

// GetOpportunities returns all opportunities
//
func (i *Insightly) GetOpportunities() ([]Opportunity, error) {
	return i.GetOpportunitiesInternal("")
}

// GetOpportunitiesUpdatedAfter returns all opportunities updated after certain date
//
func (i *Insightly) GetOpportunitiesUpdatedAfter(updatedAfter time.Time) ([]Opportunity, error) {
	from := updatedAfter.Format("2006-01-02")
	searchFilter := fmt.Sprintf("updated_after_utc=%s&", from)
	return i.GetOpportunitiesInternal(searchFilter)
}

// GetOpportunitiesFiltered returns all opportunities fulfulling the specified filter
//
func (i *Insightly) GetOpportunitiesFiltered(fieldname string, fieldvalue string) ([]Opportunity, error) {
	searchFilter := fmt.Sprintf("field_name=%s&field_value=%s&", fieldname, fieldvalue)
	return i.GetOpportunitiesInternal(searchFilter)
}

// GetOpportunitiesInternal is the generic function retrieving opportunities from Insightly
//
func (i *Insightly) GetOpportunitiesInternal(searchFilter string) ([]Opportunity, error) {
	searchString := ""

	if searchFilter != "" {
		searchString = "/Search?" + searchFilter
	} else {
		searchString = "?"
	}

	urlStr := "%sOpportunities%sskip=%s&top=%s"
	skip := 0
	top := 500
	rowCount := 1

	opportunities := []Opportunity{}

	for rowCount > 0 {
		url := fmt.Sprintf(urlStr, i.apiURL, searchString, strconv.Itoa(skip), strconv.Itoa(top))
		//fmt.Println(url)

		os := []Opportunity{}

		err := i.Get(url, &os)
		if err != nil {
			return nil, err
		}

		for _, o := range os {
			o.ParseDates()
			opportunities = append(opportunities, o)
		}

		rowCount = len(os)
		skip += top
	}

	if len(opportunities) == 0 {
		opportunities = nil
	}

	return opportunities, nil
}

func (o *Opportunity) ParseDates() {
	// parse ACTUAL_CLOSE_DATE to time.Time
	if o.ACTUAL_CLOSE_DATE != "" {
		t, _ := time.Parse("2006-01-02 15:04:05 +0000 UTC", o.ACTUAL_CLOSE_DATE+" +0000 UTC")
		//errortools.Fatal(err)
		o.ActualCloseDate = &t
	}

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

	// parse FORECAST_CLOSE_DATE to time.Time
	if o.FORECAST_CLOSE_DATE != "" {
		t, _ := time.Parse("2006-01-02 15:04:05 +0000 UTC", o.FORECAST_CLOSE_DATE+" +0000 UTC")
		//errortools.Fatal(err)
		o.ForecastCloseDate = &t
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
}
