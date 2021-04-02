package insightly

import (
	"fmt"
	"net/url"
	"time"

	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
	i_types "github.com/leapforce-libraries/go_insightly/types"
)

type Quote struct {
	QuoteID                 int64                   `json:"QUOTE_ID"`
	QuoteName               string                  `json:"QUOTATION_NAME"`
	OpportunityID           *int64                  `json:"OPPORTUNITY_ID"`
	ContactID               *int64                  `json:"CONTACT_ID"`
	OrganisationID          *int64                  `json:"ORGANISATION_ID"`
	QuotationNumber         *string                 `json:"QUOTATION_NUMBER"`
	QuotationDescription    *string                 `json:"QUOTATION_DESCRIPTION"`
	QuotationPhone          *string                 `json:"QUOTATION_PHONE"`
	QuotationEmail          *string                 `json:"QUOTATION_EMAIL"`
	QuotationFax            *string                 `json:"QUOTATION_FAX"`
	QuoteStatus             string                  `json:"QUOTE_STATUS"`
	QuotationExpirationDate *i_types.DateTimeString `json:"QUOTATION_EXPIRATION_DATE"`
	LineItemCount           int64                   `json:"LINE_ITEM_COUNT"`
	IsSyncing               bool                    `json:"IS_SYNCING"`
	QuotationCurrencyCode   *string                 `json:"QUOTATION_CURRENCY_CODE"`
	Subtotal                *float64                `json:"SUBTOTAL"`
	Discount                *float64                `json:"DISCOUNT"`
	TotalPrice              *float64                `json:"TOTAL_PRICE"`
	ShappingHandling        *float64                `json:"SHIPPING_HANDLING"`
	Tax                     *float64                `json:"TAX"`
	GrandTotal              *float64                `json:"GRAND_TOTAL"`
	AddressBillingName      *string                 `json:"ADDRESS_BILLING_NAME"`
	AddressBillingStreet    *string                 `json:"ADDRESS_BILLING_STREET"`
	AddressBillingCity      *string                 `json:"ADDRESS_BILLING_CITY"`
	AddressBillingState     *string                 `json:"ADDRESS_BILLING_STATE"`
	AddressBillingCountry   *string                 `json:"ADDRESS_BILLING_COUNTRY"`
	AddressBillingPostcode  *string                 `json:"ADDRESS_BILLING_POSTCODE"`
	AddressShippingName     *string                 `json:"ADDRESS_SHIPPING_NAME"`
	AddressShippingStreet   *string                 `json:"ADDRESS_SHIPPING_STREET"`
	AddressShippingCity     *string                 `json:"ADDRESS_SHIPPING_CITY"`
	AddressShippingState    *string                 `json:"ADDRESS_SHIPPING_STATE"`
	AddressShippingCountry  *string                 `json:"ADDRESS_SHIPPING_COUNTRY"`
	AddressShippingPostcode *string                 `json:"ADDRESS_SHIPPING_POSTCODE"`
	OwnerUserID             int64                   `json:"OWNER_USER_ID"`
	DateCreatedUTC          i_types.DateTimeString  `json:"DATE_CREATED_UTC"`
	DateUpdatedUTC          i_types.DateTimeString  `json:"DATE_UPDATED_UTC"`
	CreatedUserID           int64                   `json:"CREATED_USER_ID"`
	CustomFields            *CustomFields           `json:"CUSTOMFIELDS"`
}

// GetQuote returns a specific quote
//
func (service *Service) GetQuote(quoteID int64) (*Quote, *errortools.Error) {
	quote := Quote{}

	requestConfig := go_http.RequestConfig{
		URL:           service.url(fmt.Sprintf("Quotation/%v", quoteID)),
		ResponseModel: &quote,
	}
	_, _, e := service.get(&requestConfig)
	if e != nil {
		return nil, e
	}

	return &quote, nil
}

type GetQuotesConfig struct {
	Skip         *uint64
	Top          *uint64
	Brief        *bool
	CountTotal   *bool
	UpdatedAfter *time.Time
	FieldFilter  *FieldFilter
}

// GetQuotes returns all quotes
//
func (service *Service) GetQuotes(config *GetQuotesConfig) (*[]Quote, *errortools.Error) {
	params := url.Values{}

	endpoint := "Quotation"
	quotes := []Quote{}
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
		quotesBatch := []Quote{}

		requestConfig := go_http.RequestConfig{
			URL:           service.url(fmt.Sprintf("%s?%s", endpoint, params.Encode())),
			ResponseModel: &quotesBatch,
		}
		_, _, e := service.get(&requestConfig)
		if e != nil {
			return nil, e
		}

		quotes = append(quotes, quotesBatch...)

		if len(quotesBatch) < int(top) {
			delete(service.nextSkips, endpoint)
			break
		}

		service.nextSkips[endpoint] += top
		rowCount += top

		if rowCount >= service.maxRowCount {
			return &quotes, nil
		}
	}

	return &quotes, nil
}
