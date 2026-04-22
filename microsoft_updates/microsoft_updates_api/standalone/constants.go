package standalone

// Application IDs used in the Microsoft Office CDN URL path.
// CDN URL format: {ChannelBaseURL}{ApplicationID}.xml
const (
	AppIDWord           = "MSWD2019"
	AppIDExcel          = "XCEL2019"
	AppIDPowerPoint     = "PPT32019"
	AppIDOutlook        = "OPIM2019"
	AppIDOneNote        = "ONMC2019"
	AppIDTeams          = "TEAMS21"
	AppIDSkypeForBiz    = "MSFB16"
	AppIDDefenderEP     = "WDAV00"
	AppIDDefenderCons   = "WDAVC00"
	AppIDDefenderShim   = "WDAVS00"
	AppIDCompanyPortal  = "IMCP01"
	AppIDAutoUpdate     = "MSau04"
	AppIDWindowsApp     = "MSRD10"
	AppIDCopilot        = "MSCP10"
	AppIDQuickAssist    = "MSQA01"
	AppIDRemoteHelp     = "MSRH01"
	AppIDLicensing      = "OLIC02"
)

// Human-readable display names for each application ID.
var AppNames = map[string]string{
	AppIDWord:          "Microsoft Word",
	AppIDExcel:         "Microsoft Excel",
	AppIDPowerPoint:    "Microsoft PowerPoint",
	AppIDOutlook:       "Microsoft Outlook",
	AppIDOneNote:       "Microsoft OneNote",
	AppIDTeams:         "Microsoft Teams",
	AppIDSkypeForBiz:   "Skype for Business",
	AppIDDefenderEP:    "Microsoft Defender (Endpoint)",
	AppIDDefenderCons:  "Microsoft Defender (Consumer)",
	AppIDDefenderShim:  "Microsoft Defender (Shim)",
	AppIDCompanyPortal: "Intune Company Portal",
	AppIDAutoUpdate:    "Microsoft AutoUpdate",
	AppIDWindowsApp:    "Windows App",
	AppIDCopilot:       "Microsoft 365 Copilot",
	AppIDQuickAssist:   "Quick Assist",
	AppIDRemoteHelp:    "Remote Help",
	AppIDLicensing:     "Licensing Helper Tool",
}

// AllAppIDs is the ordered list of all production standalone application IDs.
var AllAppIDs = []string{
	AppIDWord,
	AppIDExcel,
	AppIDPowerPoint,
	AppIDOutlook,
	AppIDOneNote,
	AppIDTeams,
	AppIDSkypeForBiz,
	AppIDDefenderEP,
	AppIDDefenderCons,
	AppIDDefenderShim,
	AppIDCompanyPortal,
	AppIDAutoUpdate,
	AppIDWindowsApp,
	AppIDCopilot,
	AppIDQuickAssist,
	AppIDRemoteHelp,
	AppIDLicensing,
}

// Bundle ID constants for all standalone macOS Microsoft applications.
const (
	BundleIDWord           = "com.microsoft.word"
	BundleIDExcel          = "com.microsoft.excel"
	BundleIDPowerPoint     = "com.microsoft.powerpoint"
	BundleIDOutlook        = "com.microsoft.outlook"
	BundleIDOneNote        = "com.microsoft.onenote.mac"
	BundleIDTeams          = "com.microsoft.teams2"
	BundleIDSkypeForBiz    = "com.microsoft.skypeforbusiness"
	BundleIDDefenderEP     = "com.microsoft.wdav"
	BundleIDDefenderCons   = "com.microsoft.wdav.tray"
	BundleIDDefenderShim   = "com.microsoft.wdav.epsext"
	BundleIDCompanyPortal  = "com.microsoft.CompanyPortalMac"
	BundleIDAutoUpdate     = "com.microsoft.autoupdate2"
	BundleIDWindowsApp     = "com.microsoft.rdc.macos"
	BundleIDCopilot        = "com.microsoft.m365copilot"
	BundleIDQuickAssist    = "com.microsoft.quickassist"
	BundleIDRemoteHelp     = "com.microsoft.remotehelp"
	BundleIDLicensing      = "com.microsoft.office.licensingV2.helper"
)

// AppIDBundleMap maps application ID to its primary macOS bundle identifier.
var AppIDBundleMap = map[string]string{
	AppIDWord:          BundleIDWord,
	AppIDExcel:         BundleIDExcel,
	AppIDPowerPoint:    BundleIDPowerPoint,
	AppIDOutlook:       BundleIDOutlook,
	AppIDOneNote:       BundleIDOneNote,
	AppIDTeams:         BundleIDTeams,
	AppIDSkypeForBiz:   BundleIDSkypeForBiz,
	AppIDDefenderEP:    BundleIDDefenderEP,
	AppIDDefenderCons:  BundleIDDefenderCons,
	AppIDDefenderShim:  BundleIDDefenderShim,
	AppIDCompanyPortal: BundleIDCompanyPortal,
	AppIDAutoUpdate:    BundleIDAutoUpdate,
	AppIDWindowsApp:    BundleIDWindowsApp,
	AppIDCopilot:       BundleIDCopilot,
	AppIDQuickAssist:   BundleIDQuickAssist,
	AppIDRemoteHelp:    BundleIDRemoteHelp,
	AppIDLicensing:     BundleIDLicensing,
}
