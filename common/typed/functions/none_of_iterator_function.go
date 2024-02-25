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

func NoneOf(it typed.Iterable, p predicates.Predicate, args typed.Valuable) (bool, error) {
	match, err := AnyOf(it, p, args)
	return !match, err
}
