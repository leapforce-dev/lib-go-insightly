package insightly

import (
	"fmt"
	"strconv"
	"time"
)

// Team stores Team from Insightly
//
type Team struct {
	TEAM_ID          int          `json:"TEAM_ID"`
	TEAM_NAME        string       `json:"TEAM_NAME"`
	ANONYMOUS_TEAM   bool         `json:"ANONYMOUS_TEAM"`
	DATE_CREATED_UTC string       `json:"DATE_CREATED_UTC"`
	DATE_UPDATED_UTC string       `json:"DATE_UPDATED_UTC"`
	TEAMMEMBERS      []TeamMember `json:"TEAMMEMBERS"`
	DateCreated      *time.Time
	DateUpdated      *time.Time
}

// GetTeams returns all teams
//
func (i *Insightly) GetTeams() ([]Team, error) {
	return i.GetTeamsInternal()
}

// GetTeamsInternal is the generic function retrieving teams from Insightly
//
func (i *Insightly) GetTeamsInternal() ([]Team, error) {
	urlStr := "%sTeams?skip=%s&top=%s"
	skip := 0
	top := 500
	rowCount := top

	teams := []Team{}

	for rowCount >= top {
		url := fmt.Sprintf(urlStr, i.apiURL, strconv.Itoa(skip), strconv.Itoa(top))
		//fmt.Println(url)

		ls := []Team{}

		err := i.Get(url, &ls)
		if err != nil {
			return nil, err
		}

		for _, l := range ls {
			l.ParseDates()
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

func (l *Team) ParseDates() {
	// parse DATE_CREATED_UTC to time.Time
	if l.DATE_CREATED_UTC != "" {
		t, _ := time.Parse("2006-01-02 15:04:05 +0000 UTC", l.DATE_CREATED_UTC+" +0000 UTC")
		//errortools.Fatal(err)
		l.DateCreated = &t
	}

	// parse DATE_UPDATED_UTC to time.Time
	if l.DATE_UPDATED_UTC != "" {
		t, _ := time.Parse("2006-01-02 15:04:05 +0000 UTC", l.DATE_UPDATED_UTC+" +0000 UTC")
		//errortools.Fatal(err)
		l.DateUpdated = &t
	}
}
