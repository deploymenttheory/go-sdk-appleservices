// Package device_management is a generated Go SDK for Apple MDM and
// Declarative Device Management (DDM), built from Apple's canonical schema
// repo (github.com/apple/device-management).
//
// Typed values go in, spec-validated Apple config comes out:
//
//	cmd, _ := mdm.NewCommand(&commands.DeviceLock{PIN: ptr.To("123456")})
//	prof, _ := mdm.NewProfile("com.example.wifi",
//	    mdm.WithPayload(&profiles.WiFiManaged{...}))
//	decl, _ := ddm.BuildDeclaration("com.example.passcode",
//	    &configurations.PasscodeSettings{MinimumLength: ptr.To(int64(8))})
//
// The generated packages (mdm/commands, mdm/profiles, ddm/configurations,
// ddm/assets, ddm/activations, ddm/management) honour Apple's spec:
// required keys are value fields, optional keys are pointers, and every
// payload's Validate method enforces allowed values, ranges, formats and
// nested payload keys. Regeneration is driven by cmd/fetchspec (pinned
// upstream commit → metadata/specs snapshots) and cmd/gendm (offline
// codegen); cmd/specdiff renders semantic schema diffs between drops.
package device_management
