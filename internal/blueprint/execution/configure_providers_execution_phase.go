// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package execution

type ConfigureProvidersExecutionPhase struct {
	steps []ExecutionStep
}

func NewConfigureProvidersExecutionPhase() *ConfigureProvidersExecutionPhase {
	return &ConfigureProvidersExecutionPhase{
		steps: make([]ExecutionStep, 0, 10),
	}
}

func (phase *ConfigureProvidersExecutionPhase) AddStep(step ExecutionStep) {
	phase.steps = append(phase.steps, step)
}

func (phase *ConfigureProvidersExecutionPhase) Execute(blprntExecCtx *BlueprintExecutionContext) {
	for _, step := range phase.steps {
		step.Run(blprntExecCtx)
	}
}
