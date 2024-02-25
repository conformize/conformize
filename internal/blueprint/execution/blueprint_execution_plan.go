// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package execution

import "github.com/conformize/conformize/common/diagnostics"

type BlueprintExecutionPlan struct {
	phases []ExecutionPhase
}

func (plan *BlueprintExecutionPlan) AddPhase(phase ExecutionPhase) {
	plan.phases = append(plan.phases, phase)
}

func (plan *BlueprintExecutionPlan) Execute() *diagnostics.Diagnostics {
	diags := diagnostics.NewDiagnostics()
	for _, phase := range plan.phases {
		phaseDiags := phase.Execute()
		if phaseDiags != nil {
			diags.Append(phaseDiags.Entries()...)
		}
	}
	return diags
}

func NewBlueprintExecutionPlan() *BlueprintExecutionPlan {
	return &BlueprintExecutionPlan{
		phases: make([]ExecutionPhase, 0, 10),
	}
}
