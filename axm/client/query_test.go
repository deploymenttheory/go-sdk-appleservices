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
		{
			name:  "Non-empty string",
			key:   "test",
			value: "value",
			want:  true,
		},
		{
			name:  "Empty string",
			key:   "test",
			value: "",
			want:  false,
		},
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
		{
			name:  "Positive integer",
			key:   "limit",
			value: 10,
			want:  true,
		},
		{
			name:  "Zero",
			key:   "limit",
			value: 0,
			want:  false,
		},
		{
			name:  "Negative integer",
			key:   "limit",
			value: -1,
			want:  false,
		},
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

func TestQueryBuilder_AddInt64(t *testing.T) {
	tests := []struct {
		name  string
		key   string
		value int64
		want  string
	}{
		{
			name:  "Positive int64",
			key:   "timestamp",
			value: 1234567890,
			want:  "1234567890",
		},
		{
			name:  "Zero",
			key:   "timestamp",
			value: 0,
			want:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			qb := NewQueryBuilder()
			qb.AddInt64(tt.key, tt.value)
			
			got := qb.Get(tt.key)
			if got != tt.want {
				t.Errorf("Get() = %v, want %v", got, tt.want)
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
		{
			name:  "True",
			key:   "active",
			value: true,
			want:  "true",
		},
		{
			name:  "False",
			key:   "active",
			value: false,
			want:  "false",
		},
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
		{
			name:  "Non-zero time",
			key:   "created",
			value: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			want:  true,
		},
		{
			name:  "Zero time",
			key:   "created",
			value: time.Time{},
			want:  false,
		},
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
		{
			name:   "Multiple values",
			key:    "fields",
			values: []string{"name", "id", "status"},
			want:   "name,id,status",
		},
		{
			name:   "Single value",
			key:    "fields",
			values: []string{"name"},
			want:   "name",
		},
		{
			name:   "Empty slice",
			key:    "fields",
			values: []string{},
			want:   "",
		},
		{
			name:   "Values with empty strings",
			key:    "fields",
			values: []string{"name", "", "status"},
			want:   "name,status",
		},
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

func TestQueryBuilder_AddIntSlice(t *testing.T) {
	tests := []struct {
		name   string
		key    string
		values []int
		want   string
	}{
		{
			name:   "Multiple values",
			key:    "ids",
			values: []int{1, 2, 3},
			want:   "1,2,3",
		},
		{
			name:   "Single value",
			key:    "ids",
			values: []int{42},
			want:   "42",
		},
		{
			name:   "Empty slice",
			key:    "ids",
			values: []int{},
			want:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			qb := NewQueryBuilder()
			qb.AddIntSlice(tt.key, tt.values)
			
			got := qb.Get(tt.key)
			if got != tt.want {
				t.Errorf("Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQueryBuilder_AddCustom(t *testing.T) {
	qb := NewQueryBuilder()
	qb.AddCustom("custom", "value")
	
	if !qb.Has("custom") {
		t.Error("AddCustom did not add parameter")
	}
	
	got := qb.Get("custom")
	if got != "value" {
		t.Errorf("Get() = %v, want %v", got, "value")
	}
}

func TestQueryBuilder_AddIfNotEmpty(t *testing.T) {
	qb := NewQueryBuilder()
	
	qb.AddIfNotEmpty("key1", "value")
	qb.AddIfNotEmpty("key2", "")
	
	if !qb.Has("key1") {
		t.Error("AddIfNotEmpty did not add non-empty value")
	}
	
	if qb.Has("key2") {
		t.Error("AddIfNotEmpty added empty value")
	}
}

func TestQueryBuilder_AddIfTrue(t *testing.T) {
	qb := NewQueryBuilder()
	
	qb.AddIfTrue(true, "key1", "value1")
	qb.AddIfTrue(false, "key2", "value2")
	
	if !qb.Has("key1") {
		t.Error("AddIfTrue did not add parameter when condition is true")
	}
	
	if qb.Has("key2") {
		t.Error("AddIfTrue added parameter when condition is false")
	}
}

func TestQueryBuilder_Merge(t *testing.T) {
	qb := NewQueryBuilder()
	qb.AddString("existing", "value1")
	
	other := map[string]string{
		"new": "value2",
		"another": "value3",
	}
	
	qb.Merge(other)
	
	if !qb.Has("existing") {
		t.Error("Merge removed existing parameter")
	}
	
	if !qb.Has("new") || !qb.Has("another") {
		t.Error("Merge did not add new parameters")
	}
}

func TestQueryBuilder_Remove(t *testing.T) {
	qb := NewQueryBuilder()
	qb.AddString("key", "value")
	
	if !qb.Has("key") {
		t.Fatal("Parameter was not added")
	}
	
	qb.Remove("key")
	
	if qb.Has("key") {
		t.Error("Remove did not remove parameter")
	}
}

func TestQueryBuilder_Clear(t *testing.T) {
	qb := NewQueryBuilder()
	qb.AddString("key1", "value1")
	qb.AddString("key2", "value2")
	
	if qb.Count() != 2 {
		t.Fatalf("Expected 2 parameters, got %d", qb.Count())
	}
	
	qb.Clear()
	
	if qb.Count() != 0 {
		t.Errorf("Clear did not remove all parameters, count = %d", qb.Count())
	}
	
	if !qb.IsEmpty() {
		t.Error("IsEmpty returned false after Clear")
	}
}

func TestQueryBuilder_Count(t *testing.T) {
	qb := NewQueryBuilder()
	
	if qb.Count() != 0 {
		t.Errorf("Initial count = %d, want 0", qb.Count())
	}
	
	qb.AddString("key1", "value1")
	if qb.Count() != 1 {
		t.Errorf("Count after one add = %d, want 1", qb.Count())
	}
	
	qb.AddString("key2", "value2")
	if qb.Count() != 2 {
		t.Errorf("Count after two adds = %d, want 2", qb.Count())
	}
}

func TestQueryBuilder_IsEmpty(t *testing.T) {
	qb := NewQueryBuilder()
	
	if !qb.IsEmpty() {
		t.Error("IsEmpty returned false for new builder")
	}
	
	qb.AddString("key", "value")
	
	if qb.IsEmpty() {
		t.Error("IsEmpty returned true after adding parameter")
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
	
	// Verify it's a copy (modification shouldn't affect builder)
	result["key3"] = "value3"
	if qb.Has("key3") {
		t.Error("Modifying Build result affected builder")
	}
}

func TestQueryBuilder_BuildString(t *testing.T) {
	tests := []struct {
		name     string
		setup    func(*QueryBuilder)
		expected []string // Multiple possible orders due to map iteration
	}{
		{
			name: "Empty builder",
			setup: func(qb *QueryBuilder) {},
			expected: []string{""},
		},
		{
			name: "Single parameter",
			setup: func(qb *QueryBuilder) {
				qb.AddString("key", "value")
			},
			expected: []string{"key=value"},
		},
		{
			name: "Multiple parameters",
			setup: func(qb *QueryBuilder) {
				qb.AddString("key1", "value1")
				qb.AddString("key2", "value2")
			},
			expected: []string{
				"key1=value1&key2=value2",
				"key2=value2&key1=value1",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			qb := NewQueryBuilder()
			tt.setup(qb)
			
			got := qb.BuildString()
			
			found := false
			for _, exp := range tt.expected {
				if got == exp {
					found = true
					break
				}
			}
			
			if !found {
				t.Errorf("BuildString() = %v, expected one of %v", got, tt.expected)
			}
		})
	}
}

func TestQueryBuilder_FluentInterface(t *testing.T) {
	// Test that methods can be chained
	qb := NewQueryBuilder().
		AddString("key1", "value1").
		AddInt("key2", 42).
		AddBool("key3", true).
		AddStringSlice("key4", []string{"a", "b"})
	
	if qb.Count() != 4 {
		t.Errorf("Fluent interface resulted in %d parameters, want 4", qb.Count())
	}
}

func TestQueryBuilder_BuildString_NoParameters(t *testing.T) {
	qb := NewQueryBuilder()
	result := qb.BuildString()
	
	if result != "" {
		t.Errorf("BuildString() for empty builder = %v, want empty string", result)
	}
}

func TestQueryBuilder_BuildString_Format(t *testing.T) {
	qb := NewQueryBuilder()
	qb.AddString("name", "test")
	qb.AddInt("limit", 10)
	
	result := qb.BuildString()
	
	// Check that it contains both parameters with & separator
	if !strings.Contains(result, "=") {
		t.Error("BuildString does not contain '='")
	}
	
	// For multiple params, should contain &
	if qb.Count() > 1 && !strings.Contains(result, "&") {
		t.Error("BuildString with multiple params does not contain '&'")
	}
}
