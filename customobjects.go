package insightly

import (
	"fmt"
	"time"
)

// CustomObject stores CustomObject from Insightly
//
type CustomObject struct {
	OBJECT_NAME                string `json:"OBJECT_NAME"`
	SINGULAR_LABEL             string `json:"SINGULAR_LABEL"`
	PLURAL_LABEL               string `json:"PLURAL_LABEL"`
	DESCRIPTION                string `json:"DESCRIPTION"`
	RECORD_NAME_LABEL          string `json:"RECORD_NAME_LABEL"`
	RECORD_NAME_TYPE           string `json:"RECORD_NAME_TYPE"`
	RECORD_NAME_DISPLAY_FORMAT string `json:"RECORD_NAME_DISPLAY_FORMAT"`
	ENABLE_NAVBAR              bool   `json:"ENABLE_NAVBAR"`
	ENABLE_WORKFLOWS           bool   `json:"ENABLE_WORKFLOWS"`
	CREATED_USER_ID            int    `json:"CREATED_USER_ID"`
	DATE_CREATED_UTC           string `json:"DATE_CREATED_UTC"`
	DATE_UPDATED_UTC           string `json:"DATE_UPDATED_UTC"`
	DateCreated                *time.Time
	DateUpdated                *time.Time
}

// GetCustomObjects returns all customobjects
//
func (i *Insightly) GetCustomObjects() ([]CustomObject, error) {
	urlStr := "%sCustomObjects"

	customobjects := []CustomObject{}

	url := fmt.Sprintf(urlStr, i.apiURL)
	//fmt.Println(url)

	os := []CustomObject{}

	err := i.Get(url, &os)
	if err != nil {
		return nil, err
	}

	for _, o := range os {
		o.ParseDates()
		customobjects = append(customobjects, o)
	}

	if len(customobjects) == 0 {
		customobjects = nil
	}

	return customobjects, nil
}

func (o *CustomObject) ParseDates() {
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
}
