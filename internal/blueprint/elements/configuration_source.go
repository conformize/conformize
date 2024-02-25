// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package elements

import "gopkg.in/yaml.v3"

type ConfigurationSource struct {
	Provider     string                 `json:"provider" yaml:"provider"`
	Config       map[string]interface{} `json:"config,omitempty" yaml:"config"`
	ConfigFile   *string                `json:"configFile,omitempty" yaml:"configFile"`
	QueryOptions map[string]interface{} `json:"queryOptions,omitempty" yaml:"queryOptions"`
}

func (cs ConfigurationSource) MarshalYAML() (interface{}, error) {
	raw := make(map[string]interface{})
	raw["provider"] = cs.Provider

	if len(cs.Config) == 0 {
		raw["config"] = &yaml.Node{
			Kind:  yaml.ScalarNode,
			Style: yaml.FlowStyle,
		}
	} else {
		raw["config"] = cs.Config
	}

	if cs.ConfigFile != nil {
		raw["configFile"] = *cs.ConfigFile
	}

	if len(cs.QueryOptions) > 0 {
		raw["queryOptions"] = cs.QueryOptions
	}
	return raw, nil
}
