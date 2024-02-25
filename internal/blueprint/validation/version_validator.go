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

type BlueprintVersionValidator struct{}

func (blprntVerVld *BlueprintVersionValidator) Validate(blueprint *blueprint.Blueprint) *diagnostics.Diagnostics {
	diags := diagnostics.NewDiagnostics()
	if blueprint.Version == 0 {
		diags.Append(diagnostics.Builder().
			Error().
			Summary("Version is not set").
			Build(),
		)
	}
	return diags
}
