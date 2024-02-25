// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package typed

type BooleanTyped struct{}

var booleanHint = &simpleTypeHint{kind: Boolean}

func (bt *BooleanTyped) Name() NamedType {
	return BooleanType
}

func (bt *BooleanTyped) Hint() TypeHinter {
	return booleanHint
}
