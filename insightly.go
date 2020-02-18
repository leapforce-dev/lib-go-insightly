package insightly

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	types "github.com/Leapforce-nl/go_types"
)

// type
//
type Insightly struct {
	//RelationTypes RelationTypes
	//Organisations []Organisation
	//Contacts      []Contact
	Token  string
	ApiUrl string
	//OnlyPushToEO  bool
	//FromTimestamp time.Time
	// geo
	//Geo               *geo.Geo
	//BigQuery          *bigquerytools.BigQuery
	//BigQueryDataset   string
	//BigQueryTableName string
	//IsLive            bool
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

	return nil
}

// UpdateOrganisationRemovePushToEO remove PushToEo ( = true) custom value for specified organisation
//
/*func (i *Insightly) UpdateOrganisationRemovePushToEO(o *Organisation) error {
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
}*/

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

// generic Put method
//
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

	return nil
}
