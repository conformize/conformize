// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package typed

import (
	"fmt"
	"reflect"
)

type VariantValue struct {
	Value         Valuable
	VariantsTypes []Typeable
}

func (vv *VariantValue) Type() Typeable {
	return &VariantTyped{VariantsTypes: vv.VariantsTypes}
}

func (vv *VariantValue) Assign(val Valuable) error {
	varVal, ok := val.(*VariantValue)
	if !ok {
		return fmt.Errorf("cannot apply %v to %v", val, vv)
	}
	*vv = *varVal
	return nil
}

func (vv *VariantValue) As(dst any) error {
	if dst == nil {
		return fmt.Errorf("destination cannot be nil")
	}

	targetVal := reflect.ValueOf(dst)
	if targetVal.Kind() != reflect.Ptr {
		return fmt.Errorf("destination must be a pointer")
	}

	targetValTypeHint := TypeHintOf(targetVal.Elem())
	for _, variantType := range vv.VariantsTypes {
		if targetValTypeHint.TypeHint() != variantType.Hint().TypeHint() {
			continue
		}

		var err error
		var variantVal Valuable
		variantVal, err = CreateValue(variantType)
		if err != nil {
			return err
		}

		if err = variantVal.Assign(vv.Value); err != nil {
			return err
		}

		if err = variantVal.As(dst); err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("no variant type matches destination type")
}

func (vv *VariantValue) String() string {
	return "VariantValue"
}

func NewVariantValue(value Valuable, variantsTypes []Typeable) Valuable {
	return &VariantValue{Value: value, VariantsTypes: variantsTypes}
}
