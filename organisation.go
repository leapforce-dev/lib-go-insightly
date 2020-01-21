package insightly

import (
	"exactonline"
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
	RelationTypeName         string
	KvKNummer                string
	CountryId                string
	PartnerSinds             string
	PushToEO                 bool
	ExactOnlineAccount       *exactonline.Account //the matched account from exact online
	MainContact              *Contact
	ExactOnlineMainContact   *exactonline.Contact //the matched main contact from exact online
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

func (o *Organisation) GetRelationTypeName(relationTypes RelationTypes) {
	relationTypeName := ""
	relationTypeRank := 1000

	for _, cf := range o.CUSTOMFIELDS {
		//fmt.Println("cf.FIELD_NAME:")
		//fmt.Println(cf.FIELD_NAME)

		if cf.FIELD_NAME == customFieldNameRelationType {
			//fmt.Println("original:")
			//fmt.Println(cf.FieldValueString)
			//fmt.Println(cf.getFieldValues())

			for _, fv := range cf.getFieldValues() {
				//fmt.Println("found:", fv, len(relationTypes.RelationTypes))
				rt := relationTypes.findRelationType(fv)
				if rt != nil {
					//fmt.Println("rt", rt)
					if rt.Rank < relationTypeRank {
						relationTypeName = rt.Name
						relationTypeRank = rt.Rank
					}
				} else {
					//fmt.Println("(", fv, ")")
				}
			}
		}
	}
	//fmt.Println("inside:")
	//fmt.Println(relationTypeName)

	o.RelationTypeName = relationTypeName
}
