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
	"strconv"
	"strings"
	"sync"

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
	blprntExecCtxBuilder := NewBlueprintExecutionContextBuilder().
		WithBlueprint(blueprint).
		WithDiagnostics(diags)

	blueprintExecutionCtx := blprntExecCtxBuilder.Build()
	if diags.HasErrors() {
		return
	}

	ruleSetEvaluation := evaluation.NewRuleSetEvaluation()
	results, ruleSetOk := ruleSetEvaluation.Evaluate(blueprintExecutionCtx.ruleSet)
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
	wg.Add(resultsCount + 1)

	go func() {
		defer wg.Done()
		processResults(resultSignalChan, resultDiags, processedResults, diags, resultsCount)
	}()

	for resultIdx, result := range results {
		resultTasks <- struct{}{}
		go func() {
			defer wg.Done()
			processResult(resultIdx, &result, blueprintExecutionCtx.rulesMeta, resultDiags, processedResults)
			<-resultTasks
			resultSignalChan <- struct{}{}
		}()
	}
	wg.Wait()
}

func processResults(resultSignalChan <-chan struct{}, resultDiags []diagnostics.Diagnostic,
	processedResults *ds.BitSet, diags diagnostics.Diagnosable, resultsCount int) {
	nextResultIdx := 0
	for nextResultIdx < resultsCount {
		done, _ := processedResults.IsSet(nextResultIdx)
		for !done {
			<-resultSignalChan
			done, _ = processedResults.IsSet(nextResultIdx)
		}

		if diag := resultDiags[nextResultIdx]; diag != nil {
			diags.Append(diag)
		}
		nextResultIdx++
	}
}

func processResult(resultIdx int, result *evaluation.RuleEvaluationResult, rulesMeta []*elements.RuleMeta, diags []diagnostics.Diagnostic, processed *ds.BitSet) {
	if !result.OK {
		ruleMeta := rulesMeta[resultIdx]

		var errMsgBldr strings.Builder
		errMsgBldr.WriteString(ruleNotSatisfiedErrorMessage(ruleMeta, resultIdx))

		for argIdx, argMeta := range ruleMeta.ArgumentsMeta {
			errMsgBldr.WriteString(fmt.Sprintf("\n[%d]:%s", argIdx+1, argMeta.String()))
		}
		diags[resultIdx] = diagnostics.Builder().Error().Details(errMsgBldr.String()).Build()
	}
	processed.Set(resultIdx)
}

func ruleNotSatisfiedErrorMessage(ruleMeta *elements.RuleMeta, ruleIdx int) string {
	var ruleDescription string
	if len(ruleMeta.Name) > 0 {
		ruleDescription = fmt.Sprintf("Rule \"%s\"", ruleMeta.Name)
	} else {
		ruleDescription = fmt.Sprintf("Rule %s", strconv.Itoa(ruleIdx))
	}
	return fmt.Sprintf("\n%s not satisfied - predicate: %s, value path: %s, arguments:",
		ruleDescription, ruleMeta.Predicate, ruleMeta.ValuePath)
}
