package insightly

import (
	"strings"

	exactonline "github.com/Leapforce-nl/go_exactonline"
)

// types
//

type Article struct {
	AantalMedewerkers string
	Name              string
	ExactOnlineItem   *exactonline.Item
}

type RelationType struct {
	Name                            string
	Rank                            int
	ExactOnlineSubscriptionTypeCode string
	Articles                        []Article
	ExactOnlineSubscriptionType     *exactonline.SubscriptionType //the matched subscriptiontype from exact online
}

type RelationTypes struct {
	RelationTypes []RelationType
}

// methods
//

func (rts *RelationTypes) Append(name string, rank int, exactOnlineSubscriptionTypeCode string, articles []Article) {
	rts.RelationTypes = append(rts.RelationTypes, RelationType{name, rank, exactOnlineSubscriptionTypeCode, articles, nil})
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
