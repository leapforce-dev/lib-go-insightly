package insightly

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	types "github.com/Leapforce-nl/go_types"

	bigquerytools "github.com/Leapforce-nl/go_bigquerytools"
	geo "github.com/Leapforce-nl/go_geo"
)

const (
	customFieldNameRelationType            = "Relatietype__c"
	customFieldNameKvKNummer               = "KVKnummer__c"
	customFieldNameOrganizationOwner       = "Organization_Owner__c"
	customFieldNamePushToEO                = "Push_to_EO__c"
	customFieldNamePartnerSinds            = "Partner_sinds__c"
	customFieldNameBeeindigingPartnerschap = "Beindiging_partnerschap__c"
	customFieldNameOpzegdatum              = "Opzegdatum__c"
	customFieldNameAantalMedewerkers       = "Aantal_medewerkers__c"
	customFieldNameMainContactPerson       = "Main_contactperson__c"
	customFieldNameInitials                = "initialen__c"
	customFieldNameGender                  = "Gender__c"
)

// type
//
type Insightly struct {
	RelationTypes RelationTypes
	Organisations []Organisation
	Contacts      []Contact
	Token         string
	ApiUrl        string
	OnlyPushToEO  bool
	FromTimestamp time.Time
	// geo
	Geo               *geo.Geo
	BigQuery          *bigquerytools.BigQuery
	BigQueryDataset   string
	BigQueryTableName string
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

	// SamPartner
	articles := []Article{
		Article{"1-10", "S1", nil},
		Article{"11-100", "S2", nil},
		Article{"101-250", "S3", nil},
		Article{"251-500", "S4", nil},
		Article{"500+", "S5", nil},
	}
	i.RelationTypes.Append("In kind partners", 1, "SamPartner", articles)
	// Koploper
	articles = []Article{
		Article{"1-10", "K1", nil},
		Article{"11-100", "K2", nil},
		Article{"101-250", "K3", nil},
		Article{"251-500", "K4", nil},
		Article{"500+", "K5", nil},
	}
	i.RelationTypes.Append("Koploper", 2, "Koploper", articles)
	// GBN
	articles = []Article{
		Article{"500+", "GBN", nil},
	}
	i.RelationTypes.Append("GBN", 3, "GBN", articles)
	// Netwerkpartner
	articles = []Article{
		Article{"1-10", "N1", nil},
		Article{"11-100", "N2", nil},
		Article{"101-250", "N3", nil},
		Article{"251-500", "N4", nil},
		Article{"500+", "N5", nil},
	}
	i.RelationTypes.Append("Netwerkpartner", 4, "Netwerk", articles)
	articles = []Article{
		Article{"500+", "GBN", nil},
	}
	// Partner
	articles = []Article{
		Article{"1-10", "P1", nil},
		Article{"11-100", "P2", nil},
		Article{"101-250", "P3", nil},
		Article{"251-500", "P4", nil},
		Article{"500+", "P5", nil},
	}
	i.RelationTypes.Append("Partner", 5, "Partner", articles)
	// Opgezegd
	i.RelationTypes.Append("Opgezegd", 6, "", nil)

	i.Geo = new(geo.Geo)
	i.Geo.BigQuery = i.BigQuery
	i.Geo.BigQueryDataset = i.BigQueryDataset
	i.Geo.BigQueryTablenameCountries = i.BigQueryTableName

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
		fmt.Println("ERROR in UpdateOrganisationRemovePushToEO:", err)
		fmt.Println("url:", urlStr)
		return err
	}

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

/*
func (i *Insightly) ToExactOnline(o *Organisation) bool {
	if o.RelationType == nil {
		return false
	}
	if o.KvKNummer == "" {
		return false
	}
	if i.OnlyPushToEO {
		if o.PushToEO {
			//fmt.Println("ToExactOnline 1")
			return true
		} else {
			return false
		}
	}
	if o.DateUpdated.After(i.FromTimestamp) {
		//fmt.Println("ToExactOnline 2", o.DateUpdated, i.FromTimestamp)
		return true
	}
	if o.MainContact != nil {
		if o.MainContact.DateUpdated.After(i.FromTimestamp) {
			//fmt.Println("ToExactOnline 3", o.MainContact.DateUpdated, i.FromTimestamp)
			return true
		}
	}

	return false
}*/
