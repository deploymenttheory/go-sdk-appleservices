package constants

const (
	// ipsw.me API — third-party aggregator of Apple CDN firmware metadata.
	// All endpoint constants are full absolute URLs because this SDK spans
	// three distinct external hosts.

	// EndpointFirmwaresAllV3 returns all Mac device models and their complete
	// firmware history in a single condensed JSON response.
	// GET https://api.ipsw.me/v3/firmwares.json/condensed
	EndpointFirmwaresAllV3 = "https://api.ipsw.me/v3/firmwares.json/condensed"

	// EndpointFirmwaresByDeviceV4 is the base path for device-specific firmware
	// queries. Append "/{identifier}?type=ipsw" before use.
	// GET https://api.ipsw.me/v4/device/{identifier}?type=ipsw
	EndpointFirmwaresByDeviceV4 = "https://api.ipsw.me/v4/device"

	// Apple GDMF API — Apple's official Global Device Management Feed.

	// EndpointGDMFVersions returns all currently-signed firmware versions across
	// macOS, iOS, and visionOS, with posting/expiration dates and supported devices.
	// GET https://gdmf.apple.com/v2/pmv
	EndpointGDMFVersions = "https://gdmf.apple.com/v2/pmv"

	// Apple CDN — firmware file delivery.

	// AppleCDNBaseURL is the base URL for Apple's firmware CDN.
	// IPSW files live under: {AppleCDNBaseURL}/{catalogRelease}/fullrestores/{assetID}/{uuid}/{filename}.ipsw
	AppleCDNBaseURL = "https://updates.cdn-apple.com"
)
