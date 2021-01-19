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
	APIURL                string = "https://api.%s.insightly.com/v3.1"
	DateFormat            string = "2006-01-02T15:04:05Z"
	DateFormatCustomField string = "2006-01-02 15:04:05"
)

// type
//
type Service struct {
	pod                string
	token              string
	httpService        *go_http.Service
	client             http.Client
	rateLimitRemaining *int64
	retryAt            *time.Time
}

type ServiceConfig struct {
	Pod                   string
	APIKey                string
	MaxRetries            *uint
	SecondsBetweenRetries *uint32
}

func NewService(config ServiceConfig) (*Service, *errortools.Error) {
	if config.Pod == "" {
		return nil, errortools.ErrorMessage("Service Pod not provided")
	}

	if config.APIKey == "" {
		return nil, errortools.ErrorMessage("Service API Key not provided")
	}

	httpServiceConfig := go_http.ServiceConfig{
		MaxRetries:            config.MaxRetries,
		SecondsBetweenRetries: config.SecondsBetweenRetries,
	}

	return &Service{
		pod:         config.Pod,
		token:       base64.URLEncoding.EncodeToString([]byte(config.APIKey)),
		httpService: go_http.NewService(httpServiceConfig),
		client:      http.Client{},
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

	request, response, e := service.httpService.HTTPRequest(httpMethod, requestConfig)

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
