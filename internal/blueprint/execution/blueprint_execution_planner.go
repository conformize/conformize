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
	"github.com/conformize/conformize/common/typed"
	"github.com/conformize/conformize/internal/blueprint"

	"github.com/conformize/conformize/internal/providers"
	"github.com/conformize/conformize/internal/providers/api/schema"
)

type BlueprintExecutionPlanner struct{}

func (blprntExecPlan *BlueprintExecutionPlanner) Plan(diags *diagnostics.Diagnostics, blueprint *blueprint.Blueprint) (*BlueprintExecutionPlan, error) {
	plan := NewBlueprintExecutionPlan()

	plan.AddPhase(NewBlueprintValidationPhase(blueprint))
	providersConfigurationPhase := NewBlueprintProvidersConfigurationPhase()

	sourcesDependencyGraph := ds.NewDependencyGraph[string]()
	for alias, source := range blueprint.Sources {
		providersConfigurationPhase.AddStep(NewConfigureProviderExecutionStep(alias, &source))
		if providers.ProviderName(source.Provider) != providers.Aggregate {
			sourcesDependencyGraph.AddEdge(alias, "")
			continue
		}

		prvdr, err := providers.ProviderFactory().Provider(source.Provider)
		if err != nil {
			return nil, err
		}

		var sourceConfig = schema.NewData(prvdr.ConfigurationSchema())
		if err = sourceConfig.Set(&source.Config); err != nil {
			return nil, err
		}

		var sourcesVal typed.Valuable
		sourcesVal, err = sourceConfig.GetAtPath("sources")
		if err != nil {
			return nil, err
		}

		var aggregateSources []string
		err = sourcesVal.As(&aggregateSources)
		if err != nil {
			return nil, err
		}

		for _, aggregateSource := range aggregateSources {
			sourcesDependencyGraph.AddEdge(alias, aggregateSource)
		}
	}
	sourcesDependencyGraph.Run()

	providersReadDataPhase := NewBlueprintReadSourcesPhase()
	for _, alias := range sourcesDependencyGraph.GetOrder() {
		if alias == "" {
			continue
		}

		sourceConfig, exists := blueprint.Sources[alias]
		if !exists {
			return nil, fmt.Errorf("source with alias '%s' not found", alias)
		}

		providersReadDataPhase.AddStep(NewReadSourceExecutionStep(alias, &sourceConfig))
	}

	plan.AddPhase(providersConfigurationPhase)
	plan.AddPhase(providersReadDataPhase)
	plan.AddPhase(NewBlueprintReferencesResolutionPhase(blueprint.References))
	plan.AddPhase(NewBlueprintRulesetEvaluationPhase(diags, &blueprint.Ruleset))

	return plan, nil
}

func NewBlueprintExecutionPlanner() *BlueprintExecutionPlanner {
	return &BlueprintExecutionPlanner{}
}
