// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package validation

import (
	"fmt"

	"github.com/conformize/conformize/common/diagnostics"
	"github.com/conformize/conformize/internal/blueprint"
)

type BlueprintRulesValidator struct {
	ruleValidator
}

func (rlsVld *BlueprintRulesValidator) Validate(blueprint *blueprint.Blueprint) *diagnostics.Diagnostics {
	diags := diagnostics.NewDiagnostics()
	if len(blueprint.Ruleset) == 0 {
		diags.Append(diagnostics.Builder().
			Error().
			Summary("No rules specified in blueprint.").
			Build(),
		)
		return diags
	}

	for rIdx, rule := range blueprint.Ruleset {
		if err := rlsVld.ruleValidator.Validate(&rule, blueprint.Sources, blueprint.References); err != nil {
			diags.Append(diagnostics.Builder().
				Error().
				Details(fmt.Sprintf("\nRule [%d]: %s", rIdx+1, err.Error())).
				Build(),
			)
		}
	}
	return diags
}
