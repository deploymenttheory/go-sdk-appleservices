package blueprints

// Field constants for fields[blueprints] query parameter.
const (
	FieldName                = "name"
	FieldDescription         = "description"
	FieldStatus              = "status"
	FieldCreatedDateTime     = "createdDateTime"
	FieldUpdatedDateTime     = "updatedDateTime"
	FieldAppLicenseDeficient = "appLicenseDeficient"
	FieldApps                = "apps"
	FieldPackages            = "packages"
	FieldConfigurations      = "configurations"
	FieldOrgDevices          = "orgDevices"
	FieldUsers               = "users"
	FieldUserGroups          = "userGroups"
)

// Include constants for the include query parameter on GetByBlueprintIDV1.
const (
	IncludeApps           = "apps"
	IncludePackages       = "packages"
	IncludeConfigurations = "configurations"
	IncludeOrgDevices     = "orgDevices"
	IncludeUsers          = "users"
	IncludeUserGroups     = "userGroups"
)

// BlueprintStatus constants for the status field.
const (
	BlueprintStatusActive   = "ACTIVE"
	BlueprintStatusInactive = "INACTIVE"
)
