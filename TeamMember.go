package insightly

import (
	"fmt"
	"net/url"

	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
)

// TeamMember stores TeamMember from Service
//
type TeamMember struct {
	PermissionID int64 `json:"PERMISSION_ID"`
	TeamID       int64 `json:"TEAM_ID"`
	MemberUserID int64 `json:"MEMBER_USER_ID"`
}

type GetTeamMembersConfig struct {
	Skip       *uint64
	Top        *uint64
	CountTotal *bool
}

// GetTeamMembers returns all teamMembers
//
func (service *Service) GetTeamMembers(config *GetTeamMembersConfig) (*[]TeamMember, *errortools.Error) {
	params := url.Values{}

	endpoint := "TeamMembers"
	teamMembers := []TeamMember{}
	rowCount := uint64(0)
	top := defaultTop

	if config != nil {
		if config.Top != nil {
			top = *config.Top
		}
		if config.Skip != nil {
			service.nextSkips[endpoint] = *config.Skip
		}
		if config.CountTotal != nil {
			params.Set("count_total", fmt.Sprintf("%v", *config.CountTotal))
		}
	}

	params.Set("top", fmt.Sprintf("%v", top))

	for {
		params.Set("skip", fmt.Sprintf("%v", service.nextSkips[endpoint]))

		teamMembersBatch := []TeamMember{}

		requestConfig := go_http.RequestConfig{
			URL:           service.url(fmt.Sprintf("%s?%s", endpoint, params.Encode())),
			ResponseModel: &teamMembersBatch,
		}

		_, _, e := service.get(&requestConfig)
		if e != nil {
			return nil, e
		}

		teamMembers = append(teamMembers, teamMembersBatch...)

		if len(teamMembersBatch) < int(top) {
			delete(service.nextSkips, endpoint)
			break
		}

		service.nextSkips[endpoint] += top
		rowCount += top

		if rowCount >= service.maxRowCount {
			return &teamMembers, nil
		}
	}

	return &teamMembers, nil
}
