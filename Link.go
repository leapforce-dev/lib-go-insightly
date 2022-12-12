package insightly

type Link struct {
	LinkID         *int64  `json:"LINK_ID,omitempty"`
	ObjectName     *string `json:"OBJECT_NAME,omitempty"`
	ObjectID       *int64  `json:"OBJECT_ID,omitempty"`
	LinkObjectName *string `json:"LINK_OBJECT_NAME,omitempty"`
	LinkObjectID   *int64  `json:"LINK_OBJECT_ID,omitempty"`
	Role           *string `json:"ROLE,omitempty"`
	Details        *string `json:"DETAILS,omitempty"`
	RelationshipID *int64  `json:"RELATIONSHIP_ID,omitempty"`
	IsForward      *bool   `json:"IS_FORWARD,omitempty"`
}
