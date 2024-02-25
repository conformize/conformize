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
	"github.com/conformize/conformize/internal/providers"
)

type ConfigureProviderExecutionStep struct {
	alias        string
	sourceConfig *elements.ConfigurationSource
}

func NewConfigureProviderExecutionStep(alias string, sourceConfig *elements.ConfigurationSource) *ConfigureProviderExecutionStep {
	return &ConfigureProviderExecutionStep{
		alias:        alias,
		sourceConfig: sourceConfig,
	}
}

func (step *ConfigureProviderExecutionStep) Run() *diagnostics.Diagnostics {
	diags := diagnostics.NewDiagnostics()

	providerFactory := providers.ProviderFactory()
	formatter := format.Formatter()
	provider, err := providerFactory.Provider(step.sourceConfig.Provider)
	if err == nil {
		providerConfigurer := blueprint.ProviderConfigurer()
		if err = providerConfigurer.Configure(provider, step.sourceConfig); err != nil {
			line := formatter.
				Detail(format.Failure).
				Color(colors.Red).
				Dimmed().
				Format(fmt.Sprintf(" %-12s %-10s", step.alias, fmt.Sprintf("[%s]", step.sourceConfig.Provider)))

			line += formatter.Dimmed().Color(colors.Red).Format(fmt.Sprintf("error: %s", err.Error()))
			diags.Append(diagnostics.Builder().Error().Summary(line).Build())
			return diags
		}

		configuredProvidersRegistry := ConfiguredProvidersRegistry()
		err = configuredProvidersRegistry.Register(step.alias, provider)
		if err != nil {
			line := formatter.
				Detail(format.Failure).
				Color(colors.Red).
				Dimmed().
				Format(fmt.Sprintf(" %-12s %-10s", step.alias, fmt.Sprintf("[%s]", step.sourceConfig.Provider)))

			line += formatter.Dimmed().Color(colors.Red).Format(fmt.Sprintf("error: %s", err.Error()))
			diags.Append(diagnostics.Builder().Error().Summary(line).Build())
			return diags
		}

		return diags
	}

	line := formatter.
		Detail(format.Failure).
		Color(colors.Red).
		Dimmed().
		Format(fmt.Sprintf(" %-12s %-10s", step.alias, fmt.Sprintf("[%s]", step.sourceConfig.Provider)))

	line += formatter.Dimmed().Color(colors.Red).Format(fmt.Sprintf("error: %s", err.Error()))
	diags.Append(diagnostics.Builder().Error().Summary(line).Build())
	return diags
}
