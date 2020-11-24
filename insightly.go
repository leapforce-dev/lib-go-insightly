package insightly

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	errortools "github.com/leapforce-libraries/go_errortools"
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
func New(apiURL string, token string) (*Insightly, *errortools.Error) {
	i := new(Insightly)

	if apiURL == "" {
		return nil, errortools.ErrorMessage("Insightly ApiUrl not provided")
	}
	if token == "" {
		return nil, errortools.ErrorMessage("Insightly Token not provided")
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
func (i *Insightly) Get(url string, model interface{}) *errortools.Error {
	client := &http.Client{}

	e := new(errortools.Error)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	e.SetRequest(req)
	if err != nil {
		e.SetMessage(err)
		return e
	}
	req.Header.Set("accept", "application/json")
	req.Header.Set("authorization", fmt.Sprintf("Basic %s", i.token))

	// Send out the HTTP request
	res, err := client.Do(req)
	e.SetResponse(res)
	if err != nil {
		e.SetMessage(err)
		return e
	}

	defer res.Body.Close()

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		e.SetMessage(err)
		return e
	}

	err = json.Unmarshal(b, &model)
	if err != nil {
		e.SetMessage(err)
		return e
	}

	return nil
}

// generic Put method
//
func (i *Insightly) Put(url string, json []byte) *errortools.Error {
	client := &http.Client{}

	e := new(errortools.Error)

	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(json))
	e.SetRequest(req)
	if err != nil {
		e.SetMessage(err)
		return e
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("authorization", "Basic "+i.token)

	// Send out the HTTP request
	res, err := client.Do(req)
	e.SetResponse(res)
	if err != nil {
		e.SetMessage(err)
		return e
	}

	// Check HTTP StatusCode
	if res.StatusCode < 200 || res.StatusCode > 299 {
		e.SetMessage(fmt.Sprintf("Server returned statuscode %v: %s", res.StatusCode, err.Error()))
		return e
	}

	return nil
}
