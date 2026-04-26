package usergroups

// Field constants for fields[userGroups] query parameter.
const (
	FieldOuId             = "ouId"
	FieldName             = "name"
	FieldType             = "type"
	FieldTotalMemberCount = "totalMemberCount"
	FieldCreatedDateTime  = "createdDateTime"
	FieldUpdatedDateTime  = "updatedDateTime"
	FieldStatus           = "status"
	FieldUsers            = "users"
)

// UserGroupStatus constants for status field values.
const (
	UserGroupStatusActive   = "ACTIVE"
	UserGroupStatusInactive = "INACTIVE"
)

// UserGroupType constants for type field values.
const (
	UserGroupTypeStandard = "STANDARD"
)
