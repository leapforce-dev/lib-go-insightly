package insightly

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strconv"
	"time"

	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
)

const (
	apiName                   string = "Insightly"
	apiURL                    string = "https://api.%s.insightly.com/v3.1"
	dateTimeFormat            string = "2006-01-02T15:04:05Z"
	dateTimeFormatCustomField string = "2006-01-02 15:04:05"
	dateFormat                string = "2006-01-02"
	defaultMaxRowCount        uint64 = ^uint64(0)
	defaultTop                uint64 = 500 //max 500, see: https://api.insightly.com/v3.1/Help#!/Overview/Introduction
)

type RateLimit struct {
	Limit     *int64
	Remaining *int64
	RetryAt   *time.Time
}

type Service struct {
	pod          string
	token        string
	maxRowCount  uint64
	httpService  *go_http.Service
	rateLimit    RateLimit
	apiCallCount int64
	nextSkips    map[string]uint64
}

type ServiceConfig struct {
	Pod         string
	APIKey      string
	MaxRowCount *uint64
}

func NewService(serviceConfig *ServiceConfig) (*Service, *errortools.Error) {
	if serviceConfig == nil {
		return nil, errortools.ErrorMessage("ServiceConfig must not be a nil pointer")
	}

	if serviceConfig.Pod == "" {
		return nil, errortools.ErrorMessage("Service Pod not provided")
	}

	if serviceConfig.APIKey == "" {
		return nil, errortools.ErrorMessage("Service API Key not provided")
	}

	httpService, e := go_http.NewService(&go_http.ServiceConfig{})
	if e != nil {
		return nil, e
	}

	maxRowCount := defaultMaxRowCount
	if serviceConfig.MaxRowCount != nil {
		maxRowCount = *serviceConfig.MaxRowCount
	}

	return &Service{
		pod:         serviceConfig.Pod,
		token:       base64.URLEncoding.EncodeToString([]byte(serviceConfig.APIKey)),
		maxRowCount: maxRowCount,
		httpService: httpService,
		nextSkips:   make(map[string]uint64),
	}, nil
}

func (service *Service) httpRequest(httpMethod string, requestConfig *go_http.RequestConfig) (*http.Request, *http.Response, *errortools.Error) {
	// check rate limit
	if service.rateLimit.Remaining != nil {
		if *service.rateLimit.Remaining == 0 {
			if service.rateLimit.RetryAt == nil {
				return nil, nil, errortools.ErrorMessage("Rate limit exceeded but RetryAt unknown.")
			}

			duration := service.rateLimit.RetryAt.Sub(time.Now())

			if duration > 0 {
				errortools.CaptureInfo(fmt.Sprintf("Rate limit exceeded, waiting %v ms.", duration.Milliseconds()))
				time.Sleep(duration)
			}
		}
	}

	// add authentication header
	header := http.Header{}
	header.Set("Authorization", fmt.Sprintf("Basic %s", service.token))
	(*requestConfig).NonDefaultHeaders = &header

	// add error model
	errorResponse := ErrorResponse{}
	(*requestConfig).ErrorModel = &errorResponse

	request, response, e := service.httpService.HTTPRequest(httpMethod, requestConfig)
	if errorResponse.Message != "" {
		e.SetMessage(errorResponse.Message)
	}

	if response != nil {
		// Read RateLimit headers
		rateLimitLimit, err := strconv.ParseInt(response.Header.Get("X-RateLimit-Limit"), 10, 64)
		if err == nil {
			service.rateLimit.Limit = &rateLimitLimit
		} else {
			service.rateLimit.Limit = nil
		}
		rateLimitRemaining, err := strconv.ParseInt(response.Header.Get("X-RateLimit-Remaining"), 10, 64)
		if err == nil {
			service.rateLimit.Remaining = &rateLimitRemaining
		} else {
			service.rateLimit.Remaining = nil
		}
		retryAfter, err := strconv.ParseInt(response.Header.Get("Retry-After"), 10, 64)
		if err == nil {
			retryAt := time.Now().Add(time.Duration(retryAfter) * time.Second)
			service.rateLimit.RetryAt = &retryAt
		} else {
			service.rateLimit.RetryAt = nil
		}
	}

	return request, response, e
}

func (service *Service) url(path string) string {
	return fmt.Sprintf("%s/%s", fmt.Sprintf(apiURL, service.pod), path)
}

func (service *Service) get(requestConfig *go_http.RequestConfig) (*http.Request, *http.Response, *errortools.Error) {
	return service.httpRequest(http.MethodGet, requestConfig)
}

func (service *Service) post(requestConfig *go_http.RequestConfig) (*http.Request, *http.Response, *errortools.Error) {
	return service.httpRequest(http.MethodPost, requestConfig)
}

func (service *Service) put(requestConfig *go_http.RequestConfig) (*http.Request, *http.Response, *errortools.Error) {
	return service.httpRequest(http.MethodPut, requestConfig)
}

func (service *Service) delete(requestConfig *go_http.RequestConfig) (*http.Request, *http.Response, *errortools.Error) {
	return service.httpRequest(http.MethodDelete, requestConfig)
}

func (service *Service) RateLimit() RateLimit {
	return service.rateLimit
}

func (service *Service) APIName() string {
	return apiName
}

func (service *Service) APICallCount() int64 {
	return service.httpService.RequestCount()
}
