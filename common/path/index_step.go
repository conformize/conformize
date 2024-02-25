// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package path

import "strconv"

type IndexStep int

func (a IndexStep) String() string {
	return strconv.FormatInt(int64(a), 10)
}

func (a IndexStep) Equal(o PathStep) bool {
	other, ok := o.(IndexStep)
	if ok {
		return int(a) == int(other)
	}
	return false
}
