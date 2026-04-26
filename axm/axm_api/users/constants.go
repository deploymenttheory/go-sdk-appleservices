package users

// Field constants for fields[users] query parameter.
const (
	FieldFirstName           = "firstName"
	FieldLastName            = "lastName"
	FieldMiddleName          = "middleName"
	FieldStatus              = "status"
	FieldManagedAppleAccount = "managedAppleAccount"
	FieldIsExternalUser      = "isExternalUser"
	FieldRoleOuList          = "roleOuList"
	FieldEmail               = "email"
	FieldEmployeeNumber      = "employeeNumber"
	FieldCostCenter          = "costCenter"
	FieldDivision            = "division"
	FieldDepartment          = "department"
	FieldJobTitle            = "jobTitle"
	FieldStartDateTime       = "startDateTime"
	FieldCreatedDateTime     = "createdDateTime"
	FieldUpdatedDateTime     = "updatedDateTime"
	FieldPhoneNumbers        = "phoneNumbers"
)

// UserStatus constants for status field values.
const (
	UserStatusActive   = "ACTIVE"
	UserStatusInactive = "INACTIVE"
)
