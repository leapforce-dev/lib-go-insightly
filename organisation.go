package insightly

import (
	exactonline "exactonline"
	"strings"
	"time"
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
	KvKNummer string
	CountryId string
	//PartnerSinds                string
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

func (o *Organisation) Process(i *Insightly) bool {
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
	if o.DateUpdated.After(i.FromTimestamp) {
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
