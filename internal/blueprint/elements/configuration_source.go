// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package elements

import (
	"encoding/json"
	"fmt"

	"gopkg.in/yaml.v3"
)

type ConfigurationSource struct {
	Provider     string
	Config       map[string]any
	ConfigFile   *string
	QueryOptions map[string]any
}

func (cs ConfigurationSource) MarshalYAML() (any, error) {
	raw := make(map[string]map[string]any)
	raw[cs.Provider] = make(map[string]any)

	if len(cs.Config) == 0 {
		raw[cs.Provider]["config"] = &yaml.Node{
			Kind:  yaml.ScalarNode,
			Style: yaml.FlowStyle,
		}
	} else {
		raw[cs.Provider]["config"] = cs.Config
	}

	if cs.ConfigFile != nil {
		raw[cs.Provider]["configFile"] = *cs.ConfigFile
	}

	if len(cs.QueryOptions) > 0 {
		raw[cs.Provider]["queryOptions"] = cs.QueryOptions
	}
	return raw, nil
}

func (cs ConfigurationSource) MarshalJSON() ([]byte, error) {
	raw := make(map[string]map[string]any)
	raw[cs.Provider] = make(map[string]any)

	if cs.ConfigFile != nil {
		raw[cs.Provider]["configFile"] = *cs.ConfigFile
	}

	if len(cs.Config) > 0 {
		raw[cs.Provider]["config"] = cs.Config
	}

	if len(cs.QueryOptions) > 0 {
		raw[cs.Provider]["queryOptions"] = cs.QueryOptions
	}
	return json.Marshal(raw)
}

func (cs *ConfigurationSource) UnmarshalJSON(data []byte) error {
	var err error
	var raw map[string]any
	if err = json.Unmarshal(data, &raw); err != nil {
		return err
	}

	if err := rawUnmarshal(raw, cs); err != nil {
		return err
	}
	return nil
}

func (cs *ConfigurationSource) UnmarshalYAML(unmarshal func(any) error) error {
	var err error
	var raw map[string]any
	if err = unmarshal(&raw); err != nil {
		return err
	}
	if err := rawUnmarshal(raw, cs); err != nil {
		return err
	}
	return nil
}

func rawUnmarshal(raw map[string]any, configSource *ConfigurationSource) error {
	matchedProvider := false

	for provider, v := range raw {
		if matchedProvider {
			return fmt.Errorf("multiple providers specified")
		}

		configSource.Provider = provider
		matchedProvider = true

		mapVal, err := unmarshalMap(v)
		if err != nil {
			return fmt.Errorf("couldn't unmarshall provider configuration, reason:\n%v", err)
		}

		for key, value := range mapVal {
			switch key {
			case "config":
				configSource.Config, err = unmarshalMap(value)
			case "configFile":
				if configFile, ok := value.(string); ok {
					configSource.ConfigFile = &configFile
				} else {
					return fmt.Errorf("invalid type for configFile")
				}

			case "queryOptions":
				configSource.QueryOptions, err = unmarshalMap(value)
			default:
				continue
			}

			if err != nil {
				return err
			}
		}
	}

	return nil
}

func unmarshalMap(input any) (map[string]any, error) {
	unmarshalled := make(map[string]any)

	switch input := input.(type) {
	case map[any]any:
		for key, value := range input {
			strKey, ok := key.(string)
			if !ok {
				return nil, fmt.Errorf("non-string key in map")
			}

			if mapVal, ok := value.(map[any]any); ok {
				var err error
				unmarshalled[strKey], err = unmarshalMap(mapVal)
				if err != nil {
					return nil, err
				}
				continue
			}
			unmarshalled[strKey] = value
		}
	case map[string]any:
		return input, nil
	default:
		return nil, fmt.Errorf("invalid type, expected map, got: %T", input)
	}

	return unmarshalled, nil
}
