// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package path

type KeyStep string

func (k KeyStep) String() string {
	return string(k)
}

func (k KeyStep) Equal(o PathStep) bool {
	other, ok := o.(KeyStep)
	if ok {
		return string(k) == string(other)
	}
	return false
}
