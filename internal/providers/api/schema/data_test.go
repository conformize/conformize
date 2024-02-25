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
			"Test": &attributes.StringAttribute{
				Description: "Test",
				Required:    true,
			},
		}},
		Raw: typed.NewObjectValue(map[string]typed.Valuable{
			"Test": strVal,
		}, map[string]typed.Typeable{
			"Test": &typed.StringTyped{},
		}),
	}

	var test struct {
		Test string
	}

	var expected = struct {
		Test string
	}{"test"}

	if err := data.Get(&test); err != nil || !reflect.DeepEqual(test, expected) {
		t.Error(err)
	}
}

func TestStructToData(t *testing.T) {
	var test = struct {
		Host        string
		Environment string
		Sandbox     bool
	}{"localhost", "development", true}

	data := NewData(
		&Schema{Attributes: map[string]Attributeable{
			"Host": &attributes.StringAttribute{
				Description: "Host",
				Required:    true,
			},
			"Environment": &attributes.StringAttribute{
				Description: "Environment",
				Required:    true,
			},
			"Sandbox": &attributes.BooleanAttribute{
				Description: "Sandbox",
				Required:    true,
			},
		}},
	)

	if err := data.Set(test); err != nil {
		t.Error(err)
	}
}

func TestGetDataValueAtPath(t *testing.T) {
	var test = struct {
		Host        string
		Environment string
		Sandbox     bool
	}{"localhost", "development", true}

	data := &Data{
		Schema: &Schema{Attributes: map[string]Attributeable{
			"Host": &attributes.StringAttribute{
				Description: "Host",
				Required:    true,
			},
			"Environment": &attributes.StringAttribute{
				Description: "Environment",
				Required:    true,
			},
			"Sandbox": &attributes.BooleanAttribute{
				Description: "Sandbox",
				Required:    true,
			},
		}},
	}

	data.Set(test)

	expected, _ := typed.NewStringValue("localhost")

	v, err := data.GetAtPath("Host")
	if err != nil || !reflect.DeepEqual(v, expected) {
		t.Fail()
	}
}

func TestDataValueSetAtPath(t *testing.T) {
	data := NewData(
		&Schema{Attributes: map[string]Attributeable{
			"Host": &attributes.StringAttribute{
				Description: "Host",
				Required:    true,
			},
			"Environment": &attributes.StringAttribute{
				Description: "Environment",
				Required:    true,
			},
			"Sandbox": &attributes.BooleanAttribute{
				Description: "Sandbox",
				Required:    true,
			},
		}},
	)

	path := "Host"
	data.SetAtPath(path, "example.com")

	expected, _ := typed.NewStringValue("example.com")
	v, err := data.GetAtPath(path)
	if err != nil || !reflect.DeepEqual(v, expected) {
		t.Fail()
	}
}

func TestDataGenericListValueSetAtPath(t *testing.T) {
	data := NewData(
		&Schema{Attributes: map[string]Attributeable{
			"Hosts": &attributes.ListAttribute{
				Description:  "Hosts",
				Required:     true,
				ElementsType: &typed.GenericTyped{},
			},
			"Environment": &attributes.StringAttribute{
				Description: "Environment",
				Required:    true,
			},
			"Sandbox": &attributes.BooleanAttribute{
				Description: "Sandbox",
				Required:    true,
			},
		}},
	)

	path := "Hosts"
	data.SetAtPath(path, []string{"host-a", "host-b"})

	hostA, _ := typed.NewStringValue("host-a")
	hostB, _ := typed.NewStringValue("host-b")
	expected := typed.NewListValue(
		[]typed.Valuable{hostA, hostB},
		&typed.GenericTyped{},
	)
	v, err := data.GetAtPath(path)
	if err != nil || !reflect.DeepEqual(v, expected) {
		t.Fail()
	}
}

func TestDataObjectValueSetAtPathFromString(t *testing.T) {
	data := NewData(
		&Schema{Attributes: map[string]Attributeable{
			"Host": &attributes.ObjectAttribute{
				FieldsTypes: map[string]typed.Typeable{
					"Address": &typed.StringTyped{},
				},
			},
		}},
	)

	path := "Host.Address"
	data.SetAtPath(path, "example.com")

	expected, _ := typed.NewStringValue("example.com")
	v, err := data.GetAtPath(path)
	if err != nil || !reflect.DeepEqual(v, expected) {
		t.Fail()
	}
}

func TestDataMapValueSetAtPathFromString(t *testing.T) {
	data := NewData(
		&Schema{Attributes: map[string]Attributeable{
			"Host": &attributes.MapAttribute{
				ElementsType: &typed.MapTyped{
					ElementsType: &typed.StringTyped{},
				},
			},
		}},
	)

	path := "Host.Address.Hostname"
	data.SetAtPath(path, "example.com")

	expected, _ := typed.NewStringValue("example.com")
	v, err := data.GetAtPath(path)
	if err != nil || !reflect.DeepEqual(v, expected) {
		t.Fail()
	}
}

func TestDataMapValueGetAtPathThatDoesNotExist(t *testing.T) {
	data := NewData(
		&Schema{Attributes: map[string]Attributeable{
			"Host": &attributes.MapAttribute{
				ElementsType: &typed.MapTyped{
					ElementsType: &typed.StringTyped{},
				},
			},
		}},
	)

	path := "Host.Address.Hostname"
	v, err := data.GetAtPath(path)
	if v != nil || err == nil {
		t.Fail()
	}
}
