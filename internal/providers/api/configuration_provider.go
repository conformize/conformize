// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package sdk

import (
	"github.com/conformize/conformize/common/diagnostics"
	"github.com/conformize/conformize/common/ds"
	"github.com/conformize/conformize/internal/providers/api/schema"
)

type ConfigurationProvider interface {
	Alias() string
	ConfigurationSchema() *schema.Schema
	ProvisionDataRequestSchema() *schema.Schema
	Configure(req *ConfigurationRequest) error
	Provide(req *ProviderDataRequest) (*ds.Node[string, any], *diagnostics.Diagnostics)
}
