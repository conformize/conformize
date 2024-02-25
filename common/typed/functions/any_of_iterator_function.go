// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package functions

import (
	"github.com/conformize/conformize/common/typed"
	"github.com/conformize/conformize/predicates"
)

func AnyOf(it typed.Iterable, p predicates.Predicate, args typed.Valuable) (bool, error) {
	for it.Next() {
		var match bool
		var err error

		it.Element(func(v typed.Valuable) {
			match, err = p.Test(v, args)
		})

		if err != nil {
			return false, err
		}

		if match {
			return true, nil
		}
	}
	return false, nil
}
