package insightly

type RecordType string

const (
	RecordTypeContacts      RecordType = "contacts"
	RecordTypeLeads         RecordType = "leads"
	RecordTypeOpportunities RecordType = "opportunities"
	RecordTypeOrganisations RecordType = "organisations"
	RecordTypeProjects      RecordType = "projects"
	RecordTypeEmails        RecordType = "emails"
)
