package insightly

import (
	"fmt"
	"strconv"
	"strings"

	errortools "github.com/leapforce-libraries/go_errortools"
)

// TeamMember stores TeamMember from Insightly
//
type TeamMember struct {
	PermissionID int `json:"PERMISSION_ID"`
	TeamID       int `json:"TEAM_ID"`
	MemberUserID int `json:"MEMBER_USER_ID"`
}

type GetTeamMembersFilter struct {
}

// GetTeamMembers returns all teamMembers
//
func (i *Insightly) GetTeamMembers(filter *GetTeamMembersFilter) (*[]TeamMember, *errortools.Error) {
	searchString := "?"
	searchFilter := []string{}

	if filter != nil {
	}

	if len(searchFilter) > 0 {
		searchString = "/Search?" + strings.Join(searchFilter, "&")
	}

	endpointStr := "TeamMembers%sskip=%s&top=%s"
	skip := 0
	top := 100
	rowCount := top

	teamMembers := []TeamMember{}

	for rowCount >= top {
		endpoint := fmt.Sprintf(endpointStr, searchString, strconv.Itoa(skip), strconv.Itoa(top))
		//fmt.Println(endpoint)

		cs := []TeamMember{}

		_, _, e := i.get(endpoint, nil, &cs)
		if e != nil {
			return nil, e
		}

		teamMembers = append(teamMembers, cs...)

		rowCount = len(cs)
		//rowCount = 0
		skip += top
	}

	if len(teamMembers) == 0 {
		teamMembers = nil
	}

	return &teamMembers, nil
}
