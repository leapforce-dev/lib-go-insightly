package insightly

import (
	"fmt"
	"strconv"
	"time"

	errortools "github.com/leapforce-libraries/go_errortools"
)

// Project stores Project from Insightly
//
type Project struct {
	ProjectID           int           `json:"PROJECT_ID"`
	ProjectName         string        `json:"PROJECT_NAME"`
	Status              string        `json:"STATUS"`
	ProjectDetails      string        `json:"PROJECT_DETAILS"`
	StartedDate         string        `json:"STARTED_DATE"`
	CompletedDate       string        `json:"COMPLETED_DATE"`
	OpportunityID       int           `json:"OPPORTUNITY_ID"`
	CategoryID          int           `json:"CATEGORY_ID"`
	PipelineID          int           `json:"PIPELINE_ID"`
	StageID             int           `json:"STAGE_ID"`
	ImageURL            string        `json:"IMAGE_URL"`
	OwnerUserID         int           `json:"OWNER_USER_ID"`
	DateCreatedUTC      string        `json:"DATE_CREATED_UTC"`
	DateUpdatedUTC      string        `json:"DATE_UPDATED_UTC"`
	LastActivityDateUTC string        `json:"LAST_ACTIVITY_DATE_UTC"`
	NextActivityDateUTC string        `json:"NEXT_ACTIVITY_DATE_UTC"`
	CreatedUserID       int           `json:"CREATED_USER_ID"`
	ResponsibleUserID   int           `json:"RESPONSIBLE_USER_ID"`
	CustomFields        []CustomField `json:"CUSTOMFIELDS"`
	Tags                []Tag         `json:"TAGS"`
	StartedDateT        *time.Time
	CompletedDateT      *time.Time
	DateCreatedT        *time.Time
	DateUpdatedT        *time.Time
	LastActivityDateT   *time.Time
	NextActivityDateT   *time.Time
}

func (i *Insightly) GetProject(id int) (*Project, *errortools.Error) {
	urlStr := "%sProjects/%v"
	url := fmt.Sprintf(urlStr, apiURL, id)
	//fmt.Println(url)

	o := Project{}

	_, _, e := i.get(url, nil, &o)
	if e != nil {
		return nil, e
	}

	o.parseDates()

	return &o, nil
}

// GetProjects returns all projects
//
func (i *Insightly) GetProjects() ([]Project, *errortools.Error) {
	return i.GetProjectsInternal("")
}

// GetProjectsUpdatedAfter returns all projects updated after certain date
//
func (i *Insightly) GetProjectsUpdatedAfter(updatedAfter time.Time) ([]Project, *errortools.Error) {
	from := updatedAfter.Format("2006-01-02")
	searchFilter := fmt.Sprintf("updated_after_utc=%s&", from)
	return i.GetProjectsInternal(searchFilter)
}

// GetProjectsFiltered returns all projects fulfulling the specified filter
//
func (i *Insightly) GetProjectsFiltered(fieldname string, fieldvalue string) ([]Project, *errortools.Error) {
	searchFilter := fmt.Sprintf("field_name=%s&field_value=%s&", fieldname, fieldvalue)
	return i.GetProjectsInternal(searchFilter)
}

// GetProjectsInternal is the generic function retrieving projects from Insightly
//
func (i *Insightly) GetProjectsInternal(searchFilter string) ([]Project, *errortools.Error) {
	searchString := ""

	if searchFilter != "" {
		searchString = "/Search?" + searchFilter
	} else {
		searchString = "?"
	}

	urlStr := "%sProjects%sskip=%s&top=%s"
	skip := 0
	top := 500
	rowCount := top

	projects := []Project{}

	for rowCount >= top {
		url := fmt.Sprintf(urlStr, apiURL, searchString, strconv.Itoa(skip), strconv.Itoa(top))
		//fmt.Println(url)

		os := []Project{}

		_, _, err := i.get(url, nil, &os)
		if err != nil {
			return nil, err
		}

		for _, o := range os {
			o.parseDates()
			projects = append(projects, o)
		}

		rowCount = len(os)
		skip += top
	}

	if len(projects) == 0 {
		projects = nil
	}

	return projects, nil
}

func (o *Project) parseDates() {
	// parse STARTED_DATE to time.Time
	if o.StartedDate != "" {
		t, _ := time.Parse("2006-01-02 15:04:05 +0000 UTC", o.StartedDate+" +0000 UTC")
		//errortools.Fatal(err)
		o.StartedDateT = &t
	}

	// parse COMPLETED_DATE to time.Time
	if o.CompletedDate != "" {
		t, _ := time.Parse("2006-01-02 15:04:05 +0000 UTC", o.CompletedDate+" +0000 UTC")
		//errortools.Fatal(err)
		o.CompletedDateT = &t
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
