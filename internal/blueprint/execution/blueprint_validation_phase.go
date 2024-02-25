// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package execution

import (
	"github.com/conformize/conformize/common/diagnostics"
	"github.com/conformize/conformize/internal/blueprint"
	"github.com/conformize/conformize/internal/blueprint/validation"
)

type BlueprintValidationPhase struct {
	validationSteps []ExecutionStep
}

func (phase *BlueprintValidationPhase) Execute() *diagnostics.Diagnostics {
	var diags *diagnostics.Diagnostics
	for _, step := range phase.validationSteps {
		stepDiags := step.Run()
		if stepDiags != nil {
			diags.Append(stepDiags.Entries()...)
		}
	}
	return diags
}

func NewBlueprintValidationPhase(blueprint *blueprint.Blueprint) *BlueprintValidationPhase {
	return &BlueprintValidationPhase{
		validationSteps: []ExecutionStep{
			&BlueprintValidationStep{
				blueprint:  blueprint,
				validation: &validation.BlueprintVersionValidator{},
			},
			&BlueprintValidationStep{
				blueprint:  blueprint,
				validation: &validation.BlueprintSourcesValidator{},
			},
			&BlueprintValidationStep{
				blueprint:  blueprint,
				validation: &validation.BlueprintReferencesValidator{},
			},
			&BlueprintValidationStep{
				blueprint:  blueprint,
				validation: &validation.BlueprintRulesValidator{},
			},
		},
	}
}
