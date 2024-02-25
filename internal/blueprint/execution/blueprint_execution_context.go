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
	"github.com/conformize/conformize/internal/blueprint"
	"github.com/conformize/conformize/internal/blueprint/elements"
	"github.com/conformize/conformize/internal/blueprint/validation"
	"github.com/conformize/conformize/internal/providers"
	api "github.com/conformize/conformize/internal/providers/api"
	"github.com/conformize/conformize/internal/valuereferencesstore"
	"github.com/conformize/conformize/predicates/condition"
	"github.com/conformize/conformize/predicates/predicatefactory"
)

type blueprintExecutionContext struct {
	blueprint   *blueprint.Blueprint
	diags       *diagnostics.Diagnostics
	msgBldr     strings.Builder
	valRefStore *valuereferencesstore.ValueReferencesStore
}

func BlueprintExecutionContext() *blueprintExecutionContext {
	blprntExecCtx := &blueprintExecutionContext{
		msgBldr:     strings.Builder{},
		valRefStore: valuereferencesstore.Instance(),
	}
	return blprntExecCtx
}

func (blprntExecCtx *blueprintExecutionContext) WithBlueprint(blueprint *blueprint.Blueprint) *blueprintExecutionContext {
	blprntExecCtx.blueprint = blueprint
	return blprntExecCtx
}

func (blprntExecCtx *blueprintExecutionContext) WithDiagnostics(diags *diagnostics.Diagnostics) *blueprintExecutionContext {
	blprntExecCtx.diags = diags
	return blprntExecCtx
}

func (blprntExecCtx *blueprintExecutionContext) readSource(alias string, sourceConfig *elements.ConfigurationSource) {
	sourceDiags := diagnostics.NewDiagnostics()
	defer func() { blprntExecCtx.diags.Append(sourceDiags.Entries()...) }()

	providerFactory := providers.ProviderFactory()
	formatter := format.Formatter()
	provider, err := providerFactory.Provider(sourceConfig.Provider)
	if err == nil {
		providerConfigurer := blueprint.ProviderConfigurer()
		if err = providerConfigurer.Configure(provider, sourceConfig); err != nil {
			line := formatter.
				Detail(format.Failure).
				Color(colors.Red).
				Dimmed().
				Format(fmt.Sprintf(" %-12s %-10s", alias, fmt.Sprintf("[%s]", sourceConfig.Provider)))

			line += formatter.Dimmed().Color(colors.Red).Format(fmt.Sprintf("error: %s", err.Error()))
			sourceDiags.Append(diagnostics.Builder().Error().Summary(line).Build())
			return
		}

		var provisionDataReq *api.ProviderDataRequest
		if sourceConfig.QueryOptions != nil {
			if provisionDataReq, err = providerConfigurer.Query(provider, sourceConfig.QueryOptions); err != nil {
				sourceDiags.Append(diagnostics.Builder().
					Error().
					Details(
						fmt.Sprintf("\nCouldn't set query options for source '%s' with '%s' provider, reason:\n%s",
							alias, sourceConfig.Provider, err.Error()),
					).
					Build())
			}
		}

		providerData, providerDiags := provider.Provide(provisionDataReq)
		if !providerDiags.HasErrors() {
			blprntExecCtx.valRefStore.AddReference(alias, providerData)
			line := formatter.
				Color(colors.Green).
				Detail(format.Item).
				Format("")

			line += formatter.
				Bold().
				Format(fmt.Sprintf(" %-12s %-10s", alias, fmt.Sprintf("[%s]", sourceConfig.Provider)))

			sourceDiags.Append(diagnostics.Builder().Info().Summary(line).Build())
			return
		}

		line := formatter.
			Detail(format.Failure).
			Color(colors.Red).
			Dimmed().
			Format(fmt.Sprintf(" %-12s %-10s", alias, fmt.Sprintf("[%s]", sourceConfig.Provider)))

		line += formatter.Dimmed().Color(colors.Red).Format(fmt.Sprintf("error: %s", providerDiags.Entries().String()))
		sourceDiags.Append(diagnostics.Builder().Error().Summary(line).Build())
		return
	}

	line := formatter.
		Detail(format.Failure).
		Color(colors.Red).
		Dimmed().
		Format(fmt.Sprintf(" %-12s %-10s", alias, fmt.Sprintf("[%s]", sourceConfig.Provider)))

	line += formatter.Dimmed().Color(colors.Red).Format(fmt.Sprintf("error: %s", err.Error()))
	sourceDiags.Append(diagnostics.Builder().Error().Summary(line).Build())
}

func (blprntExecCtx *blueprintExecutionContext) executeRule(rIdx int, r *elements.Rule) (bool, *elements.RuleMeta, error) {
	ruleIdx := rIdx + 1
	valuePathSteps := r.Value.Steps()
	predicateCondition := condition.FromString(r.Predicate)
	predicate, err := predicatefactory.Instance().Build(predicateCondition)
	if err != nil {
		blprntExecCtx.diags.
			Append(diagnostics.Builder().Error().
				Summary(fmt.Sprintf("\nRule %d - couldn't build predicate '%s', reason:", ruleIdx, r.Predicate)).
				Details(err.Error()).
				Build(),
			)
		return false, nil, err
	}

	var argVal typed.Valuable
	var argMeta *elements.ArgumentMeta
	if r.Arguments != nil {
		argMeta = &elements.ArgumentMeta{Sensitive: r.Arguments.IsSensitive()}
		switch arg := r.Arguments.(type) {
		case *elements.RawValue:
			if arg.Value != nil {
				argMeta.Value = arg.Value
				argVal, err = functions.ParseRawValue(arg.Value)
				if err != nil {
					blprntExecCtx.diags.
						Append(diagnostics.Builder().Error().
							Summary(fmt.Sprintf("\nRule %d - couldn't parse argument value %v, reason:", ruleIdx, arg.Value)).
							Details(err.Error()).
							Build(),
						)
					return false, nil, err
				}
			}
		case *elements.PathValue:
			valNode, err := blprntExecCtx.valRefStore.GetAtPath(&arg.Path)
			if err != nil {
				blprntExecCtx.diags.
					Append(diagnostics.Builder().Error().
						Summary(fmt.Sprintf("Rule %d - Couldn't resolve value path %s for argument", ruleIdx, arg.Path.String())).
						Details(err.Error()).
						Build(),
					)
				return false, nil, err
			}
			argMeta.Path = arg.Path.String()
			argVal, err = functions.ParseRawValue(valNode.Value)
			if err != nil {
				blprntExecCtx.diags.
					Append(diagnostics.Builder().Error().
						Details(fmt.Sprintf("\nRule %d - couldn't parse argument value at path %s, reason:\n%s", ruleIdx, arg.Path.String(), err.Error())).
						Build(),
					)
				return false, nil, err
			}
		}
	}

	valuePath := path.NewPath(valuePathSteps)
	valNode, err := blprntExecCtx.valRefStore.GetAtPath(valuePath)
	if err != nil {
		blprntExecCtx.diags.
			Append(diagnostics.Builder().Error().
				Summary(fmt.Sprintf("\nRule %d - couldn't resolve value path %s, reason:", ruleIdx, r.Value.String())).
				Details(err.Error()).
				Build(),
			)
		return false, nil, err
	}

	ruleMeta := &elements.RuleMeta{
		Name:          r.Name,
		ValuePath:     r.Value.String(),
		Predicate:     r.Predicate,
		ArgumentsMeta: argMeta,
	}

	fnNode, ok := valNode.Value.(*common.IterFnNodeValue)
	if !ok {
		var val typed.Valuable
		val, err = functions.ParseRawValue(valNode.Value)
		if err != nil {
			blprntExecCtx.diags.
				Append(diagnostics.Builder().Error().
					Details(fmt.Sprintf("\nRule %d - couldn't parse value at path %s, reason:\n%s", ruleIdx, r.Value.String(), err.Error())).
					Build(),
				)
			return false, nil, err
		}
		ok, err = predicate.Test(val, argVal)
		return ok, ruleMeta, err
	}
	ok, err = fnNode.Fn(fnNode.Iter, predicate, argVal)
	return ok, ruleMeta, err
}

func (blprntExecCtx *blueprintExecutionContext) Execute() {
	blprntValidation := &validation.BlueprintValidation{}
	if blprntValidDiags := blprntValidation.Validate(blprntExecCtx.blueprint); blprntValidDiags.HasErrors() {
		blprntExecCtx.diags.Append(blprntValidDiags.Entries()...)
		return
	}

	blprntExecCtx.diags.Append(
		diagnostics.Builder().
			Info().
			Summary(format.Formatter().Detail(format.Box).Dimmed().Color(colors.Blue).
				Format(fmt.Sprintf("Reading %d sources...\n", len(blprntExecCtx.blueprint.Sources)))).
			Build(),
	)

	var remainingSources atomic.Int32
	var signal *sync.Cond = sync.NewCond(&sync.Mutex{})
	remainingSources.Store(int32(len(blprntExecCtx.blueprint.Sources)))

	for alias, source := range blprntExecCtx.blueprint.Sources {
		go func(alias string, source *elements.ConfigurationSource) {
			blprntExecCtx.readSource(alias, source)
			signal.L.Lock()
			remainingSources.Add(-1)
			signal.Signal()
			signal.L.Unlock()
		}(alias, &source)
	}

	signal.L.Lock()
	for remainingSources.Load() > 0 {
		signal.Wait()
	}
	signal.L.Unlock()

	if blprntExecCtx.diags.HasErrors() {
		return
	}

	refResolver := NewReferencesResolver()
	refResolver.Resolve(blprntExecCtx.blueprint.References, blprntExecCtx.diags)
	if blprntExecCtx.diags.HasErrors() {
		return
	}

	blprntExecCtx.diags.Append(
		diagnostics.Builder().Info().
			Summary(format.Formatter().Detail(format.Pencil).Dimmed().Color(colors.Blue).
				Format("Evaluating ruleset...")),
	)

	var ruleViolationsCount int32 = 0
	for idx, rule := range blprntExecCtx.blueprint.Ruleset {
		ok, ruleMeta, err := blprntExecCtx.executeRule(idx, &rule)
		if err != nil {
			blprntExecCtx.diags.Append(diagnostics.Builder().
				Error().
				Summary(fmt.Sprintf("\nRule %d - couldn't execute rule, reason:", idx+1)).
				Details(err.Error()).
				Build(),
			)
			continue
		}

		if !ok {
			blprntExecCtx.diags.Append(diagnostics.Builder().
				Error().
				Details(ruleViolationErrorMessage(ruleMeta, idx)).
				Build(),
			)
			ruleViolationsCount++
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
			),
	)
}
