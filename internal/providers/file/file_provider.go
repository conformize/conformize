// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package file

import (
	"github.com/conformize/conformize/common/diagnostics"
	"github.com/conformize/conformize/common/ds"
	sdk "github.com/conformize/conformize/internal/providers/api"
	"github.com/conformize/conformize/internal/providers/api/schema"
	"github.com/conformize/conformize/internal/providers/api/schema/attributes"
	"github.com/conformize/conformize/serialization"
)

type fileProvider struct {
	alias        string
	path         string `cnfrmz:"path"`
	unmarshaller serialization.SourceDataUnmarshaller
}

func (filePrvdr *fileProvider) Alias() string {
	return filePrvdr.alias
}

func (filePrvdr *fileProvider) Provide(req *sdk.ProviderDataRequest) (*ds.Node[string, any], *diagnostics.Diagnostics) {
	diags := diagnostics.NewDiagnostics()
	if req != nil {
		if pathVal, err := req.GetAtPath("path"); err == nil {
			pathVal.As(&filePrvdr.path)
		} else {
			diags.Append(diagnostics.Builder().Error().Details(err.Error()).Build())
			return nil, diags
		}
	}

	fileSrc, err := serialization.NewFileSource(filePrvdr.path)
	if err != nil {
		diags.Append(diagnostics.Builder().Error().Details(err.Error()).Build())
		return nil, diags
	}

	data, err := filePrvdr.unmarshaller.Unmarshal(fileSrc)
	if err == nil {
		return data, diags
	}
	diags.Append(diagnostics.Builder().Error().Details(err.Error()).Build())
	return nil, diags
}

func (filePrvdr *fileProvider) ProvisionDataRequestSchema() *schema.Schema {
	return &schema.Schema{
		Description: "File resource request schema",
		Version:     1,
		Attributes: map[string]schema.Attributeable{
			"path": &attributes.StringAttribute{},
		},
	}
}

func (filePrvdr *fileProvider) ConfigurationSchema() *schema.Schema {
	return &schema.Schema{
		Description: "File provider schema",
		Version:     1,
		Attributes: map[string]schema.Attributeable{
			"path": &attributes.StringAttribute{},
		},
	}
}

func (filePrvdr *fileProvider) Configure(req *sdk.ConfigurationRequest) error {
	if pathVal, err := req.GetAtPath("path"); err == nil {
		var filePath string
		var err error
		if err = pathVal.As(&filePath); err != nil {
			return err
		}

		filePrvdr.path = filePath
	} else {
		return err
	}
	return nil
}

func NewFileProvider(alias string, unmarshaller serialization.SourceDataUnmarshaller) *fileProvider {
	return &fileProvider{alias: alias, unmarshaller: unmarshaller}
}
