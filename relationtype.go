package insightly

import (
	exactonline "exactonline"
	"strings"
)

// types
//

type RelationType struct {
	Name                            string
	Rank                            int
	ExactOnlineSubscriptionTypeCode string
	ExactOnlineSubscriptionType     *exactonline.SubscriptionType //the matched subscriptiontype from exact online
}

type RelationTypes struct {
	RelationTypes []RelationType
}

// methods
//

func (rts *RelationTypes) Append(name string, rank int, exactOnlineSubscriptionTypeCode string) {
	rts.RelationTypes = append(rts.RelationTypes, RelationType{name, rank, exactOnlineSubscriptionTypeCode, nil})
}

func (rts *RelationTypes) FindRelationType(relationTypeName string) *RelationType {
	for _, rt := range rts.RelationTypes {
		//fmt.Println(strings.ToLower(rt.Name), strings.ToLower(relationTypeName))
		if strings.ToLower(rt.Name) == strings.ToLower(relationTypeName) {
			return &rt
		}
	}

	return nil
}
