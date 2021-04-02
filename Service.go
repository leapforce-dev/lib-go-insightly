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
	APIURL                    string = "https://api.%s.insightly.com/v3.1"
	DateTimeFormat            string = "2006-01-02T15:04:05Z"
	DateTimeFormatCustomField string = "2006-01-02 15:04:05"
	DateFormat                string = "2006-01-02"
	defaultMaxRowCount        uint64 = ^uint64(0)
	defaultTop                uint64 = 100
)

// type
//
type Service struct {
	pod                string
	token              string
	maxRowCount        uint64
	httpService        *go_http.Service
	rateLimitRemaining *int64
	retryAt            *time.Time
	nextSkips          map[string]uint64
}

type ServiceConfig struct {
	Pod         string
	APIKey      string
	MaxRowCount *uint64
}

func NewService(serviceConfig ServiceConfig) (*Service, *errortools.Error) {
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
	if service.rateLimitRemaining != nil {
		if *service.rateLimitRemaining == 0 {
			if service.retryAt == nil {
				return nil, nil, errortools.ErrorMessage("Rate limit exceeded but RetryAt unknown.")
			}

			duration := service.retryAt.Sub(time.Now())

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
		rateLimitRemaining, err := strconv.ParseInt(response.Header.Get("X-RateLimit-Remaining"), 10, 64)
		if err == nil {
			service.rateLimitRemaining = &rateLimitRemaining
		} else {
			service.rateLimitRemaining = nil
		}
		retryAfter, err := strconv.ParseInt(response.Header.Get("Retry-After"), 10, 64)
		if err == nil {
			retryAt := time.Now().Add(time.Duration(retryAfter) * time.Second)
			service.retryAt = &retryAt
		} else {
			service.retryAt = nil
		}
	}

	return request, response, e
}

func (service *Service) url(path string) string {
	return fmt.Sprintf("%s/%s", fmt.Sprintf(APIURL, service.pod), path)
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
