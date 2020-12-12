package insightly

import (
	"fmt"
	"strconv"
	"time"

	errortools "github.com/leapforce-libraries/go_errortools"
)

// CustomObjectRecord stores CustomObjectRecord from Insightly
//
type CustomObjectRecord struct {
	RecordID       int           `json:"RECORD_ID"`
	RecordName     string        `json:"RECORD_NAME"`
	OwnerUserID    int           `json:"OWNER_USER_ID"`
	DateCreatedUTC string        `json:"DATE_CREATED_UTC"`
	DateUpdatedUTC string        `json:"DATE_UPDATED_UTC"`
	CreatedUserID  int           `json:"CREATED_USER_ID"`
	VisibleTo      string        `json:"VISIBLE_TO"`
	VisibleTeamID  int           `json:"VISIBLE_TEAM_ID"`
	CustomFields   []CustomField `json:"CUSTOMFIELDS"`
	DateCreatedT   *time.Time    `json:"-"`
	DateUpdatedT   *time.Time    `json:"-"`
}

// GetCustomObjectRecords returns all customobjectrecords
//
func (i *Insightly) GetCustomObjectRecords(objectName string) ([]CustomObjectRecord, *errortools.Error) {
	return i.GetCustomObjectRecordsInternal(objectName, "")
}

// GetCustomObjectRecordsFiltered returns all customobjectrecords fulfulling the specified filter
//
func (i *Insightly) GetCustomObjectRecordsFiltered(objectName string, fieldname string, fieldvalue string) ([]CustomObjectRecord, *errortools.Error) {
	searchFilter := fmt.Sprintf("field_name=%s&field_value=%s&", fieldname, fieldvalue)
	return i.GetCustomObjectRecordsInternal(objectName, searchFilter)
}

// GetCustomObjectRecordsInternal is the generic function retrieving customobjectrecords from Insightly
//
func (i *Insightly) GetCustomObjectRecordsInternal(objectName string, searchFilter string) ([]CustomObjectRecord, *errortools.Error) {
	searchString := ""

	if searchFilter != "" {
		searchString = "/Search?" + searchFilter
	} else {
		searchString = "?"
	}

	endpointStr := "%s%sskip=%s&top=%s"
	skip := 0
	top := 500
	rowCount := top

	customobjectrecords := []CustomObjectRecord{}

	for rowCount >= top {
		endpoint := fmt.Sprintf(endpointStr, objectName, searchString, strconv.Itoa(skip), strconv.Itoa(top))
		//fmt.Println(endpoint)

		os := []CustomObjectRecord{}

		_, _, err := i.get(endpoint, nil, &os)
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
func (i *Insightly) UpdateCustomObjectRecords(customObjectName string, customObjectRecord CustomObjectRecord) *errortools.Error {
	endpoint := customObjectName
	//fmt.Println(endpoint)

	_, _, e := i.put(endpoint, customObjectRecord, nil)
	if e != nil {
		return e
	}

	return nil
}

func (o *CustomObjectRecord) ParseDates() {
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
