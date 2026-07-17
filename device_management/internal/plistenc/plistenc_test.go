package plistenc

import (
	"strings"
	"testing"
	"time"
)

type demo struct {
	Name     string    `plist:"Name"`
	Count    int64     `plist:"Count"`
	Ratio    float64   `plist:"Ratio"`
	Enabled  bool      `plist:"Enabled"`
	Blob     []byte    `plist:"Blob,omitempty"`
	When     time.Time `plist:"When"`
	Note     *string   `plist:"Note,omitempty"`
	Zero     *int64    `plist:"Zero,omitempty"`
	Tags     []string  `plist:"Tags,omitempty"`
	Ignored  string    `plist:"-"`
	internal string    //nolint:unused // exercises unexported skipping
}

func TestDocumentStruct(t *testing.T) {
	note := "hi & <bye>"
	d := demo{
		Name:    "x",
		Count:   3,
		Ratio:   1.5,
		Enabled: true,
		Blob:    []byte{0x01, 0x02},
		When:    time.Date(2026, 7, 17, 12, 0, 0, 0, time.UTC),
		Note:    &note,
		Tags:    []string{"a", "b"},
		Ignored: "nope",
	}
	fields, err := Fields(&d)
	if err != nil {
		t.Fatal(err)
	}
	doc, err := Document(fields)
	if err != nil {
		t.Fatal(err)
	}
	out := string(doc)

	for _, want := range []string{
		`<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN"`,
		"<key>Name</key>",
		"<string>x</string>",
		"<integer>3</integer>",
		"<real>1.5</real>",
		"<true/>",
		"<data>AQI=</data>",
		"<date>2026-07-17T12:00:00Z</date>",
		"<string>hi &amp; &lt;bye&gt;</string>",
		"<array>",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("missing %q in:\n%s", want, out)
		}
	}
	// Omitted: nil pointer field, "-" tag, unexported.
	for _, banned := range []string{"Zero", "Ignored", "nope", "internal"} {
		if strings.Contains(out, banned) {
			t.Errorf("unexpected %q in output", banned)
		}
	}
	// Field order preserved: Name before Count before Ratio.
	if strings.Index(out, "Name") > strings.Index(out, "Count") || strings.Index(out, "Count") > strings.Index(out, "Ratio") {
		t.Error("field order not preserved")
	}
}

func TestZeroScalarsAreKept(t *testing.T) {
	// omitempty never drops zero scalars: optional scalars are pointers.
	type s struct {
		N int64  `plist:"N"`
		B bool   `plist:"B"`
		S string `plist:"S"`
	}
	fields, err := Fields(s{})
	if err != nil {
		t.Fatal(err)
	}
	doc, _ := Document(fields)
	out := string(doc)
	for _, want := range []string{"<integer>0</integer>", "<false/>", "<string></string>"} {
		if !strings.Contains(out, want) {
			t.Errorf("missing %q in:\n%s", want, out)
		}
	}
}

func TestNestedDictAndMap(t *testing.T) {
	type inner struct {
		A string `plist:"A"`
	}
	type outer struct {
		In  inner            `plist:"In"`
		Ptr *inner           `plist:"Ptr,omitempty"`
		M   map[string]int64 `plist:"M"`
	}
	fields, err := Fields(outer{In: inner{A: "x"}, M: map[string]int64{"b": 2, "a": 1}})
	if err != nil {
		t.Fatal(err)
	}
	doc, err := Document(fields)
	if err != nil {
		t.Fatal(err)
	}
	out := string(doc)
	// Map keys sorted.
	if strings.Index(out, "<key>a</key>") > strings.Index(out, "<key>b</key>") {
		t.Errorf("map keys not sorted:\n%s", out)
	}
	if !strings.Contains(out, "<key>A</key>") {
		t.Errorf("nested struct not encoded:\n%s", out)
	}
}

func TestDeterminism(t *testing.T) {
	fields, _ := Fields(demo{Name: "x", When: time.Unix(0, 0)})
	a, _ := Document(fields)
	b, _ := Document(fields)
	if string(a) != string(b) {
		t.Fatal("output not deterministic")
	}
}
