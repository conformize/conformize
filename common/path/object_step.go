// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package path

type ObjectStep string

func (obj ObjectStep) String() string {
	return string(obj)
}

func (obj ObjectStep) Equal(o PathStep) bool {
	other, ok := o.(ObjectStep)
	if ok {
		return string(obj) == string(other)
	}
	return false
}
