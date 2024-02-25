// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package schema

import (
	"reflect"

	"github.com/conformize/conformize/common/path"
	"github.com/conformize/conformize/common/reflected"
	"github.com/conformize/conformize/common/typed"
)

type Data struct {
	Schema Schemable
	Raw    typed.Valuable
}

func (d *Data) Get(target any) error {
	return d.Raw.As(target)
}

func (d *Data) GetAtPath(pathStr string) (typed.Valuable, error) {
	p, err := path.NewFromString(pathStr)
	if err != nil {
		return nil, err
	}
	w := path.ValuePathWalk{
		Destination: d.Raw,
	}
	return w.Walk(p)
}

func (d *Data) Set(value any) error {
	val, err := reflected.Value(reflect.ValueOf(value), d.Schema.Type())
	d.Raw = val
	return err
}

func (d *Data) SetAtPath(pathStr string, value any) error {
	p, err := path.NewFromString(pathStr)
	if err != nil {
		return err
	}

	w := path.ValuePathWalk{
		Destination:       d.Raw,
		CreateValueAtPath: true,
	}

	val, err := w.Walk(p)
	if err != nil {
		return err
	}
	newVal, err := reflected.Value(reflect.ValueOf(value), val.Type())
	if err != nil {
		return err
	}
	return val.Assign(newVal)
}

func NewData(schema Schemable) *Data {
	objVal, _ := typed.CreateValue(schema.Type())
	data := &Data{
		Schema: schema,
		Raw:    objVal,
	}

	attr := schema.GetAttributes()
	for name, attribute := range attr {
		if attribute.GetDefaultValue() != nil {
			data.SetAtPath(name, attribute.GetDefaultValue())
			continue
		}

		if attribute.GetDefaultValueFn() != nil {
			data.SetAtPath(name, attribute.GetDefaultValueFn()())
		}
	}
	return data
}
