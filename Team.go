package insightly

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	errortools "github.com/leapforce-libraries/go_errortools"
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
	endpoint := fmt.Sprintf("Teams/%v", teamID)

	team := Team{}

	_, _, e := service.get(endpoint, nil, &team)
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

	endpointStr := "Teams%sskip=%s&top=%s"
	skip := 0
	top := 100
	rowCount := top

	teams := []Team{}

	for rowCount >= top {
		endpoint := fmt.Sprintf(endpointStr, searchString, strconv.Itoa(skip), strconv.Itoa(top))
		//fmt.Println(endpoint)

		cs := []Team{}

		_, _, e := service.get(endpoint, nil, &cs)
		if e != nil {
			return nil, e
		}

		teams = append(teams, cs...)

		rowCount = len(cs)
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

	endpoint := "Teams"

	teamNew := Team{}

	_, _, e := service.post(endpoint, team.prepareMarshal(), &teamNew)
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

	endpoint := "Teams"

	teamUpdated := Team{}

	_, _, e := service.put(endpoint, team.prepareMarshal(), &teamUpdated)
	if e != nil {
		return nil, e
	}

	return &teamUpdated, nil
}

// DeleteTeam deletes a specific team
//
func (service *Service) DeleteTeam(teamID int) *errortools.Error {
	endpoint := fmt.Sprintf("Teams/%v", teamID)

	_, _, e := service.delete(endpoint, nil, nil)
	if e != nil {
		return e
	}

	return nil
}
