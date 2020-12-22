package insightly

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	errortools "github.com/leapforce-libraries/go_errortools"
)

type Project struct {
	ProjectID            int          `json:"CONTACT_ID"`
	Salutation           string       `json:"SALUTATION"`
	FirstName            string       `json:"FIRST_NAME"`
	LastName             string       `json:"LAST_NAME"`
	ImageURL             string       `json:"IMAGE_URL"`
	Background           string       `json:"BACKGROUND"`
	OwnerUserID          *int         `json:"OWNER_USER_ID"`
	DateCreatedUTC       DateUTC      `json:"DATE_CREATED_UTC"`
	DateUpdatedUTC       DateUTC      `json:"DATE_UPDATED_UTC"`
	SocialLinkedin       string       `json:"SOCIAL_LINKEDIN"`
	SocialFacebook       string       `json:"SOCIAL_FACEBOOK"`
	SocialTwitter        string       `json:"SOCIAL_TWITTER"`
	DateOfBirth          DateUTC      `json:"DATE_OF_BIRTH"`
	Phone                string       `json:"PHONE"`
	PhoneHome            string       `json:"PHONE_HOME"`
	PhoneMobile          string       `json:"PHONE_MOBILE"`
	PhoneOther           string       `json:"PHONE_OTHER"`
	PhoneAssistant       string       `json:"PHONE_ASSISTANT"`
	PhoneFax             string       `json:"PHONE_FAX"`
	EmailAddress         string       `json:"EMAIL_ADDRESS"`
	AssistantName        string       `json:"ASSISTANT_NAME"`
	AddressMailStreet    string       `json:"ADDRESS_MAIL_STREET"`
	AddressMailCity      string       `json:"ADDRESS_MAIL_CITY"`
	AddressMailState     string       `json:"ADDRESS_MAIL_STATE"`
	AddressMailPostcode  string       `json:"ADDRESS_MAIL_POSTCODE"`
	AddressMailCountry   string       `json:"ADDRESS_MAIL_COUNTRY"`
	AddressOtherStreet   string       `json:"ADDRESS_OTHER_STREET"`
	AddressOtherCity     string       `json:"ADDRESS_OTHER_CITY"`
	AddressOtherState    string       `json:"ADDRESS_OTHER_STATE"`
	AddressOtherPostcode string       `json:"ADDRESS_OTHER_POSTCODE"`
	LastActivityDateUTC  DateUTC      `json:"LAST_ACTIVITY_DATE_UTC"`
	NextActivityDateUTC  DateUTC      `json:"NEXT_ACTIVITY_DATE_UTC"`
	CreatedUserID        *int         `json:"CREATED_USER_ID"`
	OrganisationID       *int         `json:"ORGANISATION_ID"`
	Title                string       `json:"TITLE"`
	EmailOptedOut        bool         `json:"EMAIL_OPTED_OUT"`
	CustomFields         CustomFields `json:"CUSTOMFIELDS"`
	Tags                 []Tag        `json:"TAGS"`
	Dates                []Date       `json:"DATES"`
}

func (p *Project) prepareMarshal() interface{} {
	if p == nil {
		return nil
	}

	return &struct {
		ProjectID            int           `json:"CONTACT_ID"`
		Salutation           string        `json:"SALUTATION"`
		FirstName            string        `json:"FIRST_NAME"`
		LastName             string        `json:"LAST_NAME"`
		ImageURL             string        `json:"IMAGE_URL"`
		Background           string        `json:"BACKGROUND"`
		OwnerUserID          *int          `json:"OWNER_USER_ID"`
		SocialLinkedin       string        `json:"SOCIAL_LINKEDIN"`
		SocialFacebook       string        `json:"SOCIAL_FACEBOOK"`
		SocialTwitter        string        `json:"SOCIAL_TWITTER"`
		DateOfBirth          DateUTC       `json:"DATE_OF_BIRTH"`
		Phone                string        `json:"PHONE"`
		PhoneHome            string        `json:"PHONE_HOME"`
		PhoneMobile          string        `json:"PHONE_MOBILE"`
		PhoneOther           string        `json:"PHONE_OTHER"`
		PhoneAssistant       string        `json:"PHONE_ASSISTANT"`
		PhoneFax             string        `json:"PHONE_FAX"`
		EmailAddress         string        `json:"EMAIL_ADDRESS"`
		AssistantName        string        `json:"ASSISTANT_NAME"`
		AddressMailStreet    string        `json:"ADDRESS_MAIL_STREET"`
		AddressMailCity      string        `json:"ADDRESS_MAIL_CITY"`
		AddressMailState     string        `json:"ADDRESS_MAIL_STATE"`
		AddressMailPostcode  string        `json:"ADDRESS_MAIL_POSTCODE"`
		AddressMailCountry   string        `json:"ADDRESS_MAIL_COUNTRY"`
		AddressOtherStreet   string        `json:"ADDRESS_OTHER_STREET"`
		AddressOtherCity     string        `json:"ADDRESS_OTHER_CITY"`
		AddressOtherState    string        `json:"ADDRESS_OTHER_STATE"`
		AddressOtherPostcode string        `json:"ADDRESS_OTHER_POSTCODE"`
		OrganisationID       *int          `json:"ORGANISATION_ID"`
		Title                string        `json:"TITLE"`
		EmailOptedOut        bool          `json:"EMAIL_OPTED_OUT"`
		CustomFields         []CustomField `json:"CUSTOMFIELDS"`
	}{
		p.ProjectID,
		p.Salutation,
		p.FirstName,
		p.LastName,
		p.ImageURL,
		p.Background,
		p.OwnerUserID,
		p.SocialLinkedin,
		p.SocialFacebook,
		p.SocialTwitter,
		p.DateOfBirth,
		p.Phone,
		p.PhoneHome,
		p.PhoneMobile,
		p.PhoneOther,
		p.PhoneAssistant,
		p.PhoneFax,
		p.EmailAddress,
		p.AssistantName,
		p.AddressMailStreet,
		p.AddressMailCity,
		p.AddressMailState,
		p.AddressMailPostcode,
		p.AddressMailCountry,
		p.AddressOtherStreet,
		p.AddressOtherCity,
		p.AddressOtherState,
		p.AddressOtherPostcode,
		p.OrganisationID,
		p.Title,
		p.EmailOptedOut,
		p.CustomFields,
	}
}

// GetProject returns a specific project
//
func (i *Insightly) GetProject(projectID int) (*Project, *errortools.Error) {
	endpoint := fmt.Sprintf("Projects/%v", projectID)

	project := Project{}

	_, _, e := i.get(endpoint, nil, &project)
	if e != nil {
		return nil, e
	}

	return &project, nil
}

type GetProjectsFilter struct {
	UpdatedAfter *time.Time
	Field        *struct {
		FieldName  string
		FieldValue string
	}
}

// GetProjects returns all projects
//
func (i *Insightly) GetProjects(filter *GetProjectsFilter) (*[]Project, *errortools.Error) {
	searchString := "?"
	searchFilter := []string{}

	if filter != nil {
		if filter.UpdatedAfter != nil {
			from := filter.UpdatedAfter.Format(ISO8601Format)
			searchFilter = append(searchFilter, fmt.Sprintf("updated_after_utc=%s&", from))
		}

		if filter.Field != nil {
			searchFilter = append(searchFilter, fmt.Sprintf("field_name=%s&field_value=%s&", filter.Field.FieldName, filter.Field.FieldValue))
		}
	}

	if len(searchFilter) > 0 {
		searchString = "/Search?" + strings.Join(searchFilter, "&")
	}

	endpointStr := "Projects%sskip=%s&top=%s"
	skip := 0
	top := 100
	rowCount := top

	projects := []Project{}

	for rowCount >= top {
		endpoint := fmt.Sprintf(endpointStr, searchString, strconv.Itoa(skip), strconv.Itoa(top))
		//fmt.Println(endpoint)

		cs := []Project{}

		_, _, e := i.get(endpoint, nil, &cs)
		if e != nil {
			return nil, e
		}

		projects = append(projects, cs...)

		rowCount = len(cs)
		//rowCount = 0
		skip += top
	}

	if len(projects) == 0 {
		projects = nil
	}

	return &projects, nil
}

// CreateProject creates a new contract
//
func (i *Insightly) CreateProject(project *Project) (*Project, *errortools.Error) {
	if project == nil {
		return nil, nil
	}

	endpoint := "Projects"

	projectNew := Project{}

	_, _, e := i.post(endpoint, project.prepareMarshal(), &projectNew)
	if e != nil {
		return nil, e
	}

	return &projectNew, nil
}

// UpdateProject updates an existing contract
//
func (i *Insightly) UpdateProject(project *Project) (*Project, *errortools.Error) {
	if project == nil {
		return nil, nil
	}

	endpoint := "Projects"

	projectUpdated := Project{}

	_, _, e := i.put(endpoint, project.prepareMarshal(), &projectUpdated)
	if e != nil {
		return nil, e
	}

	return &projectUpdated, nil
}

// DeleteProject deletes a specific project
//
func (i *Insightly) DeleteProject(projectID int) *errortools.Error {
	endpoint := fmt.Sprintf("Projects/%v", projectID)

	_, _, e := i.delete(endpoint, nil, nil)
	if e != nil {
		return e
	}

	return nil
}
