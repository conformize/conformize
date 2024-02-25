// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package typed

type TupleTyped struct {
	ElementsTypes []Typeable
}

func (tup *TupleTyped) Hint() TypeHint {
	return Tuple
}

func (tup *TupleTyped) Name() NamedType {
	return TupleType
}

func (tup *TupleTyped) GetElementsTypes() []Typeable {
	return tup.ElementsTypes
}

func (tup *TupleTyped) ValueType() Valuable {
	return &TupleValue{ElementsTypes: tup.ElementsTypes}
}
