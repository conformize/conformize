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
	sdk "github.com/conformize/conformize/internal/providers/api"
)

type BlueprintProvidersAggregationPhase struct{}

func NewBlueprintProvidersAggregationPhase() *BlueprintProvidersAggregationPhase {
	return &BlueprintProvidersAggregationPhase{}
}

func (phase *BlueprintProvidersAggregationPhase) Execute(blprntExecCtx *BlueprintExecutionContext) {
	if blprntExecCtx.providersDependenciesGraph == nil || blprntExecCtx.providersDependenciesGraph.IsEmpty() {
		return
	}

	diags := diagnostics.NewDiagnostics()
	blprntExecCtx.providersDependenciesGraph.Run()
	if blprntExecCtx.providersDependenciesGraph.HasCycles() {
		cycles := blprntExecCtx.providersDependenciesGraph.GetCycles()
		for _, cycle := range cycles {
			ref := cycle[0]
			otherRef := cycle[1]
			blprntExecCtx.diags.Append(
				diagnostics.Builder().Error().
					Summary(fmt.Sprintf("Cyclic dependency detected between providers '%s' and '%s'", ref, otherRef)).
					Build(),
			)
		}
		return
	}

	depsOrdered := blprntExecCtx.providersDependenciesGraph.GetOrder()
	for _, providerAlias := range depsOrdered {
		prvdr, exists := blprntExecCtx.providersRegistry.Get(providerAlias)
		if !exists {
			diags.Append(
				diagnostics.Builder().Error().
					Summary(fmt.Sprintf("Provider '%s' not found in the configured providers registry.", providerAlias)).
					Build(),
			)
			continue
		}

		data, prvdDiags := prvdr.Provide(sdk.NewProviderDataRequest(prvdr.ProvisionDataRequestSchema()))
		blprntExecCtx.diags.Append(prvdDiags.Entries()...)
		if prvdDiags.HasErrors() {
			continue
		}

		blprntExecCtx.valueReferencesStore.AddReference(providerAlias, data)
	}

	if diags.HasErrors() {
		blprntExecCtx.diags.Append(diags.Entries()...)
	}
}
