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
	OrganisationName       string                  `json:"ORGANISATION_NAME"`
	Background             *string                 `json:"BACKGROUND"`
	ImageURL               *string                 `json:"IMAGE_URL"`
	OwnerUserID            int64                   `json:"OWNER_USER_ID"`
	DateCreatedUTC         i_types.DateTimeString  `json:"DATE_CREATED_UTC"`
	DateUpdatedUTC         i_types.DateTimeString  `json:"DATE_UPDATED_UTC"`
	LastActivityDateUTC    *i_types.DateTimeString `json:"LAST_ACTIVITY_DATE_UTC"`
	NextActivityDateUTC    *i_types.DateTimeString `json:"NEXT_ACTIVITY_DATE_UTC"`
	CreatedUserID          int64                   `json:"CREATED_USER_ID"`
	Phone                  *string                 `json:"PHONE"`
	PhoneFax               *string                 `json:"PHONE_FAX"`
	Website                *string                 `json:"WEBSITE"`
	AddressBillingStreet   *string                 `json:"ADDRESS_BILLING_STREET"`
	AddressBillingCity     *string                 `json:"ADDRESS_BILLING_CITY"`
	AddressBillingState    *string                 `json:"ADDRESS_BILLING_STATE"`
	AddressBillingCountry  *string                 `json:"ADDRESS_BILLING_COUNTRY"`
	AddressBillingPostcode *string                 `json:"ADDRESS_BILLING_POSTCODE"`
	AddressShipStreet      *string                 `json:"ADDRESS_SHIP_STREET"`
	AddressShipCity        *string                 `json:"ADDRESS_SHIP_CITY"`
	AddressShipState       *string                 `json:"ADDRESS_SHIP_STATE"`
	AddressShipCountry     *string                 `json:"ADDRESS_SHIP_COUNTRY"`
	AddressShipPostcode    *string                 `json:"ADDRESS_SHIP_POSTCODE"`
	SocialLinkedin         *string                 `json:"SOCIAL_LINKEDIN"`
	SocialFacebook         *string                 `json:"SOCIAL_FACEBOOK"`
	SocialTwitter          *string                 `json:"SOCIAL_TWITTER"`
	CustomFields           *CustomFields           `json:"CUSTOMFIELDS"`
	Tags                   *[]Tag                  `json:"TAGS"`
	Dates                  *[]Date                 `json:"DATES"`
	EmailDomains           *[]EmailDomain          `json:"EMAILDOMAINS"`
}

func (o *Organisation) prepareMarshal() interface{} {
	if o == nil {
		return nil
	}

	return &struct {
		OrganisationID         *int64        `json:"ORGANISATION_ID"`
		OrganisationName       *string       `json:"ORGANISATION_NAME"`
		Background             *string       `json:"BACKGROUND"`
		ImageURL               *string       `json:"IMAGE_URL"`
		OwnerUserID            *int64        `json:"OWNER_USER_ID"`
		Phone                  *string       `json:"PHONE"`
		PhoneFax               *string       `json:"PHONE_FAX"`
		Website                *string       `json:"WEBSITE"`
		AddressBillingStreet   *string       `json:"ADDRESS_BILLING_STREET"`
		AddressBillingCity     *string       `json:"ADDRESS_BILLING_CITY"`
		AddressBillingState    *string       `json:"ADDRESS_BILLING_STATE"`
		AddressBillingCountry  *string       `json:"ADDRESS_BILLING_COUNTRY"`
		AddressBillingPostcode *string       `json:"ADDRESS_BILLING_POSTCODE"`
		AddressShipStreet      *string       `json:"ADDRESS_SHIP_STREET"`
		AddressShipCity        *string       `json:"ADDRESS_SHIP_CITY"`
		AddressShipState       *string       `json:"ADDRESS_SHIP_STATE"`
		AddressShipCountry     *string       `json:"ADDRESS_SHIP_COUNTRY"`
		AddressShipPostcode    *string       `json:"ADDRESS_SHIP_POSTCODE"`
		SocialLinkedin         *string       `json:"SOCIAL_LINKEDIN"`
		SocialFacebook         *string       `json:"SOCIAL_FACEBOOK"`
		SocialTwitter          *string       `json:"SOCIAL_TWITTER"`
		CustomFields           *CustomFields `json:"CUSTOMFIELDS"`
	}{
		&o.OrganisationID,
		&o.OrganisationName,
		o.Background,
		o.ImageURL,
		&o.OwnerUserID,
		o.Phone,
		o.PhoneFax,
		o.Website,
		o.AddressBillingStreet,
		o.AddressBillingCity,
		o.AddressBillingState,
		o.AddressBillingCountry,
		o.AddressBillingPostcode,
		o.AddressShipStreet,
		o.AddressShipCity,
		o.AddressShipState,
		o.AddressShipCountry,
		o.AddressShipPostcode,
		o.SocialLinkedin,
		o.SocialFacebook,
		o.SocialTwitter,
		o.CustomFields,
	}
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
		BodyModel:     organisation.prepareMarshal(),
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
		BodyModel:     organisation.prepareMarshal(),
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
func (service *Service) DeleteOrganisation(organisationID int) *errortools.Error {
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
func (service *Service) GetOrganisationLinks(organisationID int) (*[]Link, *errortools.Error) {
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
