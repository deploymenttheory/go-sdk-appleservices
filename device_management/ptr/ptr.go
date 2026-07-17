// Package ptr provides helpers for the SDK's optional fields. The
// generated structs honour Apple's spec presence rules: required keys are
// value fields, optional keys are pointers — so callers set optional
// scalars with ptr.To.
package ptr

// To returns a pointer to v.
func To[T any](v T) *T { return &v }

// Value returns the value p points to, or the zero value when p is nil.
func Value[T any](p *T) T {
	if p == nil {
		var zero T
		return zero
	}
	return *p
}
