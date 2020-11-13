package insightly

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	types "github.com/leapforce-libraries/go_types"
)

// type
//
type Insightly struct {
	token  string
	apiURL string
}

// Init initializes all settings in the Insightly struct
//
/*func (i *Insightly) Init() error {
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
}*/
func New(apiURL string, token string) (*Insightly, error) {
	i := new(Insightly)

	if apiURL == "" {
		return nil, &types.ErrorString{"Insightly ApiUrl not provided"}
	}
	if token == "" {
		return nil, &types.ErrorString{"Insightly Token not provided"}
	}

	i.apiURL = apiURL
	i.token = token

	if !strings.HasSuffix(i.apiURL, "/") {
		i.apiURL = i.apiURL + "/"
	}

	return i, nil
}

// generic Get method
//
func (i *Insightly) Get(url string, model interface{}) error {
	client := &http.Client{}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("accept", "application/json")
	req.Header.Set("authorization", "Basic "+i.token)

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
	req.Header.Set("authorization", "Basic "+i.token)

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
