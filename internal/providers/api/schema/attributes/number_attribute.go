// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package attributes

import "github.com/conformize/conformize/common/typed"

type NumberAttribute struct {
	Required        bool
	Description     string
	Deprecated      bool
	DeprecationHint string
	DefaultValue    any
	DefaultValueFn  func() any
}

func (numAttr *NumberAttribute) IsRequired() bool {
	return numAttr.Required
}

func (numAttr *NumberAttribute) GetDescription() string {
	return numAttr.Description
}

func (numAttr *NumberAttribute) IsDeprecated() bool {
	return numAttr.Deprecated
}

func (numAttr *NumberAttribute) GetDefaultValue() any {
	return numAttr.DefaultValue
}

func (numAttr *NumberAttribute) GetDefaultValueFn() func() any {
	return numAttr.DefaultValueFn
}

func (numAttr *NumberAttribute) GetDeprecationHint() string {
	return numAttr.DeprecationHint
}

func (numAttr *NumberAttribute) Type() typed.Typeable {
	return &typed.NumberTyped{}
}
