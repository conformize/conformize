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

type ConfigureProviderExecutionStep struct {
	alias  string
	config *elements.ConfigurationSource
}

func NewConfigureProviderExecutionStep(alias string, sourceConfig *elements.ConfigurationSource) *ConfigureProviderExecutionStep {
	return &ConfigureProviderExecutionStep{
		alias:  alias,
		config: sourceConfig,
	}
}

func (step *ConfigureProviderExecutionStep) Run(blprntExecCtx *BlueprintExecutionContext) {
	formatter := format.Formatter()

	var err error
	providerConfigurer := ProviderConfigurer()

	prvdrConfigCtx := &ProviderConfigurationContext{
		Alias:                      step.alias,
		Config:                     step.config,
		ValueReferencesStore:       blprntExecCtx.valueReferencesStore,
		ProvidersRegistry:          blprntExecCtx.providersRegistry,
		ProvidersDependenciesGraph: blprntExecCtx.providersDependencyGraph,
	}

	if err = providerConfigurer.Configure(prvdrConfigCtx); err != nil {
		line := formatter.
			Detail(format.Failure).
			Color(colors.Red).
			Dimmed().
			Format(fmt.Sprintf(" %-12s %-10s", step.alias, fmt.Sprintf("[%s]", step.config.Provider)))

		line += formatter.Dimmed().Color(colors.Red).Format(fmt.Sprintf("error: %s", err.Error()))
		blprntExecCtx.diags.Append(diagnostics.Builder().Error().Summary(line).Build())
	}
}
