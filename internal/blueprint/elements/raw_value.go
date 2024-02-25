// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package elements

import (
	"encoding/json"

	"gopkg.in/yaml.v3"
)

type RawValue struct {
	Value     interface{}
	Sensitive bool
}

func (r *RawValue) GetValue() interface{} {
	return r.Value
}

func (r *RawValue) IsSensitive() bool {
	return r.Sensitive
}

func (r *RawValue) MarkSensitive() {
	r.Sensitive = true
}

func (r RawValue) MarshalYAML() (interface{}, error) {
	mapVal := r.asMap()
	if val := mapVal["value"]; isEmpty(val) {
		mapVal["value"] = &yaml.Node{
			Kind:  yaml.ScalarNode,
			Style: yaml.FlowStyle,
		}
	}
	return mapVal, nil
}

func (r *RawValue) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.asMap())
}

func (r *RawValue) asMap() map[string]interface{} {
	val := map[string]interface{}{
		"value": r.Value,
	}

	if r.Sensitive {
		sensitiveVal := make(map[string]interface{})
		sensitiveVal["sensitive"] = val
		return sensitiveVal
	}
	return val
}

func isEmpty(val interface{}) bool {
	if val == nil {
		return true
	}

	switch v := val.(type) {
	case string:
		return v == ""
	case []interface{}:
		return len(v) == 0
	case map[string]interface{}:
		return len(v) == 0
	default:
		return false
	}
}
