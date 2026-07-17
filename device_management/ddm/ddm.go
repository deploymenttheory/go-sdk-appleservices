// Package ddm turns validated, generated declaration structs into
// syntactically correct Declarative Device Management JSON. It is the DDM
// half of the SDK's workflow engine — typed values in, spec-validated
// declarations out.
package ddm

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// DeclarationPayload is implemented by every generated declaration struct
// (configurations, assets, activations, management).
type DeclarationPayload interface {
	// DeclarationType is the declaration's wire identifier,
	// e.g. "com.apple.configuration.passcode.settings".
	DeclarationType() string
	// Validate checks the payload against Apple's spec.
	Validate() error
}

// Declaration is the DDM envelope: what a declarations endpoint serves.
type Declaration struct {
	Type        string `json:"Type"`
	Identifier  string `json:"Identifier"`
	ServerToken string `json:"ServerToken,omitempty"`
	Payload     any    `json:"Payload"`
}

// DeclarationOption customizes declaration construction.
type DeclarationOption func(*Declaration)

// WithServerToken sets the declaration's ServerToken, used by clients to
// detect changed declarations.
func WithServerToken(token string) DeclarationOption {
	return func(d *Declaration) { d.ServerToken = token }
}

// NewDeclaration validates payload and builds the declaration envelope.
func NewDeclaration(identifier string, payload DeclarationPayload, opts ...DeclarationOption) (*Declaration, error) {
	if identifier == "" {
		return nil, fmt.Errorf("ddm: declaration identifier is required")
	}
	if payload == nil {
		return nil, fmt.Errorf("ddm: nil declaration payload")
	}
	if err := payload.Validate(); err != nil {
		return nil, fmt.Errorf("ddm: invalid %s declaration: %w", payload.DeclarationType(), err)
	}
	d := &Declaration{
		Type:       payload.DeclarationType(),
		Identifier: identifier,
		Payload:    payload,
	}
	for _, opt := range opts {
		opt(d)
	}
	return d, nil
}

// JSON renders the declaration as indented DDM JSON.
func (d *Declaration) JSON() ([]byte, error) {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")
	if err := enc.Encode(d); err != nil {
		return nil, fmt.Errorf("ddm: encode %s: %w", d.Identifier, err)
	}
	return buf.Bytes(), nil
}

// BuildDeclaration is the one-call form: validate, wrap and render JSON.
func BuildDeclaration(identifier string, payload DeclarationPayload, opts ...DeclarationOption) ([]byte, error) {
	d, err := NewDeclaration(identifier, payload, opts...)
	if err != nil {
		return nil, err
	}
	return d.JSON()
}
