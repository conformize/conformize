package attributes

import "github.com/conformize/conformize/common/typed"

type TupleAttribute struct {
	Required        bool
	Description     string
	Deprecated      bool
	DeprecationHint string
	DefaultValue    interface{}
	DefaultValueFn  func() interface{}
	ElementsTypes   []typed.Typeable
}

func (tupAttr *TupleAttribute) Type() typed.Typeable {
	return &typed.TupleTyped{ElementsTypes: tupAttr.ElementsTypes}
}

func (tupAttr *TupleAttribute) IsRequired() bool {
	return tupAttr.Required
}

func (tupAttr *TupleAttribute) GetDescription() string {
	return tupAttr.Description
}

func (tupAttr *TupleAttribute) IsDeprecated() bool {
	return tupAttr.Deprecated
}

func (tupAttr *TupleAttribute) GetDefaultValue() interface{} {
	return tupAttr.DefaultValue
}

func (tupAttr *TupleAttribute) GetDefaultValueFn() func() interface{} {
	return tupAttr.DefaultValueFn
}

func (tupAttr *TupleAttribute) GetDeprecationHint() string {
	return tupAttr.DeprecationHint
}

func (tupAttr *TupleAttribute) GetElementsTypes() []typed.Typeable {
	return tupAttr.ElementsTypes
}
