// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package execution

import (
	"sync"
	"sync/atomic"
)

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
	var remainingSteps atomic.Int32
	var signal *sync.Cond = sync.NewCond(&sync.Mutex{})
	remainingSteps.Store(int32(len(phase.steps)))

	for _, step := range phase.steps {
		go func() {
			step.Run(blprntExecCtx)
			remainingSteps.Add(-1)

			signal.L.Lock()
			signal.Signal()
			signal.L.Unlock()
		}()
	}

	signal.L.Lock()
	for remainingSteps.Load() > 0 {
		signal.Wait()
	}
	signal.L.Unlock()
}
