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
	maxRetries                int    = 10
)

type RateLimit struct {
	Limit     *int64
	Remaining *int64
	RetryAt   *time.Time
}

type Service struct {
	pod         string
	apiKey      string
	token       string
	maxRowCount uint64
	httpService *go_http.Service
	rateLimit   RateLimit
	nextSkips   map[string]uint64
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
		apiKey:      serviceConfig.APIKey,
		token:       base64.URLEncoding.EncodeToString([]byte(serviceConfig.APIKey)),
		maxRowCount: maxRowCount,
		httpService: httpService,
		nextSkips:   make(map[string]uint64),
	}, nil
}

func (service *Service) httpRequest(requestConfig *go_http.RequestConfig) (*http.Request, *http.Response, *errortools.Error) {
	retries := 0

retry:
	// check rate limit
	if service.rateLimit.Remaining != nil {
		if *service.rateLimit.Remaining <= 0 {
			if service.rateLimit.RetryAt == nil {
				return nil, nil, errortools.ErrorMessage("Rate limit exceeded but RetryAt unknown.")
			}

			duration := time.Until(*service.rateLimit.RetryAt)

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

	request, response, e := service.httpService.HTTPRequest(requestConfig)
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

		if response.StatusCode == http.StatusTooManyRequests {
			if retryAfter > 0 {
				if retries < maxRetries {
					retryAfter++
					fmt.Printf("waiting %v seconds...\n", retryAfter)
					// wait 2 seconds
					time.Sleep(time.Duration(retryAfter) * time.Second)
					retries++
					goto retry
				}
			}
		}
	}

	return request, response, e
}

func (service *Service) url(path string) string {
	return fmt.Sprintf("%s/%s", fmt.Sprintf(apiURL, service.pod), path)
}

func (service *Service) RateLimit() RateLimit {
	return service.rateLimit
}

func (service *Service) APIName() string {
	return apiName
}

func (service *Service) APIKey() string {
	return service.apiKey
}

func (service *Service) APICallCount() int64 {
	return service.httpService.RequestCount()
}

func (service *Service) APIReset() {
	service.httpService.ResetRequestCount()
}
