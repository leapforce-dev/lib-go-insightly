package insightly

import (
	"fmt"
	"net/http"
	"net/url"

	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
)

// Relationship stores Relationship from Service
//
type Relationship struct {
	RelationshipID   int64  `json:"RELATIONSHIP_ID"`
	ForwardTitle     string `json:"FORWARD_TITLE"`
	Forward          string `json:"FORWARD"`
	ReverseTitle     string `json:"REVERSE_TITLE"`
	Reverse          string `json:"REVERSE"`
	ForContacts      bool   `json:"FOR_CONTACTS"`
	ForOrganisations bool   `json:"FOR_ORGANISATIONS"`
}

type GetRelationshipsConfig struct {
	Skip       *uint64
	Top        *uint64
	CountTotal *bool
}

// GetRelationships returns all relationships
//
func (service *Service) GetRelationships(config *GetRelationshipsConfig) (*[]Relationship, *errortools.Error) {
	params := url.Values{}

	endpoint := "Relationships"
	relationships := []Relationship{}
	rowCount := uint64(0)
	top := defaultTop

	if config != nil {
		if config.Top != nil {
			top = *config.Top
		}
		if config.Skip != nil {
			service.nextSkips[endpoint] = *config.Skip
		}
		if config.CountTotal != nil {
			params.Set("count_total", fmt.Sprintf("%v", *config.CountTotal))
		}
	}

	params.Set("top", fmt.Sprintf("%v", top))

	for {
		params.Set("skip", fmt.Sprintf("%v", service.nextSkips[endpoint]))

		relationshipsBatch := []Relationship{}

		requestConfig := go_http.RequestConfig{
			Method:        http.MethodGet,
			Url:           service.url(fmt.Sprintf("%s?%s", endpoint, params.Encode())),
			ResponseModel: &relationshipsBatch,
		}

		_, _, e := service.httpRequest(&requestConfig)
		if e != nil {
			return nil, e
		}

		relationships = append(relationships, relationshipsBatch...)

		if len(relationshipsBatch) < int(top) {
			delete(service.nextSkips, endpoint)
			break
		}

		service.nextSkips[endpoint] += top
		rowCount += top

		if rowCount >= service.maxRowCount {
			return &relationships, nil
		}
	}

	return &relationships, nil
}
