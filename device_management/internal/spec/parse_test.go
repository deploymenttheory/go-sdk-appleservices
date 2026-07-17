package spec

import (
	"testing"
)

// fixture is a condensed spec exercising the schema features the parser
// and codegen depend on: payload metadata, presence, rangelist, range,
// nested subkeys, arrays, defaults and per-OS support.
const fixture = `title: Demo Command
description: A demo command.
payload:
  requesttype: DemoCommand
  supportedOS:
    iOS:
      introduced: '4.0'
      accessrights: AllowDemo
    macOS:
      introduced: '10.7'
    tvOS:
      introduced: n/a
  content: Demo content.
payloadkeys:
- key: Message
  type: <string>
  presence: optional
  content: The message.
- key: Level
  type: <integer>
  presence: required
  rangelist:
  - 1
  - 2
  - 3
  default: 1
  range:
    min: 1
    max: 3
  content: The level.
- key: Settings
  type: <dictionary>
  presence: optional
  subkeys:
  - key: URL
    type: <string>
    subtype: <url>
    presence: required
    content: Where to go.
  - key: Retries
    type: <integer>
    presence: optional
    default: 3
    content: How many times.
- key: Tags
  type: <array>
  presence: optional
  repetition:
    min: 1
    max: 10
  subkeys:
  - key: Tag
    type: <string>
    presence: required
    format: '[a-z]+'
responsekeys:
- key: Status
  type: <string>
  presence: required
`

func TestParse(t *testing.T) {
	s, err := Parse([]byte(fixture), `mdm\commands\demo.command.yaml`)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if s.Category != "mdm/commands" || s.Name != "demo.command" {
		t.Fatalf("category/name = %q %q", s.Category, s.Name)
	}
	if s.Title != "Demo Command" || s.Payload.RequestType != "DemoCommand" {
		t.Fatalf("title/requesttype = %q %q", s.Title, s.Payload.RequestType)
	}
	if s.TypeIdentifier() != "DemoCommand" {
		t.Fatalf("TypeIdentifier = %q", s.TypeIdentifier())
	}
	if got := s.Payload.SupportedOS["iOS"].Introduced; got != "4.0" {
		t.Errorf("iOS introduced = %q", got)
	}
	if got := s.Payload.SupportedOS["tvOS"].Introduced; got != "n/a" {
		t.Errorf("tvOS introduced = %q", got)
	}

	if len(s.PayloadKeys) != 4 {
		t.Fatalf("payload keys = %d", len(s.PayloadKeys))
	}
	msg, level, settings, tags := s.PayloadKeys[0], s.PayloadKeys[1], s.PayloadKeys[2], s.PayloadKeys[3]

	if msg.Required() || msg.Type != "<string>" {
		t.Errorf("Message = %+v", msg)
	}
	if !level.Required() || level.Type != "<integer>" {
		t.Errorf("Level = %+v", level)
	}
	if len(level.RangeList) != 3 {
		t.Fatalf("Level rangelist = %v", level.RangeList)
	}
	if v, ok := level.RangeList[0].(int64); !ok || v != 1 {
		t.Errorf("rangelist[0] = %T %v, want int64 1", level.RangeList[0], level.RangeList[0])
	}
	if d, ok := level.Default.(int64); !ok || d != 1 {
		t.Errorf("default = %T %v", level.Default, level.Default)
	}
	if level.Range == nil || *level.Range.Min != 1 || *level.Range.Max != 3 {
		t.Errorf("range = %+v", level.Range)
	}

	if len(settings.Subkeys) != 2 {
		t.Fatalf("Settings subkeys = %d", len(settings.Subkeys))
	}
	if settings.Subkeys[0].Subtype != "<url>" || !settings.Subkeys[0].Required() {
		t.Errorf("Settings.URL = %+v", settings.Subkeys[0])
	}

	if tags.Repetition == nil || tags.Repetition.Min != 1 || tags.Repetition.Max != 10 {
		t.Errorf("Tags repetition = %+v", tags.Repetition)
	}
	if len(tags.Subkeys) != 1 || tags.Subkeys[0].Format != "[a-z]+" {
		t.Errorf("Tags item = %+v", tags.Subkeys)
	}

	if len(s.ResponseKeys) != 1 || s.ResponseKeys[0].Key != "Status" {
		t.Errorf("response keys = %+v", s.ResponseKeys)
	}
}

// TestParseRecursiveAnchor mirrors upstream specs (safari.bookmarks) whose
// subkeys reference themselves via a YAML anchor: the cycle must be cut,
// not error, with everything before the cycle intact.
func TestParseRecursiveAnchor(t *testing.T) {
	const recursive = `title: Recursive Demo
payload:
  declarationtype: com.example.recursive
payloadkeys:
- key: Groups
  type: <array>
  presence: optional
  subkeys: &grp
  - key: Group
    type: <dictionary>
    presence: required
    subkeys:
    - key: Name
      type: <string>
      presence: required
    - key: Children
      type: <array>
      presence: optional
      subkeys: *grp
`
	s, err := Parse([]byte(recursive), "declarative/declarations/configurations/recursive.yaml")
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	groups := s.PayloadKeys[0]
	if len(groups.Subkeys) != 1 {
		t.Fatalf("Groups subkeys = %+v", groups.Subkeys)
	}
	group := groups.Subkeys[0]
	if len(group.Subkeys) != 2 || group.Subkeys[0].Key != "Name" || group.Subkeys[1].Key != "Children" {
		t.Fatalf("Group subkeys = %+v", group.Subkeys)
	}
	// The cycle is cut: Children keeps its type but has no subkeys.
	children := group.Subkeys[1]
	if children.Type != "<array>" || len(children.Subkeys) != 0 {
		t.Fatalf("Children = %+v", children)
	}
}

// TestParseSharedAnchor: non-recursive anchors (the common case upstream)
// must still expand normally everywhere they are referenced.
func TestParseSharedAnchor(t *testing.T) {
	const shared = `title: Shared Demo
payload:
  requesttype: SharedDemo
payloadkeys:
- key: A
  type: <dictionary>
  presence: required
  subkeys: &common
  - key: Inner
    type: <string>
    presence: required
- key: B
  type: <dictionary>
  presence: optional
  subkeys: *common
`
	s, err := Parse([]byte(shared), "mdm/commands/shared.yaml")
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if len(s.PayloadKeys) != 2 {
		t.Fatalf("keys = %d", len(s.PayloadKeys))
	}
	for _, k := range s.PayloadKeys {
		if len(k.Subkeys) != 1 || k.Subkeys[0].Key != "Inner" {
			t.Fatalf("%s subkeys = %+v", k.Key, k.Subkeys)
		}
	}
}

func TestParseRejectsGarbage(t *testing.T) {
	if _, err := Parse([]byte(":\tnot yaml"), "x.yaml"); err == nil {
		t.Fatal("expected YAML error")
	}
	if _, err := Parse([]byte("description: no title"), "x.yaml"); err == nil {
		t.Fatal("expected missing-title error")
	}
}
