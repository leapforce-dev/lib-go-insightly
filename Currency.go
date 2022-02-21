package insightly

import (
	"net/http"

	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
)

// Currency stores Currency from Service
//
type Currency struct {
	CurrencyCode   string `json:"CURRENCY_CODE"`
	CurrencySymbol string `json:"CURRENCY_SYMBOL"`
}

// GetCurrencies returns all currencies
//

func (service *Service) GetCurrencies() (*[]Currency, *errortools.Error) {
	currencies := []Currency{}
	requestConfig := go_http.RequestConfig{
		Method:        http.MethodGet,
		Url:           service.url("Currencies"),
		ResponseModel: &currencies,
	}

	_, _, e := service.httpRequest(&requestConfig)
	if e != nil {
		return nil, e
	}

	return &currencies, nil
}
