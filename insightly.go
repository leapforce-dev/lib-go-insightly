package insightly

import (
	"encoding/json"
	"fmt"
	"geo"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"types"
)

const (
	customFieldNameRelationType      = "Relatietype__c"
	customFieldNameKvKNummer         = "KVKnummer__c"
	customFieldNameOrganizationOwner = "Organization_Owner__c"
	customFieldNameMainContactPerson = "Main_contactperson__c"
	customFieldNameInitials          = "initialen__c"
	customFieldNameGender            = "Gender__c"
	customFieldPushToEO              = "Push_to_EO__c"
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
}

// methods
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

//
// get methods
//
func (i *Insightly) GetAll() error {
	//
	// get iContacts
	//

	errContacts := i.getContacts()
	if errContacts != nil {
		return errContacts
	}

	fmt.Println("#iContacts: ", len(i.Contacts))
	//jsonString, _ := json.Marshal(Insightly.contacts)
	//fmt.Println(string(jsonString))

	//
	// get iOrganisations
	//
	err := i.getOrganisations()
	if err != nil {
		return err
	}
	fmt.Println("#iOrganisations: ", len(i.Organisations))

	//jsonString, _ := json.Marshal(Insightly.Organisations)
	//fmt.Println(string(jsonString))

	return nil
}

func (i *Insightly) getOrganisations() error {
	urlStr := "%sOrganisations?skip=%s&top=%s"
	skip := 0
	top := 500
	rowCount := 1
	pushToEOCount := 1

	for rowCount > 0 {
		url := fmt.Sprintf(urlStr, i.ApiUrl, strconv.Itoa(skip), strconv.Itoa(top))
		//fmt.Printf(url)

		os := []Organisation{}

		err := i.get(url, &os, i.Token)
		if err != nil {
			return err
		}

		for _, o := range os {
			for i := range o.CUSTOMFIELDS {
				o.CUSTOMFIELDS[i].UnmarshalValue()
			}
			o.getRelationTypeName(i.RelationTypes)
			//fmt.Println("outside:", o.RelationTypeName)
			o.KvKNummer = i.FindCustomFieldValue(o.CUSTOMFIELDS, customFieldNameKvKNummer)

			o.PushToEO = i.FindCustomFieldValueBool(o.CUSTOMFIELDS, customFieldPushToEO)
			if o.PushToEO {
				pushToEOCount++
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

		//i.Organisations.Organisations = append(i.Organisations.Organisations, os...)
	}

	/*
		for _, o := range i.Organisations {
			fmt.Println("KvK:")
			fmt.Println(o.KvKNummer)
		}*/
	fmt.Println("pushToEOCount (Organisation):", pushToEOCount)

	return nil
}

func (i *Insightly) getContacts() error {
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

		err := i.get(url, &cs, i.Token)
		if err != nil {
			return err
		}

		for _, c := range cs {
			for ii := range c.CUSTOMFIELDS {
				c.CUSTOMFIELDS[ii].UnmarshalValue()
				c.Initials = i.FindCustomFieldValue(c.CUSTOMFIELDS, customFieldNameInitials)
				c.Gender = c.iGenderToGender(i.FindCustomFieldValue(c.CUSTOMFIELDS, customFieldNameGender))
				c.Title = c.iGenderToTitle(i.FindCustomFieldValue(c.CUSTOMFIELDS, customFieldNameGender))
			}

			//fmt.Println(c.CUSTOMFIELDS)
			//jsonString, _ := json.Marshal(c.CUSTOMFIELDS)
			//fmt.Println(string(jsonString))
			/*b, err := strconv.ParseBool(findCustomFieldValue(c.CUSTOMFIELDS, customFieldNameMainContactPerson))
			if err == nil {
				c.IsMainContact = (b == true)
			} else {
				c.IsMainContact = false
			}*/
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

			c.PushToEO = i.FindCustomFieldValueBool(c.CUSTOMFIELDS, customFieldPushToEO)
			if c.PushToEO {
				pushToEOCount++
			}

			i.Contacts = append(i.Contacts, c)
			//fmt.Println(c.CONTACT_ID, c.LAST_NAME, "initials:", c.Initials, "gender:", c.Gender, "title:", c.Title)
		}

		rowCount = len(cs)
		skip += top

		//i.Organisations.Organisations = append(i.Organisations.Organisations, os...)
	}

	/*
		for _, o := range i.Organisations.Organisations {
			fmt.Println("KvK:")
			fmt.Println(findCustomFieldValue(o.CUSTOMFIELDS, customFieldNameKvKNummer))
		}
	*/

	fmt.Println("isMainContactCount:", isMainContactCount)
	fmt.Println("pushToEOCount (Contact):", pushToEOCount)

	return nil
}

//
// generic Get method
//

func (i *Insightly) get(url string, model interface{}, basicAuthorizationToken string) error {
	client := &http.Client{}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("accept", "application/json")
	req.Header.Set("authorization", "Basic "+basicAuthorizationToken)

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
