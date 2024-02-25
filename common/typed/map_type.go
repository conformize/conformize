// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package typed

type MapTyped struct {
	ElementsType Typeable
}

func (mt *MapTyped) Name() NamedType {
	return MapType
}

func (mt *MapTyped) Hint() TypeHinter {
	return &complexTypeHint{
		kind:     Map,
		elemType: mt.ElementsType.Hint(),
	}
}

func (mt *MapTyped) ElementType() Typeable {
	return mt.ElementsType
}
