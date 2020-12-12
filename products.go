package insightly

import (
	"fmt"
	"strconv"
	"time"

	errortools "github.com/leapforce-libraries/go_errortools"
)

// Product stores Product from Insightly
//
type Product struct {
	ProductID       int           `json:"PRODUCT_ID"`
	ProductName     string        `json:"PRODUCT_NAME"`
	ProductCode     string        `json:"PRODUCT_CODE"`
	ProductSKU      string        `json:"PRODUCT_SKU"`
	Description     string        `json:"DESCRIPTION"`
	ProductFamily   string        `json:"PRODUCT_FAMILY"`
	ProductImageURL string        `json:"PRODUCT_IMAGE_URL"`
	CurrencyCode    string        `json:"CURRENCY_CODE"`
	DefaultPrice    int           `json:"DEFAULT_PRICE"`
	DateCreatedUTC  string        `json:"DATE_CREATED_UTC"`
	DateUpdatedUTC  string        `json:"DATE_UPDATED_UTC"`
	CreatedUserID   int           `json:"CREATED_USER_ID"`
	OwnerUserID     int           `json:"OWNER_USER_ID"`
	Active          bool          `json:"ACTIVE"`
	CustomFields    []CustomField `json:"CUSTOMFIELDS"`
	DateCreatedT    *time.Time
	DateUpdatedT    *time.Time
}

func (i *Insightly) GetProduct(id int) (*Product, *errortools.Error) {
	urlStr := "%sProduct/%v"
	url := fmt.Sprintf(urlStr, apiURL, id)
	//fmt.Println(url)

	o := Product{}

	_, _, e := i.get(url, nil, &o)
	if e != nil {
		return nil, e
	}

	o.parseDates()

	return &o, nil
}

// GetProducts returns all products
//
func (i *Insightly) GetProducts() ([]Product, *errortools.Error) {
	return i.GetProductsInternal("")
}

// GetProductsUpdatedAfter returns all products updated after certain date
//
func (i *Insightly) GetProductsUpdatedAfter(updatedAfter time.Time) ([]Product, *errortools.Error) {
	from := updatedAfter.Format("2006-01-02")
	searchFilter := fmt.Sprintf("updated_after_utc=%s&", from)
	return i.GetProductsInternal(searchFilter)
}

// GetProductsFiltered returns all products fulfulling the specified filter
//
func (i *Insightly) GetProductsFiltered(fieldname string, fieldvalue string) ([]Product, *errortools.Error) {
	searchFilter := fmt.Sprintf("field_name=%s&field_value=%s&", fieldname, fieldvalue)
	return i.GetProductsInternal(searchFilter)
}

// GetProductsInternal is the generic function retrieving products from Insightly
//
func (i *Insightly) GetProductsInternal(searchFilter string) ([]Product, *errortools.Error) {
	searchString := ""

	if searchFilter != "" {
		searchString = "/Search?" + searchFilter
	} else {
		searchString = "?"
	}

	urlStr := "%sProduct%sskip=%s&top=%s"
	skip := 0
	top := 500
	rowCount := top

	products := []Product{}

	for rowCount >= top {
		url := fmt.Sprintf(urlStr, apiURL, searchString, strconv.Itoa(skip), strconv.Itoa(top))
		//fmt.Println(url)

		os := []Product{}

		_, _, e := i.get(url, nil, &os)
		if e != nil {
			return nil, e
		}

		for _, o := range os {
			o.parseDates()
			products = append(products, o)
		}

		rowCount = len(os)
		skip += top
	}

	if len(products) == 0 {
		products = nil
	}

	return products, nil
}

func (o *Product) parseDates() {
	// parse DATE_CREATED_UTC to time.Time
	if o.DateCreatedUTC != "" {
		t, _ := time.Parse("2006-01-02 15:04:05 +0000 UTC", o.DateCreatedUTC+" +0000 UTC")
		//errortools.Fatal(err)
		o.DateCreatedT = &t
	}

	// parse DATE_UPDATED_UTC to time.Time
	if o.DateUpdatedUTC != "" {
		t, _ := time.Parse("2006-01-02 15:04:05 +0000 UTC", o.DateUpdatedUTC+" +0000 UTC")
		//errortools.Fatal(err)
		o.DateUpdatedT = &t
	}
}
