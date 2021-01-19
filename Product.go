package insightly

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
)

// Product stores Product from Service
//
type Product struct {
	ProductID       int          `json:"PRODUCT_ID"`
	ProductName     string       `json:"PRODUCT_NAME"`
	ProductCode     string       `json:"PRODUCT_CODE"`
	ProductSKU      string       `json:"PRODUCT_SKU"`
	Description     string       `json:"DESCRIPTION"`
	ProductFamily   string       `json:"PRODUCT_FAMILY"`
	ProductImageURL string       `json:"PRODUCT_IMAGE_URL"`
	CurrencyCode    string       `json:"CURRENCY_CODE"`
	DefaultPrice    int          `json:"DEFAULT_PRICE"`
	DateCreatedUTC  DateUTC      `json:"DATE_CREATED_UTC"`
	DateUpdatedUTC  DateUTC      `json:"DATE_UPDATED_UTC"`
	CreatedUserID   *int         `json:"CREATED_USER_ID"`
	OwnerUserID     *int         `json:"OWNER_USER_ID"`
	Active          bool         `json:"ACTIVE"`
	CustomFields    CustomFields `json:"CUSTOMFIELDS"`
}

func (p *Product) prepareMarshal() interface{} {
	if p == nil {
		return nil
	}

	return &struct {
		ProductID       int           `json:"PRODUCT_ID"`
		ProductName     string        `json:"PRODUCT_NAME"`
		ProductCode     string        `json:"PRODUCT_CODE"`
		ProductSKU      string        `json:"PRODUCT_SKU"`
		Description     string        `json:"DESCRIPTION"`
		ProductFamily   string        `json:"PRODUCT_FAMILY"`
		ProductImageURL string        `json:"PRODUCT_IMAGE_URL"`
		CurrencyCode    string        `json:"CURRENCY_CODE"`
		DefaultPrice    int           `json:"DEFAULT_PRICE"`
		OwnerUserID     *int          `json:"OWNER_USER_ID"`
		Active          bool          `json:"ACTIVE"`
		CustomFields    []CustomField `json:"CUSTOMFIELDS"`
	}{
		p.ProductID,
		p.ProductName,
		p.ProductCode,
		p.ProductSKU,
		p.Description,
		p.ProductFamily,
		p.ProductImageURL,
		p.CurrencyCode,
		p.DefaultPrice,
		p.OwnerUserID,
		p.Active,
		p.CustomFields,
	}
}

// GetProduct returns a specific product
//
func (service *Service) GetProduct(productID int) (*Product, *errortools.Error) {
	product := Product{}

	requestConfig := go_http.RequestConfig{
		URL:           service.url(fmt.Sprintf("Products/%v", productID)),
		ResponseModel: &product,
	}
	_, _, e := service.get(&requestConfig)
	if e != nil {
		return nil, e
	}

	return &product, nil
}

type GetProductsFilter struct {
	UpdatedAfter *time.Time
	Field        *struct {
		FieldName  string
		FieldValue string
	}
}

// GetProducts returns all products
//
func (service *Service) GetProducts(filter *GetProductsFilter) (*[]Product, *errortools.Error) {
	searchString := "?"
	searchFilter := []string{}

	if filter != nil {
		if filter.UpdatedAfter != nil {
			from := filter.UpdatedAfter.Format(time.RFC3339)
			searchFilter = append(searchFilter, fmt.Sprintf("updated_after_utc=%s&", from))
		}

		if filter.Field != nil {
			searchFilter = append(searchFilter, fmt.Sprintf("field_name=%s&field_value=%s&", filter.Field.FieldName, filter.Field.FieldValue))
		}
	}

	if len(searchFilter) > 0 {
		searchString = "/Search?" + strings.Join(searchFilter, "&")
	}

	endpointStr := "Products%sskip=%s&top=%s"
	skip := 0
	top := 100
	rowCount := top

	products := []Product{}

	for rowCount >= top {
		_products := []Product{}

		requestConfig := go_http.RequestConfig{
			URL:           service.url(fmt.Sprintf(endpointStr, searchString, strconv.Itoa(skip), strconv.Itoa(top))),
			ResponseModel: &_products,
		}
		_, _, e := service.get(&requestConfig)
		if e != nil {
			return nil, e
		}

		products = append(products, _products...)

		rowCount = len(_products)
		//rowCount = 0
		skip += top
	}

	if len(products) == 0 {
		products = nil
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
		URL:           service.url("Products"),
		BodyModel:     product.prepareMarshal(),
		ResponseModel: &productNew,
	}
	_, _, e := service.post(&requestConfig)
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
		URL:           service.url("Products"),
		BodyModel:     product.prepareMarshal(),
		ResponseModel: &productUpdated,
	}
	_, _, e := service.put(&requestConfig)
	if e != nil {
		return nil, e
	}

	return &productUpdated, nil
}

// DeleteProduct deletes a specific product
//
func (service *Service) DeleteProduct(productID int) *errortools.Error {
	requestConfig := go_http.RequestConfig{
		URL: service.url(fmt.Sprintf("Products/%v", productID)),
	}
	_, _, e := service.delete(&requestConfig)
	if e != nil {
		return e
	}

	return nil
}
