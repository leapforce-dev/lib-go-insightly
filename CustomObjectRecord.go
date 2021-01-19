package insightly

import (
	"fmt"
	"strconv"

	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
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
func (service *Service) GetCustomObjectRecord(customObjectName string, customObjectRecordID int) (*CustomObjectRecord, *errortools.Error) {
	customObjectRecord := CustomObjectRecord{}

	requestConfig := go_http.RequestConfig{
		URL:           service.url(fmt.Sprintf("%s/%v", customObjectName, customObjectRecordID)),
		ResponseModel: &customObjectRecord,
	}
	_, _, e := service.get(&requestConfig)
	if e != nil {
		return nil, e
	}

	return &customObjectRecord, nil
}

// GetCustomObjectRecords returns all customObjectRecords
//
func (service *Service) GetCustomObjectRecords(customObjectName string, filter *FieldFilter) (*[]CustomObjectRecord, *errortools.Error) {
	endpointStr := "%s%sskip=%s&top=%s"
	skip := 0
	top := 100
	rowCount := top

	customObjectRecords := []CustomObjectRecord{}

	for rowCount >= top {
		_customObjectRecords := []CustomObjectRecord{}

		requestConfig := go_http.RequestConfig{
			URL:           service.url(fmt.Sprintf(endpointStr, customObjectName, filter.Search(), strconv.Itoa(skip), strconv.Itoa(top))),
			ResponseModel: &_customObjectRecords,
		}
		_, _, e := service.get(&requestConfig)
		if e != nil {
			return nil, e
		}

		customObjectRecords = append(customObjectRecords, _customObjectRecords...)

		rowCount = len(_customObjectRecords)
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
func (service *Service) CreateCustomObjectRecord(customObjectName string, customObjectRecord *CustomObjectRecord) (*CustomObjectRecord, *errortools.Error) {
	if customObjectRecord == nil {
		return nil, nil
	}

	customObjectRecordNew := CustomObjectRecord{}

	requestConfig := go_http.RequestConfig{
		URL:           service.url(customObjectName),
		BodyModel:     customObjectRecord.prepareMarshal(),
		ResponseModel: &customObjectRecordNew,
	}
	_, _, e := service.post(&requestConfig)
	if e != nil {
		return nil, e
	}

	return &customObjectRecordNew, nil
}

// UpdateCustomObjectRecord updates an existing contract
//
func (service *Service) UpdateCustomObjectRecord(customObjectName string, customObjectRecord *CustomObjectRecord) (*CustomObjectRecord, *errortools.Error) {
	if customObjectRecord == nil {
		return nil, nil
	}

	customObjectRecordUpdated := CustomObjectRecord{}

	requestConfig := go_http.RequestConfig{
		URL:           service.url(customObjectName),
		BodyModel:     customObjectRecord.prepareMarshal(),
		ResponseModel: &customObjectRecordUpdated,
	}
	_, _, e := service.put(&requestConfig)
	if e != nil {
		return nil, e
	}

	return &customObjectRecordUpdated, nil
}

// DeleteCustomObjectRecord deletes a specific customObjectRecord
//
func (service *Service) DeleteCustomObjectRecord(customObjectName string, customObjectRecordID int) *errortools.Error {
	requestConfig := go_http.RequestConfig{
		URL: service.url(fmt.Sprintf("%s/%v", customObjectName, customObjectRecordID)),
	}
	_, _, e := service.delete(&requestConfig)
	if e != nil {
		return e
	}

	return nil
}
