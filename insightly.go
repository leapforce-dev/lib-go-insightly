package insightly

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"

	errortools "github.com/leapforce-libraries/go_errortools"
	utilities "github.com/leapforce-libraries/go_utilities"
)

const (
	apiName string = "Insightly"
	apiURL  string = "https://api.insightly.com/v3.1"
)

// type
//
type Insightly struct {
	token                 string
	client                http.Client
	maxRetries            uint
	secondsBetweenRetries uint32
}

type InsightlyConfig struct {
	Token                 string
	MaxRetries            *uint
	SecondsBetweenRetries *uint32
}

func NewInsightly(config InsightlyConfig) (*Insightly, *errortools.Error) {
	i := new(Insightly)

	if config.Token == "" {
		return nil, errortools.ErrorMessage("Insightly Token not provided")
	}

	i.token = config.Token

	if config.MaxRetries != nil {
		i.maxRetries = *config.MaxRetries
	} else {
		i.maxRetries = 0
	}

	if config.SecondsBetweenRetries != nil {
		i.secondsBetweenRetries = *config.SecondsBetweenRetries
	} else {
		i.secondsBetweenRetries = 3
	}
	i.client = http.Client{}

	return i, nil
}

func (ins *Insightly) httpRequest(httpMethod string, url string, requestBody interface{}, responseModel interface{}, errorModel interface{}) (*http.Request, *http.Response, *errortools.Error) {
	e := new(errortools.Error)

	request, err := func() (*http.Request, error) {
		if requestBody == nil {
			return http.NewRequest(httpMethod, url, nil)
		} else {
			b, err := json.Marshal(requestBody)
			if err != nil {
				return nil, err
			}

			return http.NewRequest(httpMethod, url, bytes.NewBuffer(b))
		}
	}()

	e.SetRequest(request)

	if err != nil {
		e.SetMessage(err)
		return request, nil, e
	}

	// Add authorization token to header
	request.Header.Set("Authorization", fmt.Sprintf("Basic %s", ins.token))
	request.Header.Set("Accept", "application/json")
	if requestBody != nil {
		request.Header.Set("Content-Type", "application/json")
	}

	// Send out the HTTP request
	response, e := utilities.DoWithRetry(&ins.client, request, ins.maxRetries, ins.secondsBetweenRetries)

	if response != nil {
		// Check HTTP StatusCode
		if response.StatusCode < 200 || response.StatusCode > 299 {
			if e == nil {
				e = new(errortools.Error)
				e.SetRequest(request)
				e.SetResponse(response)
			}

			e.SetMessage(fmt.Sprintf("Server returned statuscode %v", response.StatusCode))
		}
	}

	if e != nil {
		if errorModel != nil {
			err2 := unmarshalError(response, errorModel)
			errortools.CaptureInfo(err2)
		}

		return request, response, e
	}

	if responseModel != nil {
		defer response.Body.Close()

		b, err := ioutil.ReadAll(response.Body)
		if err != nil {
			e.SetMessage(err)
			return request, response, e
		}

		err = json.Unmarshal(b, &responseModel)
		if err != nil {
			e.SetMessage(err)
			return request, response, e
		}
	}

	return request, response, nil
}

// generic Get method
//

func (ins *Insightly) get(url string, requestBody interface{}, responseModel interface{}) (*http.Request, *http.Response, *errortools.Error) {
	return ins.httpRequest(http.MethodGet, url, requestBody, responseModel, nil)
}

func (ins *Insightly) post(url string, requestBody interface{}, responseModel interface{}) (*http.Request, *http.Response, *errortools.Error) {
	return ins.httpRequest(http.MethodPost, url, requestBody, responseModel, nil)
}

func (ins *Insightly) put(url string, requestBody interface{}, responseModel interface{}) (*http.Request, *http.Response, *errortools.Error) {
	return ins.httpRequest(http.MethodPut, url, requestBody, responseModel, nil)
}

func (ins *Insightly) delete(url string, requestBody interface{}, responseModel interface{}) (*http.Request, *http.Response, *errortools.Error) {
	return ins.httpRequest(http.MethodDelete, url, requestBody, responseModel, nil)
}

func unmarshalError(response *http.Response, errorModel interface{}) *errortools.Error {
	if response == nil {
		return nil
	}
	if reflect.TypeOf(errorModel).Kind() != reflect.Ptr {
		return errortools.ErrorMessage("Type of errorModel must be a pointer.")
	}
	if reflect.ValueOf(errorModel).IsNil() {
		return nil
	}

	defer response.Body.Close()

	b, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return errortools.ErrorMessage(err)
	}

	err = json.Unmarshal(b, &errorModel)
	if err != nil {
		return errortools.ErrorMessage(err)
	}

	return nil
}
