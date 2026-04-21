package firmware

const (
	// PlatformMac is the platform string for Mac computers in the ipsw.me API.
	PlatformMac = "REALBRIDGE"

	// FilenamePlatformUniversalMac is the platform prefix used in Universal Mac IPSW filenames.
	// Universal IPSW files run on all Apple Silicon Mac models.
	FilenamePlatformUniversalMac = "UniversalMac"

	// RestoreTypeRestore is the standard restore type for full IPSW restore files.
	RestoreTypeRestore = "Restore"

	// FirmwareTypeIPSW is the query parameter value for requesting IPSW firmware
	// from the ipsw.me v4 device endpoint.
	FirmwareTypeIPSW = "ipsw"

	// IdentifierPrefixIPhone is the model identifier prefix for iPhone devices.
	IdentifierPrefixIPhone = "iPhone"

	// IdentifierPrefixIPad is the model identifier prefix for iPad devices.
	IdentifierPrefixIPad = "iPad"

	// IdentifierPrefixIPod is the model identifier prefix for iPod touch devices.
	IdentifierPrefixIPod = "iPod"
)
