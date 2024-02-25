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
)

type BlueprintExecutor struct{}

func (blprntExec *BlueprintExecutor) Execute(blueprint *blueprint.Blueprint, diags *diagnostics.Diagnostics) {
	blueprintExecutionCtx := BlueprintExecutionContext().
		WithBlueprint(blueprint).
		WithDiagnostics(diags)

	blueprintExecutionCtx.Execute()
}
