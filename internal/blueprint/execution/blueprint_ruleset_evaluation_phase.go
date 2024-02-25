// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package execution

import (
	"fmt"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/conformize/conformize/common"
	"github.com/conformize/conformize/common/diagnostics"
	"github.com/conformize/conformize/common/format"
	"github.com/conformize/conformize/common/format/colors"
	"github.com/conformize/conformize/common/functions"
	"github.com/conformize/conformize/common/path"
	"github.com/conformize/conformize/common/typed"
	"github.com/conformize/conformize/internal/blueprint/elements"
	"github.com/conformize/conformize/predicates/condition"
	"github.com/conformize/conformize/predicates/predicatefactory"
)

type ruleEvaluationResult struct {
	ok       bool
	ruleMeta *elements.RuleMeta
}

type BlueprintRulesetEvaluationPhase struct {
	ruleset *[]elements.Rule
}

func NewBlueprintRulesetEvaluationPhase(ruleset *[]elements.Rule) *BlueprintRulesetEvaluationPhase {
	return &BlueprintRulesetEvaluationPhase{
		ruleset: ruleset,
	}
}

func (phase *BlueprintRulesetEvaluationPhase) Execute(blprntExecCtx *BlueprintExecutionContext) {
	var ruleViolationsCount int32 = 0

	blprntExecCtx.diags.Append(
		diagnostics.Builder().Info().
			Summary(format.Formatter().Detail(format.Pencil).Color(colors.Blue).
				Format("Evaluating ruleset...")).Build(),
	)

	var remainingEvaluations atomic.Int32
	var signal *sync.Cond = sync.NewCond(&sync.Mutex{})
	remainingEvaluations.Store(int32(len(*phase.ruleset)))

	evaluationResults := make([]ruleEvaluationResult, len(*phase.ruleset))
	for idx, rule := range *phase.ruleset {
		go func() {
			ok, ruleMeta := evaluateRule(blprntExecCtx, idx, &rule)
			evaluationResults[idx] = ruleEvaluationResult{
				ok:       ok,
				ruleMeta: ruleMeta,
			}
			remainingEvaluations.Add(-1)

			signal.L.Lock()
			signal.Signal()
			signal.L.Unlock()
		}()
	}

	signal.L.Lock()
	for remainingEvaluations.Load() > 0 {
		signal.Wait()
	}
	signal.L.Unlock()

	for idx, evalResult := range evaluationResults {
		if evalResult.ruleMeta != nil {
			if evalResult.ruleMeta != nil && evalResult.ruleMeta.Diagnostics.HasErrors() {
				blprntExecCtx.diags.Append(evalResult.ruleMeta.Diagnostics.Entries()...)
				continue
			}

			if !evalResult.ok {
				ruleViolationsCount++
				blprntExecCtx.diags.Append(diagnostics.Builder().
					Error().
					Details(ruleViolationErrorMessage(evalResult.ruleMeta, idx)).
					Build(),
				)
			}
		}
	}

	if ruleViolationsCount == 0 {
		return
	}

	var errMsg string
	if ruleViolationsCount == 1 {
		errMsg = "1 rule assertion failed."
	} else {
		errMsg = fmt.Sprintf("%d rule assertions failed.", ruleViolationsCount)
	}
	blprntExecCtx.diags.Append(
		diagnostics.Builder().
			Error().
			Summary(
				format.Formatter().
					Color(colors.Red).
					Bold().
					Detail(format.FailureWarning).
					Format(errMsg),
			).Build(),
	)
}

func evaluateRule(blprntExecCtx *BlueprintExecutionContext, rIdx int, r *elements.Rule) (bool, *elements.RuleMeta) {
	ruleIdx := rIdx + 1
	valuePathSteps := r.Value.Steps()
	predicateCondition := condition.FromString(r.Predicate)
	predicate, err := predicatefactory.Instance().Build(predicateCondition)
	diags := diagnostics.NewDiagnostics()

	ruleMeta := &elements.RuleMeta{
		Name:        r.Name,
		ValuePath:   r.Value.String(),
		Predicate:   r.Predicate,
		Diagnostics: diags,
	}

	if err != nil {
		diags.
			Append(diagnostics.Builder().Error().
				Summary(fmt.Sprintf("\nRule %d - couldn't build predicate '%s', reason:", ruleIdx, r.Predicate)).
				Details(err.Error()).
				Build(),
			)
		return false, ruleMeta
	}

	var argVal typed.Valuable
	var argMeta *elements.ArgumentMeta
	var valRefStore = blprntExecCtx.valueReferencesStore
	if r.Arguments != nil {
		argMeta = &elements.ArgumentMeta{Sensitive: r.Arguments.IsSensitive()}
		ruleMeta.ArgumentsMeta = argMeta
		switch arg := r.Arguments.(type) {
		case *elements.RawValue:
			if arg.Value != nil {
				argMeta.Value = arg.Value
				argVal, err = functions.ParseRawValue(arg.Value)
				if err != nil {
					diags.
						Append(diagnostics.Builder().Error().
							Summary(fmt.Sprintf("\nRule %d - couldn't parse argument value %v, reason:", ruleIdx, arg.Value)).
							Details(err.Error()).
							Build(),
						)
					ruleMeta.ArgumentsMeta = argMeta
					return false, ruleMeta
				}
			}
		case *elements.PathValue:
			valNode, err := valRefStore.GetAtPath(&arg.Path)
			if err != nil {
				diags.
					Append(diagnostics.Builder().Error().
						Summary(fmt.Sprintf("Rule %d - Couldn't resolve value path %s for argument", ruleIdx, arg.Path.String())).
						Details(err.Error()).
						Build(),
					)
				return false, ruleMeta
			}
			argMeta.Path = arg.Path.String()
			argVal, err = functions.ParseRawValue(valNode.Value)
			if err != nil {
				blprntExecCtx.diags.
					Append(diagnostics.Builder().Error().
						Details(fmt.Sprintf("\nRule %d - couldn't parse argument value at path %s, reason:\n%s", ruleIdx, arg.Path.String(), err.Error())).
						Build(),
					)
				return false, ruleMeta
			}
		}
	}

	valuePath := path.NewPath(valuePathSteps)
	valNode, err := valRefStore.GetAtPath(valuePath)
	if err != nil {
		diags.
			Append(diagnostics.Builder().Error().
				Summary(fmt.Sprintf("\nRule %d - couldn't resolve value path %s, reason:", ruleIdx, r.Value.String())).
				Details(err.Error()).
				Build(),
			)
		return false, ruleMeta
	}

	fnNode, ok := valNode.Value.(*common.IterFnNodeValue)
	if !ok {
		var val typed.Valuable
		val, err = functions.ParseRawValue(valNode.Value)
		if err != nil {
			diags.
				Append(diagnostics.Builder().Error().
					Details(fmt.Sprintf("\nRule %d - couldn't parse value at path %s, reason:\n%s", ruleIdx, r.Value.String(), err.Error())).
					Build(),
				)
			return false, ruleMeta
		}
		ok, err = predicate.Test(val, argVal)
		return ok, ruleMeta
	}
	ok, err = fnNode.Fn(fnNode.Iter, predicate, argVal)
	return ok, ruleMeta
}

const bulletIndent = " "
const labelWidth = 12

func ruleViolationErrorMessage(ruleMeta *elements.RuleMeta, ruleIdx int) string {
	var msgBldr strings.Builder
	msgBldr.Grow(256)

	var ruleHeader string
	if len(ruleMeta.Name) > 0 {
		ruleHeader = fmt.Sprintf("Rule '%s':\n\n", ruleMeta.Name)
	} else {
		ruleHeader = fmt.Sprintf("Rule %d:\n\n", ruleIdx+1)
	}

	msgBldr.WriteString(format.Formatter().
		Detail(format.Failure).
		Color(colors.Red).
		Bold().
		Format(ruleHeader))

	writeLine := func(label, value string) {
		line := fmt.Sprintf("%-*s%s\n", labelWidth, label+":", value)
		msgBldr.WriteString(
			strings.Repeat(" ", 4) +
				format.Formatter().
					Detail(format.Bullet).
					Color(colors.Red).
					Format(line),
		)
	}

	writeLine("$value", ruleMeta.ValuePath)
	writeLine("predicate", ruleMeta.Predicate)
	if ruleMeta.ArgumentsMeta.Value != nil {
		if args, ok := ruleMeta.ArgumentsMeta.Value.([]any); ok {
			writeLine("arguments", fmt.Sprintf("%v", args))
		} else {
			writeLine("argument", ruleMeta.ArgumentsMeta.String())
		}
	}

	return msgBldr.String()
}
