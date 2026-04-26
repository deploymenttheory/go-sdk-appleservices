package configurations

// Field constants for fields[configurations] query parameter.
const (
	FieldType                   = "type"
	FieldName                   = "name"
	FieldConfiguredForPlatforms = "configuredForPlatforms"
	FieldCustomSettingsValues   = "customSettingsValues"
	FieldCreatedDateTime        = "createdDateTime"
	FieldUpdatedDateTime        = "updatedDateTime"
)

// ConfigurationType constants for type field values.
const (
	ConfigurationTypeCustomSetting = "CUSTOM_SETTING"
	ConfigurationTypeAirDrop       = "AIR_DROP"
)

// Platform constants for configuredForPlatforms field values.
const (
	PlatformIOS   = "PLATFORM_IOS"
	PlatformMacOS = "PLATFORM_MACOS"
	PlatformTvOS  = "PLATFORM_TVOS"
)
