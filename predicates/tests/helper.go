// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package tests

import (
	"reflect"

	"github.com/conformize/conformize/common/reflected"
	"github.com/conformize/conformize/common/typed"
)

func PrimVal(val any, vType typed.Typeable) typed.Valuable {
	v, _ := reflected.Primitive(reflect.ValueOf(val), vType)
	return v
}
