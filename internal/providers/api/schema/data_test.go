// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package schema

import (
	"reflect"
	"testing"

	"github.com/conformize/conformize/common/typed"
	"github.com/conformize/conformize/internal/providers/api/schema/attributes"
)

func TestDataAsStruct(t *testing.T) {
	strVal, _ := typed.NewStringValue("test")
	data := &Data{
		Schema: &Schema{Attributes: map[string]Attributeable{
			"test": &attributes.StringAttribute{
				Description: "Test attribute",
				Required:    true,
			},
		}},
		Raw: typed.NewObjectValue(map[string]typed.Valuable{
			"test": strVal,
		}, map[string]typed.Typeable{
			"test": &typed.StringTyped{},
		}),
	}

	var test struct {
		Test string `cnfrmz:"test"`
	}

	var expected = struct {
		Test string `cnfrmz:"test"`
	}{"test"}

	if err := data.Get(&test); err != nil || !reflect.DeepEqual(test, expected) {
		t.Error(err)
	}
}

func TestDataAsNestedStruct(t *testing.T) {
	data := NewData(&Schema{
		Attributes: map[string]Attributeable{
			"test": &attributes.StringAttribute{},
			"nested": &attributes.ObjectAttribute{
				FieldsTypes: map[string]typed.Typeable{
					"testNested": &typed.BooleanTyped{},
				},
			},
		},
	})

	data.SetAtPath("test", "test")
	data.SetAtPath("nested.testNested", true)

	var test struct {
		Test   string `cnfrmz:"test"`
		Nested *struct {
			TestNested bool `cnfrmz:"testNested"`
		} `cnfrmz:"nested"`
	}

	var expected = struct {
		Test   string `cnfrmz:"test"`
		Nested *struct {
			TestNested bool `cnfrmz:"testNested"`
		} `cnfrmz:"nested"`
	}{"test", &struct {
		TestNested bool `cnfrmz:"testNested"`
	}{TestNested: true}}

	if err := data.Get(&test); err != nil || !reflect.DeepEqual(test, expected) {
		t.Error(err.Error())
	}
}

func TestStructToData(t *testing.T) {
	var test = struct {
		Host        string `cnfrmz:"host"`
		Environment string `cnfrmz:"environment"`
		Sandbox     bool   `cnfrmz:"sandbox"`
	}{"localhost", "development", true}

	data := NewData(
		&Schema{Attributes: map[string]Attributeable{
			"host": &attributes.StringAttribute{
				Description: "Host attribute",
				Required:    true,
			},
			"environment": &attributes.StringAttribute{
				Description: "Environment attribute",
				Required:    true,
			},
			"sandbox": &attributes.BooleanAttribute{
				Description: "Sandbox environment",
				Required:    true,
			},
		}},
	)

	if err := data.Set(&test); err != nil {
		t.Error(err)
	}
}

func TestGetDataValueAtPath(t *testing.T) {
	var test = struct {
		Host        string `cnfrmz:"host"`
		Environment string `cnfrmz:"environment"`
		Sandbox     bool   `cnfrmz:"sandbox"`
	}{"localhost", "development", true}

	data := NewData(
		&Schema{
			Attributes: map[string]Attributeable{
				"host": &attributes.StringAttribute{
					Description: "Host attribute",
					Required:    true,
				},
				"environment": &attributes.StringAttribute{
					Description: "Environment attribute",
					Required:    true,
				},
				"sandbox": &attributes.BooleanAttribute{
					Description: "Sandbox attribute",
					Required:    true,
				},
			},
		},
	)

	data.Set(&test)

	expected, _ := typed.NewStringValue("localhost")

	v, err := data.GetAtPath("host")
	if err != nil || !reflect.DeepEqual(v, expected) {
		t.Fail()
	}
}

func TestDataValueSetAtPath(t *testing.T) {
	data := NewData(
		&Schema{Attributes: map[string]Attributeable{
			"host": &attributes.StringAttribute{
				Description: "Host attribute",
				Required:    true,
			},
			"environment": &attributes.StringAttribute{
				Description: "Environment attribute",
				Required:    true,
			},
			"sandbox": &attributes.BooleanAttribute{
				Description: "Sandbox attribute",
				Required:    true,
			},
		}},
	)

	path := "host"
	data.SetAtPath(path, "example.com")

	expected, _ := typed.NewStringValue("example.com")
	v, err := data.GetAtPath(path)
	if err != nil || !reflect.DeepEqual(v, expected) {
		t.Error(err.Error())
	}
}

func TestDataGenericListValueSetAtPath(t *testing.T) {
	data := NewData(
		&Schema{Attributes: map[string]Attributeable{
			"hosts": &attributes.ListAttribute{
				Description:  "Hosts",
				Required:     true,
				ElementsType: &typed.StringTyped{},
			},
			"environment": &attributes.StringAttribute{
				Description: "Environment",
				Required:    true,
			},
			"sandbox": &attributes.BooleanAttribute{
				Description: "Sandbox",
				Required:    true,
			},
		}},
	)

	path := "hosts"
	data.SetAtPath(path, []string{"host-a", "host-b"})

	hostA, _ := typed.NewStringValue("host-a")
	hostB, _ := typed.NewStringValue("host-b")
	expected := typed.NewListValue(
		[]typed.Valuable{hostA, hostB},
		&typed.StringTyped{},
	)
	v, err := data.GetAtPath(path)
	if err != nil {
		t.Error(err.Error())
	}

	if !reflect.DeepEqual(v, expected) {
		t.Error("Values are not equal")
	}
}

func TestDataObjectValueSetAtPathFromString(t *testing.T) {
	data := NewData(
		&Schema{Attributes: map[string]Attributeable{
			"host": &attributes.ObjectAttribute{
				FieldsTypes: map[string]typed.Typeable{
					"address": &typed.StringTyped{},
				},
			},
		}},
	)

	path := "host.address"
	data.SetAtPath(path, "example.com")

	expected, _ := typed.NewStringValue("example.com")
	v, err := data.GetAtPath(path)
	if err != nil || !reflect.DeepEqual(v, expected) {
		t.Error(err.Error())
	}
}

func TestDataMapValueSetAtPathFromString(t *testing.T) {
	data := NewData(
		&Schema{Attributes: map[string]Attributeable{
			"host": &attributes.MapAttribute{
				ElementsType: &typed.MapTyped{
					ElementsType: &typed.StringTyped{},
				},
			},
		}},
	)

	path := "host.address.hostname"
	data.SetAtPath(path, "example.com")

	expected, _ := typed.NewStringValue("example.com")
	v, err := data.GetAtPath(path)
	if err != nil || !reflect.DeepEqual(v, expected) {
		t.Error(err.Error())
	}
}

func TestDataMapValueGetAtPathThatDoesNotExist(t *testing.T) {
	data := NewData(
		&Schema{Attributes: map[string]Attributeable{
			"host": &attributes.MapAttribute{
				ElementsType: &typed.MapTyped{
					ElementsType: &typed.StringTyped{},
				},
			},
		}},
	)

	path := "host.address.hostname"
	_, err := data.GetAtPath(path)
	if err == nil {
		t.Fail()
	}
}
