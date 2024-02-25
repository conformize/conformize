// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package xml

import (
	"github.com/conformize/conformize/common/diagnostics"
	"github.com/conformize/conformize/common/ds"
	sdk "github.com/conformize/conformize/internal/providers/api"
	"github.com/conformize/conformize/internal/providers/api/schema"
	"github.com/conformize/conformize/internal/providers/api/schema/attributes"
	"github.com/conformize/conformize/serialization"
	"github.com/conformize/conformize/serialization/unmarshal/xml"
)

type XmlFileProvider struct {
	path string
}

func (xmlFilePrvdr *XmlFileProvider) Provide(req *sdk.ProviderDataRequest) (*ds.Node[string, interface{}], *diagnostics.Diagnostics) {
	diags := diagnostics.NewDiagnostics()
	if req != nil {
		if pathVal, err := req.GetAtPath("Path"); err == nil {
			pathVal.As(&xmlFilePrvdr.path)
		} else {
			diags.Append(diagnostics.Builder().Error().Details(err.Error()).Build())
			return nil, diags
		}
	}

	fileSrc, err := serialization.NewFileSource(xmlFilePrvdr.path)
	if err != nil {
		diags.Append(diagnostics.Builder().Error().Details(err.Error()).Build())
		return nil, diags
	}

	unmarshaller := xml.NewXmlFileUnmarshal(fileSrc)
	if data, err := unmarshaller.Unmarshal(); err == nil {
		return data, diags
	} else {
		diags.Append(diagnostics.Builder().Error().Details(err.Error()).Build())
		return nil, diags
	}
}

func (xmlFilePrvdr *XmlFileProvider) ProvisionDataRequestSchema() *schema.Schema {
	return &schema.Schema{
		Description: "Xml file resource request schema",
		Version:     1,
		Attributes: map[string]schema.Attributeable{
			"Path": &attributes.StringAttribute{},
		},
	}
}

func (xmlFilePrvdr *XmlFileProvider) ConfigurationSchema() *schema.Schema {
	return &schema.Schema{
		Description: "Xml file provider schema",
		Version:     1,
		Attributes: map[string]schema.Attributeable{
			"Path": &attributes.StringAttribute{},
		},
	}
}

func (xmlFilePrvdr *XmlFileProvider) Configure(req *sdk.ConfigurationRequest) error {
	if pathVal, err := req.GetAtPath("Path"); err == nil {
		return pathVal.As(&xmlFilePrvdr.path)
	} else {
		return err
	}
}
