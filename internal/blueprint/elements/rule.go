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
)

type Rule struct {
	Value     string  `json:"$value" yaml:"$value"`
	Predicate string  `json:"predicate" yaml:"predicate"`
	Arguments []Value `json:"arguments" yaml:"arguments"`
}

func (r *Rule) UnmarshalJSON(data []byte) error {
	type Alias Rule
	rule := struct {
		Arguments []map[interface{}]interface{} `json:"arguments"`
		*Alias
	}{
		Alias: (*Alias)(r),
	}

	var err error
	if err = json.Unmarshal(data, &rule); err != nil {
		return err
	}

	if r.Arguments, err = unmarshalArgumentValues(rule.Arguments); err != nil {
		return err
	}
	return nil
}

func (r *Rule) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var rule struct {
		Value     string                        `yaml:"$value"`
		Predicate string                        `yaml:"predicate"`
		Arguments []map[interface{}]interface{} `yaml:"arguments"`
	}

	var err error
	if err = unmarshal(&rule); err != nil {
		return err
	}

	r.Value = rule.Value
	r.Predicate = rule.Predicate
	if r.Arguments, err = unmarshalArgumentValues(rule.Arguments); err != nil {
		return err
	}
	return nil
}

func unmarshalArgumentValues(arguments []map[interface{}]interface{}) ([]Value, error) {
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

func unmarshalArgumentValue(argVal map[interface{}]interface{}) (Value, error) {
	if len(argVal) > 1 {
		return nil, fmt.Errorf("invalid argument: %v, only one value can be specified", argVal)
	}

	var err error
	var valOk bool
	var rawVal interface{}
	if rawVal, valOk = argVal["sensitive"]; valOk {
		if mapVal, ok := rawVal.(map[interface{}]interface{}); ok {
			var val Value
			if val, err = unmarshalArgumentValue(mapVal); err == nil {
				val.MarkSensitive()
				return val, nil
			}
		}
		return nil, fmt.Errorf("invalid sensitive value")
	} else if rawVal, valOk = argVal["path"]; valOk {
		if pathStr, ok := rawVal.(string); ok {
			return &PathValue{Path: pathStr}, nil
		}
		err = fmt.Errorf("invalid path value: %v", rawVal)
	} else if rawVal, valOk = argVal["value"]; valOk {
		return &RawValue{Value: rawVal}, nil
	} else {
		err = fmt.Errorf("invalid argument value: %v", argVal)
	}
	return nil, err
}
