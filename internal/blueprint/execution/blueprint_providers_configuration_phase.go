// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package execution

type BlueprintProvidersConfigurationPhase struct {
	steps []ExecutionStep
}

func NewBlueprintProvidersConfigurationPhase() *BlueprintProvidersConfigurationPhase {
	return &BlueprintProvidersConfigurationPhase{
		steps: make([]ExecutionStep, 0, 10),
	}
}

func (phase *BlueprintProvidersConfigurationPhase) AddStep(step ExecutionStep) {
	phase.steps = append(phase.steps, step)
}

func (phase *BlueprintProvidersConfigurationPhase) Execute(blprntExecCtx *BlueprintExecutionContext) {
	for _, step := range phase.steps {
		step.Run(blprntExecCtx)
	}
}
