package insightly

import (
	"fmt"
	"net/url"
	"time"

	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
	i_types "github.com/leapforce-libraries/go_insightly/types"
)

// QuoteProduct stores QuoteProduct from Service
//
type QuoteProduct struct {
	QuotationItemID   int64                  `json:"QUOTATION_ITEM_ID"`
	OpportunityItemID int64                  `json:"OPPORTUNITY_ITEM_ID"`
	PricebookEntryID  int64                  `json:"PRICEBOOK_ENTRY_ID"`
	Description       string                 `json:"DESCRIPTION"`
	CurrencyCode      string                 `json:"CURRENCY_CODE"`
	Quantity          int64                  `json:"QUANTITY"`
	ListPrice         float64                `json:"LIST_PRICE"`
	UnitPrice         float64                `json:"UNIT_PRICE"`
	Subtotal          float64                `json:"SUBTOTAL"`
	Discount          float64                `json:"DISCOUNT"`
	TotalPrice        float64                `json:"TOTAL_PRICE"`
	DateCreatedUTC    i_types.DateTimeString `json:"DATE_CREATED_UTC"`
	DateUpdatedUTC    i_types.DateTimeString `json:"DATE_UPDATED_UTC"`
	SortOrder         int64                  `json:"SORT_ORDER"`
	CustomFields      *CustomFields          `json:"CUSTOMFIELDS"`
}

type GetQuoteProductsConfig struct {
	Skip         *uint64
	Top          *uint64
	Brief        *bool
	CountTotal   *bool
	UpdatedAfter *time.Time
	FieldFilter  *FieldFilter
}

// GetQuoteProducts returns all quoteProducts
//
func (service *Service) GetQuoteProducts(config *GetQuoteProductsConfig) (*[]QuoteProduct, *errortools.Error) {
	params := url.Values{}

	endpoint := "QuotationLineItem"
	quoteProducts := []QuoteProduct{}
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

		quoteProductsBatch := []QuoteProduct{}

		requestConfig := go_http.RequestConfig{
			URL:           service.url(fmt.Sprintf("%s?%s", endpoint, params.Encode())),
			ResponseModel: &quoteProductsBatch,
		}

		_, _, e := service.get(&requestConfig)
		if e != nil {
			return nil, e
		}

		quoteProducts = append(quoteProducts, quoteProductsBatch...)

		if len(quoteProductsBatch) < int(top) {
			delete(service.nextSkips, endpoint)
			break
		}

		service.nextSkips[endpoint] += top
		rowCount += top

		if rowCount >= service.maxRowCount {
			return &quoteProducts, nil
		}
	}

	return &quoteProducts, nil
}
