package insightly

import (
	"fmt"
	"strconv"
	"time"

	errortools "github.com/leapforce-libraries/go_errortools"
)

// Team stores Team from Insightly
//
type Team struct {
	TeamID         int          `json:"TEAM_ID"`
	TeamName       string       `json:"TEAM_NAME"`
	AnonymousTeam  bool         `json:"ANONYMOUS_TEAM"`
	DateCreatedUTC string       `json:"DATE_CREATED_UTC"`
	DateUpdatedUTC string       `json:"DATE_UPDATED_UTC"`
	TeamMembers    []TeamMember `json:"TEAMMEMBERS"`
	DateCreatedT   *time.Time
	DateUpdatedT   *time.Time
}

// GetTeams returns all teams
//
func (i *Insightly) GetTeams() ([]Team, *errortools.Error) {
	return i.GetTeamsInternal()
}

// GetTeamsInternal is the generic function retrieving teams from Insightly
//
func (i *Insightly) GetTeamsInternal() ([]Team, *errortools.Error) {
	urlStr := "Teams?skip=%s&top=%s"
	skip := 0
	top := 500
	rowCount := top

	teams := []Team{}

	for rowCount >= top {
		url := fmt.Sprintf(urlStr, strconv.Itoa(skip), strconv.Itoa(top))
		//fmt.Println(url)

		ls := []Team{}

		_, _, e := i.get(url, nil, &ls)
		if e != nil {
			return nil, e
		}

		for _, l := range ls {
			l.parseDates()
			teams = append(teams, l)
		}

		rowCount = len(ls)
		skip += top
	}

	if len(teams) == 0 {
		teams = nil
	}

	return teams, nil
}

func (l *Team) parseDates() {
	// parse DATE_CREATED_UTC to time.Time
	if l.DateCreatedUTC != "" {
		t, _ := time.Parse("2006-01-02 15:04:05 +0000 UTC", l.DateCreatedUTC+" +0000 UTC")
		//errortools.Fatal(err)
		l.DateCreatedT = &t
	}

	// parse DATE_UPDATED_UTC to time.Time
	if l.DateUpdatedUTC != "" {
		t, _ := time.Parse("2006-01-02 15:04:05 +0000 UTC", l.DateUpdatedUTC+" +0000 UTC")
		//errortools.Fatal(err)
		l.DateUpdatedT = &t
	}
}
