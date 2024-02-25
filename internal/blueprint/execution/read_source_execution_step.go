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
	"github.com/conformize/conformize/internal/blueprint/elements"
	api "github.com/conformize/conformize/internal/providers/api"
	"github.com/mitchellh/reflectwalk"
)

type ReadSourceExecutionStep struct {
	alias  string
	config *elements.ConfigurationSource
}

func NewReadSourceExecutionStep(alias string, config *elements.ConfigurationSource) *ReadSourceExecutionStep {
	return &ReadSourceExecutionStep{
		alias:  alias,
		config: config,
	}
}

func (step *ReadSourceExecutionStep) SourceAlias() string {
	return step.alias
}

func (step *ReadSourceExecutionStep) Provider() string {
	return step.config.Provider
}

func (step *ReadSourceExecutionStep) Run(blprntExecCtx *BlueprintExecutionContext) {
	formatter := format.Formatter()

	confiredProvidersRegistry := blprntExecCtx.providersRegistry
	provider, exists := confiredProvidersRegistry.Get(step.alias)
	if !exists {
		line := formatter.
			Detail(format.Failure).
			Color(colors.Red).
			Dimmed().
			Format(fmt.Sprintf(" %-12s %-10s", step.alias, fmt.Sprintf("[ %s ]", step.config.Provider)))

		line += formatter.Color(colors.Red).Format("error: provider not configured")
		blprntExecCtx.diags.Append(diagnostics.Builder().Error().Summary(line).Build())
		return
	}

	var provisionDataReq *api.ProviderDataRequest

	var queryOptions = step.config.QueryOptions
	if queryOptions != nil {
		provider, found := blprntExecCtx.providersRegistry.Get(step.alias)
		if !found {
			line := formatter.
				Detail(format.Failure).
				Color(colors.Red).
				Dimmed().
				Format(fmt.Sprintf(" %-12s %-10s", step.alias, fmt.Sprintf("[ %s ]", step.config.Provider)))

			line += formatter.Color(colors.Red).Format("error: provider not found")
			blprntExecCtx.diags.Append(diagnostics.Builder().Error().Summary(line).Build())
			return
		}

		provisionDataReqSchema := provider.ProvisionDataRequestSchema()
		provisionDataReq = api.NewProviderDataRequest(provisionDataReqSchema)

		if err := reflectwalk.Walk(&queryOptions, &mapValueWalker{}); err != nil {
			blprntExecCtx.diags.
				Append(diagnostics.Builder().
					Error().
					Details(
						fmt.Sprintf("\nCouldn't set query options for source '%s' with '%s' provider, reason:\n%s",
							step.alias, step.config.Provider, err.Error()),
					).
					Build())
			return
		}

		if err := provisionDataReq.Set(queryOptions); err != nil {
			blprntExecCtx.diags.
				Append(diagnostics.Builder().
					Error().
					Details(
						fmt.Sprintf("\nCouldn't set query options for source '%s' with '%s' provider, reason:\n%s",
							step.alias, step.config.Provider, err.Error()),
					).
					Build())
			return
		}

	}

	valRefStore := blprntExecCtx.valueReferencesStore
	providerData, providerDiags := provider.Provide(provisionDataReq)
	if !providerDiags.HasErrors() {
		valRefStore.AddReference(step.alias, providerData)
		line := formatter.
			Color(colors.Green).
			Detail(format.Item).
			Format("")

		line += formatter.
			Bold().
			Format(fmt.Sprintf(" %-12s %-10s", step.alias, fmt.Sprintf("[ %s ]", step.config.Provider)))

		blprntExecCtx.diags.Append(diagnostics.Builder().Info().Summary(line).Build())
		return
	}

	line := formatter.
		Detail(format.Failure).
		Color(colors.Red).
		Dimmed().
		Format(fmt.Sprintf(" %-12s %-10s", step.alias, fmt.Sprintf("[ %s ]", step.config.Provider)))

	line += formatter.Color(colors.Red).Format(fmt.Sprintf("error: %s", providerDiags.Entries().String()))
	blprntExecCtx.diags.Append(diagnostics.Builder().Error().Summary(line).Build())
}
