package insightly

import (
	"fmt"
	"net/http"
	"net/url"

	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
	i_types "github.com/leapforce-libraries/go_insightly/types"
)

// ActivitySet stores ActivitySet from Service
//
type ActivitySet struct {
	ActivitySetID    int64      `json:"ACTIVITYSET_ID"`
	Name             string     `json:"NAME"`
	ForContacts      bool       `json:"FOR_CONTACTS"`
	ForOrganisations bool       `json:"FOR_ORGANISATIONS"`
	ForOpportunities bool       `json:"FOR_OPPORTUNITIES"`
	ForProjects      bool       `json:"FOR_PROJECTS"`
	ForLeads         bool       `json:"FOR_LEADS"`
	OwnerUserID      int64      `json:"OWNER_USER_ID"`
	Activities       []Activity `json:"ACTIVITIES"`
}

type Activity struct {
	ActivityID          int64               `json:"ACTIVITY_ID"`
	ActivitySetID       int64               `json:"ACTIVITYSET_ID"`
	ActivityName        string              `json:"ACTIVITY_NAME"`
	ActivityDetails     string              `json:"ACTIVITY_DETAILS"`
	ActivityType        string              `json:"ACTIVITY_TYPE"`
	CategoryID          int64               `json:"CATEGORY_ID"`
	Reminder            bool                `json:"REMINDER"`
	ReminderTime        *i_types.TimeString `json:"REMINDER_TIME"`
	PubliclyVisible     bool                `json:"PUBLICLY_VISIBLE"`
	OwnerVisible        bool                `json:"OWNER_VISIBLE"`
	OwnerUserID         int64               `json:"OWNER_USER_ID"`
	ResponsibleUserID   *int64              `json:"RESPONSIBLE_USER_ID"`
	AssignedTeamID      *int64              `json:"ASSIGNED_TEAM_ID"`
	SkipSunday          bool                `json:"SKIP_SUN"`
	SkipMonday          bool                `json:"SKIP_MON"`
	SkipTuesday         bool                `json:"SKIP_TUE"`
	SkipWednesday       bool                `json:"SKIP_WED"`
	SkipThirsday        bool                `json:"SKIP_THU"`
	SkipFriday          bool                `json:"SKIP_FRI"`
	SkipSaturday        bool                `json:"SKIP_SAT"`
	DueDaysAfterStart   *int64              `json:"DUE_DAYS_AFTER_START"`
	DueDaysBeforeEnd    *int64              `json:"DUE_DAYS_BEFORE_END"`
	EventDaysAfterStart *int64              `json:"EVENT_DAYS_AFTER_START"`
	EventDaysBeforeEnd  *int64              `json:"EVENT_DAYS_BEFORE_END"`
	EventTime           *i_types.TimeString `json:"EVENT_TIME"`
	AllDay              *bool               `json:"ALL_DAY"`
	Duration            *int64              `json:"DURATION"`
}

type GetActivitySetsConfig struct {
	Skip       *uint64
	Top        *uint64
	Brief      *bool
	CountTotal *bool
}

// GetActivitySets returns all activitySets
//

func (service *Service) GetActivitySets(config *GetActivitySetsConfig) (*[]ActivitySet, *errortools.Error) {
	params := url.Values{}

	endpoint := "ActivitySets"
	activitySets := []ActivitySet{}
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

		activitySetsBatch := []ActivitySet{}

		requestConfig := go_http.RequestConfig{
			Method:        http.MethodGet,
			Url:           service.url(fmt.Sprintf("%s?%s", endpoint, params.Encode())),
			ResponseModel: &activitySetsBatch,
		}

		_, _, e := service.httpRequest(&requestConfig)
		if e != nil {
			return nil, e
		}

		activitySets = append(activitySets, activitySetsBatch...)

		if len(activitySetsBatch) < int(top) {
			delete(service.nextSkips, endpoint)
			break
		}

		service.nextSkips[endpoint] += top
		rowCount += top

		if rowCount >= service.maxRowCount {
			return &activitySets, nil
		}
	}

	return &activitySets, nil
}
