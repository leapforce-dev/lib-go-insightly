package insightly

import (
	"strings"
)

// types
//

type RelationType struct {
	Name string
	Rank int
}

type RelationTypes struct {
	RelationTypes []RelationType
}

// methods
//

func (rts *RelationTypes) Append(name string, rank int) {
	rts.RelationTypes = append(rts.RelationTypes, RelationType{name, rank})
}

func (rts *RelationTypes) findRelationType(relationTypeName string) *RelationType {
	for _, rt := range rts.RelationTypes {
		//fmt.Println(strings.ToLower(rt.Name), strings.ToLower(relationTypeName))
		if strings.ToLower(rt.Name) == strings.ToLower(relationTypeName) {
			return &rt
		}
	}

	return nil
}
