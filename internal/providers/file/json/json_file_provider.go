// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package json

import (
	"github.com/conformize/conformize/common/diagnostics"
	"github.com/conformize/conformize/common/ds"
	sdk "github.com/conformize/conformize/internal/providers/api"
	"github.com/conformize/conformize/internal/providers/api/schema"
	"github.com/conformize/conformize/internal/providers/api/schema/attributes"
	"github.com/conformize/conformize/serialization"
	"github.com/conformize/conformize/serialization/unmarshal/yaml"
)

type JsonFileProvider struct {
	path string
}

func (jsonFilePrvdr *JsonFileProvider) Provide(req *sdk.ProviderDataRequest) (*ds.Node[string, interface{}], *diagnostics.Diagnostics) {
	diags := diagnostics.NewDiagnostics()
	if req != nil {
		if pathVal, err := req.GetAtPath("Path"); err == nil {
			pathVal.As(&jsonFilePrvdr.path)
		} else {
			diags.Append(diagnostics.Builder().Error().Details(err.Error()).Build())
			return nil, diags
		}
	}

	fileSrc, err := serialization.NewFileSource(jsonFilePrvdr.path)
	if err != nil {
		diags.Append(diagnostics.Builder().Error().Details(err.Error()).Build())
		return nil, diags
	}

	unmarshaller := yaml.NewYamlUnmarshal(fileSrc)
	if data, err := unmarshaller.Unmarshal(); err == nil {
		return data, diags
	} else {
		diags.Append(diagnostics.Builder().Error().Details(err.Error()).Build())
		return nil, diags
	}
}

func (jsonFilePrvdr *JsonFileProvider) ProvisionDataRequestSchema() *schema.Schema {
	return &schema.Schema{
		Description: "Json file resource request schema",
		Version:     1,
		Attributes: map[string]schema.Attributeable{
			"Path": &attributes.StringAttribute{},
		},
	}
}

func (jsonFilePrvdr *JsonFileProvider) ConfigurationSchema() *schema.Schema {
	return &schema.Schema{
		Description: "Json file provider schema",
		Version:     1,
		Attributes: map[string]schema.Attributeable{
			"Path": &attributes.StringAttribute{},
		},
	}
}

func (jsonFilePrvdr *JsonFileProvider) Configure(req *sdk.ConfigurationRequest) error {
	if pathVal, err := req.GetAtPath("Path"); err == nil {
		return pathVal.As(&jsonFilePrvdr.path)
	} else {
		return err
	}
}
