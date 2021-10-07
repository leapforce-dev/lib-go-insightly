package insightly

import (
	"net/http"

	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
)

// Country stores Country from Service
//
type Country struct {
	CountryName string `json:"COUNTRY_NAME"`
}

// GetCountries returns all countries
//

func (service *Service) GetCountries() (*[]Country, *errortools.Error) {
	countries := []Country{}
	requestConfig := go_http.RequestConfig{
		Method:        http.MethodGet,
		URL:           service.url("Countries"),
		ResponseModel: &countries,
	}

	_, _, e := service.httpRequest(&requestConfig)
	if e != nil {
		return nil, e
	}

	return &countries, nil
}
