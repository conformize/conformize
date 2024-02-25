// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package path

type PropertyStep string

func (a PropertyStep) String() string {
	return string(a)
}

func (a PropertyStep) Equal(o PathStep) bool {
	other, ok := o.(PropertyStep)
	if ok {
		return string(a) == string(other)
	}
	return false
}
