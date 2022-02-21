package insightly

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
	i_types "github.com/leapforce-libraries/go_insightly/types"
)

// Product stores Product from Service
//
type Product struct {
	ProductID       int64                  `json:"PRODUCT_ID"`
	ProductName     string                 `json:"PRODUCT_NAME"`
	ProductCode     *string                `json:"PRODUCT_CODE"`
	ProductSKU      *string                `json:"PRODUCT_SKU"`
	Description     *string                `json:"DESCRIPTION"`
	ProductFamily   *string                `json:"PRODUCT_FAMILY"`
	ProductImageUrl *string                `json:"PRODUCT_IMAGE_Url"`
	CurrencyCode    string                 `json:"CURRENCY_CODE"`
	DefaultPrice    float64                `json:"DEFAULT_PRICE"`
	DateCreatedUTC  i_types.DateTimeString `json:"DATE_CREATED_UTC"`
	DateUpdatedUTC  i_types.DateTimeString `json:"DATE_UPDATED_UTC"`
	CreatedUserID   int64                  `json:"CREATED_USER_ID"`
	OwnerUserID     int64                  `json:"OWNER_USER_ID"`
	Active          bool                   `json:"ACTIVE"`
	CustomFields    *CustomFields          `json:"CUSTOMFIELDS"`
}

func (p *Product) prepareMarshal() interface{} {
	if p == nil {
		return nil
	}

	return &struct {
		ProductID       *int64        `json:"PRODUCT_ID"`
		ProductName     *string       `json:"PRODUCT_NAME"`
		ProductCode     *string       `json:"PRODUCT_CODE"`
		ProductSKU      *string       `json:"PRODUCT_SKU"`
		Description     *string       `json:"DESCRIPTION"`
		ProductFamily   *string       `json:"PRODUCT_FAMILY"`
		ProductImageUrl *string       `json:"PRODUCT_IMAGE_Url"`
		CurrencyCode    *string       `json:"CURRENCY_CODE"`
		DefaultPrice    *float64      `json:"DEFAULT_PRICE"`
		OwnerUserID     *int64        `json:"OWNER_USER_ID"`
		Active          *bool         `json:"ACTIVE"`
		CustomFields    *CustomFields `json:"CUSTOMFIELDS"`
	}{
		&p.ProductID,
		&p.ProductName,
		p.ProductCode,
		p.ProductSKU,
		p.Description,
		p.ProductFamily,
		p.ProductImageUrl,
		&p.CurrencyCode,
		&p.DefaultPrice,
		&p.OwnerUserID,
		&p.Active,
		p.CustomFields,
	}
}

// GetProduct returns a specific product
//
func (service *Service) GetProduct(productID int64) (*Product, *errortools.Error) {
	product := Product{}

	requestConfig := go_http.RequestConfig{
		Method:        http.MethodGet,
		Url:           service.url(fmt.Sprintf("Product/%v", productID)),
		ResponseModel: &product,
	}
	_, _, e := service.httpRequest(&requestConfig)
	if e != nil {
		return nil, e
	}

	return &product, nil
}

type GetProductsConfig struct {
	Skip         *uint64
	Top          *uint64
	Brief        *bool
	CountTotal   *bool
	UpdatedAfter *time.Time
	FieldFilter  *FieldFilter
}

// GetProducts returns all products
//
func (service *Service) GetProducts(config *GetProductsConfig) (*[]Product, *errortools.Error) {
	params := url.Values{}

	endpoint := "Product"
	products := []Product{}
	rowCount := uint64(0)
	top := defaultTop
	isSearch := false

	if config != nil {
		if config.Top != nil {
			top = *config.Top
		}
		if config.Skip != nil {
			service.nextSkips[endpoint] = *config.Skip
		}
		if config.Brief != nil {
			params.Set("brief", fmt.Sprintf("%v", *config.Brief))
		}
		if config.CountTotal != nil {
			params.Set("count_total", fmt.Sprintf("%v", *config.CountTotal))
		}
		if config.UpdatedAfter != nil {
			isSearch = true
			params.Set("updated_after_utc", fmt.Sprintf("%v", config.UpdatedAfter.Format(dateTimeFormat)))
		}
		if config.FieldFilter != nil {
			isSearch = true
			params.Set("field_name", config.FieldFilter.FieldName)
			params.Set("field_value", config.FieldFilter.FieldValue)
		}
	}

	if isSearch {
		endpoint += "/Search"
	}

	params.Set("top", fmt.Sprintf("%v", top))

	for {
		params.Set("skip", fmt.Sprintf("%v", service.nextSkips[endpoint]))
		productsBatch := []Product{}

		requestConfig := go_http.RequestConfig{
			Method:        http.MethodGet,
			Url:           service.url(fmt.Sprintf("%s?%s", endpoint, params.Encode())),
			ResponseModel: &productsBatch,
		}
		_, _, e := service.httpRequest(&requestConfig)
		if e != nil {
			return nil, e
		}

		products = append(products, productsBatch...)

		if len(productsBatch) < int(top) {
			delete(service.nextSkips, endpoint)
			break
		}

		service.nextSkips[endpoint] += top
		rowCount += top

		if rowCount >= service.maxRowCount {
			return &products, nil
		}
	}

	return &products, nil
}

// CreateProduct creates a new contract
//
func (service *Service) CreateProduct(product *Product) (*Product, *errortools.Error) {
	if product == nil {
		return nil, nil
	}

	productNew := Product{}

	requestConfig := go_http.RequestConfig{
		Method:        http.MethodPost,
		Url:           service.url("Products"),
		BodyModel:     product.prepareMarshal(),
		ResponseModel: &productNew,
	}
	_, _, e := service.httpRequest(&requestConfig)
	if e != nil {
		return nil, e
	}

	return &productNew, nil
}

// UpdateProduct updates an existing contract
//
func (service *Service) UpdateProduct(product *Product) (*Product, *errortools.Error) {
	if product == nil {
		return nil, nil
	}

	productUpdated := Product{}

	requestConfig := go_http.RequestConfig{
		Method:        http.MethodPut,
		Url:           service.url("Products"),
		BodyModel:     product.prepareMarshal(),
		ResponseModel: &productUpdated,
	}
	_, _, e := service.httpRequest(&requestConfig)
	if e != nil {
		return nil, e
	}

	return &productUpdated, nil
}

// DeleteProduct deletes a specific product
//
func (service *Service) DeleteProduct(productID int64) *errortools.Error {
	requestConfig := go_http.RequestConfig{
		Method: http.MethodDelete,
		Url:    service.url(fmt.Sprintf("Products/%v", productID)),
	}
	_, _, e := service.httpRequest(&requestConfig)
	if e != nil {
		return e
	}

	return nil
}
