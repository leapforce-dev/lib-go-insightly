package insightly

import (
	"fmt"
	"strconv"

	errortools "github.com/leapforce-libraries/go_errortools"
)

// TeamMember stores TeamMember from Insightly
//
type TeamMember struct {
	PermissionID int `json:"PERMISSION_ID"`
	TeamID       int `json:"TEAM_ID"`
	MemberUserID int `json:"MEMBER_USER_ID"`
}

// GetTeamMembers returns all TeamMembers
//
func (i *Insightly) GetTeamMembers() ([]TeamMember, *errortools.Error) {
	return i.GetTeamMembersInternal()
}

// GetTeamMembersInternal is the generic function retrieving TeamMembers from Insightly
//
func (i *Insightly) GetTeamMembersInternal() ([]TeamMember, *errortools.Error) {
	urlStr := "TeamMembers?skip=%s&top=%s"
	skip := 0
	top := 500
	rowCount := top

	teamMembers := []TeamMember{}

	for rowCount >= top {
		url := fmt.Sprintf(urlStr, strconv.Itoa(skip), strconv.Itoa(top))
		//fmt.Println(url)

		oc := []TeamMember{}

		_, _, e := i.get(url, nil, &oc)
		if e != nil {
			return nil, e
		}

		for _, o := range oc {
			teamMembers = append(teamMembers, o)
		}

		rowCount = len(oc)
		skip += top
	}

	if len(teamMembers) == 0 {
		teamMembers = nil
	}

	return teamMembers, nil
}
