// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package execution

import (
	"fmt"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/conformize/conformize/common/diagnostics"
	"github.com/conformize/conformize/common/ds"
	"github.com/conformize/conformize/internal/blueprint"
	"github.com/conformize/conformize/internal/blueprint/elements"
	"github.com/conformize/conformize/predicates/evaluation"
)

type BlueprintExecutor struct{}

func NewBlueprintExecutor() *BlueprintExecutor {
	return &BlueprintExecutor{}
}

func (blprntExec *BlueprintExecutor) Execute(blueprint *blueprint.Blueprint, diags diagnostics.Diagnosable) {
	blprntExecBuilder := NewBlueprintExecutionBuilder().
		WithBlueprint(blueprint).
		WithDiagnostics(diags)

	blueprintExecution := blprntExecBuilder.Build()
	if diags.HasErrors() {
		return
	}

	ruleSetEvaluation := evaluation.NewRuleSetEvaluation()
	results, ruleSetOk := ruleSetEvaluation.Evaluate(blueprintExecution.ruleSet)
	if ruleSetOk {
		return
	}

	availableCPUs := max(1, runtime.NumCPU()-1)
	cpus := runtime.GOMAXPROCS(availableCPUs)
	defer runtime.GOMAXPROCS(cpus)

	resultsCount := len(results)
	resultTasks := make(chan struct{}, min(availableCPUs, resultsCount))
	defer close(resultTasks)

	resultSignalChan := make(chan struct{}, resultsCount)
	defer close(resultSignalChan)

	resultDiags := make([]diagnostics.Diagnostic, resultsCount)
	processedResults := ds.NewBitSet(resultsCount)

	var wg sync.WaitGroup
	wg.Add(resultsCount)

	wg.Add(1)
	go func() {
		defer wg.Done()
		processResults(resultSignalChan, resultDiags, processedResults, diags, resultsCount)
	}()

	for resultIdx, result := range results {
		resultTasks <- struct{}{}
		go func() {
			defer wg.Done()
			processResult(resultIdx, &result, blueprintExecution.rulesMeta,
				resultDiags, processedResults, resultSignalChan, resultTasks)
		}()
	}

	wg.Wait()
}

func processResults(resultSignalChan <-chan struct{}, resultDiags []diagnostics.Diagnostic,
	processedResults *ds.BitSet, diags diagnostics.Diagnosable, resultsCount int) {
	nextResultIdx := 0
	for resultsCount > 0 {
		done, _ := processedResults.IsSet(nextResultIdx)
		for !done {
			<-resultSignalChan
			if done, _ = processedResults.IsSet(nextResultIdx); !done {
				time.Sleep(5 * time.Millisecond)
			}
		}

		if diag := resultDiags[nextResultIdx]; diag != nil {
			diags.Append(diag)
		}
		resultsCount--
		nextResultIdx++
	}
}

func processResult(resultIdx int, result *evaluation.RuleEvaluationResult, rulesMeta []*elements.RuleMeta,
	resultDiags []diagnostics.Diagnostic, processedResults *ds.BitSet,
	resultSignalChan chan<- struct{}, resultTasks <-chan struct{}) {
	if !result.OK {
		ruleMeta := rulesMeta[resultIdx]

		var errMsgBldr strings.Builder
		errMsgBldr.WriteString(
			fmt.Sprintf("\nRule %d not satisfied - predicate: %s, value path: %s, arguments:",
				resultIdx+1, ruleMeta.Predicate, ruleMeta.ValuePath),
		)

		for argIdx, argMeta := range ruleMeta.ArgumentsMeta {
			errMsgBldr.WriteString(fmt.Sprintf("\n[%d]:%s", argIdx+1, argMeta.String()))
		}
		resultDiags[resultIdx] = diagnostics.Builder().Error().Details(errMsgBldr.String()).Build()
	}
	processedResults.Set(resultIdx)
	resultSignalChan <- struct{}{}
	<-resultTasks
}
