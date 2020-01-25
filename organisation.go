package insightly

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	exactonline "github.com/leapforce-nl/go_exactonline"

	errortools "github.com/leapforce-nl/go_errortools"
)

// Organisation store Organisation from Insightly
//
type Organisation struct {
	ORGANISATION_ID          int           `json:"ORGANISATION_ID"`
	ORGANISATION_NAME        string        `json:"ORGANISATION_NAME"`
	ADDRESS_BILLING_STREET   string        `json:"ADDRESS_BILLING_STREET"`
	ADDRESS_BILLING_CITY     string        `json:"ADDRESS_BILLING_CITY"`
	ADDRESS_BILLING_STATE    string        `json:"ADDRESS_BILLING_STATE"`
	ADDRESS_BILLING_COUNTRY  string        `json:"ADDRESS_BILLING_COUNTRY"`
	ADDRESS_BILLING_POSTCODE string        `json:"ADDRESS_BILLING_POSTCODE"`
	DATE_UPDATED_UTC         string        `json:"DATE_UPDATED_UTC"`
	CUSTOMFIELDS             []CustomField `json:"CUSTOMFIELDS"`
	DateUpdated              time.Time
	//RelationTypeName         string
	//RelationType *RelationType
	KvKNummer                   string
	CountryId                   string
	PartnerSinds                string
	PartnerSindsTime            *time.Time
	BeeindigingPartnerschap     string
	BeeindigingPartnerschapTime *time.Time
	AantalMedewerkers           string
	PushToEO                    bool
	Opgezegd                    bool
	MainContact                 *Contact
	ExactOnlineAccount          *exactonline.Account          //the matched account from exact online
	ExactOnlineMainContact      *exactonline.Contact          //the matched main contact from exact online
	ExactOnlineSubscriptionType *exactonline.SubscriptionType //the matched SubscriptionType from exact online
	ExactOnlineItem             *exactonline.Item             //the matched Item from exact online
}

/*
type iOrganisations struct {
	Organisations []Organisation
}*/

// ToExactOnline return whether an organisation should be copied to ExactOnline or not
//
/*func (o *Organisation) ToExactOnline(onlyPushToEO bool, maxDateModified time.Time) bool {
	return o.RelationTypeName != "" && o.KvKNummer != "" && (o.PushToEO || !onlyPushToEO)
}*/

func (o *Organisation) GetExactOnlineSubscriptionTypeAndItem(relationTypes *RelationTypes) {
	var relationType *RelationType = nil
	relationTypeRank := 1000
	//value := ""

	for _, cf := range o.CUSTOMFIELDS {

		if cf.FIELD_NAME == customFieldNameRelationType {
			//value = cf.FieldValueString

			for _, fv := range cf.GetFieldValues() {
				rt := relationTypes.FindRelationType(fv)
				if rt != nil {
					if rt.Rank < relationTypeRank {
						relationTypeRank = rt.Rank
						relationType = rt
					}
				}
			}
		}
	}

	if relationType != nil {
		o.Opgezegd = strings.ToLower(relationType.Name) == "opgezegd"
		o.ExactOnlineSubscriptionType = relationType.ExactOnlineSubscriptionType

		if !o.Opgezegd {
			for ii, a := range relationType.Articles {
				if a.AantalMedewerkers == o.AantalMedewerkers {
					o.ExactOnlineItem = relationType.Articles[ii].ExactOnlineItem
				}
			}
		}
	}

	//fmt.Println("RelationType", value, "AantalMedewerkers", o.AantalMedewerkers, "\nExactOnlineSubscriptionType", o.ExactOnlineSubscriptionType, "ExactOnlineItem", o.ExactOnlineItem)
}

func (o *Organisation) Updated(i *Insightly) bool {
	return o.DateUpdated.After(i.FromTimestamp)
}

func (o *Organisation) Process(i *Insightly) bool {
	//fmt.Println(o.ExactOnlineSubscriptionType == nil, o.Opgezegd, o.KvKNummer, i.OnlyPushToEO, o.DateUpdated)
	if o.ExactOnlineSubscriptionType == nil && !o.Opgezegd {
		return false
	}
	if o.KvKNummer == "" {
		return false
	}
	if i.OnlyPushToEO {
		if o.PushToEO {
			//fmt.Println("ToExactOnline 1")
			return true
		} else {
			return false
		}
	}
	if o.Updated(i) {
		//fmt.Println("ToExactOnline 2", o.DateUpdated, i.FromTimestamp)
		return true
	}
	/*if o.MainContact != nil {
		if o.MainContact.DateUpdated.After(i.FromTimestamp) {
			//fmt.Println("ToExactOnline 3", o.MainContact.DateUpdated, i.FromTimestamp)
			return true
		}
	}*/

	return false
}

func (o *Organisation) ProcessSubscriptions() bool {
	return (o.ExactOnlineSubscriptionType != nil && o.ExactOnlineItem != nil) || o.Opgezegd
}

func (i *Insightly) GetOrganisation(id int) (*Organisation, error) {
	urlStr := "%sOrganisations/%v"
	url := fmt.Sprintf(urlStr, i.ApiUrl, id)
	fmt.Println(url)

	o := Organisation{}

	err := i.Get(url, &o)
	if err != nil {
		return nil, err
	}

	err = i.PrepareOrganisation(&o)
	if err != nil {
		return nil, err
	}

	//fmt.Println(o)
	//i.Organisations = append(i.Organisations, o)

	return &o, nil
}

func (i *Insightly) GetOrganisations() error {
	urlStr := "%sOrganisations/Search?updated_after_utc=%s&skip=%s&top=%s"
	skip := 0
	top := 500
	rowCount := 1

	from := i.FromTimestamp.Format("2006-01-02")

	for rowCount > 0 {
		url := fmt.Sprintf(urlStr, i.ApiUrl, from, strconv.Itoa(skip), strconv.Itoa(top))
		fmt.Println(url)

		os := []Organisation{}

		err := i.Get(url, &os)
		if err != nil {
			return err
		}

		for _, o := range os {
			err = i.PrepareOrganisation(&o)
			errortools.Fatal(err)

			//fmt.Println(o)
			i.Organisations = append(i.Organisations, o)
		}

		rowCount = len(os)
		skip += top
	}

	return nil
}

func (i *Insightly) PrepareOrganisation(o *Organisation) error {
	// unmarshal custom fields
	for ii := range o.CUSTOMFIELDS {
		o.CUSTOMFIELDS[ii].UnmarshalValue()
	}

	// parse DATE_UPDATED_UTC to time.Time
	t, err := time.Parse("2006-01-02 15:04:05 +0000 UTC", o.DATE_UPDATED_UTC+" +0000 UTC")
	errortools.Fatal(err)
	o.DateUpdated = t

	// get KvKNummer from custom field
	o.KvKNummer = i.FindCustomFieldValue(o.CUSTOMFIELDS, customFieldNameKvKNummer)

	// get Aantal Medewerkers from custom field
	o.AantalMedewerkers = i.FindCustomFieldValue(o.CUSTOMFIELDS, customFieldNameAantalMedewerkers)

	// get RelationTypeName from custom field
	o.GetExactOnlineSubscriptionTypeAndItem(&i.RelationTypes)

	// get PushToEO from custom field
	o.PushToEO = i.FindCustomFieldValueBool(o.CUSTOMFIELDS, customFieldNamePushToEO)

	// get PartnerSinds from custom field
	o.PartnerSinds = i.FindCustomFieldValue(o.CUSTOMFIELDS, customFieldNamePartnerSinds)
	if o.PartnerSinds != "" {
		t1, err := time.Parse("2006-01-02 15:04:05", o.PartnerSinds)
		errortools.Fatal(err)
		o.PartnerSindsTime = &t1

		//fmt.Println("o.PartnerSindsTime", t1)
	}

	// get BeeindigingPartnerschap from custom field
	o.BeeindigingPartnerschap = i.FindCustomFieldValue(o.CUSTOMFIELDS, customFieldNameBeeindigingPartnerschap)
	if o.BeeindigingPartnerschap != "" {
		t1, err := time.Parse("2006-01-02 15:04:05", o.BeeindigingPartnerschap)
		errortools.Fatal(err)
		o.BeeindigingPartnerschapTime = &t1

		//fmt.Println("o.BeeindigingPartnerschapTime", o.BeeindigingPartnerschap, t1)
	}

	// find CountryId
	id, err := i.Geo.FindCountryId(o.ADDRESS_BILLING_COUNTRY, "", "", "")
	if err != nil {
		return err
	}
	o.CountryId = id

	return nil
}
