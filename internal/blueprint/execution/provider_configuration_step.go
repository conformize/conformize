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
)

type ProviderConfigurationStep struct {
	alias  string
	config *elements.ConfigurationSource
}

func NewProviderConfigurationStep(alias string, sourceConfig *elements.ConfigurationSource) *ProviderConfigurationStep {
	return &ProviderConfigurationStep{
		alias:  alias,
		config: sourceConfig,
	}
}

func (step *ProviderConfigurationStep) Run(blprntExecCtx *BlueprintExecutionContext) {
	formatter := format.Formatter()

	var err error
	providerConfigurer := ProviderConfigurer()

	prvdrConfigCtx := &ProviderConfigurationContext{
		Alias:                    step.alias,
		Config:                   step.config,
		ProvidersRegistry:        blprntExecCtx.providersRegistry,
		ProvidersDependencyGraph: blprntExecCtx.providersDependenciesGraph,
	}

	_, registered := blprntExecCtx.providersRegistry.Get(step.alias)
	if !registered {
		line := formatter.
			Detail(format.Failure).
			Color(colors.Red).
			Dimmed().
			Format(fmt.Sprintf(" %-12s %-10s", step.alias, fmt.Sprintf("[ %s ]", step.config.Provider)))

		line += formatter.Color(colors.Red).Format("error: provider not initialized")
		blprntExecCtx.diags.Append(diagnostics.Builder().Error().Summary(line).Build())
		return
	}

	if err = providerConfigurer.Configure(prvdrConfigCtx); err != nil {
		line := formatter.
			Detail(format.Failure).
			Color(colors.Red).
			Dimmed().
			Format(fmt.Sprintf(" %-12s %-10s", step.alias, fmt.Sprintf("[ %s ]", step.config.Provider)))

		line += formatter.Color(colors.Red).Format(fmt.Sprintf("error: %s", err.Error()))
		blprntExecCtx.diags.Append(diagnostics.Builder().Error().Summary(line).Build())
	}
}
