// Package plistenc is a deterministic XML property-list encoder for the
// device-management SDK. It exists (instead of a third-party plist
// library) for two reasons: envelope building needs to merge common keys
// (RequestType, PayloadType, …) with generated payload fields into a
// single dict with stable ordering, and generated output must be
// byte-deterministic for the regen CI gate.
//
// Encoding rules: struct fields emit in declaration order via their
// `plist:"Name,omitempty"` tags; maps emit with sorted keys; omitempty
// skips nil pointers, nil slices and nil maps (never zero scalars —
// optional scalars are pointers in generated structs).
package plistenc

import (
	"encoding/base64"
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"
)

// Entry is one dict key/value pair. Dicts are ordered slices of entries so
// output is deterministic.
type Entry struct {
	Key   string
	Value any
}

// Dict is an explicitly ordered dictionary value.
type Dict []Entry

const (
	header = `<?xml version="1.0" encoding="UTF-8"?>` + "\n" +
		`<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">` + "\n" +
		`<plist version="1.0">` + "\n"
	footer = `</plist>` + "\n"
)

// Document renders a complete plist document with the given root dict.
func Document(root Dict) ([]byte, error) {
	var b strings.Builder
	b.WriteString(header)
	if err := encodeDict(&b, root, 0); err != nil {
		return nil, err
	}
	b.WriteString(footer)
	return []byte(b.String()), nil
}

// Fields extracts a struct's plist-tagged fields, in declaration order,
// applying omitempty. v must be a struct or non-nil pointer to one.
func Fields(v any) (Dict, error) {
	rv := reflect.ValueOf(v)
	for rv.Kind() == reflect.Pointer {
		if rv.IsNil() {
			return nil, fmt.Errorf("plistenc: nil value")
		}
		rv = rv.Elem()
	}
	if rv.Kind() != reflect.Struct {
		return nil, fmt.Errorf("plistenc: %T is not a struct", v)
	}
	return structFields(rv)
}

func structFields(rv reflect.Value) (Dict, error) {
	var out Dict
	t := rv.Type()
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if !f.IsExported() {
			continue
		}
		name, omitEmpty, skip := parseTag(f)
		if skip {
			continue
		}
		fv := rv.Field(i)
		if omitEmpty && isEmpty(fv) {
			continue
		}
		// Deref optional pointers so scalar encoding sees the value.
		for fv.Kind() == reflect.Pointer {
			if fv.IsNil() {
				fv = reflect.Value{}
				break
			}
			fv = fv.Elem()
		}
		if !fv.IsValid() {
			continue
		}
		out = append(out, Entry{Key: name, Value: fv.Interface()})
	}
	return out, nil
}

func parseTag(f reflect.StructField) (name string, omitEmpty, skip bool) {
	tag := f.Tag.Get("plist")
	if tag == "-" {
		return "", false, true
	}
	name = f.Name
	if tag != "" {
		parts := strings.Split(tag, ",")
		if parts[0] != "" {
			name = parts[0]
		}
		for _, p := range parts[1:] {
			if p == "omitempty" {
				omitEmpty = true
			}
		}
	}
	return name, omitEmpty, false
}

// isEmpty implements omitempty: only absence counts (nil pointer, nil
// slice, nil map, nil interface), never zero scalars.
func isEmpty(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Pointer, reflect.Interface, reflect.Slice, reflect.Map:
		return v.IsNil()
	}
	return false
}

func encodeDict(b *strings.Builder, d Dict, depth int) error {
	ind := strings.Repeat("\t", depth)
	if len(d) == 0 {
		b.WriteString(ind + "<dict/>\n")
		return nil
	}
	b.WriteString(ind + "<dict>\n")
	for _, e := range d {
		b.WriteString(ind + "\t<key>" + escape(e.Key) + "</key>\n")
		if err := encodeValue(b, e.Value, depth+1); err != nil {
			return fmt.Errorf("key %q: %w", e.Key, err)
		}
	}
	b.WriteString(ind + "</dict>\n")
	return nil
}

func encodeValue(b *strings.Builder, v any, depth int) error {
	ind := strings.Repeat("\t", depth)
	switch tv := v.(type) {
	case nil:
		return fmt.Errorf("cannot encode nil value")
	case Dict:
		return encodeDict(b, tv, depth)
	case []byte:
		b.WriteString(ind + "<data>" + base64.StdEncoding.EncodeToString(tv) + "</data>\n")
		return nil
	case time.Time:
		b.WriteString(ind + "<date>" + tv.UTC().Format("2006-01-02T15:04:05Z") + "</date>\n")
		return nil
	case string:
		b.WriteString(ind + "<string>" + escape(tv) + "</string>\n")
		return nil
	case bool:
		if tv {
			b.WriteString(ind + "<true/>\n")
		} else {
			b.WriteString(ind + "<false/>\n")
		}
		return nil
	}

	rv := reflect.ValueOf(v)
	for rv.Kind() == reflect.Pointer || rv.Kind() == reflect.Interface {
		if rv.IsNil() {
			return fmt.Errorf("cannot encode nil value")
		}
		rv = rv.Elem()
	}
	switch rv.Kind() {
	case reflect.String:
		b.WriteString(ind + "<string>" + escape(rv.String()) + "</string>\n")
	case reflect.Bool:
		return encodeValue(b, rv.Bool(), depth)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		b.WriteString(ind + "<integer>" + strconv.FormatInt(rv.Int(), 10) + "</integer>\n")
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		b.WriteString(ind + "<integer>" + strconv.FormatUint(rv.Uint(), 10) + "</integer>\n")
	case reflect.Float32, reflect.Float64:
		b.WriteString(ind + "<real>" + strconv.FormatFloat(rv.Float(), 'g', -1, 64) + "</real>\n")
	case reflect.Slice, reflect.Array:
		if rv.Len() == 0 {
			b.WriteString(ind + "<array/>\n")
			return nil
		}
		b.WriteString(ind + "<array>\n")
		for i := 0; i < rv.Len(); i++ {
			if err := encodeValue(b, rv.Index(i).Interface(), depth+1); err != nil {
				return fmt.Errorf("index %d: %w", i, err)
			}
		}
		b.WriteString(ind + "</array>\n")
	case reflect.Map:
		if rv.Type().Key().Kind() != reflect.String {
			return fmt.Errorf("map keys must be strings, have %s", rv.Type().Key())
		}
		keys := make([]string, 0, rv.Len())
		for _, k := range rv.MapKeys() {
			keys = append(keys, k.String())
		}
		sort.Strings(keys)
		d := make(Dict, 0, len(keys))
		for _, k := range keys {
			d = append(d, Entry{Key: k, Value: rv.MapIndex(reflect.ValueOf(k)).Interface()})
		}
		return encodeDict(b, d, depth)
	case reflect.Struct:
		fields, err := structFields(rv)
		if err != nil {
			return err
		}
		return encodeDict(b, fields, depth)
	default:
		return fmt.Errorf("unsupported type %s", rv.Type())
	}
	return nil
}

func escape(s string) string {
	r := strings.NewReplacer("&", "&amp;", "<", "&lt;", ">", "&gt;")
	return r.Replace(s)
}
