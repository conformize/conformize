// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package path

type Steps []PathStep

func (s *Steps) Add(step ...PathStep) Steps {
	*s = append(*s, step...)
	return *s
}

func (s *Steps) Next() (PathStep, bool) {
	if len(*s) == 0 {
		return nil, false
	}

	step := (*s)[0]
	*s = (*s)[1:]
	return step, true
}

func (s Steps) Clone() Steps {
	clone := make(Steps, len(s))
	copy(clone, s)

	return clone
}
