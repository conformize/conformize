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
	"github.com/conformize/conformize/common/format"
	"github.com/conformize/conformize/common/format/colors"
	"github.com/conformize/conformize/internal/blueprint"
	"github.com/conformize/conformize/internal/blueprint/elements"
	api "github.com/conformize/conformize/internal/providers/api"
	"github.com/conformize/conformize/internal/valuereferencesstore"
)

type ReadSourceExecutionStep struct {
	alias        string
	sourceConfig *elements.ConfigurationSource
}

func NewReadSourceExecutionStep(alias string, sourceConfig *elements.ConfigurationSource) *ReadSourceExecutionStep {
	return &ReadSourceExecutionStep{
		alias:        alias,
		sourceConfig: sourceConfig,
	}
}

func (step *ReadSourceExecutionStep) Run() *diagnostics.Diagnostics {
	diags := diagnostics.NewDiagnostics()

	formatter := format.Formatter()

	confiredProvidersRegistry := ConfiguredProvidersRegistry()
	provider, exists := confiredProvidersRegistry.Get(step.alias)
	if !exists {
		line := formatter.
			Detail(format.Failure).
			Color(colors.Red).
			Dimmed().
			Format(fmt.Sprintf(" %-12s %-10s", step.alias, fmt.Sprintf("[%s]", step.sourceConfig.Provider)))

		line += formatter.Dimmed().Color(colors.Red).Format("error: provider not configured")
		diags.Append(diagnostics.Builder().Error().Summary(line).Build())
		return diags
	}

	providerConfigurer := blueprint.ProviderConfigurer()

	var err error
	var provisionDataReq *api.ProviderDataRequest
	if step.sourceConfig.QueryOptions != nil {
		if provisionDataReq, err = providerConfigurer.Query(provider, step.sourceConfig.QueryOptions); err != nil {
			diags.Append(diagnostics.Builder().
				Error().
				Details(
					fmt.Sprintf("\nCouldn't set query options for source '%s' with '%s' provider, reason:\n%s",
						step.alias, step.sourceConfig.Provider, err.Error()),
				).
				Build())
			return diags
		}
	}

	valRefStore := valuereferencesstore.Instance()
	providerData, providerDiags := provider.Provide(provisionDataReq)
	if !providerDiags.HasErrors() {
		valRefStore.AddReference(step.alias, providerData)
		line := formatter.
			Color(colors.Green).
			Detail(format.Item).
			Format("")

		line += formatter.
			Bold().
			Format(fmt.Sprintf(" %-12s %-10s", step.alias, fmt.Sprintf("[%s]", step.sourceConfig.Provider)))

		diags.Append(diagnostics.Builder().Info().Summary(line).Build())
		return diags
	}

	line := formatter.
		Detail(format.Failure).
		Color(colors.Red).
		Dimmed().
		Format(fmt.Sprintf(" %-12s %-10s", step.alias, fmt.Sprintf("[%s]", step.sourceConfig.Provider)))

	line += formatter.Dimmed().Color(colors.Red).Format(fmt.Sprintf("error: %s", providerDiags.Entries().String()))
	diags.Append(diagnostics.Builder().Error().Summary(line).Build())
	return diags

}
