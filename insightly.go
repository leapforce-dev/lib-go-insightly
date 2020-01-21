package insightly

import (
	"bytes"
	"encoding/json"
	"errortools"
	"fmt"
	"geo"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
	"types"
)

const (
	customFieldNameRelationType            = "Relatietype__c"
	customFieldNameKvKNummer               = "KVKnummer__c"
	customFieldNameOrganizationOwner       = "Organization_Owner__c"
	customFieldNameMainContactPerson       = "Main_contactperson__c"
	customFieldNameInitials                = "initialen__c"
	customFieldNameGender                  = "Gender__c"
	customFieldNamePushToEO                = "Push_to_EO__c"
	customFieldNamePartnerSinds            = "Partner_sinds__c"
	customFieldNameBeeindigingPartnerschap = "Beindiging_partnerschap__c"
)

// type
//
type Insightly struct {
	RelationTypes RelationTypes
	Organisations []Organisation
	Contacts      []Contact
	Token         string
	ApiUrl        string
	Geo           *geo.Geo
	OnlyPushToEO  bool
	FromTimestamp time.Time
}

// Init initializes all settings in the Insightly struct
//
func (i *Insightly) Init() error {
	if i.ApiUrl == "" {
		return &types.ErrorString{"Insightly ApiUrl not provided"}
	}
	if i.Token == "" {
		return &types.ErrorString{"Insightly Token not provided"}
	}

	if !strings.HasSuffix(i.ApiUrl, "/") {
		i.ApiUrl = i.ApiUrl + "/"
	}

	i.OnlyPushToEO = false
	i.FromTimestamp, _ = time.Parse("2006-01-02", "1800-01-01")

	i.RelationTypes.Append("In kind partners", 1)
	i.RelationTypes.Append("Koploper", 2)
	i.RelationTypes.Append("GBN", 3)
	i.RelationTypes.Append("Netwerkpartner", 4)
	i.RelationTypes.Append("Partner", 5)
	i.RelationTypes.Append("Opgezegd", 6)

	i.Geo = new(geo.Geo)
	i.Geo.InitBigQuery()

	return nil
}

func (i *Insightly) GetOrganisations() error {
	urlStr := "%sOrganisations?skip=%s&top=%s"
	skip := 0
	top := 500
	rowCount := 1
	pushToEOCount := 1

	for rowCount > 0 {
		url := fmt.Sprintf(urlStr, i.ApiUrl, strconv.Itoa(skip), strconv.Itoa(top))
		//fmt.Printf(url)

		os := []Organisation{}

		err := i.Get(url, &os)
		if err != nil {
			return err
		}

		for _, o := range os {
			// unmarshal custom fields
			for ii := range o.CUSTOMFIELDS {
				o.CUSTOMFIELDS[ii].UnmarshalValue()
			}

			// get RelationTypeName from custom field
			o.GetRelationType(&i.RelationTypes)

			// parse DATE_UPDATED_UTC to time.Time
			t, err := time.Parse("2006-01-02 15:04:05 +0000 UTC", o.DATE_UPDATED_UTC+" +0000 UTC")
			errortools.Fatal(err)
			o.DateUpdated = t

			// get KvKNummer from custom field
			o.KvKNummer = i.FindCustomFieldValue(o.CUSTOMFIELDS, customFieldNameKvKNummer)

			// get PushToEO from custom field
			o.PushToEO = i.FindCustomFieldValueBool(o.CUSTOMFIELDS, customFieldNamePushToEO)
			if o.PushToEO {
				pushToEOCount++
			}

			// get PushToEO from custom field
			o.BeeindigingPartnerschap = i.FindCustomFieldValue(o.CUSTOMFIELDS, customFieldNameBeeindigingPartnerschap)
			if o.BeeindigingPartnerschap != "" {
				t1, err := time.Parse("2006-01-02 15:04:05", o.BeeindigingPartnerschap)
				errortools.Fatal(err)
				o.BeeindigingPartnerschapTime = &t1

				//fmt.Println("o.BeeindigingPartnerschapTime", t1)
			}

			i.Organisations = append(i.Organisations, o)

			// find CountryId
			id, err := i.Geo.FindCountryId(o.ADDRESS_BILLING_COUNTRY, "", "", "")
			if err != nil {
				return err
			}
			o.CountryId = id
		}

		rowCount = len(os)
		skip += top
	}

	fmt.Println("pushToEOCount (Organisation):", pushToEOCount)

	return nil
}

func (i *Insightly) GetContacts() error {
	urlStr := "%sContacts?skip=%s&top=%s"
	skip := 0
	top := 500
	rowCount := 1
	isMainContactCount := 1
	pushToEOCount := 1

	for rowCount > 0 {
		url := fmt.Sprintf(urlStr, i.ApiUrl, strconv.Itoa(skip), strconv.Itoa(top))
		//fmt.Printf(url)

		cs := []Contact{}

		err := i.Get(url, &cs)
		if err != nil {
			return err
		}

		for _, c := range cs {
			// unmarshal custom fields
			for ii := range c.CUSTOMFIELDS {
				c.CUSTOMFIELDS[ii].UnmarshalValue()
			}

			// get Initials from custom field
			c.Initials = i.FindCustomFieldValue(c.CUSTOMFIELDS, customFieldNameInitials)

			// get Gender from custom field
			c.Gender = c.iGenderToGender(i.FindCustomFieldValue(c.CUSTOMFIELDS, customFieldNameGender))

			// get Title from custom field
			c.Title = c.iGenderToTitle(i.FindCustomFieldValue(c.CUSTOMFIELDS, customFieldNameGender))

			// parse DATE_UPDATED_UTC to time.Time
			t, err := time.Parse("2006-01-02 15:04:05 +0000 UTC", c.DATE_UPDATED_UTC+" +0000 UTC")
			errortools.Fatal(err)
			c.DateUpdated = t

			//fmt.Println("o.DATE_UPDATED_UTC", c.DATE_UPDATED_UTC, "o.DateUpdated", c.DateUpdated, "Now", time.Now(), "Diff", time.Now().Sub(c.DateUpdated))

			// validate email
			if c.EMAIL_ADDRESS != "" {
				err := ValidateFormat(c.EMAIL_ADDRESS)
				if err != nil {
					fmt.Println("invalid emailadress:", c.EMAIL_ADDRESS)
				} else {
					c.Email = c.EMAIL_ADDRESS
				}
			}

			c.IsMainContact = i.FindCustomFieldValueBool(c.CUSTOMFIELDS, customFieldNameMainContactPerson)
			if c.IsMainContact {
				isMainContactCount++
			}

			c.PushToEO = i.FindCustomFieldValueBool(c.CUSTOMFIELDS, customFieldNamePushToEO)
			if c.PushToEO {
				pushToEOCount++
			}

			i.Contacts = append(i.Contacts, c)
			//fmt.Println(c.CONTACT_ID, c.LAST_NAME, "initials:", c.Initials, "gender:", c.Gender, "title:", c.Title)
		}

		rowCount = len(cs)
		skip += top
	}

	fmt.Println("isMainContactCount:", isMainContactCount)
	fmt.Println("pushToEOCount (Contact):", pushToEOCount)

	return nil
}

// UpdateOrganisationRemovePushToEO remove PushToEo ( = true) custom value for specified organisation
//
func (i *Insightly) UpdateOrganisationRemovePushToEO(o *Organisation) error {
	urlStr := "%sOrganisations"
	url := fmt.Sprintf(urlStr, i.ApiUrl)

	type CustomFieldDelete struct {
		FIELD_NAME      string
		CUSTOM_FIELD_ID string
	}

	type OrganisationID struct {
		ORGANISATION_ID int
		CUSTOMFIELDS    []CustomFieldDelete
	}

	o1 := OrganisationID{}
	o1.ORGANISATION_ID = o.ORGANISATION_ID
	o1.CUSTOMFIELDS = make([]CustomFieldDelete, 1)
	o1.CUSTOMFIELDS[0] = CustomFieldDelete{customFieldNamePushToEO, customFieldNamePushToEO}

	b, err := json.Marshal(o1)
	if err != nil {
		fmt.Println("error:", err)
	}

	err = i.Put(url, b)
	if err != nil {
		return err
	}

	//fmt.Println("unchecked:", o.ORGANISATION_ID)
	//time.Sleep(1 * time.Second)

	return nil
}

//
// generic Get method
//
func (i *Insightly) Get(url string, model interface{}) error {
	client := &http.Client{}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("accept", "application/json")
	req.Header.Set("authorization", "Basic "+i.Token)

	// Send out the HTTP request
	res, err := client.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	b, err := ioutil.ReadAll(res.Body)

	errr := json.Unmarshal(b, &model)
	if errr != nil {
		return err
	}

	return nil
}
func (i *Insightly) Put(url string, json []byte) error {
	client := &http.Client{}

	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(json))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("authorization", "Basic "+i.Token)

	// Send out the HTTP request
	res, err := client.Do(req)
	if err != nil {
		return err
	}

	// Check HTTP StatusCode
	if res.StatusCode < 200 || res.StatusCode > 299 {
		return &types.ErrorString{fmt.Sprintf("Server returned statuscode %v: %s", res.StatusCode, err.Error())}
	}

	//fmt.Println(res)

	return nil
}

func (i *Insightly) ToExactOnline(o *Organisation) bool {
	if o.RelationType == nil {
		return false
	}
	if o.KvKNummer == "" {
		return false
	}
	if i.OnlyPushToEO {
		if o.PushToEO {
			fmt.Println("ToExactOnline 1")
			return true
		} else {
			return false
		}
	}
	if o.DateUpdated.After(i.FromTimestamp) {
		fmt.Println("ToExactOnline 2", o.DateUpdated, i.FromTimestamp)
		return true
	}
	if o.MainContact != nil {
		if o.MainContact.DateUpdated.After(i.FromTimestamp) {
			fmt.Println("ToExactOnline 3", o.MainContact.DateUpdated, i.FromTimestamp)
			return true
		}
	}

	return false
}
