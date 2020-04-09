package insightly

import (
	"fmt"
	"strconv"
	"time"
)

// User stores User from Insightly
//
type User struct {
	USER_ID                  int    `json:"USER_ID"`
	CONTACT_ID               int    `json:"CONTACT_ID"`
	FIRST_NAME               string `json:"FIRST_NAME"`
	LAST_NAME                string `json:"LAST_NAME"`
	TIMEZONE_ID              string `json:"TIMEZONE_ID"`
	EMAIL_ADDRESS            string `json:"EMAIL_ADDRESS"`
	EMAIL_DROPBOX_IDENTIFIER string `json:"EMAIL_DROPBOX_IDENTIFIER"`
	EMAIL_DROPBOX_ADDRESS    string `json:"EMAIL_DROPBOX_ADDRESS"`
	ADMINISTRATOR            bool   `json:"ADMINISTRATOR"`
	ACCOUNT_OWNER            bool   `json:"ACCOUNT_OWNER"`
	ACTIVE                   bool   `json:"ACTIVE"`
	DATE_CREATED_UTC         string `json:"DATE_CREATED_UTC"`
	DATE_UPDATED_UTC         string `json:"DATE_UPDATED_UTC"`
	USER_CURRENCY            string `json:"USER_CURRENCY"`
	CONTACT_DISPLAY          string `json:"CONTACT_DISPLAY"`
	CONTACT_ORDER            string `json:"CONTACT_ORDER"`
	TASK_WEEK_START          int    `json:"TASK_WEEK_START"`
	INSTANCE_ID              int    `json:"INSTANCE_ID"`
	DateCreated              *time.Time
	DateUpdated              *time.Time
}

// GetUsers returns all users
//
func (i *Insightly) GetUsers() ([]User, error) {
	return i.GetUsersInternal()
}

// GetUsersInternal is the generic function retrieving users from Insightly
//
func (i *Insightly) GetUsersInternal() ([]User, error) {
	urlStr := "%sUsers?skip=%s&top=%s"
	skip := 0
	top := 500
	rowCount := top

	users := []User{}

	for rowCount >= top {
		url := fmt.Sprintf(urlStr, i.apiURL, strconv.Itoa(skip), strconv.Itoa(top))
		//fmt.Println(url)

		ls := []User{}

		err := i.Get(url, &ls)
		if err != nil {
			return nil, err
		}

		for _, l := range ls {
			l.ParseDates()
			users = append(users, l)
		}

		rowCount = len(ls)
		skip += top
	}

	if len(users) == 0 {
		users = nil
	}

	return users, nil
}

func (l *User) ParseDates() {
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
