# device_management — Implementation Plan

A generated Go SDK for Apple MDM and Declarative Device Management (DDM),
built from Apple's canonical schema repo
[apple/device-management](https://github.com/apple/device-management)
(YAML specs, one file per command / profile payload / declaration).

The lifecycle mirrors go-sdk-windowscsp: **collect → snapshot → generate →
ship**, with a semantic spec-diff pipeline to track Apple's moving target.

```
apple/device-management @ pinned commit
        │ fetchspec (download archive, verify, parse YAML)
        ▼
metadata/specs/**.json          committed, normalized snapshots + PROVENANCE
        │ gendm (offline, deterministic)
        ▼
mdm/commands, mdm/profiles,     generated structs + Validate() + enums
ddm/{configurations,assets,
     activations,management}
        │ mdm / ddm envelope builders ("the workflow engine")
        ▼
valid macOS config: plist (.mobileconfig, command plists) & DDM JSON
```

## Spec semantics honoured

Apple's schema (docs/schema.yaml upstream) drives everything:

| Spec construct | Generated Go |
|---|---|
| `presence: required` | value field (always serialized) |
| `presence: optional` | pointer / nil-able field with `omitempty` |
| `<string> <integer> <real> <boolean> <date> <data> <any>` | `string int64 float64 bool time.Time []byte any` |
| `<array>` (single item subkey) | `[]T` of the item type |
| `<dictionary>` + subkeys | named nested struct (`<Parent><Key>`) |
| `subkeytype` | canonical shared struct, emitted once per package |
| `rangelist` | typed constants + Validate() membership check |
| `range` / `format` / `repetition` | Validate() bounds / regex / cardinality checks |
| `subtype: <url> <hostname> <email>` | Validate() format checks |
| `payload.requesttype / payloadtype / declarationtype` | `RequestType()` / `PayloadType()` / `DeclarationType()` methods |
| supportedOS introduced/deprecated | doc comments |

## Phases

- **Phase 0** — scaffold `device_management/`, add `howett.net/plist`
  (struct-tag plist marshalling; stdlib has no plist encoder).
- **Phase 1** — `internal/spec`: YAML parser + normalized model.
  `cmd/fetchspec`: pinned-commit archive download (SHA-256 verified),
  snapshots to `metadata/specs/<category>/<name>.json` + `PROVENANCE.json`;
  `-ref` to re-pin, `-dir` for offline parse of a checkout.
- **Phase 2** — runtime: `mdm` (command envelope, profile/.mobileconfig
  builder), `ddm` (declaration JSON envelope), `validate` (shared checks),
  `ptr` (optional-field helpers).
- **Phase 3** — `internal/codegen`: build → view → render (embedded `.tmpl`)
  → fileasm, with the two-firewall design and DO-NOT-EDIT prune sentinel
  from the sibling generators. `cmd/gendm` drives it.
- **Phase 4** — generate all spec families, `go build ./device_management/...`,
  golden + round-trip tests (emit plist/JSON, decode back, compare).
- **Phase 5** — `cmd/specdiff` (semantic markdown diff of two snapshot
  trees: added/removed specs, key/type/presence/rangelist/OS changes), CI
  (unit tests, regen determinism gate, weekly spec-update PR whose body is
  the specdiff report), examples, README.

## Generated layout (one package per family, files separated by construct)

```
device_management/
  mdm/commands/       package commands   devicelock.go, devicelock_functions.go, …
  mdm/profiles/       package profiles   wifi.go, wifi_functions.go, wifi_enums.go, …
  ddm/configurations/ package configurations  passcodesettings.go, …
  ddm/assets/         package assets
  ddm/activations/    package activations
  ddm/management/     package management
```

Each spec emits up to three files, keeping construct kinds separate:
`<name>.go` (struct declarations: payload, response, nested subkey types),
`<name>_functions.go` (the wire-identifier method and `Validate() error`),
`<name>_enums.go` (rangelist constants, when present). The per-family
`registry.go` maps type identifiers → factories for decoding.

## The workflow engine

Typed values in → validated config out:

```go
cmd, _ := mdm.NewCommand(&commands.DeviceLock{PIN: ptr.To("123456")})      // → command plist
prof, _ := mdm.NewProfile("com.example.wifi", mdm.WithDisplayName("WiFi"),
    mdm.WithPayload(&profiles.WiFi{...}))                                   // → .mobileconfig
decl, _ := ddm.NewDeclaration("com.example.passcode",
    &configurations.PasscodeSettings{...})                                  // → DDM JSON
```

Every builder validates against the spec (required fields, key types,
rangelists, ranges, formats) before emitting; invalid input never becomes
config.

## Non-goals (v1)

- `declarative/status`, `declarative/protocol`, `mdm/checkin`, `mdm/errors`
  and `other/` specs (read-side / protocol plumbing; same generator can add
  them later).
- Transporting config to devices (that's an MDM server's job); this SDK
  produces and validates the artifacts.
