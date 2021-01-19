package insightly

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
)

// Team stores Team from Service
//
type Team struct {
	TeamID         int          `json:"TEAM_ID"`
	TeamName       string       `json:"TEAM_NAME"`
	AnonymousTeam  bool         `json:"ANONYMOUS_TEAM"`
	DateCreatedUTC DateUTC      `json:"DATE_CREATED_UTC"`
	DateUpdatedUTC DateUTC      `json:"DATE_UPDATED_UTC"`
	TeamMembers    []TeamMember `json:"TEAMMEMBERS"`
}

func (t *Team) prepareMarshal() interface{} {
	if t == nil {
		return nil
	}

	return &struct {
		TeamID        int          `json:"TEAM_ID"`
		TeamName      string       `json:"TEAM_NAME"`
		AnonymousTeam bool         `json:"ANONYMOUS_TEAM"`
		TeamMembers   []TeamMember `json:"TEAMMEMBERS"`
	}{
		t.TeamID,
		t.TeamName,
		t.AnonymousTeam,
		t.TeamMembers,
	}
}

// GetTeam returns a specific team
//
func (service *Service) GetTeam(teamID int) (*Team, *errortools.Error) {
	team := Team{}

	requestConfig := go_http.RequestConfig{
		URL:           service.url(fmt.Sprintf("Teams/%v", teamID)),
		ResponseModel: &team,
	}
	_, _, e := service.get(&requestConfig)
	if e != nil {
		return nil, e
	}

	return &team, nil
}

type GetTeamsFilter struct {
	UpdatedAfter *time.Time
	Field        *struct {
		FieldName  string
		FieldValue string
	}
}

// GetTeams returns all teams
//
func (service *Service) GetTeams(filter *GetTeamsFilter) (*[]Team, *errortools.Error) {
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

	endpointStr := "Teams%sskip=%s&top=%s"
	skip := 0
	top := 100
	rowCount := top

	teams := []Team{}

	for rowCount >= top {
		_teams := []Team{}

		requestConfig := go_http.RequestConfig{
			URL:           service.url(fmt.Sprintf(endpointStr, searchString, strconv.Itoa(skip), strconv.Itoa(top))),
			ResponseModel: &_teams,
		}
		_, _, e := service.get(&requestConfig)
		if e != nil {
			return nil, e
		}

		teams = append(teams, _teams...)

		rowCount = len(_teams)
		//rowCount = 0
		skip += top
	}

	if len(teams) == 0 {
		teams = nil
	}

	return &teams, nil
}

// CreateTeam creates a new contract
//
func (service *Service) CreateTeam(team *Team) (*Team, *errortools.Error) {
	if team == nil {
		return nil, nil
	}

	teamNew := Team{}

	requestConfig := go_http.RequestConfig{
		URL:           service.url("Teams"),
		BodyModel:     team.prepareMarshal(),
		ResponseModel: &teamNew,
	}
	_, _, e := service.post(&requestConfig)
	if e != nil {
		return nil, e
	}

	return &teamNew, nil
}

// UpdateTeam updates an existing contract
//
func (service *Service) UpdateTeam(team *Team) (*Team, *errortools.Error) {
	if team == nil {
		return nil, nil
	}

	teamUpdated := Team{}

	requestConfig := go_http.RequestConfig{
		URL:           service.url("Teams"),
		BodyModel:     team.prepareMarshal(),
		ResponseModel: &teamUpdated,
	}
	_, _, e := service.put(&requestConfig)
	if e != nil {
		return nil, e
	}

	return &teamUpdated, nil
}

// DeleteTeam deletes a specific team
//
func (service *Service) DeleteTeam(teamID int) *errortools.Error {
	requestConfig := go_http.RequestConfig{
		URL: service.url(fmt.Sprintf("Teams/%v", teamID)),
	}
	_, _, e := service.delete(&requestConfig)
	if e != nil {
		return e
	}

	return nil
}
