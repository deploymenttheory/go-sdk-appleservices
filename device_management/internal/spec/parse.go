package spec

import (
	"fmt"
	"path"
	"strings"

	"gopkg.in/yaml.v3"
)

// yamlSpec is the raw YAML shape; json-facing normalization happens in
// Parse. Payload and Key already carry yaml tags, so they decode directly.
type yamlSpec struct {
	Title        string  `yaml:"title"`
	Description  string  `yaml:"description"`
	Payload      Payload `yaml:"payload"`
	PayloadKeys  []Key   `yaml:"payloadkeys"`
	ResponseKeys []Key   `yaml:"responsekeys"`
}

// Parse parses one upstream YAML spec file. relPath is the file's path
// inside the upstream repo (e.g. "mdm/commands/device.lock.yaml") and
// provides Category and Name.
//
// A few upstream specs (safari.bookmarks, homescreenlayout, …) use
// self-referential YAML anchors to describe recursively nested structures;
// yaml.v3 refuses to expand those into structs. Parse first clones the node
// tree with cycle-cutting (decycle), so a recursive reference becomes an
// open-ended dictionary/array at the point of recursion.
func Parse(data []byte, relPath string) (*Spec, error) {
	var doc yaml.Node
	if err := yaml.Unmarshal(data, &doc); err != nil {
		return nil, fmt.Errorf("parse %s: %w", relPath, err)
	}
	sanitized := decycle(&doc, map[*yaml.Node]bool{})
	if sanitized == nil {
		return nil, fmt.Errorf("parse %s: empty document", relPath)
	}
	var raw yamlSpec
	if err := sanitized.Decode(&raw); err != nil {
		return nil, fmt.Errorf("parse %s: %w", relPath, err)
	}
	if raw.Title == "" {
		return nil, fmt.Errorf("parse %s: missing title", relPath)
	}

	rel := path.Clean(strings.ReplaceAll(relPath, `\`, "/"))
	name := strings.TrimSuffix(path.Base(rel), path.Ext(rel))
	s := &Spec{
		Category:     path.Dir(rel),
		Name:         name,
		Title:        raw.Title,
		Description:  strings.TrimSpace(raw.Description),
		Payload:      raw.Payload,
		PayloadKeys:  normalizeKeys(raw.PayloadKeys),
		ResponseKeys: normalizeKeys(raw.ResponseKeys),
	}
	s.Payload.Content = strings.TrimSpace(s.Payload.Content)
	return s, nil
}

// decycle deep-copies a YAML node tree, expanding aliases but cutting
// self-referential cycles: an alias that points at a node currently being
// expanded returns nil, and nil children are dropped from the copy. The
// original tree is never mutated (alias targets are shared nodes).
func decycle(n *yaml.Node, expanding map[*yaml.Node]bool) *yaml.Node {
	if n == nil {
		return nil
	}
	if n.Kind == yaml.AliasNode {
		// Cycle detection happens in the target's own expansion below.
		return decycle(n.Alias, expanding)
	}
	if expanding[n] {
		return nil // self-referential anchor: cut the cycle here
	}
	expanding[n] = true
	defer delete(expanding, n)

	out := *n
	out.Anchor = ""
	out.Alias = nil
	out.Content = nil

	switch n.Kind {
	case yaml.MappingNode:
		// Content is key/value pairs; drop a pair when its value is cut.
		for i := 0; i+1 < len(n.Content); i += 2 {
			k := decycle(n.Content[i], expanding)
			v := decycle(n.Content[i+1], expanding)
			if k == nil || v == nil {
				continue
			}
			out.Content = append(out.Content, k, v)
		}
	default:
		for _, c := range n.Content {
			if cc := decycle(c, expanding); cc != nil {
				out.Content = append(out.Content, cc)
			}
		}
	}
	return &out
}

// normalizeKeys trims free text and normalizes YAML scalar types so the
// JSON snapshots are stable (ints stay ints, not float drift).
func normalizeKeys(keys []Key) []Key {
	for i := range keys {
		k := &keys[i]
		k.Content = strings.TrimSpace(k.Content)
		k.Title = strings.TrimSpace(k.Title)
		for j, v := range k.RangeList {
			k.RangeList[j] = normalizeScalar(v)
		}
		k.Default = normalizeScalar(k.Default)
		k.Subkeys = normalizeKeys(k.Subkeys)
	}
	return keys
}

// normalizeScalar keeps YAML scalars in JSON-stable form: integers as
// int64, everything else as-is.
func normalizeScalar(v any) any {
	switch n := v.(type) {
	case int:
		return int64(n)
	case int32:
		return int64(n)
	case uint64:
		return int64(n)
	case float64:
		if n == float64(int64(n)) {
			return int64(n)
		}
		return n
	}
	return v
}

// TypeIdentifier returns the spec's wire type identifier: requesttype for
// commands, payloadtype for profiles, declarationtype for declarations,
// statusitemtype for status items — whichever is set.
func (s *Spec) TypeIdentifier() string {
	p := s.Payload
	for _, t := range []string{p.RequestType, p.PayloadType, p.DeclarationType, p.StatusItemType} {
		if t != "" {
			return t
		}
	}
	return ""
}
