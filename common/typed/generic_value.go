// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package typed

type GenericValue struct {
	Value Valuable
}

func (gv *GenericValue) Type() Typeable {
	return &GenericTyped{}
}

func (gv *GenericValue) Assign(val Valuable) error {
	gv.Value = val
	return nil
}

func (gv *GenericValue) String() string {
	return "GenericValue"
}

func (gv *GenericValue) As(dst any) error {
	if valuable, ok := dst.(Valuable); ok {
		return valuable.Assign(gv.Value)
	}
	return gv.Value.As(dst)
}

func NewGenericValue(value Valuable) (Valuable, error) {
	return &GenericValue{Value: value}, nil
}
