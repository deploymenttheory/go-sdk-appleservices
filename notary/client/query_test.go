package client

import (
	"strings"
	"testing"
	"time"
)

func TestNewQueryBuilder(t *testing.T) {
	qb := NewQueryBuilder()
	if qb == nil {
		t.Fatal("NewQueryBuilder returned nil")
	}
	if qb.params == nil {
		t.Error("QueryBuilder params map is nil")
	}
}

func TestQueryBuilder_AddString(t *testing.T) {
	tests := []struct {
		name  string
		key   string
		value string
		want  bool
	}{
		{"Non-empty string", "test", "value", true},
		{"Empty string", "test", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			qb := NewQueryBuilder()
			qb.AddString(tt.key, tt.value)

			has := qb.Has(tt.key)
			if has != tt.want {
				t.Errorf("Has() = %v, want %v", has, tt.want)
			}

			if has {
				got := qb.Get(tt.key)
				if got != tt.value {
					t.Errorf("Get() = %v, want %v", got, tt.value)
				}
			}
		})
	}
}

func TestQueryBuilder_AddInt(t *testing.T) {
	tests := []struct {
		name  string
		key   string
		value int
		want  bool
	}{
		{"Positive integer", "limit", 10, true},
		{"Zero", "limit", 0, false},
		{"Negative integer", "limit", -1, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			qb := NewQueryBuilder()
			qb.AddInt(tt.key, tt.value)

			has := qb.Has(tt.key)
			if has != tt.want {
				t.Errorf("Has() = %v, want %v", has, tt.want)
			}
		})
	}
}

func TestQueryBuilder_AddBool(t *testing.T) {
	tests := []struct {
		name  string
		key   string
		value bool
		want  string
	}{
		{"True", "active", true, "true"},
		{"False", "active", false, "false"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			qb := NewQueryBuilder()
			qb.AddBool(tt.key, tt.value)

			got := qb.Get(tt.key)
			if got != tt.want {
				t.Errorf("Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQueryBuilder_AddTime(t *testing.T) {
	tests := []struct {
		name  string
		key   string
		value time.Time
		want  bool
	}{
		{"Non-zero time", "created", time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), true},
		{"Zero time", "created", time.Time{}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			qb := NewQueryBuilder()
			qb.AddTime(tt.key, tt.value)

			has := qb.Has(tt.key)
			if has != tt.want {
				t.Errorf("Has() = %v, want %v", has, tt.want)
			}

			if has {
				got := qb.Get(tt.key)
				expected := tt.value.Format(time.RFC3339)
				if got != expected {
					t.Errorf("Get() = %v, want %v", got, expected)
				}
			}
		})
	}
}

func TestQueryBuilder_AddStringSlice(t *testing.T) {
	tests := []struct {
		name   string
		key    string
		values []string
		want   string
	}{
		{"Multiple values", "fields", []string{"name", "id", "status"}, "name,id,status"},
		{"Single value", "fields", []string{"name"}, "name"},
		{"Empty slice", "fields", []string{}, ""},
		{"Values with empty strings", "fields", []string{"name", "", "status"}, "name,status"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			qb := NewQueryBuilder()
			qb.AddStringSlice(tt.key, tt.values)

			got := qb.Get(tt.key)
			if got != tt.want {
				t.Errorf("Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQueryBuilder_Build(t *testing.T) {
	qb := NewQueryBuilder()
	qb.AddString("key1", "value1")
	qb.AddInt("key2", 42)

	result := qb.Build()

	if len(result) != 2 {
		t.Errorf("Build returned %d parameters, want 2", len(result))
	}

	if result["key1"] != "value1" {
		t.Errorf("Build['key1'] = %v, want 'value1'", result["key1"])
	}

	if result["key2"] != "42" {
		t.Errorf("Build['key2'] = %v, want '42'", result["key2"])
	}

	// Verify it's a copy
	result["key3"] = "value3"
	if qb.Has("key3") {
		t.Error("Modifying Build result affected builder")
	}
}

func TestQueryBuilder_Clear(t *testing.T) {
	qb := NewQueryBuilder()
	qb.AddString("key1", "value1")
	qb.AddString("key2", "value2")

	qb.Clear()

	if qb.Count() != 0 {
		t.Errorf("Clear did not remove all parameters, count = %d", qb.Count())
	}

	if !qb.IsEmpty() {
		t.Error("IsEmpty returned false after Clear")
	}
}

func TestQueryBuilder_BuildString(t *testing.T) {
	qb := NewQueryBuilder()
	qb.AddString("name", "test")
	qb.AddInt("limit", 10)

	result := qb.BuildString()

	if !strings.Contains(result, "=") {
		t.Error("BuildString does not contain '='")
	}

	if qb.Count() > 1 && !strings.Contains(result, "&") {
		t.Error("BuildString with multiple params does not contain '&'")
	}
}

func TestQueryBuilder_FluentInterface(t *testing.T) {
	qb := NewQueryBuilder().
		AddString("key1", "value1").
		AddInt("key2", 42).
		AddBool("key3", true).
		AddStringSlice("key4", []string{"a", "b"})

	if qb.Count() != 4 {
		t.Errorf("Fluent interface resulted in %d parameters, want 4", qb.Count())
	}
}
