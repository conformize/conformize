// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package execution

import (
	"github.com/conformize/conformize/common/diagnostics"
	"github.com/conformize/conformize/internal/blueprint"
	"github.com/conformize/conformize/internal/blueprint/validation"
)

type BlueprintValidationStep struct {
	blueprint  *blueprint.Blueprint
	validation validation.BlueprintValidator
}

func (step *BlueprintValidationStep) Run() *diagnostics.Diagnostics {
	return step.validation.Validate(step.blueprint)
}

func NewBlueprintValidationStep(blueprint *blueprint.Blueprint, validator validation.BlueprintValidator) *BlueprintValidationStep {
	return &BlueprintValidationStep{
		blueprint:  blueprint,
		validation: validator,
	}
}
