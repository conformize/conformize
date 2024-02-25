// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package execution

import (
	"fmt"

	"github.com/conformize/conformize/common/diagnostics"
	"github.com/conformize/conformize/common/ds"
	"github.com/conformize/conformize/internal/blueprint"
	sdk "github.com/conformize/conformize/internal/providers/api"
	"github.com/conformize/conformize/internal/valuereferencesstore"
)

type BlueprintExecutor struct{}

func (blprntExec *BlueprintExecutor) Execute(blueprint *blueprint.Blueprint, diags *diagnostics.Diagnostics) {
	blprntExecCtx := &BlueprintExecutionContext{
		diags:                      diags,
		providersRegistry:          sdk.NewConfiguredProvidersRegistry(),
		providersDependenciesGraph: ds.NewDependencyGraph[string](),
		valueReferencesStore:       valuereferencesstore.NewValueReferencesStore(),
	}

	blprntPlanner := &BlueprintExecutionPlanner{}
	blprntPlan, err := blprntPlanner.Plan(blueprint)
	if err != nil {
		diags.Append(diagnostics.
			Builder().
			Summary(fmt.Sprintf("Failed to prepare blueprint execution plan, reason: %s", err)).
			Build(),
		)
		return
	}

	blprntPlan.Execute(blprntExecCtx)
}
