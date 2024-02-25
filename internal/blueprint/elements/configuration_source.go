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
	Config       map[string]interface{}
	ConfigFile   *string
	QueryOptions map[string]interface{}
}

func (cs ConfigurationSource) MarshalYAML() (interface{}, error) {
	raw := make(map[string]map[string]interface{})
	raw[cs.Provider] = make(map[string]interface{})

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
	raw := make(map[string]map[string]interface{})
	raw[cs.Provider] = make(map[string]interface{})

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
	var raw map[string]interface{}
	if err = json.Unmarshal(data, &raw); err != nil {
		return err
	}

	if err := rawUnmarshal(raw, cs); err != nil {
		return err
	}
	return nil
}

func (cs *ConfigurationSource) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var err error
	var raw map[string]interface{}
	if err = unmarshal(&raw); err != nil {
		return err
	}
	if err := rawUnmarshal(raw, cs); err != nil {
		return err
	}
	return nil
}

func rawUnmarshal(raw map[string]interface{}, configSource *ConfigurationSource) error {
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

func unmarshalMap(input interface{}) (map[string]interface{}, error) {
	unmarshalled := make(map[string]interface{})

	switch input := input.(type) {
	case map[interface{}]interface{}:
		for key, value := range input {
			strKey, ok := key.(string)
			if !ok {
				return nil, fmt.Errorf("non-string key in map")
			}

			if mapVal, ok := value.(map[interface{}]interface{}); ok {
				var err error
				unmarshalled[strKey], err = unmarshalMap(mapVal)
				if err != nil {
					return nil, err
				}
				continue
			}
			unmarshalled[strKey] = value
		}
	case map[string]interface{}:
		return input, nil
	default:
		return nil, fmt.Errorf("invalid type, expected map, got: %T", input)
	}

	return unmarshalled, nil
}
