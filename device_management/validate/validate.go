// Package validate implements the value checks Apple's device-management
// schema declares on payload keys: closed value sets (rangelist), numeric
// bounds (range), regular-expression formats, array cardinality
// (repetition) and string subtypes (<url>, <hostname>, <email>).
//
// Generated Validate() methods call these helpers; they are exported so
// hand-written code can reuse the same semantics.
package validate

import (
	"fmt"
	"net/mail"
	"net/url"
	"regexp"
	"strings"
	"sync"
)

// InList checks membership of v in the spec's rangelist.
func InList[T comparable](field string, v T, allowed []T) error {
	for _, a := range allowed {
		if v == a {
			return nil
		}
	}
	return fmt.Errorf("%s: value %v is not in the allowed set %v", field, v, allowed)
}

// IntRange checks integer bounds from the spec's range block.
func IntRange(field string, v int64, min, max *int64) error {
	if min != nil && v < *min {
		return fmt.Errorf("%s: value %d is below the minimum %d", field, v, *min)
	}
	if max != nil && v > *max {
		return fmt.Errorf("%s: value %d is above the maximum %d", field, v, *max)
	}
	return nil
}

// FloatRange checks real-number bounds from the spec's range block.
func FloatRange(field string, v float64, min, max *float64) error {
	if min != nil && v < *min {
		return fmt.Errorf("%s: value %v is below the minimum %v", field, v, *min)
	}
	if max != nil && v > *max {
		return fmt.Errorf("%s: value %v is above the maximum %v", field, v, *max)
	}
	return nil
}

// Repetition checks array cardinality from the spec's repetition block.
func Repetition(field string, n, min, max int) error {
	if n < min {
		return fmt.Errorf("%s: %d items, need at least %d", field, n, min)
	}
	if max > 0 && n > max {
		return fmt.Errorf("%s: %d items, allows at most %d", field, n, max)
	}
	return nil
}

var (
	regexpMu    sync.Mutex
	regexpCache = map[string]*regexp.Regexp{}
)

// Format checks v against the spec's regular-expression format. Patterns
// are anchored to the full value, matching the spec's intent.
func Format(field, v, pattern string) error {
	regexpMu.Lock()
	re, ok := regexpCache[pattern]
	var err error
	if !ok {
		re, err = regexp.Compile("^(?:" + pattern + ")$")
		if err == nil {
			regexpCache[pattern] = re
		}
	}
	regexpMu.Unlock()
	if err != nil {
		return fmt.Errorf("%s: spec format %q is not a valid pattern: %w", field, pattern, err)
	}
	if !re.MatchString(v) {
		return fmt.Errorf("%s: value %q does not match the required format %q", field, v, pattern)
	}
	return nil
}

// URL checks the <url> subtype: an absolute URL with a scheme and host.
func URL(field, v string) error {
	u, err := url.Parse(v)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return fmt.Errorf("%s: value %q is not an absolute URL", field, v)
	}
	return nil
}

var hostnameRe = regexp.MustCompile(`^(?i)[a-z0-9]([a-z0-9-]{0,61}[a-z0-9])?(\.[a-z0-9]([a-z0-9-]{0,61}[a-z0-9])?)*$`)

// Hostname checks the <hostname> subtype (RFC 1123 labels).
func Hostname(field, v string) error {
	if len(v) == 0 || len(v) > 253 || !hostnameRe.MatchString(strings.TrimSuffix(v, ".")) {
		return fmt.Errorf("%s: value %q is not a valid hostname", field, v)
	}
	return nil
}

// Email checks the <email> subtype.
func Email(field, v string) error {
	a, err := mail.ParseAddress(v)
	if err != nil || a.Address != v {
		return fmt.Errorf("%s: value %q is not a valid email address", field, v)
	}
	return nil
}

// Required reports a missing required key.
func Required(field string) error {
	return fmt.Errorf("%s: required key is missing", field)
}

// Nested wraps a nested dictionary's validation error with its key.
func Nested(field string, err error) error {
	return fmt.Errorf("%s: %w", field, err)
}

// Indexed wraps an array item's validation error with its position.
func Indexed(field string, i int, err error) error {
	return fmt.Errorf("%s[%d]: %w", field, i, err)
}
