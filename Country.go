package insightly

import (
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
		URL:           service.url("Countries"),
		ResponseModel: &countries,
	}

	_, _, e := service.get(&requestConfig)
	if e != nil {
		return nil, e
	}

	return &countries, nil
}
