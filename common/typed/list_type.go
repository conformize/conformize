// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package typed

type ListTyped struct {
	ElementsType Typeable
}

func (lt *ListTyped) Name() NamedType {
	return ListType
}

func (lt *ListTyped) Hint() TypeHinter {
	return &complexTypeHint{kind: List, elemType: lt.ElementsType.Hint()}
}

func (lt *ListTyped) ElementType() Typeable {
	return lt.ElementsType
}
