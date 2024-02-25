// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package attributes

import "github.com/conformize/conformize/common/typed"

type ListAttribute struct {
	Required        bool
	Description     string
	Deprecated      bool
	DeprecationHint string
	DefaultValue    any
	DefaultValueFn  func() any
	ElementsType    typed.Typeable
}

func (lAttr *ListAttribute) Type() typed.Typeable {
	return &typed.ListTyped{ElementsType: lAttr.ElementsType}
}

func (lAttr *ListAttribute) IsRequired() bool {
	return lAttr.Required
}

func (lAttr *ListAttribute) GetDescription() string {
	return lAttr.Description
}

func (lAttr *ListAttribute) IsDeprecated() bool {
	return lAttr.Deprecated
}

func (lAttr *ListAttribute) GetDefaultValue() any {
	return lAttr.DefaultValue
}

func (lAttr *ListAttribute) GetDefaultValueFn() func() any {
	return lAttr.DefaultValueFn
}

func (lAttr *ListAttribute) GetDeprecationHint() string {
	return lAttr.DeprecationHint
}

func (lAttr *ListAttribute) ElementType() typed.Typeable {
	return lAttr.ElementsType
}
