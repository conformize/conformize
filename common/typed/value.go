package typed

import (
	"fmt"
)

func CreateValue(typ Typeable) (Valuable, error) {
	switch typ.Hint().TypeHint() {
	case Boolean:
		return NewBooleanValue(false)
	case Number:
		return NewNumberValue(0)
	case String:
		return NewStringValue("")
	case List:
		listType, ok := typ.(*ListTyped)
		if !ok {
			return nil, fmt.Errorf("expected ListTyped, got %T", typ)
		}
		return &ListValue{Elements: make([]Valuable, 0), ElementsType: listType.ElementsType}, nil
	case Map:
		mapType, ok := typ.(*MapTyped)
		if !ok {
			return nil, fmt.Errorf("expected MapTyped, got %T", typ)
		}
		return &MapValue{Elements: make(map[string]Valuable), ElementsType: mapType.ElementsType}, nil
	case Tuple:
		tupleType, ok := typ.(*TupleTyped)
		if !ok {
			return nil, fmt.Errorf("expected TupleTyped, got %T", typ)
		}
		return &TupleValue{Elements: make([]Valuable, 0), ElementsTypes: tupleType.ElementsTypes}, nil
	case Object:
		objectType, ok := typ.(*ObjectTyped)
		if !ok {
			return nil, fmt.Errorf("expected ObjectTyped, got %T", typ)
		}
		return &ObjectValue{Fields: make(map[string]Valuable), FieldsTypes: objectType.FieldsTypes}, nil
	case Variant:
		variantType, ok := typ.(*VariantTyped)
		if !ok {
			return nil, fmt.Errorf("expected VariantTyped, got %T", typ)
		}
		return &VariantValue{VariantsTypes: variantType.VariantsTypes}, nil
	case Generic:
		return &GenericValue{}, nil
	default:
		return nil, fmt.Errorf("unsupported type hint %s", typ.Name())
	}
}
