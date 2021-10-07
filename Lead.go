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

// Lead stores Lead from Service
//
type Lead struct {
	LeadID                  int64                   `json:"LEAD_ID"`
	Salutation              *string                 `json:"SALUTATION"`
	FirstName               *string                 `json:"FIRST_NAME"`
	LastName                *string                 `json:"LAST_NAME"`
	LeadSourceID            int64                   `json:"LEAD_SOURCE_ID"`
	LeadStatusID            int64                   `json:"LEAD_STATUS_ID"`
	Title                   *string                 `json:"TITLE"`
	Converted               bool                    `json:"CONVERTED"`
	ConvertedContactID      *int64                  `json:"CONVERTED_CONTACT_ID"`
	ConvertedDateUTC        *i_types.DateTimeString `json:"CONVERTED_DATE_UTC"`
	ConvertedOpportunityID  *int64                  `json:"CONVERTED_OPPORTUNITY_ID"`
	ConvertedOrganisationID *int64                  `json:"CONVERTED_ORGANISATION_ID"`
	DateCreatedUTC          i_types.DateTimeString  `json:"DATE_CREATED_UTC"`
	DateUpdatedUTC          i_types.DateTimeString  `json:"DATE_UPDATED_UTC"`
	Email                   *string                 `json:"EMAIL"`
	EmployeeCount           *int64                  `json:"EMPLOYEE_COUNT"`
	Fax                     *string                 `json:"FAX"`
	Industry                *string                 `json:"INDUSTRY"`
	LeadDescription         *string                 `json:"LEAD_DESCRIPTION"`
	LeadRating              *int64                  `json:"LEAD_RATING"`
	Mobile                  *string                 `json:"MOBILE"`
	OwnerUserID             int64                   `json:"OWNER_USER_ID"`
	Phone                   *string                 `json:"PHONE"`
	ResponsibleUserID       *int64                  `json:"RESPONSIBLE_USER_ID"`
	Website                 *string                 `json:"WEBSITE"`
	AddressStreet           *string                 `json:"ADDRESS_STREET"`
	AddressCity             *string                 `json:"ADDRESS_CITY"`
	AddressState            *string                 `json:"ADDRESS_STATE"`
	AddressPostcode         *string                 `json:"ADDRESS_POSTCODE"`
	AddressCountry          *string                 `json:"ADDRESS_COUNTRY"`
	LastActivityDateUTC     *i_types.DateTimeString `json:"LAST_ACTIVITY_DATE_UTC"`
	NextActivityDateUTC     *i_types.DateTimeString `json:"NEXT_ACTIVITY_DATE_UTC"`
	OrganisationName        *string                 `json:"ORGANISATION_NAME"`
	CreatedUserID           int64                   `json:"CREATED_USER_ID"`
	ImageURL                *string                 `json:"IMAGE_URL"`
	EmailOptedOut           bool                    `json:"EMAIL_OPTED_OUT"`
	CustomFields            *CustomFields           `json:"CUSTOMFIELDS"`
	Tags                    *[]Tag                  `json:"TAGS"`
	Links                   *[]Link                 `json:"LINKS"`
}

func (l *Lead) prepareMarshal() interface{} {
	if l == nil {
		return nil
	}

	return &struct {
		LeadID                  *int64                  `json:"LEAD_ID,omitempty"`
		Salutation              *string                 `json:"SALUTATION,omitempty"`
		FirstName               *string                 `json:"FIRST_NAME,omitempty"`
		LastName                *string                 `json:"LAST_NAME,omitempty"`
		LeadSourceID            *int64                  `json:"LEAD_SOURCE_ID,omitempty"`
		LeadStatusID            *int64                  `json:"LEAD_STATUS_ID,omitempty"`
		Title                   *string                 `json:"TITLE,omitempty"`
		Converted               *bool                   `json:"CONVERTED,omitempty"`
		ConvertedContactID      *int64                  `json:"CONVERTED_CONTACT_ID,omitempty"`
		ConvertedDateUTC        *i_types.DateTimeString `json:"CONVERTED_DATE_UTC,omitempty"`
		ConvertedOpportunityID  *int64                  `json:"CONVERTED_OPPORTUNITY_ID,omitempty"`
		ConvertedOrganisationID *int64                  `json:"CONVERTED_ORGANISATION_ID,omitempty"`
		Email                   *string                 `json:"EMAIL,omitempty"`
		EmployeeCount           *int64                  `json:"EMPLOYEE_COUNT,omitempty"`
		Fax                     *string                 `json:"FAX,omitempty"`
		Industry                *string                 `json:"INDUSTRY,omitempty"`
		LeadDescription         *string                 `json:"LEAD_DESCRIPTION,omitempty"`
		LeadRating              *int64                  `json:"LEAD_RATING,omitempty"`
		Mobile                  *string                 `json:"MOBILE,omitempty"`
		OwnerUserID             *int64                  `json:"OWNER_USER_ID,omitempty"`
		Phone                   *string                 `json:"PHONE,omitempty"`
		ResponsibleUserID       *int64                  `json:"RESPONSIBLE_USER_ID,omitempty"`
		Website                 *string                 `json:"WEBSITE,omitempty"`
		AddressStreet           *string                 `json:"ADDRESS_STREET,omitempty"`
		AddressCity             *string                 `json:"ADDRESS_CITY,omitempty"`
		AddressState            *string                 `json:"ADDRESS_STATE,omitempty"`
		AddressPostcode         *string                 `json:"ADDRESS_POSTCODE,omitempty"`
		AddressCountry          *string                 `json:"ADDRESS_COUNTRY,omitempty"`
		ImageURL                *string                 `json:"IMAGE_URL,omitempty"`
		EmailOptedOut           *bool                   `json:"EMAIL_OPTED_OUT,omitempty"`
		CustomFields            *CustomFields           `json:"CUSTOMFIELDS,omitempty"`
	}{
		&l.LeadID,
		l.Salutation,
		l.FirstName,
		l.LastName,
		&l.LeadSourceID,
		&l.LeadStatusID,
		l.Title,
		&l.Converted,
		l.ConvertedContactID,
		l.ConvertedDateUTC,
		l.ConvertedOpportunityID,
		l.ConvertedOrganisationID,
		l.Email,
		l.EmployeeCount,
		l.Fax,
		l.Industry,
		l.LeadDescription,
		l.LeadRating,
		l.Mobile,
		&l.OwnerUserID,
		l.Phone,
		l.ResponsibleUserID,
		l.Website,
		l.AddressStreet,
		l.AddressCity,
		l.AddressState,
		l.AddressPostcode,
		l.AddressCountry,
		l.ImageURL,
		&l.EmailOptedOut,
		l.CustomFields,
	}
}

// GetLead returns a specific lead
//
func (service *Service) GetLead(leadID int64) (*Lead, *errortools.Error) {
	lead := Lead{}

	requestConfig := go_http.RequestConfig{
		Method:        http.MethodGet,
		URL:           service.url(fmt.Sprintf("Leads/%v", leadID)),
		ResponseModel: &lead,
	}
	_, _, e := service.httpRequest(&requestConfig)
	if e != nil {
		return nil, e
	}

	return &lead, nil
}

type GetLeadsConfig struct {
	Skip         *uint64
	Top          *uint64
	Brief        *bool
	CountTotal   *bool
	UpdatedAfter *time.Time
	FieldFilter  *FieldFilter
}

// GetLeads returns all leads
//
func (service *Service) GetLeads(config *GetLeadsConfig) (*[]Lead, *errortools.Error) {
	params := url.Values{}

	endpoint := "Leads"
	leads := []Lead{}
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
		leadsBatch := []Lead{}

		requestConfig := go_http.RequestConfig{
			Method:        http.MethodGet,
			URL:           service.url(fmt.Sprintf("%s?%s", endpoint, params.Encode())),
			ResponseModel: &leadsBatch,
		}
		_, _, e := service.httpRequest(&requestConfig)
		if e != nil {
			return nil, e
		}

		leads = append(leads, leadsBatch...)

		if len(leadsBatch) < int(top) {
			delete(service.nextSkips, endpoint)
			break
		}

		service.nextSkips[endpoint] += top
		rowCount += top

		if rowCount >= service.maxRowCount {
			return &leads, nil
		}
	}

	return &leads, nil
}

// CreateLead creates a new contract
//
func (service *Service) CreateLead(lead *Lead) (*Lead, *errortools.Error) {
	if lead == nil {
		return nil, nil
	}

	leadNew := Lead{}

	requestConfig := go_http.RequestConfig{
		Method:        http.MethodPost,
		URL:           service.url("Leads"),
		BodyModel:     lead.prepareMarshal(),
		ResponseModel: &leadNew,
	}
	_, _, e := service.httpRequest(&requestConfig)
	if e != nil {
		return nil, e
	}

	return &leadNew, nil
}

// UpdateLead updates an existing contract
//
func (service *Service) UpdateLead(lead *Lead) (*Lead, *errortools.Error) {
	if lead == nil {
		return nil, nil
	}

	leadUpdated := Lead{}

	requestConfig := go_http.RequestConfig{
		Method:        http.MethodPut,
		URL:           service.url("Leads"),
		BodyModel:     lead.prepareMarshal(),
		ResponseModel: &leadUpdated,
	}
	_, _, e := service.httpRequest(&requestConfig)
	if e != nil {
		return nil, e
	}

	return &leadUpdated, nil
}

// DeleteLead deletes a specific lead
//
func (service *Service) DeleteLead(leadID int64) *errortools.Error {
	requestConfig := go_http.RequestConfig{
		Method: http.MethodDelete,
		URL:    service.url(fmt.Sprintf("Leads/%v", leadID)),
	}
	_, _, e := service.httpRequest(&requestConfig)
	if e != nil {
		return e
	}

	return nil
}
