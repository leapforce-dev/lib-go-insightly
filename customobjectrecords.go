package insightly

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

// CustomObjectRecord stores CustomObjectRecord from Insightly
//
type CustomObjectRecord struct {
	RECORD_ID        int           `json:"RECORD_ID"`
	RECORD_NAME      string        `json:"RECORD_NAME"`
	OWNER_USER_ID    int           `json:"OWNER_USER_ID"`
	DATE_CREATED_UTC string        `json:"DATE_CREATED_UTC"`
	DATE_UPDATED_UTC string        `json:"DATE_UPDATED_UTC"`
	CREATED_USER_ID  int           `json:"CREATED_USER_ID"`
	VISIBLE_TO       string        `json:"VISIBLE_TO"`
	VISIBLE_TEAM_ID  int           `json:"VISIBLE_TEAM_ID"`
	CUSTOMFIELDS     []CustomField `json:"CUSTOMFIELDS"`
	DateCreated      *time.Time    `json:"-"`
	DateUpdated      *time.Time    `json:"-"`
}

// GetCustomObjectRecords returns all customobjectrecords
//
func (i *Insightly) GetCustomObjectRecords(objectName string) ([]CustomObjectRecord, error) {
	return i.GetCustomObjectRecordsInternal(objectName, "")
}

// GetCustomObjectRecordsFiltered returns all customobjectrecords fulfulling the specified filter
//
func (i *Insightly) GetCustomObjectRecordsFiltered(objectName string, fieldname string, fieldvalue string) ([]CustomObjectRecord, error) {
	searchFilter := fmt.Sprintf("field_name=%s&field_value=%s&", fieldname, fieldvalue)
	return i.GetCustomObjectRecordsInternal(objectName, searchFilter)
}

// GetCustomObjectRecordsInternal is the generic function retrieving customobjectrecords from Insightly
//
func (i *Insightly) GetCustomObjectRecordsInternal(objectName string, searchFilter string) ([]CustomObjectRecord, error) {
	searchString := ""

	if searchFilter != "" {
		searchString = "/Search?" + searchFilter
	} else {
		searchString = "?"
	}

	urlStr := "%s%s%sskip=%s&top=%s"
	skip := 0
	top := 500
	rowCount := top

	customobjectrecords := []CustomObjectRecord{}

	for rowCount >= top {
		url := fmt.Sprintf(urlStr, i.apiURL, objectName, searchString, strconv.Itoa(skip), strconv.Itoa(top))
		//fmt.Println(url)

		os := []CustomObjectRecord{}

		err := i.Get(url, &os)
		if err != nil {
			return nil, err
		}

		for _, o := range os {
			o.ParseDates()
			customobjectrecords = append(customobjectrecords, o)
		}

		rowCount = len(os)
		skip += top
	}

	if len(customobjectrecords) == 0 {
		customobjectrecords = nil
	}

	return customobjectrecords, nil
}

// GetCustomObjectRecordsInternal is the generic function retrieving customobjectrecords from Insightly
//
func (i *Insightly) UpdateCustomObjectRecords(customObjectName string, customObjectRecord CustomObjectRecord) error {
	urlStr := "%s%s"

	url := fmt.Sprintf(urlStr, i.apiURL, customObjectName)
	//fmt.Println(url)

	b, err := json.Marshal(customObjectRecord)
	if err != nil {
		return err
	}

	err = i.Put(url, b)
	if err != nil {
		return err
	}

	return nil
}

func (o *CustomObjectRecord) ParseDates() {
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
