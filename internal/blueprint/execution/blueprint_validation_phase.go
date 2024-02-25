// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package execution

import (
	"github.com/conformize/conformize/internal/blueprint"
	"github.com/conformize/conformize/internal/blueprint/validation"
)

type BlueprintValidationPhase struct {
	validationSteps []ExecutionStep
}

func (phase *BlueprintValidationPhase) Execute(blprntExecCtx *BlueprintExecutionContext) {
	for _, step := range phase.validationSteps {
		step.Run(blprntExecCtx)
	}
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
