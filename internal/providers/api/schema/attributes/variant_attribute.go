// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package attributes

import "github.com/conformize/conformize/common/typed"

type VariantAttribute struct {
	Required        bool
	Description     string
	Deprecated      bool
	DeprecationHint string
	DefaultValue    any
	DefaultValueFn  func() any
	VariantsTypes   []typed.Typeable
}

func (varAttr *VariantAttribute) Type() typed.Typeable {
	return &typed.VariantTyped{VariantsTypes: varAttr.VariantsTypes}
}

func (varAttr *VariantAttribute) IsRequired() bool {
	return varAttr.Required
}

func (varAttr *VariantAttribute) GetDescription() string {
	return varAttr.Description
}

func (varAttr *VariantAttribute) IsDeprecated() bool {
	return varAttr.Deprecated
}

func (varAttr *VariantAttribute) GetDefaultValue() any {
	return varAttr.DefaultValue
}

func (varAttr *VariantAttribute) GetDefaultValueFn() func() any {
	return varAttr.DefaultValueFn
}

func (varAttr *VariantAttribute) GetDeprecationHint() string {
	return varAttr.DeprecationHint
}

func (varAttr *VariantAttribute) GetVariantsTypes() []typed.Typeable {
	return varAttr.VariantsTypes
}
