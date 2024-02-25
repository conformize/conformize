// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package attributes

import "github.com/conformize/conformize/common/typed"

type StringAttribute struct {
	Required        bool
	Description     string
	Deprecated      bool
	DefaultValue    any
	DefaultValueFn  func() any
	DeprecationHint string
}

func (sAttr *StringAttribute) IsRequired() bool {
	return sAttr.Required
}

func (sAttr *StringAttribute) GetDescription() string {
	return sAttr.Description
}

func (sAttr *StringAttribute) IsDeprecated() bool {
	return sAttr.Deprecated
}

func (sAttr *StringAttribute) GetDeprecationHint() string {
	return sAttr.DeprecationHint
}

func (sAttr *StringAttribute) GetDefaultValue() any {
	return sAttr.DefaultValue
}

func (sAttr *StringAttribute) GetDefaultValueFn() func() any {
	return sAttr.DefaultValueFn
}

func (sAttr *StringAttribute) Type() typed.Typeable {
	return &typed.StringTyped{}
}
