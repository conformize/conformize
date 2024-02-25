// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package attributes

import "github.com/conformize/conformize/common/typed"

type MapAttribute struct {
	Required        bool
	Description     string
	Deprecated      bool
	DeprecationHint string
	DefaultValue    any
	DefaultValueFn  func() any
	ElementsType    typed.Typeable
}

func (mAttr *MapAttribute) Type() typed.Typeable {
	return &typed.MapTyped{ElementsType: mAttr.ElementsType}
}

func (mAttr *MapAttribute) IsRequired() bool {
	return mAttr.Required
}

func (mAttr *MapAttribute) GetDescription() string {
	return mAttr.Description
}

func (mAttr *MapAttribute) IsDeprecated() bool {
	return mAttr.Deprecated
}

func (mAttr *MapAttribute) GetDefaultValue() any {
	return mAttr.DefaultValue
}

func (mAttr *MapAttribute) GetDefaultValueFn() func() any {
	return mAttr.DefaultValueFn
}

func (mAttr *MapAttribute) GetDeprecationHint() string {
	return mAttr.DeprecationHint
}

func (mAttr *MapAttribute) ElementType() typed.Typeable {
	return mAttr.ElementsType
}
