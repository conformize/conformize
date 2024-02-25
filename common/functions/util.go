// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package functions

import (
	"fmt"
	"os"
	"reflect"
	"regexp"

	"github.com/conformize/conformize/common/reflected"
	"github.com/conformize/conformize/common/typed"
	"github.com/conformize/conformize/serialization/unmarshal/functions"
)

var envVarPlaceHolderExp = regexp.MustCompile(`(?m)\$\{?([A-Za-z_][A-Za-z0-9_]*)\}?|\$([A-Za-z_][A-Za-z0-9_])`)

func FindEnvVarIndices(s string) ([][]int, bool) {
	matches := envVarPlaceHolderExp.FindAllStringSubmatchIndex(s, -1)
	return matches, len(matches) > 0
}

func parseEnvVars(strVal string, envVarsIdxs [][]int) (*string, error) {
	for i := len(envVarsIdxs) - 1; i >= 0; i-- {
		envVarIdxs := envVarsIdxs[i]
		envVar := strVal[envVarIdxs[2]:envVarIdxs[3]]
		if envVal, ok := os.LookupEnv(envVar); ok {
			if innerEnvVarIdxs, found := FindEnvVarIndices(envVal); !found {
				strVal = strVal[:envVarIdxs[0]] + envVal + strVal[envVarIdxs[1]:]
			} else if interpolatedInnerEnvVars, err := parseEnvVars(envVal, innerEnvVarIdxs); err == nil {
				strVal = strVal[:envVarIdxs[0]] + *interpolatedInnerEnvVars + strVal[envVarIdxs[1]:]
			} else {
				return nil, err
			}
		} else {
			return nil, fmt.Errorf("environment variable %s not set", envVar)
		}
	}
	return &strVal, nil
}

func ParseRawValue(val any) (typed.Valuable, error) {
	if val == nil {
		return nil, fmt.Errorf("cannot parse <nil> value")
	}

	reflectVal := reflect.ValueOf(val)
	valTypeHint := typed.TypeHintOf(reflectVal)
	if v, err := reflected.ValueFromTypeHint(reflectVal, valTypeHint); err == nil {
		if v.Type().Hint().TypeHint() != typed.String {
			return v, nil
		}

		var strVal string
		v.As(&strVal)
		if envVarsIdxs, found := FindEnvVarIndices(strVal); found {
			interpolatedVal, err := parseEnvVars(strVal, envVarsIdxs)
			if err != nil {
				return nil, err
			}
			strVal = *interpolatedVal
		}

		if decodedVal, err := functions.DecodeStringValue(strVal); err == nil {
			reflectDecodedVal := reflect.ValueOf(decodedVal)
			decodedValTypeHint := typed.TypeHintOf(reflectDecodedVal)
			return reflected.ValueFromTypeHint(reflectDecodedVal, decodedValTypeHint)
		}
		return v, nil
	} else {
		return nil, err
	}
}

func LookupEnvVar(envVar string) (string, error) {
	if envVal, ok := os.LookupEnv(envVar); ok {
		if innerEnvVarIdxs, found := FindEnvVarIndices(envVal); found {
			if interpolatedInnerEnvVars, err := parseEnvVars(envVal, innerEnvVarIdxs); err == nil {
				return *interpolatedInnerEnvVars, nil
			} else {
				return "", err
			}
		} else {
			return envVal, nil
		}
	}
	return "", fmt.Errorf("environment variable %s not set", envVar)
}

func InterpolateEnvVars(v string) (string, error) {
	if envVarsIdxs, found := FindEnvVarIndices(v); found {
		if interpolatedVal, err := parseEnvVars(v, envVarsIdxs); err == nil {
			return *interpolatedVal, nil
		} else {
			return "", err
		}
	}
	return v, nil
}
