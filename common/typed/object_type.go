// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package typed

type ObjectTyped struct {
	FieldsTypes map[string]Typeable
}

func (ot *ObjectTyped) Name() NamedType {
	return ObjectType
}

func (ot *ObjectTyped) Hint() TypeHinter {
	fieldHints := make(map[string]TypeHinter, len(ot.FieldsTypes))
	for key, typ := range ot.FieldsTypes {
		fieldHints[key] = typ.Hint()
	}

	return &complexObjectHint{
		kind:       Object,
		fieldsType: fieldHints,
	}
}

func (ot *ObjectTyped) GetFieldsTypes() map[string]Typeable {
	return ot.FieldsTypes
}
