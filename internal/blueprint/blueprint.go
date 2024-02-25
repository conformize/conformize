// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package blueprint

import "github.com/conformize/conformize/internal/blueprint/elements"

type Blueprint struct {
	Version    float64                                 `json:"version" yaml:"version"`
	References map[string]string                       `json:"$refs" yaml:"$refs"`
	Sources    map[string]elements.ConfigurationSource `json:"sources" yaml:"sources"`
	Ruleset    []elements.Rule                         `json:"ruleset" yaml:"ruleset"`
}
