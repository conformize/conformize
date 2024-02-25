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

	"github.com/conformize/conformize/common/pathparser"
	"gopkg.in/yaml.v3"
)

type Rule struct {
	Name      string
	Value     string
	Predicate string
	Arguments []Value
}

func (r *Rule) UnmarshalJSON(data []byte) error {
	var err error
	var raw map[string]interface{}
	if err = json.Unmarshal(data, &raw); err != nil {
		return err
	}

	if err := unmarshalRaw(raw, r); err != nil {
		return err
	}
	return nil
}

func (r *Rule) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var err error
	var raw map[string]interface{}
	if err = unmarshal(&raw); err != nil {
		return err
	}
	if err := unmarshalRaw(raw, r); err != nil {
		return err
	}
	return nil
}

func unmarshalRaw(raw map[string]interface{}, rule *Rule) error {
	matchedPredicate := false
	for k, v := range raw {
		if k == "$value" {
			if val, ok := v.(string); ok {
				rule.Value = val
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
			if predicateVal, ok := v.([]interface{}); ok {
				args, err := unmarshalArgumentValues(predicateVal)
				if err != nil {
					return err
				}
				rule.Arguments = args
				continue
			}
		} else {
			return fmt.Errorf("condition '%s' already defined", rule.Predicate)
		}
	}
	return nil
}

func unmarshalArgumentValues(arguments []interface{}) ([]Value, error) {
	var args []Value
	for _, rawArg := range arguments {
		val, err := unmarshalArgumentValue(rawArg)
		if err != nil {
			return nil, err
		}
		args = append(args, val)
	}
	return args, nil
}

func unmarshalArgumentValue(rawArgVal interface{}) (Value, error) {
	argVal, ok := rawArgVal.(map[interface{}]interface{})
	if ok {
		rawVal, valOk := argVal["sensitive"]
		if valOk {
			if val, err := unmarshalArgumentValue(rawVal); err == nil {
				val.MarkSensitive()
				return val, nil
			}
			return nil, fmt.Errorf("invalid sensitive value")
		}
	}

	if val, ok := rawArgVal.(string); ok && strings.HasPrefix(val, "$") {
		pathParser := pathparser.NewPathParser()
		if _, err := pathParser.Parse(val); err != nil {
			return &RawValue{Value: val}, nil
		}
		return &PathValue{Path: val}, nil
	}

	return &RawValue{Value: rawArgVal}, nil
}

func (r Rule) MarshalJSON() ([]byte, error) {
	raw := make(map[string]interface{})
	raw["name"] = r.Name
	raw["$value"] = r.Value

	var args interface{}
	if len(r.Arguments) > 0 {
		args = r.Arguments
	} else {
		args = make([]interface{}, 0)
	}
	raw[r.Predicate] = args
	return json.Marshal(raw)
}

func (r Rule) MarshalYAML() (interface{}, error) {
	raw := make(map[string]interface{})
	raw["name"] = &yaml.Node{
		Kind:  yaml.ScalarNode,
		Style: yaml.FlowStyle,
		Value: r.Name,
	}

	raw["$value"] = &yaml.Node{
		Kind:  yaml.ScalarNode,
		Style: yaml.FlowStyle,
		Value: r.Value,
	}

	var args interface{}
	if len(r.Arguments) > 0 {
		args = r.Arguments
	} else {
		args = &yaml.Node{
			Kind: yaml.ScalarNode,
		}
	}
	raw[r.Predicate] = args
	return raw, nil
}
