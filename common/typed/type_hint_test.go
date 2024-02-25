package typed

import (
	"reflect"
	"testing"
)

func TestIsPrimitive(t *testing.T) {
	tests := []struct {
		hint     TypeHinter
		expected bool
	}{
		{&simpleTypeHint{kind: Boolean}, true},
		{&simpleTypeHint{kind: Number}, true},
		{&simpleTypeHint{kind: String}, true},
		{&simpleTypeHint{kind: List}, false},
		{&complexTypeHint{kind: Map, elemType: &simpleTypeHint{kind: Number}}, false},
		{&simpleTypeHint{kind: Object}, false},
		{&complexTypeHint{kind: Tuple, elemType: &simpleTypeHint{kind: Generic}}, false},
		{&complexTypeMixedElementsHint{kind: Variant, elementsType: []TypeHinter{&simpleTypeHint{kind: Boolean}, &simpleTypeHint{kind: String}, &simpleTypeHint{kind: Number}}}, false},
		{&simpleTypeHint{kind: Generic}, false},
		{&simpleTypeHint{kind: Invalid}, false},
	}

	for _, tt := range tests {
		t.Run(string(NamedTypeFromTypeHint(tt.hint.TypeHint())), func(t *testing.T) {
			if match := IsPrimitive(tt.hint) != tt.expected; match {
				t.Errorf("expected %v, got %v", tt.expected, match)
			}
		})
	}
}

func TestTypeHintOf(t *testing.T) {
	tests := []struct {
		value    any
		expected TypeHinter
	}{
		{true, &simpleTypeHint{kind: Boolean}},
		{42, &simpleTypeHint{kind: Number}},
		{"hello", &simpleTypeHint{kind: String}},
		{[]int{1, 2, 3}, &complexTypeHint{kind: List, elemType: &simpleTypeHint{kind: Number}}},
		{map[string]int{"a": 1}, &complexTypeHint{kind: Map, elemType: &simpleTypeHint{kind: Number}}},
		{struct{}{}, &complexObjectHint{kind: Object, fieldsType: map[string]TypeHinter{}}},
		{[]any{1, "two"}, &complexTypeMixedElementsHint{kind: Tuple, elementsType: []TypeHinter{&simpleTypeHint{kind: Number}, &simpleTypeHint{kind: String}}}},
	}

	for _, tt := range tests {
		t.Run(reflect.TypeOf(tt.value).String(), func(t *testing.T) {
			actual := TypeHintOf(reflect.ValueOf(tt.value))
			if !reflect.DeepEqual(actual, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, actual)
			}
		})
	}
}

func TestNativeTypeForTypeHint(t *testing.T) {
	tests := []struct {
		typeHint TypeHinter
		expected reflect.Type
	}{
		{&simpleTypeHint{kind: Boolean}, reflect.TypeOf(true)},
		{&simpleTypeHint{kind: Number}, reflect.TypeOf(float64(0))},
		{&simpleTypeHint{kind: String}, reflect.TypeOf("")},
		{&complexTypeHint{kind: List, elemType: &simpleTypeHint{kind: Generic}}, reflect.TypeOf([]any{})},
		{&complexTypeHint{kind: Tuple, elemType: &simpleTypeHint{kind: Generic}}, reflect.TypeOf([]any{})},
		{&complexTypeHint{kind: Map, elemType: &simpleTypeHint{kind: Generic}}, reflect.TypeOf(map[string]any{})},
		{&simpleTypeHint{kind: Object}, reflect.TypeOf(struct{}{})},
		{&complexTypeHint{kind: Variant, elemType: &simpleTypeHint{kind: Generic}}, reflect.TypeOf((*any)(nil)).Elem()},
		{&simpleTypeHint{kind: Generic}, reflect.TypeOf((*any)(nil)).Elem()},
		{&simpleTypeHint{kind: Invalid}, nil},
	}

	for _, tt := range tests {
		t.Run(NamedTypeFromTypeHint(tt.typeHint.TypeHint()), func(t *testing.T) {
			actual, err := NativeTypeForTypeHint(tt.typeHint)
			if tt.typeHint.TypeHint() == Invalid {
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
