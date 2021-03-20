package insightly

import (
	"fmt"
	"strconv"
	"strings"

	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
	i_types "github.com/leapforce-libraries/go_insightly/types"
)

// CustomObjectRecord stores CustomObjectRecord from Service
//
type CustomObjectRecord struct {
	RecordID       int64                  `json:"RECORD_ID"`
	RecordName     string                 `json:"RECORD_NAME"`
	OwnerUserID    int64                  `json:"OWNER_USER_ID"`
	DateCreatedUTC i_types.DateTimeString `json:"DATE_CREATED_UTC"`
	DateUpdatedUTC i_types.DateTimeString `json:"DATE_UPDATED_UTC"`
	CreatedUserID  int64                  `json:"CREATED_USER_ID"`
	VisibleTo      *string                `json:"VISIBLE_TO"`
	VisibleTeamID  *int64                 `json:"VISIBLE_TEAM_ID"`
	CustomFields   *CustomFields          `json:"CUSTOMFIELDS"`
}

func (c *CustomObjectRecord) prepareMarshal() interface{} {
	if c == nil {
		return nil
	}

	return &struct {
		RecordID      *int64        `json:"RECORD_ID,omitempty"`
		RecordName    *string       `json:"RECORD_NAME,omitempty"`
		OwnerUserID   *int64        `json:"OWNER_USER_ID,omitempty"`
		VisibleTo     *string       `json:"VISIBLE_TO,omitempty"`
		VisibleTeamID *int64        `json:"VISIBLE_TEAM_ID,omitempty"`
		CustomFields  *CustomFields `json:"CUSTOMFIELDS,omitempty"`
	}{
		&c.RecordID,
		&c.RecordName,
		&c.OwnerUserID,
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

type GetCustomObjectRecordsConfig struct {
	CustomObjectName string
	FieldFilter      *FieldFilter
}

// GetCustomObjectRecords returns all customObjectRecords
//
func (service *Service) GetCustomObjectRecords(config *GetCustomObjectRecordsConfig) (*[]CustomObjectRecord, *errortools.Error) {
	searchString := "?"
	searchFilter := []string{}

	if config == nil {
		return nil, errortools.ErrorMessage("GetCustomObjectRecordsConfig must not be nil or a nil pointer")
	}

	if config.FieldFilter != nil {
		searchFilter = append(searchFilter, fmt.Sprintf("field_name=%s&field_value=%s&", config.FieldFilter.FieldName, config.FieldFilter.FieldValue))
	}

	if len(searchFilter) > 0 {
		searchString = "/Search?" + strings.Join(searchFilter, "&")
	}

	endpointStr := "%s%sskip=%s&top=%s"
	skip := 0
	top := 100
	rowCount := top

	customObjectRecords := []CustomObjectRecord{}

	for rowCount >= top {
		_customObjectRecords := []CustomObjectRecord{}

		requestConfig := go_http.RequestConfig{
			URL:           service.url(fmt.Sprintf(endpointStr, config.CustomObjectName, searchString, strconv.Itoa(skip), strconv.Itoa(top))),
			ResponseModel: &_customObjectRecords,
		}
		_, _, e := service.get(&requestConfig)
		if e != nil {
			return nil, e
		}

		customObjectRecords = append(customObjectRecords, _customObjectRecords...)

		rowCount = len(_customObjectRecords)
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
