package insightly

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strconv"
	"time"

	errortools "github.com/leapforce-libraries/go_errortools"
	utilities "github.com/leapforce-libraries/go_utilities"
)

const (
	apiName string = "Insightly"
	apiURL  string = "https://api.%s.insightly.com/v3.1"
)

// type
//
type Insightly struct {
	pod                   string
	token                 string
	client                http.Client
	maxRetries            uint
	secondsBetweenRetries uint32
	rateLimitRemaining    *int64
	retryAt               *time.Time
}

type InsightlyConfig struct {
	Pod                   string
	APIKey                string
	MaxRetries            *uint
	SecondsBetweenRetries *uint32
}

func NewInsightly(config InsightlyConfig) (*Insightly, *errortools.Error) {
	i := new(Insightly)

	if config.Pod == "" {
		return nil, errortools.ErrorMessage("Insightly Pod not provided")
	}
	i.pod = config.Pod

	if config.APIKey == "" {
		return nil, errortools.ErrorMessage("Insightly API Key not provided")
	}
	i.token = base64.URLEncoding.EncodeToString([]byte(config.APIKey))

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

func (ins *Insightly) baseURL() string {
	return fmt.Sprintf(apiURL, ins.pod)
}

func (ins *Insightly) httpRequest(httpMethod string, endpoint string, requestBody interface{}, responseModel interface{}) (*http.Request, *http.Response, *errortools.Error) {
	// check rate limit
	if ins.rateLimitRemaining != nil {
		if *ins.rateLimitRemaining == 0 {
			if ins.retryAt == nil {
				return nil, nil, errortools.ErrorMessage("Rate limit exceeded but RetryAt unknown.")
			}

			duration := ins.retryAt.Sub(time.Now())

			if duration > 0 {
				errortools.CaptureInfo(fmt.Sprintf("Rate limit exceeded, waiting %v ms.", duration.Milliseconds()))
				time.Sleep(duration)
			}
		}
	}

	e := new(errortools.Error)

	url := fmt.Sprintf("%s/%s", ins.baseURL(), endpoint)

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
		// Read RateLimit headers
		rateLimitRemaining, err := strconv.ParseInt(response.Header.Get("X-RateLimit-Remaining"), 10, 64)
		if err == nil {
			ins.rateLimitRemaining = &rateLimitRemaining
		} else {
			ins.rateLimitRemaining = nil
		}
		retryAfter, err := strconv.ParseInt(response.Header.Get("Retry-After"), 10, 64)
		if err == nil {
			retryAt := time.Now().Add(time.Duration(retryAfter) * time.Second)
			ins.retryAt = &retryAt
		} else {
			ins.retryAt = nil
		}

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
		if response != nil {

			defer response.Body.Close()

			b, err := ioutil.ReadAll(response.Body)
			if err == nil {
				e.SetMessage(string(b))
			}
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
			fmt.Println(err)
			e.SetMessage(err)
			return request, response, e
		}
	}

	return request, response, nil
}

// generic Get method
//

func (ins *Insightly) get(endpoint string, requestBody interface{}, responseModel interface{}) (*http.Request, *http.Response, *errortools.Error) {
	return ins.httpRequest(http.MethodGet, endpoint, requestBody, responseModel)
}

func (ins *Insightly) post(endpoint string, requestBody interface{}, responseModel interface{}) (*http.Request, *http.Response, *errortools.Error) {
	return ins.httpRequest(http.MethodPost, endpoint, requestBody, responseModel)
}

func (ins *Insightly) put(endpoint string, requestBody interface{}, responseModel interface{}) (*http.Request, *http.Response, *errortools.Error) {
	return ins.httpRequest(http.MethodPut, endpoint, requestBody, responseModel)
}

func (ins *Insightly) delete(endpoint string, requestBody interface{}, responseModel interface{}) (*http.Request, *http.Response, *errortools.Error) {
	return ins.httpRequest(http.MethodDelete, endpoint, requestBody, responseModel)
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
	fmt.Println(string(b))
	if err != nil {
		return errortools.ErrorMessage(err)
	}

	err = json.Unmarshal(b, &errorModel)
	if err != nil {
		return errortools.ErrorMessage(err)
	}

	return nil
}
