package insightly

import (
	"fmt"
	"net/url"
	"time"

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
	Skip             *uint64
	Top              *uint64
	Brief            *bool
	CountTotal       *bool
	CustomObjectName string
	UpdatedAfter     *time.Time
	FieldFilter      *FieldFilter
}

// GetCustomObjectRecords returns all customObjectRecords
//
func (service *Service) GetCustomObjectRecords(config *GetCustomObjectRecordsConfig) (*[]CustomObjectRecord, *errortools.Error) {
	if config == nil {
		return nil, nil
	}

	params := url.Values{}

	customObjectRecords := []CustomObjectRecord{}

	endpoint := config.CustomObjectName
	rowCount := uint64(0)
	top := defaultTop
	isSearch := false

	if config != nil {
		if config.Top != nil {
			top = *config.Top
		}
		if config.Skip != nil {
			service.nextSkips[endpoint] = *config.Skip
		}
		if config.Brief != nil {
			params.Set("brief", fmt.Sprintf("%v", *config.Brief))
		}
		if config.CountTotal != nil {
			params.Set("count_total", fmt.Sprintf("%v", *config.CountTotal))
		}
		if config.UpdatedAfter != nil {
			isSearch = true
			params.Set("updated_after_utc", fmt.Sprintf("%v", config.UpdatedAfter.Format(DateTimeFormat)))
		}
		if config.FieldFilter != nil {
			isSearch = true
			params.Set("field_name", config.FieldFilter.FieldName)
			params.Set("field_value", config.FieldFilter.FieldValue)
		}
	}

	if isSearch {
		endpoint += "/Search"
	}

	params.Set("top", fmt.Sprintf("%v", top))

	for true {
		params.Set("skip", fmt.Sprintf("%v", service.nextSkips[endpoint]))
		customObjectRecordsBatch := []CustomObjectRecord{}

		requestConfig := go_http.RequestConfig{
			URL:           service.url(fmt.Sprintf("%s?%s", endpoint, params.Encode())),
			ResponseModel: &customObjectRecordsBatch,
		}
		_, _, e := service.get(&requestConfig)
		if e != nil {
			return nil, e
		}

		customObjectRecords = append(customObjectRecords, customObjectRecordsBatch...)

		if len(customObjectRecordsBatch) < int(top) {
			service.nextSkips[endpoint] = 0
			break
		}

		service.nextSkips[endpoint] += top
		rowCount += top

		if rowCount >= service.maxRowCount {
			return &customObjectRecords, nil
		}
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
