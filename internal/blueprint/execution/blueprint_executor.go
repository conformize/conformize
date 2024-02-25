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

	"github.com/conformize/conformize/common/diagnostics"
	"github.com/conformize/conformize/common/ds"
	"github.com/conformize/conformize/common/format"
	"github.com/conformize/conformize/common/format/colors"
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
	diags.Append(
		diagnostics.
			Builder().
			Info().
			Summary(
				fmt.Sprintf(
					"%s\n",
					format.Formatter().
						Dimmed().
						Color(colors.Blue).
						Detail(format.TestTube).
						Format("Evaluating ruleset..."),
				),
			),
	)

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
		ruleErrMsg := ruleNotSatisfiedErrorMessage(rulesMeta[resultIdx], resultIdx)
		diags[resultIdx] = diagnostics.Builder().Error().Details(ruleErrMsg).Build()
	}
	processed.Set(resultIdx)
}

const bulletIndent = "   "
const labelWidth = -14

func ruleNotSatisfiedErrorMessage(ruleMeta *elements.RuleMeta, ruleIdx int) string {
	var bldr strings.Builder

	var ruleHeader string
	if len(ruleMeta.Name) > 0 {
		ruleHeader = fmt.Sprintf("Rule '%s':\n\n", ruleMeta.Name)
	} else {
		ruleHeader = fmt.Sprintf("Rule %d:\n\n", ruleIdx+1)
	}

	formatter := format.Formatter()
	bldr.WriteString(formatter.
		Detail(format.Failure).
		Color(colors.Red).
		Bold().
		Format(ruleHeader))

	writeField := func(label, value string) {
		line := fmt.Sprintf("%-*s%s\n", labelWidth, label+":", value)
		bldr.WriteString(
			bulletIndent +
				formatter.
					Detail(format.Bullet).
					Dimmed().
					Color(colors.Red).
					Format(line),
		)
	}

	writeField("$value", ruleMeta.ValuePath)
	writeField("Predicate", ruleMeta.Predicate)
	writeField("Argument(s)", ruleMeta.ArgumentsMeta.String())

	return bldr.String()
}
