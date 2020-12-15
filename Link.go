package insightly

// types
//
type Link struct {
	LinkID         int    `json:"LINK_ID"`
	ObjectName     string `json:"OBJECT_NAME"`
	ObjectID       int    `json:"OBJECT_ID"`
	LinkObjectName string `json:"LINK_OBJECT_NAME"`
	LinkObjectID   int    `json:"LINK_OBJECT_ID"`
	Role           string `json:"ROLE"`
	Details        string `json:"DETAILS"`
	RelationshipID int    `json:"RELATIONSHIP_ID"`
	IsForward      bool   `json:"IS_FORWARD"`
}
