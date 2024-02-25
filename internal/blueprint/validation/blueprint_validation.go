// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package validation

import (
	"github.com/conformize/conformize/common/diagnostics"
	"github.com/conformize/conformize/internal/blueprint"
)

var validators = []BlueprintValidator{
	&BlueprintVersionValidator{},
	&BlueprintSourcesValidator{},
	&BlueprintReferencesValidator{},
	&BlueprintRulesValidator{},
}

type BlueprintValidation struct{}

func (blprntVld *BlueprintValidation) Validate(blueprint *blueprint.Blueprint) *diagnostics.Diagnostics {
	validationDiags := diagnostics.NewDiagnostics()
	if blueprint == nil {
		validationDiags.Append(diagnostics.Builder().Error().Summary("Blueprint not specified!"))
	} else {
		for _, validator := range validators {
			if validatorDiags := validator.Validate(blueprint); validatorDiags.HasErrors() {
				validationDiags.Append(validatorDiags.Entries()...)
			}
		}
	}
	return validationDiags
}
