package insightly

import (
	"fmt"
	"net/url"
	"time"

	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
	i_types "github.com/leapforce-libraries/go_insightly/types"
)

// Organisation stores Organisation from Service
//
type Organisation struct {
	OrganisationID         int64                   `json:"ORGANISATION_ID"`
	OrganisationName       *string                 `json:"ORGANISATION_NAME,omitempty"`
	Background             *string                 `json:"BACKGROUND,omitempty"`
	ImageURL               *string                 `json:"IMAGE_URL,omitempty"`
	OwnerUserID            *int64                  `json:"OWNER_USER_ID,omitempty"`
	DateCreatedUTC         *i_types.DateTimeString `json:"DATE_CREATED_UTC,omitempty"`
	DateUpdatedUTC         *i_types.DateTimeString `json:"DATE_UPDATED_UTC,omitempty"`
	LastActivityDateUTC    *i_types.DateTimeString `json:"LAST_ACTIVITY_DATE_UTC,omitempty"`
	NextActivityDateUTC    *i_types.DateTimeString `json:"NEXT_ACTIVITY_DATE_UTC,omitempty"`
	CreatedUserID          *int64                  `json:"CREATED_USER_ID,omitempty"`
	Phone                  *string                 `json:"PHONE,omitempty"`
	PhoneFax               *string                 `json:"PHONE_FAX,omitempty"`
	Website                *string                 `json:"WEBSITE,omitempty"`
	AddressBillingStreet   *string                 `json:"ADDRESS_BILLING_STREET,omitempty"`
	AddressBillingCity     *string                 `json:"ADDRESS_BILLING_CITY,omitempty"`
	AddressBillingState    *string                 `json:"ADDRESS_BILLING_STATE,omitempty"`
	AddressBillingCountry  *string                 `json:"ADDRESS_BILLING_COUNTRY,omitempty"`
	AddressBillingPostcode *string                 `json:"ADDRESS_BILLING_POSTCODE,omitempty"`
	AddressShipStreet      *string                 `json:"ADDRESS_SHIP_STREET,omitempty"`
	AddressShipCity        *string                 `json:"ADDRESS_SHIP_CITY,omitempty"`
	AddressShipState       *string                 `json:"ADDRESS_SHIP_STATE,omitempty"`
	AddressShipCountry     *string                 `json:"ADDRESS_SHIP_COUNTRY,omitempty"`
	AddressShipPostcode    *string                 `json:"ADDRESS_SHIP_POSTCODE,omitempty"`
	SocialLinkedin         *string                 `json:"SOCIAL_LINKEDIN,omitempty"`
	SocialFacebook         *string                 `json:"SOCIAL_FACEBOOK,omitempty"`
	SocialTwitter          *string                 `json:"SOCIAL_TWITTER,omitempty"`
	CustomFields           *CustomFields           `json:"CUSTOMFIELDS,omitempty"`
	Tags                   *[]Tag                  `json:"TAGS,omitempty"`
	Dates                  *[]Date                 `json:"DATES,omitempty"`
	EmailDomains           *[]EmailDomain          `json:"EMAILDOMAINS,omitempty"`
	Links                  *[]Link                 `json:"LINKS,omitempty"`
}

// GetOrganisation returns a specific organisation
//
func (service *Service) GetOrganisation(organisationID int64) (*Organisation, *errortools.Error) {
	organisation := Organisation{}

	requestConfig := go_http.RequestConfig{
		URL:           service.url(fmt.Sprintf("Organisations/%v", organisationID)),
		ResponseModel: &organisation,
	}
	_, _, e := service.get(&requestConfig)
	if e != nil {
		return nil, e
	}

	return &organisation, nil
}

type GetOrganisationsConfig struct {
	Skip         *uint64
	Top          *uint64
	Brief        *bool
	CountTotal   *bool
	UpdatedAfter *time.Time
	FieldFilter  *FieldFilter
}

// GetOrganisations returns all organisations
//
func (service *Service) GetOrganisations(config *GetOrganisationsConfig) (*[]Organisation, *errortools.Error) {
	params := url.Values{}

	endpoint := "Organisations"
	organisations := []Organisation{}
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
		organisationsBatch := []Organisation{}

		requestConfig := go_http.RequestConfig{
			URL:           service.url(fmt.Sprintf("%s?%s", endpoint, params.Encode())),
			ResponseModel: &organisationsBatch,
		}
		_, _, e := service.get(&requestConfig)
		if e != nil {
			return nil, e
		}

		organisations = append(organisations, organisationsBatch...)

		if len(organisationsBatch) < int(top) {
			delete(service.nextSkips, endpoint)
			break
		}

		service.nextSkips[endpoint] += top
		rowCount += top

		if rowCount >= service.maxRowCount {
			return &organisations, nil
		}
	}

	return &organisations, nil
}

// CreateOrganisation creates a new contract
//
func (service *Service) CreateOrganisation(organisation *Organisation) (*Organisation, *errortools.Error) {
	if organisation == nil {
		return nil, nil
	}

	organisationNew := Organisation{}

	requestConfig := go_http.RequestConfig{
		URL:           service.url("Organisations"),
		BodyModel:     organisation,
		ResponseModel: &organisationNew,
	}
	_, _, e := service.post(&requestConfig)
	if e != nil {
		return nil, e
	}

	return &organisationNew, nil
}

// UpdateOrganisation updates an existing contract
//
func (service *Service) UpdateOrganisation(organisation *Organisation) (*Organisation, *errortools.Error) {
	if organisation == nil {
		return nil, nil
	}

	organisationUpdated := Organisation{}

	requestConfig := go_http.RequestConfig{
		URL:           service.url("Organisations"),
		BodyModel:     organisation,
		ResponseModel: &organisationUpdated,
	}
	_, _, e := service.put(&requestConfig)
	if e != nil {
		return nil, e
	}

	return &organisationUpdated, nil
}

// DeleteOrganisation deletes a specific organisation
//
func (service *Service) DeleteOrganisation(organisationID int64) *errortools.Error {
	requestConfig := go_http.RequestConfig{
		URL: service.url(fmt.Sprintf("Organisations/%v", organisationID)),
	}
	_, _, e := service.delete(&requestConfig)
	if e != nil {
		return e
	}

	return nil
}

// GetOrganisationLinks returns links for a specific organisation
//
func (service *Service) GetOrganisationLinks(organisationID int64) (*[]Link, *errortools.Error) {
	links := []Link{}

	requestConfig := go_http.RequestConfig{
		URL:           service.url(fmt.Sprintf("Organisations/%v/Links", organisationID)),
		ResponseModel: &links,
	}
	_, _, e := service.get(&requestConfig)
	if e != nil {
		return nil, e
	}

	return &links, nil
}
