package insightly

import (
	"fmt"
	"net/http"
	"net/url"

	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
	i_types "github.com/leapforce-libraries/go_insightly/types"
)

// Team stores Team from Service
//
type Team struct {
	TeamID         int64                  `json:"TEAM_ID"`
	TeamName       string                 `json:"TEAM_NAME"`
	AnonymousTeam  bool                   `json:"ANONYMOUS_TEAM"`
	DateCreatedUTC i_types.DateTimeString `json:"DATE_CREATED_UTC"`
	DateUpdatedUTC i_types.DateTimeString `json:"DATE_UPDATED_UTC"`
	TeamMembers    *[]TeamMember          `json:"TEAMMEMBERS"`
}

func (t *Team) prepareMarshal() interface{} {
	if t == nil {
		return nil
	}

	return &struct {
		TeamID        *int64        `json:"TEAM_ID"`
		TeamName      *string       `json:"TEAM_NAME"`
		AnonymousTeam *bool         `json:"ANONYMOUS_TEAM"`
		TeamMembers   *[]TeamMember `json:"TEAMMEMBERS"`
	}{
		&t.TeamID,
		&t.TeamName,
		&t.AnonymousTeam,
		t.TeamMembers,
	}
}

// GetTeam returns a specific team
//
func (service *Service) GetTeam(teamID int64) (*Team, *errortools.Error) {
	team := Team{}

	requestConfig := go_http.RequestConfig{
		Method:        http.MethodGet,
		URL:           service.url(fmt.Sprintf("Teams/%v", teamID)),
		ResponseModel: &team,
	}
	_, _, e := service.httpRequest(&requestConfig)
	if e != nil {
		return nil, e
	}

	return &team, nil
}

type GetTeamsConfig struct {
	Skip       *uint64
	Top        *uint64
	Brief      *bool
	CountTotal *bool
}

// GetTeams returns all teams
//
func (service *Service) GetTeams(config *GetTeamsConfig) (*[]Team, *errortools.Error) {
	params := url.Values{}

	endpoint := "Teams"
	teams := []Team{}
	rowCount := uint64(0)
	top := defaultTop

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
	}

	params.Set("top", fmt.Sprintf("%v", top))

	for {
		params.Set("skip", fmt.Sprintf("%v", service.nextSkips[endpoint]))

		teamsBatch := []Team{}

		requestConfig := go_http.RequestConfig{
			Method:        http.MethodGet,
			URL:           service.url(fmt.Sprintf("%s?%s", endpoint, params.Encode())),
			ResponseModel: &teamsBatch,
		}

		_, _, e := service.httpRequest(&requestConfig)
		if e != nil {
			return nil, e
		}

		teams = append(teams, teamsBatch...)

		if len(teamsBatch) < int(top) {
			delete(service.nextSkips, endpoint)
			break
		}

		service.nextSkips[endpoint] += top
		rowCount += top

		if rowCount >= service.maxRowCount {
			return &teams, nil
		}
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
		Method:        http.MethodPost,
		URL:           service.url("Teams"),
		BodyModel:     team.prepareMarshal(),
		ResponseModel: &teamNew,
	}
	_, _, e := service.httpRequest(&requestConfig)
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
		Method:        http.MethodPut,
		URL:           service.url("Teams"),
		BodyModel:     team.prepareMarshal(),
		ResponseModel: &teamUpdated,
	}
	_, _, e := service.httpRequest(&requestConfig)
	if e != nil {
		return nil, e
	}

	return &teamUpdated, nil
}

// DeleteTeam deletes a specific team
//
func (service *Service) DeleteTeam(teamID int) *errortools.Error {
	requestConfig := go_http.RequestConfig{
		Method: http.MethodDelete,
		URL:    service.url(fmt.Sprintf("Teams/%v", teamID)),
	}
	_, _, e := service.httpRequest(&requestConfig)
	if e != nil {
		return e
	}

	return nil
}
