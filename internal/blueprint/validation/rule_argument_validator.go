// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package validation

import (
	"fmt"

	"github.com/conformize/conformize/internal/blueprint/elements"
)

type ruleArgumentValidator struct{}

func (ruleArgValidator *ruleArgumentValidator) Validate(value elements.Value, configSources map[string]elements.ConfigurationSource, refs map[string]string) error {
	switch argVal := value.(type) {
	case *elements.RawValue:
		return nil
	case *elements.PathValue:
		pathSteps := argVal.Path.Steps()
		root := pathSteps[0].String()
		if _, sourceRefOk := configSources[root]; !sourceRefOk {
			if _, aliasRefOk := refs[root]; !aliasRefOk {
				return fmt.Errorf("couldn't resolve root '%s' in path %s", root, argVal.Path.String())
			}
		}
	default:
		return fmt.Errorf("invalid argument type")
	}
	return nil
}
