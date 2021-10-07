package insightly

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
	i_types "github.com/leapforce-libraries/go_insightly/types"
)

// OpportunityProduct stores OpportunityProduct from Service
//
type OpportunityProduct struct {
	OpportunityItemID int64                   `json:"OPPORTUNITY_ITEM_ID"`
	OpportunityID     int64                   `json:"OPPORTUNITY_ID"`
	PricebookEntryID  int64                   `json:"PRICEBOOK_ENTRY_ID"`
	CurrencyCode      string                  `json:"CURRENCY_CODE"`
	UnitPrice         float64                 `json:"UNIT_PRICE"`
	Description       string                  `json:"DESCRIPTION"`
	Quantity          int64                   `json:"QUANTITY"`
	ServiceDate       *i_types.DateTimeString `json:"SERVICE_DATE"`
	TotalPrice        float64                 `json:"TOTAL_PRICE"`
	DateCreatedUTC    i_types.DateTimeString  `json:"DATE_CREATED_UTC"`
	DateUpdatedUTC    i_types.DateTimeString  `json:"DATE_UPDATED_UTC"`
	ListPrice         float64                 `json:"LIST_PRICE"`
	Subtotal          float64                 `json:"SUBTOTAL"`
	Discount          float64                 `json:"DISCOUNT"`
	CustomFields      *CustomFields           `json:"CUSTOMFIELDS"`
}

type GetOpportunityProductsConfig struct {
	Skip         *uint64
	Top          *uint64
	Brief        *bool
	CountTotal   *bool
	UpdatedAfter *time.Time
	FieldFilter  *FieldFilter
}

// GetOpportunityProducts returns all opportunityProducts
//
func (service *Service) GetOpportunityProducts(config *GetOpportunityProductsConfig) (*[]OpportunityProduct, *errortools.Error) {
	params := url.Values{}

	endpoint := "OpportunityLineItem"
	opportunityProducts := []OpportunityProduct{}
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

	for {
		params.Set("skip", fmt.Sprintf("%v", service.nextSkips[endpoint]))

		opportunityProductsBatch := []OpportunityProduct{}

		requestConfig := go_http.RequestConfig{
			Method:        http.MethodGet,
			URL:           service.url(fmt.Sprintf("%s?%s", endpoint, params.Encode())),
			ResponseModel: &opportunityProductsBatch,
		}

		_, _, e := service.httpRequest(&requestConfig)
		if e != nil {
			return nil, e
		}

		opportunityProducts = append(opportunityProducts, opportunityProductsBatch...)

		if len(opportunityProductsBatch) < int(top) {
			delete(service.nextSkips, endpoint)
			break
		}

		service.nextSkips[endpoint] += top
		rowCount += top

		if rowCount >= service.maxRowCount {
			return &opportunityProducts, nil
		}
	}

	return &opportunityProducts, nil
}
