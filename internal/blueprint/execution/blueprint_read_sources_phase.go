// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package execution

import (
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/conformize/conformize/common/diagnostics"
)

type BlueprintReadSourcesPhase struct {
	steps map[string]ExecutionStep
}

func (phase *BlueprintReadSourcesPhase) AddStep(step ExecutionStep) {
	if readSourceStep, ok := step.(*ReadSourceExecutionStep); ok {
		phase.steps[readSourceStep.SourceAlias()] = step
	}
}

func (phase *BlueprintReadSourcesPhase) Execute(blprntExecCtx *BlueprintExecutionContext) {
	blprntExecCtx.providersDependencyGraph.Run()
	sourcesDepsOrder := blprntExecCtx.providersDependencyGraph.GetOrder()
	if blprntExecCtx.providersDependencyGraph.HasCycles() {
		blprntExecCtx.diags.Append(diagnostics.Builder().
			Error().
			Summary("Blueprint execution failed: providers dependency graph has cycles.").
			Build())
		return
	}

	var remainingSteps atomic.Int32
	var signal *sync.Cond = sync.NewCond(&sync.Mutex{})
	remainingSteps.Store(int32(len(phase.steps)))

	for _, sourceAlias := range sourcesDepsOrder {
		if sourceAlias == "" {
			continue
		}

		currentStep, exists := phase.steps[sourceAlias]
		if !exists {
			blprntExecCtx.diags.Append(diagnostics.Builder().
				Error().
				Summary(fmt.Sprintf("Blueprint execution failed: no read source step found for source alias '%s'.", sourceAlias)).
				Build())
			continue
		}

		go func(step ExecutionStep) {
			signal.L.Lock()

			step.Run(blprntExecCtx)
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
}

func NewBlueprintReadSourcesPhase() *BlueprintReadSourcesPhase {
	return &BlueprintReadSourcesPhase{
		steps: make(map[string]ExecutionStep),
	}
}
