package insightly

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
)

// Lead stores Lead from Service
//
type Lead struct {
	LeadID                  int          `json:"LEAD_ID"`
	Salutation              string       `json:"SALUTATION"`
	FirstName               string       `json:"FIRST_NAME"`
	LastName                string       `json:"LAST_NAME"`
	LeadSourceID            *int         `json:"LEAD_SOURCE_ID"`
	LeadStatusID            *int         `json:"LEAD_STATUS_ID"`
	Title                   string       `json:"TITLE"`
	Converted               bool         `json:"CONVERTED"`
	ConvertedContactID      *int         `json:"CONVERTED_CONTACT_ID"`
	ConvertedDateUTC        DateUTC      `json:"CONVERTED_DATE_UTC"`
	ConvertedOpportunityID  *int         `json:"CONVERTED_OPPORTUNITY_ID"`
	ConvertedOrganisationID *int         `json:"CONVERTED_ORGANISATION_ID"`
	DateCreateUTC           DateUTC      `json:"DATE_CREATED_UTC"`
	DateUpdatedUTC          DateUTC      `json:"DATE_UPDATED_UTC"`
	Email                   string       `json:"EMAIL"`
	EmployeeCount           int          `json:"EMPLOYEE_COUNT"`
	Fax                     string       `json:"FAX"`
	Industry                string       `json:"INDUSTRY"`
	LeadDescription         string       `json:"LEAD_DESCRIPTION"`
	LeadRating              int          `json:"LEAD_RATING"`
	Mobile                  string       `json:"MOBILE"`
	OwnerUserID             *int         `json:"OWNER_USER_ID"`
	Phone                   string       `json:"PHONE"`
	ResponsibleUserID       *int         `json:"RESPONSIBLE_USER_ID"`
	Website                 string       `json:"WEBSITE"`
	AddressStreet           string       `json:"ADDRESS_STREET"`
	AddressCity             string       `json:"ADDRESS_CITY"`
	AddressState            string       `json:"ADDRESS_STATE"`
	AddressPostcode         string       `json:"ADDRESS_POSTCODE"`
	AddressCountry          string       `json:"ADDRESS_COUNTRY"`
	LastActivityDateUTC     DateUTC      `json:"LAST_ACTIVITY_DATE_UTC"`
	NextActivityDateUTC     DateUTC      `json:"NEXT_ACTIVITY_DATE_UTC"`
	OrganisationName        string       `json:"ORGANISATION_NAME"`
	CreatedUserID           *int         `json:"CREATED_USER_ID"`
	ImageURL                string       `json:"IMAGE_URL"`
	EmailOptedOut           bool         `json:"EMAIL_OPTED_OUT"`
	CustomFields            CustomFields `json:"CUSTOMFIELDS"`
	Tags                    []Tag        `json:"TAGS"`
}

func (l *Lead) prepareMarshal() interface{} {
	if l == nil {
		return nil
	}

	return &struct {
		LeadID                  int           `json:"LEAD_ID"`
		Salutation              string        `json:"SALUTATION"`
		FirstName               string        `json:"FIRST_NAME"`
		LastName                string        `json:"LAST_NAME"`
		LeadSourceID            *int          `json:"LEAD_SOURCE_ID"`
		LeadStatusID            *int          `json:"LEAD_STATUS_ID"`
		Title                   string        `json:"TITLE"`
		Converted               bool          `json:"CONVERTED"`
		ConvertedContactID      *int          `json:"CONVERTED_CONTACT_ID"`
		ConvertedDateUTC        DateUTC       `json:"CONVERTED_DATE_UTC"`
		ConvertedOpportunityID  *int          `json:"CONVERTED_OPPORTUNITY_ID"`
		ConvertedOrganisationID *int          `json:"CONVERTED_ORGANISATION_ID"`
		Email                   string        `json:"EMAIL"`
		EmployeeCount           int           `json:"EMPLOYEE_COUNT"`
		Fax                     string        `json:"FAX"`
		Industry                string        `json:"INDUSTRY"`
		LeadDescription         string        `json:"LEAD_DESCRIPTION"`
		LeadRating              int           `json:"LEAD_RATING"`
		Mobile                  string        `json:"MOBILE"`
		OwnerUserID             *int          `json:"OWNER_USER_ID"`
		Phone                   string        `json:"PHONE"`
		ResponsibleUserID       *int          `json:"RESPONSIBLE_USER_ID"`
		Website                 string        `json:"WEBSITE"`
		AddressStreet           string        `json:"ADDRESS_STREET"`
		AddressCity             string        `json:"ADDRESS_CITY"`
		AddressState            string        `json:"ADDRESS_STATE"`
		AddressPostcode         string        `json:"ADDRESS_POSTCODE"`
		AddressCountry          string        `json:"ADDRESS_COUNTRY"`
		ImageURL                string        `json:"IMAGE_URL"`
		EmailOptedOut           bool          `json:"EMAIL_OPTED_OUT"`
		CustomFields            []CustomField `json:"CUSTOMFIELDS"`
	}{
		l.LeadID,
		l.Salutation,
		l.FirstName,
		l.LastName,
		l.LeadSourceID,
		l.LeadStatusID,
		l.Title,
		l.Converted,
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
		l.OwnerUserID,
		l.Phone,
		l.ResponsibleUserID,
		l.Website,
		l.AddressStreet,
		l.AddressCity,
		l.AddressState,
		l.AddressPostcode,
		l.AddressCountry,
		l.ImageURL,
		l.EmailOptedOut,
		l.CustomFields,
	}
}

// GetLead returns a specific lead
//
func (service *Service) GetLead(leadID int) (*Lead, *errortools.Error) {
	lead := Lead{}

	requestConfig := go_http.RequestConfig{
		URL:           service.url(fmt.Sprintf("Leads/%v", leadID)),
		ResponseModel: &lead,
	}
	_, _, e := service.get(&requestConfig)
	if e != nil {
		return nil, e
	}

	return &lead, nil
}

type GetLeadsFilter struct {
	UpdatedAfter *time.Time
	Field        *struct {
		FieldName  string
		FieldValue string
	}
}

// GetLeads returns all leads
//
func (service *Service) GetLeads(filter *GetLeadsFilter) (*[]Lead, *errortools.Error) {
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

	endpointStr := "Leads%sskip=%s&top=%s"
	skip := 0
	top := 100
	rowCount := top

	leads := []Lead{}

	for rowCount >= top {
		_leads := []Lead{}

		requestConfig := go_http.RequestConfig{
			URL:           service.url(fmt.Sprintf(endpointStr, searchString, strconv.Itoa(skip), strconv.Itoa(top))),
			ResponseModel: &_leads,
		}
		_, _, e := service.get(&requestConfig)
		if e != nil {
			return nil, e
		}

		leads = append(leads, _leads...)

		rowCount = len(_leads)
		//rowCount = 0
		skip += top
	}

	if len(leads) == 0 {
		leads = nil
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
		URL:           service.url("Leads"),
		BodyModel:     lead.prepareMarshal(),
		ResponseModel: &leadNew,
	}
	_, _, e := service.post(&requestConfig)
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
		URL:           service.url("Leads"),
		BodyModel:     lead.prepareMarshal(),
		ResponseModel: &leadUpdated,
	}
	_, _, e := service.put(&requestConfig)
	if e != nil {
		return nil, e
	}

	return &leadUpdated, nil
}

// DeleteLead deletes a specific lead
//
func (service *Service) DeleteLead(leadID int) *errortools.Error {
	requestConfig := go_http.RequestConfig{
		URL: service.url(fmt.Sprintf("Leads/%v", leadID)),
	}
	_, _, e := service.delete(&requestConfig)
	if e != nil {
		return e
	}

	return nil
}
