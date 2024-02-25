// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package execution

import (
	"github.com/conformize/conformize/internal/blueprint"
	"github.com/conformize/conformize/internal/providers"
)

type BlueprintExecutionPlanner struct{}

func (blprntExecPlan *BlueprintExecutionPlanner) Plan(blueprint *blueprint.Blueprint) (*BlueprintExecutionPlan, error) {
	plan := NewBlueprintExecutionPlan()
	plan.AddPhase(NewBlueprintValidationPhase(blueprint))

	providersInitializationPhase := NewProvidersInitializationPhase()
	providersConfigurationPhase := NewBlueprintProvidersConfigurationPhase()
	readSourcesPhase := NewBlueprintReadSourcesPhase()

	for sourceAlias, sourceConfig := range blueprint.Sources {
		providersInitializationPhase.AddStep(NewProviderInitializationStep(sourceAlias, &sourceConfig))
		providersConfigurationPhase.AddStep(NewProviderConfigurationStep(sourceAlias, &sourceConfig))

		if providers.ProviderName(sourceConfig.Provider) == providers.Aggregate {
			continue
		}

		readSourcesPhase.AddStep(NewReadSourceExecutionStep(sourceAlias, &sourceConfig))
	}

	plan.AddPhase(providersInitializationPhase)
	plan.AddPhase(providersConfigurationPhase)
	plan.AddPhase(readSourcesPhase)
	plan.AddPhase(NewBlueprintReferencesResolutionPhase(&blueprint.References))
	plan.AddPhase(NewBlueprintProvidersAggregationPhase())
	plan.AddPhase(NewBlueprintRulesetEvaluationPhase(&blueprint.Ruleset))

	return plan, nil
}

func NewBlueprintExecutionPlanner() *BlueprintExecutionPlanner {
	return &BlueprintExecutionPlanner{}
}
