package appstore_macos

// iTunes Search API entity value for macOS software.
const iTunesEntityMacOS = "macSoftware"

// iTunes Search API country code.
const iTunesCountry = "us"

// macOS App Store search terms for Microsoft applications.
// These are the search strings passed to the iTunes Search API.
var AppSearchTerms = []string{
	"Microsoft Word",
	"Microsoft Excel",
	"Microsoft PowerPoint",
	"Microsoft Outlook",
	"Microsoft OneNote",
	"Microsoft OneDrive",
	"Microsoft Teams",
	"Microsoft Defender",
	"Microsoft Edge",
	"Microsoft Copilot",
	"Azure VPN Client",
	"Microsoft Intune Company Portal",
	"Microsoft Remote Desktop",
}

// macOS App Store bundle ID constants for Microsoft applications.
const (
	BundleIDWord           = "com.microsoft.Word"
	BundleIDExcel          = "com.microsoft.Excel"
	BundleIDPowerPoint     = "com.microsoft.Powerpoint"
	BundleIDOutlook        = "com.microsoft.Outlook"
	BundleIDOneNote        = "com.microsoft.onenote.mac"
	BundleIDOneDrive       = "com.microsoft.OneDrive"
	BundleIDTeams          = "com.microsoft.teams2"
	BundleIDDefender       = "com.microsoft.wdav"
	BundleIDEdge           = "com.microsoft.edgemac"
	BundleIDCopilot        = "com.microsoft.m365copilot"
	BundleIDCompanyPortal  = "com.microsoft.CompanyPortalMac"
	BundleIDRemoteDesktop  = "com.microsoft.rdc.macos"
	BundleIDAutoUpdate     = "com.microsoft.autoupdate2"
	BundleIDIntuneAgent    = "com.microsoft.intuneMDMAgent"
)
