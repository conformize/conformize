// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package evaluation

import (
	"runtime"
	"sync"
	"sync/atomic"

	"github.com/conformize/conformize/predicates/rule"
	"github.com/conformize/conformize/predicates/ruleset"
)

type RuleEvaluationResult struct {
	OK    bool
	Error error
}

type RuleSetEvaluation struct {
	ruleEvaluation *RuleEvaluation
}

func (ruleSetEval *RuleSetEvaluation) Evaluate(ruleSet *ruleset.RuleSet) ([]RuleEvaluationResult, bool) {
	results := make([]RuleEvaluationResult, len(ruleSet.Rules))
	ruleSetOk := atomic.Bool{}
	ruleSetOk.Store(true)

	availableCPUs := max(1, runtime.NumCPU()-1)
	cpus := runtime.GOMAXPROCS(availableCPUs)
	defer runtime.GOMAXPROCS(cpus)

	var wg sync.WaitGroup
	wg.Add(len(ruleSet.Rules))

	maxParallelTasksCount := min(availableCPUs, len(results))
	tasks := make(chan struct{}, maxParallelTasksCount)
	defer close(tasks)
	for idx, r := range ruleSet.Rules {
		tasks <- struct{}{}
		go func(rule *rule.Rule, idx int) {
			defer wg.Done()
			ruleOk, err := ruleSetEval.ruleEvaluation.Evaluate(rule)
			if !ruleOk {
				ruleSetOk.Store(false)
			}
			results[idx] = RuleEvaluationResult{OK: ruleOk, Error: err}
			<-tasks
		}(r, idx)
	}
	wg.Wait()
	return results, ruleSetOk.Load()
}

func NewRuleSetEvaluation() *RuleSetEvaluation {
	return &RuleSetEvaluation{NewRuleEvaluation()}
}
