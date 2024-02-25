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
	"strings"

	"github.com/conformize/conformize/common/path"
	"github.com/conformize/conformize/common/pathparser"
	"gopkg.in/yaml.v3"
)

type Rule struct {
	Name      string
	Value     path.Path
	Predicate string
	Arguments Value
}

func (r *Rule) UnmarshalJSON(data []byte) error {
	var raw map[string]any
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	return unmarshalRaw(raw, r)
}

func (r *Rule) UnmarshalYAML(unmarshal func(any) error) error {
	var raw map[string]any
	if err := unmarshal(&raw); err != nil {
		return err
	}
	return unmarshalRaw(raw, r)
}

func unmarshalRaw(raw map[string]any, rule *Rule) error {
	matchedPredicate := false
	pathParser := pathparser.NewPathParser()
	for k, v := range raw {
		if k == "$value" {
			if val, ok := v.(string); ok {
				pathSteps, err := pathParser.Parse(val)
				if err != nil {
					return fmt.Errorf("couldn't parse value path '%s': %w", val, err)
				}
				rule.Value = *path.NewPath(pathSteps)
				continue
			}
			return fmt.Errorf(" \"$value\" attribute is not valid, expected value to be a string")
		}

		if k == "name" {
			if name, ok := v.(string); ok {
				rule.Name = name
				continue
			}
			return fmt.Errorf(" \"name\" attribute is not valid, expected value to be a string")
		}

		if !matchedPredicate {
			rule.Predicate = k
			matchedPredicate = true
			argVal, err := unmarshalArgumentValue(v)
			if err != nil {
				return err
			}
			rule.Arguments = argVal
			continue
		}

		return fmt.Errorf("condition '%s' already defined", rule.Predicate)
	}
	return nil
}

func unmarshalArgumentValue(rawArgVal any) (Value, error) {
	if argMap, ok := rawArgVal.(map[any]any); ok {
		if rawVal, valOk := argMap["sensitive"]; valOk {
			val, err := unmarshalArgumentValue(rawVal)
			if err != nil {
				return nil, fmt.Errorf("invalid sensitive value")
			}
			val.MarkSensitive()
			return val, nil
		}
	}

	if val, ok := rawArgVal.(string); ok && strings.HasPrefix(val, "$") {
		pathSteps, err := pathparser.NewPathParser().Parse(val)
		if err != nil {
			return nil, fmt.Errorf("couldn't parse path value '%s': %w", val, err)
		}
		return &PathValue{Path: *path.NewPath(pathSteps)}, nil
	}

	return &RawValue{Value: rawArgVal}, nil
}

func (r Rule) MarshalJSON() ([]byte, error) {
	raw := make(map[string]any)
	raw["name"] = r.Name
	raw["$value"] = r.Value

	if r.Arguments != nil {
		raw[r.Predicate] = r.Arguments
	} else {
		raw[r.Predicate] = nil
	}

	return json.Marshal(raw)
}

func (r Rule) MarshalYAML() (any, error) {
	raw := make(map[string]any)

	raw["name"] = &yaml.Node{
		Kind:  yaml.ScalarNode,
		Style: yaml.FlowStyle,
		Value: r.Name,
	}

	raw["$value"] = &yaml.Node{
		Kind:  yaml.ScalarNode,
		Style: yaml.FlowStyle,
		Value: r.Value.String(),
	}

	if r.Arguments != nil {
		raw[r.Predicate] = r.Arguments
	} else {
		raw[r.Predicate] = &yaml.Node{
			Kind: yaml.ScalarNode,
		}
	}

	return raw, nil
}
