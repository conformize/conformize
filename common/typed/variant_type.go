// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package typed

type VariantTyped struct {
	VariantsTypes []Typeable
}

func (vt *VariantTyped) Hint() TypeHinter {
	variantHints := make([]TypeHinter, len(vt.VariantsTypes))
	for i, variantType := range vt.VariantsTypes {
		variantHints[i] = variantType.Hint()
	}

	return &complexTypeMixedElementsHint{
		kind:         Variant,
		elementsType: variantHints,
	}
}

func (vt *VariantTyped) Name() NamedType {
	return VariantType
}

func (vt *VariantTyped) GetVariantsTypes() []Typeable {
	return vt.VariantsTypes
}
