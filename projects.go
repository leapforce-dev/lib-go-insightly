package insightly

import (
	"fmt"
	"strconv"
	"time"
)

// Project stores Project from Insightly
//
type Project struct {
	PROJECT_ID             int           `json:"PROJECT_ID"`
	PROJECT_NAME           string        `json:"PROJECT_NAME"`
	STATUS                 string        `json:"STATUS"`
	PROJECT_DETAILS        string        `json:"PROJECT_DETAILS"`
	STARTED_DATE           string        `json:"STARTED_DATE"`
	COMPLETED_DATE         string        `json:"COMPLETED_DATE"`
	OPPORTUNITY_ID         int           `json:"OPPORTUNITY_ID"`
	CATEGORY_ID            int           `json:"CATEGORY_ID"`
	PIPELINE_ID            int           `json:"PIPELINE_ID"`
	STAGE_ID               int           `json:"STAGE_ID"`
	IMAGE_URL              string        `json:"IMAGE_URL"`
	OWNER_USER_ID          int           `json:"OWNER_USER_ID"`
	DATE_CREATED_UTC       string        `json:"DATE_CREATED_UTC"`
	DATE_UPDATED_UTC       string        `json:"DATE_UPDATED_UTC"`
	LAST_ACTIVITY_DATE_UTC string        `json:"LAST_ACTIVITY_DATE_UTC"`
	NEXT_ACTIVITY_DATE_UTC string        `json:"NEXT_ACTIVITY_DATE_UTC"`
	CREATED_USER_ID        int           `json:"CREATED_USER_ID"`
	RESPONSIBLE_USER_ID    int           `json:"RESPONSIBLE_USER_ID"`
	CUSTOMFIELDS           []CustomField `json:"CUSTOMFIELDS"`
	TAGS                   []Tag         `json:"TAGS"`
	StartedDate            *time.Time
	CompletedDate          *time.Time
	DateCreated            *time.Time
	DateUpdated            *time.Time
	LastActivityDate       *time.Time
	NextActivityDate       *time.Time
}

func (i *Insightly) GetProject(id int) (*Project, error) {
	urlStr := "%sProjects/%v"
	url := fmt.Sprintf(urlStr, i.apiURL, id)
	//fmt.Println(url)

	o := Project{}

	err := i.Get(url, &o)
	if err != nil {
		return nil, err
	}

	o.ParseDates()

	return &o, nil
}

// GetProjects returns all projects
//
func (i *Insightly) GetProjects() ([]Project, error) {
	return i.GetProjectsInternal("")
}

// GetProjectsUpdatedAfter returns all projects updated after certain date
//
func (i *Insightly) GetProjectsUpdatedAfter(updatedAfter time.Time) ([]Project, error) {
	from := updatedAfter.Format("2006-01-02")
	searchFilter := fmt.Sprintf("updated_after_utc=%s&", from)
	return i.GetProjectsInternal(searchFilter)
}

// GetProjectsFiltered returns all projects fulfulling the specified filter
//
func (i *Insightly) GetProjectsFiltered(fieldname string, fieldvalue string) ([]Project, error) {
	searchFilter := fmt.Sprintf("field_name=%s&field_value=%s&", fieldname, fieldvalue)
	return i.GetProjectsInternal(searchFilter)
}

// GetProjectsInternal is the generic function retrieving projects from Insightly
//
func (i *Insightly) GetProjectsInternal(searchFilter string) ([]Project, error) {
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
		url := fmt.Sprintf(urlStr, i.apiURL, searchString, strconv.Itoa(skip), strconv.Itoa(top))
		//fmt.Println(url)

		os := []Project{}

		err := i.Get(url, &os)
		if err != nil {
			return nil, err
		}

		for _, o := range os {
			o.ParseDates()
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

func (o *Project) ParseDates() {
	// parse STARTED_DATE to time.Time
	if o.STARTED_DATE != "" {
		t, _ := time.Parse("2006-01-02 15:04:05 +0000 UTC", o.STARTED_DATE+" +0000 UTC")
		//errortools.Fatal(err)
		o.StartedDate = &t
	}

	// parse COMPLETED_DATE to time.Time
	if o.COMPLETED_DATE != "" {
		t, _ := time.Parse("2006-01-02 15:04:05 +0000 UTC", o.COMPLETED_DATE+" +0000 UTC")
		//errortools.Fatal(err)
		o.CompletedDate = &t
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
