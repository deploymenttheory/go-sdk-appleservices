package standalone_beta

import "encoding/xml"

// StandaloneBetaResponse holds all packages fetched from the beta CDN channel.
type StandaloneBetaResponse struct {
	Packages []*Package
}

// Package represents a single Microsoft application beta update entry.
type Package struct {
	ApplicationID     string
	Title             string
	ShortVersion      string
	FullVersion       string
	UpdateVersion     string
	MinimumOS         string
	Location          string
	AppOnlyLocation   string
	Hash              string
	HashSHA256        string
	AppOnlyHash       string
	AppOnlyHashSHA256 string
	Date              string
}

type plistArray struct {
	XMLName xml.Name    `xml:"plist"`
	Items   []plistDict `xml:"array>dict"`
}

type plistDict struct {
	Children []plistNode `xml:",any"`
}

type plistNode struct {
	XMLName xml.Name
	Value   string `xml:",chardata"`
}

func (d *plistDict) toPackage(appID string) *Package {
	p := &Package{ApplicationID: appID}
	children := d.Children
	for i := 0; i+1 < len(children); i += 2 {
		key := children[i].Value
		val := children[i+1].Value
		switch key {
		case "Title":
			p.Title = val
		case "Update Version":
			p.UpdateVersion = val
		case "Short Version":
			p.ShortVersion = val
		case "Full Version":
			p.FullVersion = val
		case "Minimum OS":
			p.MinimumOS = val
		case "Location":
			p.Location = val
		case "App Only Location":
			p.AppOnlyLocation = val
		case "Hash":
			p.Hash = val
		case "Hash SHA-256":
			p.HashSHA256 = val
		case "App Only Hash":
			p.AppOnlyHash = val
		case "App Only Hash SHA-256":
			p.AppOnlyHashSHA256 = val
		case "Date":
			p.Date = val
		}
	}
	if p.FullVersion == "" {
		p.FullVersion = p.UpdateVersion
	}
	return p
}
