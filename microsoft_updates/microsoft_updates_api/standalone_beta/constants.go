package standalone_beta

// Beta channel application IDs.
// The beta (Insider Fast) channel hosts the same application IDs as production.
// See standalone package for the full list and bundle ID mappings.
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

// AllAppIDs is the ordered list of all beta channel application IDs.
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

// AppNames maps application ID to human-readable display name.
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
