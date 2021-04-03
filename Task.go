package insightly

import (
	"fmt"
	"net/url"
	"time"

	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
	i_types "github.com/leapforce-libraries/go_insightly/types"
)

// Task stores Task from Service
//
type Task struct {
	TaskID            int64                   `json:"TASK_ID"`
	Title             string                  `json:"TITLE"`
	CategoryID        *int64                  `json:"CATEGORY_ID"`
	DueDate           *i_types.DateTimeString `json:"DUE_DATE"`
	CompletedDateUTC  *i_types.DateTimeString `json:"COMPLETED_DATE_UTC"`
	Completed         bool                    `json:"COMPLETED"`
	Details           *string                 `json:"DETAILS"`
	Status            string                  `json:"STATUS"`
	Priority          int64                   `json:"PRIORITY"`
	PercentComplete   int64                   `json:"PERCENT_COMPLETE"`
	StartDate         i_types.DateTimeString  `json:"START_DATE"`
	MilestoneID       *int64                  `json:"MILESTONE_ID"`
	ResponsibleUserID int64                   `json:"RESPONSIBLE_USER_ID"`
	OwnerUserID       int64                   `json:"OWNER_USER_ID"`
	DateCreatedUTC    i_types.DateTimeString  `json:"DATE_CREATED_UTC"`
	DateUpdatedUTC    i_types.DateTimeString  `json:"DATE_UPDATED_UTC"`
	EmailID           *int64                  `json:"EMAIL_ID"`
	ProjectID         *int64                  `json:"PROJECT_ID"`
	ReminderDateUTC   *i_types.DateTimeString `json:"REMINDER_DATE_UTC"`
	ReminderSent      bool                    `json:"REMINDER_SENT"`
	OwnerVisible      bool                    `json:"OWNER_VISIBLE"`
	StageID           *int64                  `json:"STAGE_ID"`
	AssignedByUserID  *int64                  `json:"ASSIGNED_BY_USER_ID"`
	ParentTaskID      *int64                  `json:"PARENT_TASK_ID"`
	Recurrence        *string                 `json:"RECURRENCE"`
	OpportunityID     *int64                  `json:"OPPORTUNITY_ID"`
	AssignedTeamID    *int64                  `json:"ASSIGNED_TEAM_ID"`
	AssignedDateUTC   *i_types.DateTimeString `json:"ASSIGNED_DATE_UTC"`
	CreatedUserID     int64                   `json:"CREATED_USER_ID"`
	CustomFields      *CustomFields           `json:"CUSTOMFIELDS"`
}

// GetTask returns a specific task
//
func (service *Service) GetTask(taskID int64) (*Task, *errortools.Error) {
	task := Task{}

	requestConfig := go_http.RequestConfig{
		URL:           service.url(fmt.Sprintf("Tasks/%v", taskID)),
		ResponseModel: &task,
	}
	_, _, e := service.get(&requestConfig)
	if e != nil {
		return nil, e
	}

	return &task, nil
}

type GetTasksConfig struct {
	Skip         *uint64
	Top          *uint64
	Brief        *bool
	CountTotal   *bool
	UpdatedAfter *time.Time
	FieldFilter  *FieldFilter
}

// GetTasks returns all tasks
//
func (service *Service) GetTasks(config *GetTasksConfig) (*[]Task, *errortools.Error) {
	params := url.Values{}

	endpoint := "Tasks"
	tasks := []Task{}
	rowCount := uint64(0)
	top := defaultTop
	isSearch := false

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
		if config.UpdatedAfter != nil {
			isSearch = true
			params.Set("updated_after_utc", fmt.Sprintf("%v", config.UpdatedAfter.Format(DateTimeFormat)))
		}
		if config.FieldFilter != nil {
			isSearch = true
			params.Set("field_name", config.FieldFilter.FieldName)
			params.Set("field_value", config.FieldFilter.FieldValue)
		}
	}

	if isSearch {
		endpoint += "/Search"
	}

	params.Set("top", fmt.Sprintf("%v", top))

	for true {
		params.Set("skip", fmt.Sprintf("%v", service.nextSkips[endpoint]))
		tasksBatch := []Task{}

		requestConfig := go_http.RequestConfig{
			URL:           service.url(fmt.Sprintf("%s?%s", endpoint, params.Encode())),
			ResponseModel: &tasksBatch,
		}
		_, _, e := service.get(&requestConfig)
		if e != nil {
			return nil, e
		}

		tasks = append(tasks, tasksBatch...)

		if len(tasksBatch) < int(top) {
			delete(service.nextSkips, endpoint)
			break
		}

		service.nextSkips[endpoint] += top
		rowCount += top

		if rowCount >= service.maxRowCount {
			return &tasks, nil
		}
	}

	return &tasks, nil
}
