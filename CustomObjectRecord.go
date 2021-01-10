package insightly

import (
	"fmt"
	"strconv"

	errortools "github.com/leapforce-libraries/go_errortools"
)

// CustomObjectRecord stores CustomObjectRecord from Service
//
type CustomObjectRecord struct {
	RecordID       int          `json:"RECORD_ID"`
	RecordName     string       `json:"RECORD_NAME"`
	OwnerUserID    *int         `json:"OWNER_USER_ID"`
	DateCreatedUTC DateUTC      `json:"DATE_CREATED_UTC"`
	DateUpdatedUTC DateUTC      `json:"DATE_UPDATED_UTC"`
	CreatedUserID  *int         `json:"CREATED_USER_ID"`
	VisibleTo      string       `json:"VISIBLE_TO"`
	VisibleTeamID  *int         `json:"VISIBLE_TEAM_ID"`
	CustomFields   CustomFields `json:"CUSTOMFIELDS"`
}

func (c *CustomObjectRecord) prepareMarshal() interface{} {
	if c == nil {
		return nil
	}

	return &struct {
		RecordID      int           `json:"RECORD_ID"`
		RecordName    string        `json:"RECORD_NAME"`
		OwnerUserID   *int          `json:"OWNER_USER_ID"`
		VisibleTo     string        `json:"VISIBLE_TO"`
		VisibleTeamID *int          `json:"VISIBLE_TEAM_ID"`
		CustomFields  []CustomField `json:"CUSTOMFIELDS"`
	}{
		c.RecordID,
		c.RecordName,
		c.OwnerUserID,
		c.VisibleTo,
		c.VisibleTeamID,
		c.CustomFields,
	}
}

// GetCustomObjectRecord returns a specific customObjectRecord
//
func (i *Service) GetCustomObjectRecord(customObjectName string, customObjectRecordID int) (*CustomObjectRecord, *errortools.Error) {
	endpoint := fmt.Sprintf("%s/%v", customObjectName, customObjectRecordID)

	customObjectRecord := CustomObjectRecord{}

	_, _, e := i.get(endpoint, nil, &customObjectRecord)
	if e != nil {
		return nil, e
	}

	return &customObjectRecord, nil
}

// GetCustomObjectRecords returns all customObjectRecords
//
func (i *Service) GetCustomObjectRecords(customObjectName string, filter *FieldFilter) (*[]CustomObjectRecord, *errortools.Error) {
	endpointStr := "%s%sskip=%s&top=%s"
	skip := 0
	top := 100
	rowCount := top

	customObjectRecords := []CustomObjectRecord{}

	for rowCount >= top {
		endpoint := fmt.Sprintf(endpointStr, customObjectName, filter.Search(), strconv.Itoa(skip), strconv.Itoa(top))
		//fmt.Println(endpoint)

		cs := []CustomObjectRecord{}

		_, _, e := i.get(endpoint, nil, &cs)
		if e != nil {
			return nil, e
		}

		customObjectRecords = append(customObjectRecords, cs...)

		rowCount = len(cs)
		//rowCount = 0
		skip += top
	}

	if len(customObjectRecords) == 0 {
		customObjectRecords = nil
	}

	return &customObjectRecords, nil
}

// CreateCustomObjectRecord creates a new contract
//
func (i *Service) CreateCustomObjectRecord(customObjectName string, customObjectRecord *CustomObjectRecord) (*CustomObjectRecord, *errortools.Error) {
	if customObjectRecord == nil {
		return nil, nil
	}

	endpoint := customObjectName

	customObjectRecordNew := CustomObjectRecord{}

	_, _, e := i.post(endpoint, customObjectRecord.prepareMarshal(), &customObjectRecordNew)
	if e != nil {
		return nil, e
	}

	return &customObjectRecordNew, nil
}

// UpdateCustomObjectRecord updates an existing contract
//
func (i *Service) UpdateCustomObjectRecord(customObjectName string, customObjectRecord *CustomObjectRecord) (*CustomObjectRecord, *errortools.Error) {
	if customObjectRecord == nil {
		return nil, nil
	}

	endpoint := customObjectName

	customObjectRecordUpdated := CustomObjectRecord{}

	_, _, e := i.put(endpoint, customObjectRecord.prepareMarshal(), &customObjectRecordUpdated)
	if e != nil {
		return nil, e
	}

	return &customObjectRecordUpdated, nil
}

// DeleteCustomObjectRecord deletes a specific customObjectRecord
//
func (i *Service) DeleteCustomObjectRecord(customObjectName string, customObjectRecordID int) *errortools.Error {
	endpoint := fmt.Sprintf("%s/%v", customObjectName, customObjectRecordID)

	_, _, e := i.delete(endpoint, nil, nil)
	if e != nil {
		return e
	}

	return nil
}
