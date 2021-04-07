package insightly

import (
	"fmt"
	"net/url"
	"time"

	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
	i_types "github.com/leapforce-libraries/go_insightly/types"
)

// Pricebook stores Pricebook from Service
//
type Pricebook struct {
	PricebookID    int64                  `json:"PRICEBOOK_ID"`
	Name           string                 `json:"NAME"`
	Description    string                 `json:"DESCRIPTION"`
	CurrencyCode   string                 `json:"CURRENCY_CODE"`
	IsStandard     bool                   `json:"IS_STANDARD"`
	Active         bool                   `json:"ACTIVE"`
	OwnerUserID    int64                  `json:"OWNER_USER_ID"`
	CreatedUserID  int64                  `json:"CREATED_USER_ID"`
	DateCreatedUTC i_types.DateTimeString `json:"DATE_CREATED_UTC"`
	DateUpdatedUTC i_types.DateTimeString `json:"DATE_UPDATED_UTC"`
}

// GetPricebook returns a specific pricebook
//
func (service *Service) GetPricebook(pricebookID int64) (*Pricebook, *errortools.Error) {
	pricebook := Pricebook{}

	requestConfig := go_http.RequestConfig{
		URL:           service.url(fmt.Sprintf("Pricebook/%v", pricebookID)),
		ResponseModel: &pricebook,
	}
	_, _, e := service.get(&requestConfig)
	if e != nil {
		return nil, e
	}

	return &pricebook, nil
}

type GetPricebooksConfig struct {
	Skip         *uint64
	Top          *uint64
	Brief        *bool
	CountTotal   *bool
	UpdatedAfter *time.Time
	FieldFilter  *FieldFilter
}

// GetPricebooks returns all pricebooks
//
func (service *Service) GetPricebooks(config *GetPricebooksConfig) (*[]Pricebook, *errortools.Error) {
	params := url.Values{}

	endpoint := "Pricebook"
	pricebooks := []Pricebook{}
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
			params.Set("updated_after_utc", fmt.Sprintf("%v", config.UpdatedAfter.Format(dateTimeFormat)))
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
		pricebooksBatch := []Pricebook{}

		requestConfig := go_http.RequestConfig{
			URL:           service.url(fmt.Sprintf("%s?%s", endpoint, params.Encode())),
			ResponseModel: &pricebooksBatch,
		}
		_, _, e := service.get(&requestConfig)
		if e != nil {
			return nil, e
		}

		pricebooks = append(pricebooks, pricebooksBatch...)

		if len(pricebooksBatch) < int(top) {
			delete(service.nextSkips, endpoint)
			break
		}

		service.nextSkips[endpoint] += top
		rowCount += top

		if rowCount >= service.maxRowCount {
			return &pricebooks, nil
		}
	}

	return &pricebooks, nil
}
