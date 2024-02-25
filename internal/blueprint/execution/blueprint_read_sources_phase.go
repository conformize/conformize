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

	"github.com/conformize/conformize/common/diagnostics"
)

type BlueprintReadSourcesPhase struct {
	steps []ExecutionStep
}

func (phase *BlueprintReadSourcesPhase) Execute() *diagnostics.Diagnostics {
	diags := diagnostics.NewDiagnostics()

	var remainingSteps atomic.Int32
	var signal *sync.Cond = sync.NewCond(&sync.Mutex{})
	remainingSteps.Store(int32(len(phase.steps)))

	for _, currentStep := range phase.steps {
		go func(step ExecutionStep) {
			signal.L.Lock()

			stepDiags := step.Run()
			if stepDiags != nil {
				diags.Append(stepDiags.Entries()...)
			}
			remainingSteps.Add(-1)

			signal.Signal()
			signal.L.Unlock()
		}(currentStep)
	}

	signal.L.Lock()
	for remainingSteps.Load() > 0 {
		signal.Wait()
	}
	signal.L.Unlock()

	return diags
}

func (phase *BlueprintReadSourcesPhase) AddStep(step ExecutionStep) {
	phase.steps = append(phase.steps, step)
}

func NewBlueprintReadSourcesPhase() *BlueprintReadSourcesPhase {
	return &BlueprintReadSourcesPhase{
		steps: make([]ExecutionStep, 0, 10),
	}
}
