package insightly

import (
	"fmt"
	"strconv"
	"time"

	errortools "github.com/leapforce-libraries/go_errortools"
)

// User stores User from Insightly
//
type User struct {
	UserID                 int    `json:"USER_ID"`
	ContactID              int    `json:"CONTACT_ID"`
	FirstName              string `json:"FIRST_NAME"`
	LastName               string `json:"LAST_NAME"`
	TimezoneID             string `json:"TIMEZONE_ID"`
	EmailAddress           string `json:"EMAIL_ADDRESS"`
	EmailDropboxIdentifier string `json:"EMAIL_DROPBOX_IDENTIFIER"`
	EmailDropboxAddress    string `json:"EMAIL_DROPBOX_ADDRESS"`
	Administrator          bool   `json:"ADMINISTRATOR"`
	AccountOwner           bool   `json:"ACCOUNT_OWNER"`
	Active                 bool   `json:"ACTIVE"`
	DateCreatedUTC         string `json:"DATE_CREATED_UTC"`
	DateUpdatedUTC         string `json:"DATE_UPDATED_UTC"`
	UserCurrency           string `json:"USER_CURRENCY"`
	ContactDisplay         string `json:"CONTACT_DISPLAY"`
	ContactOrder           string `json:"CONTACT_ORDER"`
	TaskWeekStart          int    `json:"TASK_WEEK_START"`
	InstanceID             int    `json:"INSTANCE_ID"`
	DateCreatedT           *time.Time
	DateUpdatedT           *time.Time
}

// GetUsers returns all users
//
func (i *Insightly) GetUsers() ([]User, *errortools.Error) {
	return i.GetUsersInternal()
}

// GetUsersInternal is the generic function retrieving users from Insightly
//
func (i *Insightly) GetUsersInternal() ([]User, *errortools.Error) {
	urlStr := "Users?skip=%s&top=%s"
	skip := 0
	top := 500
	rowCount := top

	users := []User{}

	for rowCount >= top {
		url := fmt.Sprintf(urlStr, strconv.Itoa(skip), strconv.Itoa(top))
		//fmt.Println(url)

		ls := []User{}

		_, _, e := i.get(url, nil, &ls)
		if e != nil {
			return nil, e
		}

		for _, l := range ls {
			l.parseDates()
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

func (l *User) parseDates() {
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
