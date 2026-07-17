// Package spec models Apple's device-management schema files (the YAML
// specs in github.com/apple/device-management) and parses them into the
// normalized form committed as JSON snapshots under metadata/specs.
//
// The model mirrors upstream docs/schema.yaml: a spec has payload metadata
// (request/payload/declaration type, OS support) and a recursive list of
// payload keys carrying type, presence, allowed values and constraints.
package spec

// Provenance records where a snapshot tree came from.
type Provenance struct {
	Source  string `json:"source"`
	Ref     string `json:"ref"`
	Commit  string `json:"commit"`
	SHA256  string `json:"sha256"`
	Fetched string `json:"fetched"`
}

// Spec is one parsed schema file.
type Spec struct {
	// Category is the upstream directory, e.g. "mdm/commands",
	// "declarative/declarations/configurations".
	Category string `json:"category"`
	// Name is the upstream file base name, e.g. "device.lock".
	Name string `json:"name"`

	Title       string  `json:"title"`
	Description string  `json:"description,omitempty"`
	Payload     Payload `json:"payload"`
	PayloadKeys []Key   `json:"payloadkeys,omitempty"`
	// ResponseKeys describe MDM command responses.
	ResponseKeys []Key `json:"responsekeys,omitempty"`
}

// Payload is the spec's payload metadata block.
type Payload struct {
	RequestType     string `json:"requesttype,omitempty" yaml:"requesttype"`
	PayloadType     string `json:"payloadtype,omitempty" yaml:"payloadtype"`
	DeclarationType string `json:"declarationtype,omitempty" yaml:"declarationtype"`
	StatusItemType  string `json:"statusitemtype,omitempty" yaml:"statusitemtype"`
	CredentialType  string `json:"credentialtype,omitempty" yaml:"credentialtype"`
	// Apply is single, multiple or combined (DDM configurations).
	Apply   string `json:"apply,omitempty" yaml:"apply"`
	Beta    bool   `json:"beta,omitempty" yaml:"beta"`
	Content string `json:"content,omitempty" yaml:"content"`
	// SupportedOS is keyed by OS name (iOS, macOS, tvOS, visionOS, watchOS).
	SupportedOS map[string]OSSupport `json:"supportedOS,omitempty" yaml:"supportedOS"`
}

// OSSupport captures per-OS availability. "n/a" for Introduced means the
// payload is not supported on that OS.
type OSSupport struct {
	Introduced   string `json:"introduced,omitempty" yaml:"introduced"`
	Deprecated   string `json:"deprecated,omitempty" yaml:"deprecated"`
	Removed      string `json:"removed,omitempty" yaml:"removed"`
	AccessRights string `json:"accessrights,omitempty" yaml:"accessrights"`
	Beta         bool   `json:"beta,omitempty" yaml:"beta"`
}

// Key is one payload key. Subkeys make it recursive: dictionaries describe
// their members, arrays describe their single item shape.
type Key struct {
	Key   string `json:"key" yaml:"key"`
	Title string `json:"title,omitempty" yaml:"title"`
	// Type is the spec type in angle brackets: <string>, <integer>, <real>,
	// <boolean>, <date>, <data>, <array>, <dictionary>, <any>.
	Type string `json:"type" yaml:"type"`
	// Subtype refines string values: <url>, <hostname>, <email>.
	Subtype string `json:"subtype,omitempty" yaml:"subtype"`
	// Presence is "required" or "optional"; absent means optional.
	Presence string `json:"presence,omitempty" yaml:"presence"`
	// RangeList is the closed set of allowed values (strings or numbers).
	RangeList []any  `json:"rangelist,omitempty" yaml:"rangelist"`
	Range     *Range `json:"range,omitempty" yaml:"range"`
	Default   any    `json:"default,omitempty" yaml:"default"`
	// Format is a regular expression the value must match.
	Format string `json:"format,omitempty" yaml:"format"`
	// Repetition bounds array cardinality.
	Repetition  *Repetition `json:"repetition,omitempty" yaml:"repetition"`
	CombineType string      `json:"combinetype,omitempty" yaml:"combinetype"`
	Content     string      `json:"content,omitempty" yaml:"content"`
	AssetTypes  []string    `json:"assettypes,omitempty" yaml:"assettypes"`
	// SubkeyType names a shared structured type reused across specs.
	SubkeyType  string               `json:"subkeytype,omitempty" yaml:"subkeytype"`
	Subkeys     []Key                `json:"subkeys,omitempty" yaml:"subkeys"`
	SupportedOS map[string]OSSupport `json:"supportedOS,omitempty" yaml:"supportedOS"`
}

// Required reports whether the key must be present.
func (k Key) Required() bool { return k.Presence == "required" }

// Range bounds a numeric value.
type Range struct {
	Min *float64 `json:"min,omitempty" yaml:"min"`
	Max *float64 `json:"max,omitempty" yaml:"max"`
}

// Repetition bounds array cardinality.
type Repetition struct {
	Min int `json:"min" yaml:"min"`
	Max int `json:"max" yaml:"max"`
}
