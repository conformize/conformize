// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package sdk

import "github.com/conformize/conformize/internal/providers/api/schema"

type ConfigurationRequest struct {
	Schema schema.Schemable
	*schema.Data
}

func NewConfigurationRequest(s schema.Schemable) *ConfigurationRequest {
	return &ConfigurationRequest{
		Schema: s,
		Data:   schema.NewData(s),
	}
}
