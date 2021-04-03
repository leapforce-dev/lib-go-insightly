package insightly

import (
	"fmt"
	"net/url"
	"time"

	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
	i_types "github.com/leapforce-libraries/go_insightly/types"
)

// PricebookEntry stores PricebookEntry from Service
//
type PricebookEntry struct {
	PricebookEntryID int64                  `json:"PRICEBOOK_ENTRY_ID"`
	PricebookID      int64                  `json:"PRICEBOOK_ID"`
	ProductID        int64                  `json:"PRODUCT_ID"`
	CurrencyCode     string                 `json:"CURRENCY_CODE"`
	Price            float64                `json:"PRICE"`
	UseStandardPrice bool                   `json:"USE_STANDARD_PRICE"`
	Active           bool                   `json:"ACTIVE"`
	CreatedUserID    int64                  `json:"CREATED_USER_ID"`
	DateCreatedUTC   i_types.DateTimeString `json:"DATE_CREATED_UTC"`
	DateUpdatedUTC   i_types.DateTimeString `json:"DATE_UPDATED_UTC"`
	CustomFields     *CustomFields          `json:"CUSTOMFIELDS"`
}

// GetPricebookEntry returns a specific pricebookEntry
//
func (service *Service) GetPricebookEntry(pricebookEntryID int64) (*PricebookEntry, *errortools.Error) {
	pricebookEntry := PricebookEntry{}

	requestConfig := go_http.RequestConfig{
		URL:           service.url(fmt.Sprintf("PricebookEntry/%v", pricebookEntryID)),
		ResponseModel: &pricebookEntry,
	}
	_, _, e := service.get(&requestConfig)
	if e != nil {
		return nil, e
	}

	return &pricebookEntry, nil
}

type GetPricebookEntriesConfig struct {
	Skip         *uint64
	Top          *uint64
	Brief        *bool
	CountTotal   *bool
	UpdatedAfter *time.Time
	FieldFilter  *FieldFilter
}

// GetPricebookEntries returns all PricebookEntries
//
func (service *Service) GetPricebookEntries(config *GetPricebookEntriesConfig) (*[]PricebookEntry, *errortools.Error) {
	params := url.Values{}

	endpoint := "PricebookEntry"
	pricebookEntries := []PricebookEntry{}
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
		pricebookEntriesBatch := []PricebookEntry{}

		requestConfig := go_http.RequestConfig{
			URL:           service.url(fmt.Sprintf("%s?%s", endpoint, params.Encode())),
			ResponseModel: &pricebookEntriesBatch,
		}
		_, _, e := service.get(&requestConfig)
		if e != nil {
			return nil, e
		}

		pricebookEntries = append(pricebookEntries, pricebookEntriesBatch...)

		if len(pricebookEntriesBatch) < int(top) {
			delete(service.nextSkips, endpoint)
			break
		}

		service.nextSkips[endpoint] += top
		rowCount += top

		if rowCount >= service.maxRowCount {
			return &pricebookEntries, nil
		}
	}

	return &pricebookEntries, nil
}
