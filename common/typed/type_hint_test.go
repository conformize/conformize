package typed

import (
	"reflect"
	"testing"
)

func TestIsPrimitive(t *testing.T) {
	tests := []struct {
		hint     TypeHint
		expected bool
	}{
		{Boolean, true},
		{Number, true},
		{String, true},
		{List, false},
		{Map, false},
		{Object, false},
		{Tuple, false},
		{Variant, false},
		{Generic, false},
		{Invalid, false},
	}

	for _, tt := range tests {
		t.Run(string(NamedTypeFromTypeHint(tt.hint)), func(t *testing.T) {
			if match := IsPrimitive(tt.hint) != tt.expected; match {
				t.Errorf("expected %v, got %v", tt.expected, match)
			}
		})
	}
}

func TestTypeHintOf(t *testing.T) {
	tests := []struct {
		value    interface{}
		expected TypeHint
	}{
		{true, Boolean},
		{42, Number},
		{"hello", String},
		{[]int{1, 2, 3}, List},
		{map[string]int{"a": 1}, Map},
		{struct{}{}, Object},
		{[]interface{}{1, "two"}, Tuple},
	}

	for _, tt := range tests {
		t.Run(reflect.TypeOf(tt.value).String(), func(t *testing.T) {
			actual := TypeHintOf(reflect.ValueOf(tt.value))
			if actual != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, actual)
			}
		})
	}
}

func TestIsTuple(t *testing.T) {
	tests := []struct {
		name     string
		value    interface{}
		expected bool
	}{
		{
			"[]int is not tuple",
			[]int{1, 2, 3},
			false,
		},
		{
			"[]interface{} with mixed int and string elements is tuple",
			[]interface{}{1, "two"},
			true,
		},
		{
			"[]string is not tuple",
			[]interface{}{"one", "two"},
			false,
		},
		{
			"[]int with single element is not tuple",
			[]interface{}{42}, false,
		},
		{
			"[]interface{} with mixed int, string and bool is tuple",
			[]interface{}{1, "two", false},
			true,
		},
		{
			"[]interface{} with no elements is tuple",
			[]interface{}{},
			true,
		},
		{
			"string is not tuple",
			"not a slice",
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := isTuple(reflect.ValueOf(tt.value))
			if actual != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, actual)
			}
		})
	}
}

func TestNativeTypeForTypeHint(t *testing.T) {
	tests := []struct {
		typeHint TypeHint
		expected reflect.Type
	}{
		{Boolean, reflect.TypeOf(true)},
		{Number, reflect.TypeOf(float64(0))},
		{String, reflect.TypeOf("")},
		{List, reflect.TypeOf([]interface{}{})},
		{Tuple, reflect.TypeOf([]interface{}{})},
		{Map, reflect.TypeOf(map[interface{}]interface{}{})},
		{Object, reflect.TypeOf(struct{}{})},
		{Variant, reflect.TypeOf((*interface{})(nil)).Elem()},
		{Generic, reflect.TypeOf((*interface{})(nil)).Elem()},
		{Invalid, nil},
	}

	for _, tt := range tests {
		t.Run(NamedTypeFromTypeHint(tt.typeHint), func(t *testing.T) {
			actual, err := NativeTypeForTypeHint(tt.typeHint)
			if tt.typeHint == Invalid {
				if err == nil {
					t.Errorf("expected an error for invalid type hint, got nil")
				}
			} else {
				if actual != tt.expected {
					t.Errorf("expected %v, got %v", tt.expected, actual)
				}
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}
