package insightly

import (
	"fmt"
	"strconv"
	"time"

	errortools "github.com/leapforce-libraries/go_errortools"
)

// Opportunity stores Opportunity from Insightly
//
type Opportunity struct {
	OpportunityID       int           `json:"OPPORTUNITY_ID"`
	OpportunityName     string        `json:"OPPORTUNITY_NAME"`
	OpportunityDetails  string        `json:"OPPORTUNITY_DETAILS"`
	OpportunityState    string        `json:"OPPORTUNITY_STATE"`
	ResponsibleUserID   int           `json:"RESPONSIBLE_USER_ID"`
	CategoryID          int           `json:"CATEGORY_ID"`
	ImageURL            string        `json:"IMAGE_URL"`
	BidCurrency         string        `json:"BID_CURRENCY"`
	BidAmount           float32       `json:"BID_AMOUNT"`
	BidType             string        `json:"BID_TYPE"`
	BidDuration         int           `json:"BID_DURATION"`
	ActualCloseDate     string        `json:"ACTUAL_CLOSE_DATE"`
	DateCreatedUTC      string        `json:"DATE_CREATED_UTC"`
	DateUpdatedUTC      string        `json:"DATE_UPDATED_UTC"`
	OpportunityValue    float32       `json:"OPPORTUNITY_VALUE"`
	Probability         int           `json:"PROBABILITY"`
	ForecastCloseDate   string        `json:"FORECAST_CLOSE_DATE"`
	OwnerUserID         int           `json:"OWNER_USER_ID"`
	LastActivityDateUTC string        `json:"LAST_ACTIVITY_DATE_UTC"`
	NextActivityDateUTC string        `json:"NEXT_ACTIVITY_DATE_UTC"`
	PipelineID          int           `json:"PIPELINE_ID"`
	StageID             int           `json:"STAGE_ID"`
	CreatedUserID       int           `json:"CREATED_USER_ID"`
	OrganisationID      int           `json:"ORGANISATION_ID"`
	CustomFields        []CustomField `json:"CUSTOMFIELDS"`
	Tags                []Tag         `json:"TAGS"`
	ActualCloseDateT    *time.Time
	DateCreatedT        *time.Time
	DateUpdatedT        *time.Time
	ForecastCloseDateT  *time.Time
	LastActivityDateT   *time.Time
	NextActivityDateT   *time.Time
}

func (i *Insightly) GetOpportunity(id int) (*Opportunity, *errortools.Error) {
	endpointStr := "%sOpportunities/%v"
	endpoint := fmt.Sprintf(endpointStr, apiURL, id)
	//fmt.Println(endpoint)

	o := Opportunity{}

	_, _, e := i.get(endpoint, nil, &o)
	if e != nil {
		return nil, e
	}

	o.parseDates()

	return &o, nil
}

// GetOpportunities returns all opportunities
//
func (i *Insightly) GetOpportunities() ([]Opportunity, *errortools.Error) {
	return i.GetOpportunitiesInternal("")
}

// GetOpportunitiesUpdatedAfter returns all opportunities updated after certain date
//
func (i *Insightly) GetOpportunitiesUpdatedAfter(updatedAfter time.Time) ([]Opportunity, *errortools.Error) {
	from := updatedAfter.Format("2006-01-02")
	searchFilter := fmt.Sprintf("updated_after_utc=%s&", from)
	return i.GetOpportunitiesInternal(searchFilter)
}

// GetOpportunitiesFiltered returns all opportunities fulfulling the specified filter
//
func (i *Insightly) GetOpportunitiesFiltered(fieldname string, fieldvalue string) ([]Opportunity, *errortools.Error) {
	searchFilter := fmt.Sprintf("field_name=%s&field_value=%s&", fieldname, fieldvalue)
	return i.GetOpportunitiesInternal(searchFilter)
}

// GetOpportunitiesInternal is the generic function retrieving opportunities from Insightly
//
func (i *Insightly) GetOpportunitiesInternal(searchFilter string) ([]Opportunity, *errortools.Error) {
	searchString := ""

	if searchFilter != "" {
		searchString = "/Search?" + searchFilter
	} else {
		searchString = "?"
	}

	endpointStr := "Opportunities%sskip=%s&top=%s"
	skip := 0
	top := 500
	rowCount := top

	opportunities := []Opportunity{}

	for rowCount >= top {
		endpoint := fmt.Sprintf(endpointStr, searchString, strconv.Itoa(skip), strconv.Itoa(top))
		//fmt.Println(endpoint)

		os := []Opportunity{}

		_, _, e := i.get(endpoint, nil, &os)
		if e != nil {
			return nil, e
		}

		for _, o := range os {
			o.parseDates()
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

func (o *Opportunity) parseDates() {
	// parse ACTUAL_CLOSE_DATE to time.Time
	if o.ActualCloseDate != "" {
		t, _ := time.Parse("2006-01-02 15:04:05 +0000 UTC", o.ActualCloseDate+" +0000 UTC")
		//errortools.Fatal(err)
		o.ActualCloseDateT = &t
	}

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

	// parse FORECAST_CLOSE_DATE to time.Time
	if o.ForecastCloseDate != "" {
		t, _ := time.Parse("2006-01-02 15:04:05 +0000 UTC", o.ForecastCloseDate+" +0000 UTC")
		//errortools.Fatal(err)
		o.ForecastCloseDateT = &t
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
}
