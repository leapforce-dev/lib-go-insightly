package insightly

import (
	"fmt"
	"strconv"
	"time"
)

// Product stores Product from Insightly
//
type Product struct {
	PRODUCT_ID        int           `json:"PRODUCT_ID"`
	PRODUCT_NAME      string        `json:"PRODUCT_NAME"`
	PRODUCT_CODE      string        `json:"PRODUCT_CODE"`
	PRODUCT_SKU       string        `json:"PRODUCT_SKU"`
	DESCRIPTION       string        `json:"DESCRIPTION"`
	PRODUCT_FAMILY    string        `json:"PRODUCT_FAMILY"`
	PRODUCT_IMAGE_URL string        `json:"PRODUCT_IMAGE_URL"`
	CURRENCY_CODE     string        `json:"CURRENCY_CODE"`
	DEFAULT_PRICE     int           `json:"DEFAULT_PRICE"`
	DATE_CREATED_UTC  string        `json:"DATE_CREATED_UTC"`
	DATE_UPDATED_UTC  string        `json:"DATE_UPDATED_UTC"`
	CREATED_USER_ID   int           `json:"CREATED_USER_ID"`
	OWNER_USER_ID     int           `json:"OWNER_USER_ID"`
	ACTIVE            bool          `json:"ACTIVE"`
	CUSTOMFIELDS      []CustomField `json:"CUSTOMFIELDS"`
	DateCreated       *time.Time
	DateUpdated       *time.Time
}

func (i *Insightly) GetProduct(id int) (*Product, error) {
	urlStr := "%sProduct/%v"
	url := fmt.Sprintf(urlStr, i.apiURL, id)
	//fmt.Println(url)

	o := Product{}

	err := i.Get(url, &o)
	if err != nil {
		return nil, err
	}

	o.ParseDates()

	return &o, nil
}

// GetProducts returns all products
//
func (i *Insightly) GetProducts() ([]Product, error) {
	return i.GetProductsInternal("")
}

// GetProductsUpdatedAfter returns all products updated after certain date
//
func (i *Insightly) GetProductsUpdatedAfter(updatedAfter time.Time) ([]Product, error) {
	from := updatedAfter.Format("2006-01-02")
	searchFilter := fmt.Sprintf("updated_after_utc=%s&", from)
	return i.GetProductsInternal(searchFilter)
}

// GetProductsFiltered returns all products fulfulling the specified filter
//
func (i *Insightly) GetProductsFiltered(fieldname string, fieldvalue string) ([]Product, error) {
	searchFilter := fmt.Sprintf("field_name=%s&field_value=%s&", fieldname, fieldvalue)
	return i.GetProductsInternal(searchFilter)
}

// GetProductsInternal is the generic function retrieving products from Insightly
//
func (i *Insightly) GetProductsInternal(searchFilter string) ([]Product, error) {
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
		url := fmt.Sprintf(urlStr, i.apiURL, searchString, strconv.Itoa(skip), strconv.Itoa(top))
		//fmt.Println(url)

		os := []Product{}

		err := i.Get(url, &os)
		if err != nil {
			return nil, err
		}

		for _, o := range os {
			o.ParseDates()
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

func (o *Product) ParseDates() {
	// parse DATE_CREATED_UTC to time.Time
	if o.DATE_CREATED_UTC != "" {
		t, _ := time.Parse("2006-01-02 15:04:05 +0000 UTC", o.DATE_CREATED_UTC+" +0000 UTC")
		//errortools.Fatal(err)
		o.DateCreated = &t
	}

	// parse DATE_UPDATED_UTC to time.Time
	if o.DATE_UPDATED_UTC != "" {
		t, _ := time.Parse("2006-01-02 15:04:05 +0000 UTC", o.DATE_UPDATED_UTC+" +0000 UTC")
		//errortools.Fatal(err)
		o.DateUpdated = &t
	}
}
