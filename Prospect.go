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

type Prospect struct {
	ProspectID            int64                   `json:"PROSPECT_ID"`
	LeadID                *int64                  `json:"LEAD_ID"`
	ContactID             *int64                  `json:"CONTACT_ID"`
	OrganisationID        *int64                  `json:"ORGANISATION_ID"`
	Salutation            *string                 `json:"SALUTATION"`
	FirstName             string                  `json:"FIRST_NAME"`
	LastName              string                  `json:"LAST_NAME"`
	OrganisationName      *string                 `json:"ORGANISATION_NAME"`
	Title                 *string                 `json:"TITLE"`
	EmailAddress          *string                 `json:"EMAIL_ADDRESS"`
	Phone                 *string                 `json:"PHONE"`
	Mobile                *string                 `json:"MOBILE"`
	Fax                   *string                 `json:"FAX"`
	Website               *string                 `json:"WEBSITE"`
	AddressStreet         *string                 `json:"ADDRESS_STREET"`
	AddressCity           *string                 `json:"ADDRESS_CITY"`
	AddressState          *string                 `json:"ADDRESS_STATE"`
	AddressPostcode       *string                 `json:"ADDRESS_POSTCODE"`
	AddressCountry        *string                 `json:"ADDRESS_COUNTRY"`
	Industry              *string                 `json:"INDUSTRY"`
	EmployeeCount         *int64                  `json:"EMPLOYEE_COUNT"`
	Score                 int64                   `json:"SCORE"`
	Grade                 string                  `json:"GRADE"`
	Description           *string                 `json:"DESCRIPTION"`
	DoNotEmail            bool                    `json:"DO_NOT_EMAIL"`
	DoNotCall             bool                    `json:"DO_NOT_CALL"`
	OptedOut              bool                    `json:"OPTED_OUT"`
	LastActivityDateUTC   *i_types.DateTimeString `json:"LAST_ACTIVITY_DATE_UTC"`
	CreatedUserID         int64                   `json:"CREATED_USER_ID"`
	OwnerUserID           int64                   `json:"OWNER_USER_ID"`
	VisibleTo             string                  `json:"VISIBLE_TO"`
	VisibleTeamID         *int64                  `json:"VISIBLE_TEAM_ID"`
	DateCreatedUTC        i_types.DateTimeString  `json:"DATE_CREATED_UTC"`
	DateUpdatedUTC        i_types.DateTimeString  `json:"DATE_UPDATED_UTC"`
	DoNotSync             bool                    `json:"DO_NOT_SYNC"`
	LeadConversionDateUTC *i_types.DateTimeString `json:"LEAD_CONVERSION_DATE_UTC"`
	GradeProfileID        *int64                  `json:"GRADE_PROFILE_ID"`
	CustomFields          *CustomFields           `json:"CUSTOMFIELDS"`
	Tags                  *[]Tag                  `json:"TAGS"`
}

// GetProspect returns a specific prospect
//
func (service *Service) GetProspect(prospectID int64) (*Prospect, *errortools.Error) {
	prospect := Prospect{}

	requestConfig := go_http.RequestConfig{
		Method:        http.MethodGet,
		Url:           service.url(fmt.Sprintf("Prospect/%v", prospectID)),
		ResponseModel: &prospect,
	}
	_, _, e := service.httpRequest(&requestConfig)
	if e != nil {
		return nil, e
	}

	return &prospect, nil
}

type GetProspectsConfig struct {
	Skip         *uint64
	Top          *uint64
	Brief        *bool
	CountTotal   *bool
	UpdatedAfter *time.Time
	FieldFilter  *FieldFilter
}

// GetProspects returns all prospects
//
func (service *Service) GetProspects(config *GetProspectsConfig) (*[]Prospect, *errortools.Error) {
	params := url.Values{}

	endpoint := "Prospect"
	prospects := []Prospect{}
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
		prospectsBatch := []Prospect{}

		requestConfig := go_http.RequestConfig{
			Method:        http.MethodGet,
			Url:           service.url(fmt.Sprintf("%s?%s", endpoint, params.Encode())),
			ResponseModel: &prospectsBatch,
		}
		_, _, e := service.httpRequest(&requestConfig)
		if e != nil {
			return nil, e
		}

		prospects = append(prospects, prospectsBatch...)

		if len(prospectsBatch) < int(top) {
			delete(service.nextSkips, endpoint)
			break
		}

		service.nextSkips[endpoint] += top
		rowCount += top

		if rowCount >= service.maxRowCount {
			return &prospects, nil
		}
	}

	return &prospects, nil
}
