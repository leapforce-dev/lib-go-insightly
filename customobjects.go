package insightly

import (
	"fmt"
	"time"

	errortools "github.com/leapforce-libraries/go_errortools"
)

// CustomObject stores CustomObject from Insightly
//
type CustomObject struct {
	ObjectName              string `json:"OBJECT_NAME"`
	SingularLabel           string `json:"SINGULAR_LABEL"`
	PluralLabel             string `json:"PLURAL_LABEL"`
	Description             string `json:"DESCRIPTION"`
	RecordNameLabel         string `json:"RECORD_NAME_LABEL"`
	RecordNameType          string `json:"RECORD_NAME_TYPE"`
	RecordNameDisplayFormat string `json:"RECORD_NAME_DISPLAY_FORMAT"`
	EnableNavbar            bool   `json:"ENABLE_NAVBAR"`
	EnableWorkflows         bool   `json:"ENABLE_WORKFLOWS"`
	CreatedUserID           int    `json:"CREATED_USER_ID"`
	DateCreatedUTC          string `json:"DATE_CREATED_UTC"`
	DateUpdatedUTC          string `json:"DATE_UPDATED_UTC"`
	DateCreatedT            *time.Time
	DateUpdatedT            *time.Time
}

// GetCustomObjects returns all customobjects
//
func (i *Insightly) GetCustomObjects() ([]CustomObject, *errortools.Error) {
	urlStr := "%sCustomObjects"

	customobjects := []CustomObject{}

	url := fmt.Sprintf(urlStr, apiURL)
	//fmt.Println(url)

	os := []CustomObject{}

	_, _, err := i.get(url, nil, &os)
	if err != nil {
		return nil, err
	}

	for _, o := range os {
		o.parseDates()
		customobjects = append(customobjects, o)
	}

	if len(customobjects) == 0 {
		customobjects = nil
	}

	return customobjects, nil
}

func (o *CustomObject) parseDates() {
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
}
