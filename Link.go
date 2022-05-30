package insightly

type Link struct {
	LinkID         int64   `json:"LINK_ID"`
	ObjectName     string  `json:"OBJECT_NAME"`
	ObjectID       int64   `json:"OBJECT_ID"`
	LinkObjectName string  `json:"LINK_OBJECT_NAME"`
	LinkObjectID   int64   `json:"LINK_OBJECT_ID"`
	Role           *string `json:"ROLE,omitempty"`
	Details        *string `json:"DETAILS,omitempty"`
	RelationshipID *int64  `json:"RELATIONSHIP_ID,omitempty"`
	IsForward      *bool   `json:"IS_FORWARD,omitempty"`
}
