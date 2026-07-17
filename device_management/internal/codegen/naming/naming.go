// Package naming holds the spec-to-Go identifier rules shared by every
// emitter. All decisions about how a spec, key or enum value becomes a Go
// name live here.
package naming

import (
	"strings"
	"unicode"
)

// ExportName converts an arbitrary spec name into an exported Go
// identifier. Word boundaries are non-alphanumeric runs; existing
// capitalization inside a word is preserved ("Hash-SHA-256" -> "HashSHA256").
func ExportName(s string) string {
	var b strings.Builder
	newWord := true
	for _, r := range s {
		switch {
		case r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || r >= '0' && r <= '9':
			if newWord && r >= 'a' && r <= 'z' {
				r = unicode.ToUpper(r)
			}
			b.WriteRune(r)
			newWord = false
		default:
			newWord = true
		}
	}
	out := b.String()
	if out == "" {
		return "X"
	}
	if out[0] >= '0' && out[0] <= '9' {
		out = "N" + out
	}
	return out
}

// typePrefixes are stripped from wire type identifiers before deriving a
// struct name, most specific first.
var typePrefixes = []string{
	"com.apple.configuration.",
	"com.apple.activation.",
	"com.apple.asset.",
	"com.apple.management.",
	"com.apple.MCX.",
	"com.apple.",
}

// StructName derives a Go struct name from a spec's wire type identifier
// (requesttype, payloadtype or declarationtype), falling back to the spec
// file name. "DeviceLock" stays "DeviceLock";
// "com.apple.configuration.passcode.settings" becomes "PasscodeSettings".
func StructName(typeIdentifier, specName string) string {
	id := typeIdentifier
	if id == "" {
		id = specName
	}
	for _, p := range typePrefixes {
		if strings.HasPrefix(id, p) {
			id = strings.TrimPrefix(id, p)
			break
		}
	}
	return ExportName(id)
}

// FileName derives the generated file base name for a spec.
func FileName(specName string) string {
	s := strings.TrimPrefix(specName, "com.apple.")
	var b strings.Builder
	for _, r := range s {
		switch {
		case r >= 'a' && r <= 'z' || r >= '0' && r <= '9':
			b.WriteRune(r)
		case r >= 'A' && r <= 'Z':
			b.WriteRune(unicode.ToLower(r))
		}
	}
	if b.Len() == 0 {
		return "spec"
	}
	return b.String()
}

// ConstName derives the member part of an allowed-value constant.
func ConstName(value string) string {
	name := ExportName(value)
	if name == "X" {
		return "Empty"
	}
	return name
}

// FirstLine returns the first line of free text, trimmed.
func FirstLine(s string) string {
	s = strings.ReplaceAll(s, "\r", "")
	if i := strings.IndexByte(s, '\n'); i >= 0 {
		s = s[:i]
	}
	return strings.TrimSpace(s)
}

// WrapComment wraps s into comment-friendly lines of at most width runes.
func WrapComment(s string, width int) []string {
	words := strings.Fields(s)
	if len(words) == 0 {
		return nil
	}
	var lines []string
	line := words[0]
	for _, w := range words[1:] {
		if len(line)+1+len(w) > width {
			lines = append(lines, line)
			line = w
			continue
		}
		line += " " + w
	}
	return append(lines, line)
}
