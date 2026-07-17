// Package view defines the pure-data models the templates render. Every
// field is a fully-resolved fragment: templates only branch and
// interpolate, they never make naming or type decisions.
package view

// File is one generated Go file (one upstream spec).
type File struct {
	Name         string // file base name without extension
	PackageName  string
	CommentLines []string // file-level context comment (spec title, OS support)
	// MainStruct is the payload struct carrying the wire-identifier method;
	// registries reference it.
	MainStruct    string
	Structs       []Struct
	Enums         []EnumBlock
	NeedsTime     bool // imports time
	NeedsValidate bool
}

// Struct is one generated payload/response/subkey struct.
type Struct struct {
	Name         string
	CommentLines []string
	Fields       []Field
	// TypeMethod/TypeValue emit the wire-identifier method, e.g.
	// RequestType() -> "DeviceLock". Empty TypeMethod means none.
	TypeMethod string
	TypeValue  string
	// ValidateLines is the fully-resolved body of Validate(); empty slice
	// still emits a Validate that returns nil (interface conformance).
	ValidateLines []string
}

// Field is one struct field.
type Field struct {
	Name         string
	Type         string
	Tag          string // full backquoted tag content (without backquotes)
	CommentLines []string
}

// EnumBlock is one allowed-values enum: a named type, its typed constants
// and a String method.
type EnumBlock struct {
	TypeName string
	BaseType string // "string" or "int64"
	Comment  string
	Members  []EnumMember
}

// EnumMember is one allowed-value constant.
type EnumMember struct {
	Name    string
	Literal string // rendered literal (already quoted for strings)
	// Dup marks members whose literal repeats an earlier member's value;
	// they are excluded from the String switch.
	Dup bool
}

// Registry is a family's identifier->factory registry file.
type Registry struct {
	PackageName  string
	CommentLines []string
	MapName      string // e.g. "ByRequestType"
	MapDoc       string
	IfaceType    string // e.g. "mdm.CommandPayload"
	Entries      []RegistryEntry
}

// RegistryEntry maps one wire identifier to its struct.
type RegistryEntry struct {
	Identifier string
	StructName string
}
