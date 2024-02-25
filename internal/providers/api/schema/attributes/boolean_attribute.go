// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package attributes

import "github.com/conformize/conformize/common/typed"

type BooleanAttribute struct {
	Required        bool
	Description     string
	Deprecated      bool
	DeprecationHint string
	DefaultValue    interface{}
	DefaultValueFn  func() interface{}
}

func (bAttr *BooleanAttribute) Type() typed.Typeable {
	return &typed.BooleanTyped{}
}

func (bAttr *BooleanAttribute) IsRequired() bool {
	return bAttr.Required
}

func (bAttr *BooleanAttribute) GetDescription() string {
	return bAttr.Description
}

func (bAttr *BooleanAttribute) IsDeprecated() bool {
	return bAttr.Deprecated
}

func (bAttr *BooleanAttribute) GetDefaultValue() interface{} {
	return bAttr.DefaultValue
}

func (bAttr *BooleanAttribute) GetDefaultValueFn() func() interface{} {
	return bAttr.DefaultValueFn
}

func (bAttr *BooleanAttribute) GetDeprecationHint() string {
	return bAttr.DeprecationHint
}
