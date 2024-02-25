// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package attributes

import "github.com/conformize/conformize/common/typed"

type ObjectAttribute struct {
	Required        bool
	Description     string
	Deprecated      bool
	DeprecationHint string
	DefaultValue    interface{}
	DefaultValueFn  func() interface{}
	FieldsTypes     map[string]typed.Typeable
}

func (oAttr *ObjectAttribute) Type() typed.Typeable {
	return &typed.ObjectTyped{FieldsTypes: oAttr.FieldsTypes}
}

func (oAttr *ObjectAttribute) IsRequired() bool {
	return oAttr.Required
}

func (oAttr *ObjectAttribute) GetDescription() string {
	return oAttr.Description
}

func (oAttr *ObjectAttribute) IsDeprecated() bool {
	return oAttr.Deprecated
}

func (oAttr *ObjectAttribute) GetDefaultValue() interface{} {
	return oAttr.DefaultValue
}

func (oAttr *ObjectAttribute) GetDefaultValueFn() func() interface{} {
	return oAttr.DefaultValueFn
}

func (oAttr *ObjectAttribute) GetDeprecationHint() string {
	return oAttr.DeprecationHint
}

func (oAttr *ObjectAttribute) GetFieldsTypes() map[string]typed.Typeable {
	return oAttr.FieldsTypes
}
