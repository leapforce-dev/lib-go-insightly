package insightly

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
)

// Organisation stores Organisation from Service
//
type Organisation struct {
	OrganisationID         int           `json:"ORGANISATION_ID"`
	OrganisationName       string        `json:"ORGANISATION_NAME"`
	Background             string        `json:"BACKGROUND"`
	ImageURL               string        `json:"IMAGE_URL"`
	OwnerUserID            *int          `json:"OWNER_USER_ID"`
	DateCreatedUTC         DateUTC       `json:"DATE_CREATED_UTC"`
	DateUpdatedUTC         DateUTC       `json:"DATE_UPDATED_UTC"`
	LastActivityDateUTC    DateUTC       `json:"LAST_ACTIVITY_DATE_UTC"`
	NextActivityDateUTC    DateUTC       `json:"NEXT_ACTIVITY_DATE_UTC"`
	CreatedUserID          *int          `json:"CREATED_USER_ID"`
	Phone                  string        `json:"PHONE"`
	PhoneFax               string        `json:"PHONE_FAX"`
	Website                string        `json:"WEBSITE"`
	AddressBillingStreet   string        `json:"ADDRESS_BILLING_STREET"`
	AddressBillingCity     string        `json:"ADDRESS_BILLING_CITY"`
	AddressBillingState    string        `json:"ADDRESS_BILLING_STATE"`
	AddressBillingCountry  string        `json:"ADDRESS_BILLING_COUNTRY"`
	AddressBillingPostcode string        `json:"ADDRESS_BILLING_POSTCODE"`
	AddressShipStreet      string        `json:"ADDRESS_SHIP_STREET"`
	AddressShipCity        string        `json:"ADDRESS_SHIP_CITY"`
	AddressShipState       string        `json:"ADDRESS_SHIP_STATE"`
	AddressShipCountry     string        `json:"ADDRESS_SHIP_COUNTRY"`
	AddressShipPostcode    string        `json:"ADDRESS_SHIP_POSTCODE"`
	SocialLinkedin         string        `json:"SOCIAL_LINKEDIN"`
	SocialFacebook         string        `json:"SOCIAL_FACEBOOK"`
	SocialTwitter          string        `json:"SOCIAL_TWITTER"`
	CustomFields           CustomFields  `json:"CUSTOMFIELDS"`
	Tags                   []Tag         `json:"TAGS"`
	Dates                  []Date        `json:"DATES"`
	EmailDomains           []EmailDomain `json:"EMAILDOMAINS"`
}

func (o *Organisation) prepareMarshal() interface{} {
	if o == nil {
		return nil
	}

	return &struct {
		OrganisationID         int          `json:"ORGANISATION_ID"`
		OrganisationName       string       `json:"ORGANISATION_NAME"`
		Background             string       `json:"BACKGROUND"`
		ImageURL               string       `json:"IMAGE_URL"`
		OwnerUserID            *int         `json:"OWNER_USER_ID"`
		Phone                  string       `json:"PHONE"`
		PhoneFax               string       `json:"PHONE_FAX"`
		Website                string       `json:"WEBSITE"`
		AddressBillingStreet   string       `json:"ADDRESS_BILLING_STREET"`
		AddressBillingCity     string       `json:"ADDRESS_BILLING_CITY"`
		AddressBillingState    string       `json:"ADDRESS_BILLING_STATE"`
		AddressBillingCountry  string       `json:"ADDRESS_BILLING_COUNTRY"`
		AddressBillingPostcode string       `json:"ADDRESS_BILLING_POSTCODE"`
		AddressShipStreet      string       `json:"ADDRESS_SHIP_STREET"`
		AddressShipCity        string       `json:"ADDRESS_SHIP_CITY"`
		AddressShipState       string       `json:"ADDRESS_SHIP_STATE"`
		AddressShipCountry     string       `json:"ADDRESS_SHIP_COUNTRY"`
		AddressShipPostcode    string       `json:"ADDRESS_SHIP_POSTCODE"`
		SocialLinkedin         string       `json:"SOCIAL_LINKEDIN"`
		SocialFacebook         string       `json:"SOCIAL_FACEBOOK"`
		SocialTwitter          string       `json:"SOCIAL_TWITTER"`
		CustomFields           CustomFields `json:"CUSTOMFIELDS"`
	}{
		o.OrganisationID,
		o.OrganisationName,
		o.Background,
		o.ImageURL,
		o.OwnerUserID,
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
func (service *Service) GetOrganisation(organisationID int) (*Organisation, *errortools.Error) {
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

type GetOrganisationsFilter struct {
	UpdatedAfter *time.Time
	Field        *struct {
		FieldName  string
		FieldValue string
	}
}

// GetOrganisations returns all organisations
//
func (service *Service) GetOrganisations(filter *GetOrganisationsFilter) (*[]Organisation, *errortools.Error) {
	searchString := "?"
	searchFilter := []string{}

	if filter != nil {
		if filter.UpdatedAfter != nil {
			from := filter.UpdatedAfter.Format(time.RFC3339)
			searchFilter = append(searchFilter, fmt.Sprintf("updated_after_utc=%s&", from))
		}

		if filter.Field != nil {
			searchFilter = append(searchFilter, fmt.Sprintf("field_name=%s&field_value=%s&", filter.Field.FieldName, filter.Field.FieldValue))
		}
	}

	if len(searchFilter) > 0 {
		searchString = "/Search?" + strings.Join(searchFilter, "&")
	}

	endpointStr := "Organisations%sskip=%s&top=%s"
	skip := 0
	top := 100
	rowCount := top

	organisations := []Organisation{}

	for rowCount >= top {
		_organisations := []Organisation{}

		requestConfig := go_http.RequestConfig{
			URL:           service.url(fmt.Sprintf(endpointStr, searchString, strconv.Itoa(skip), strconv.Itoa(top))),
			ResponseModel: &_organisations,
		}
		_, _, e := service.get(&requestConfig)
		if e != nil {
			return nil, e
		}

		organisations = append(organisations, _organisations...)

		rowCount = len(_organisations)
		//rowCount = 0
		skip += top
	}

	if len(organisations) == 0 {
		organisations = nil
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
