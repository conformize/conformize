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
	"sync"

	"github.com/conformize/conformize/common"
	"github.com/conformize/conformize/common/diagnostics"
	"github.com/conformize/conformize/common/functions"
	"github.com/conformize/conformize/common/path"
	"github.com/conformize/conformize/common/pathparser"
	"github.com/conformize/conformize/internal/blueprint"
	"github.com/conformize/conformize/internal/blueprint/elements"
	"github.com/conformize/conformize/internal/blueprint/validation"
	"github.com/conformize/conformize/internal/providers"
	api "github.com/conformize/conformize/internal/providers/api"
	"github.com/conformize/conformize/predicates/condition"
	"github.com/conformize/conformize/predicates/predicatefactory"
	"github.com/conformize/conformize/predicates/rule"
	"github.com/conformize/conformize/predicates/ruleset"
)

type builder struct {
	blueprint   *blueprint.Blueprint
	valRefStore *ValueReferencesStore
	rulesMeta   []*elements.RuleMeta
	ruleSet     []*rule.Rule
	diags       diagnostics.Diagnosable
}

func NewBlueprintExecutionContextBuilder() *builder {
	return &builder{
		valRefStore: NewValueReferencesStore(),
	}
}

func (blprntExecCtxBuilder *builder) WithBlueprint(blueprint *blueprint.Blueprint) *builder {
	blprntExecCtxBuilder.blueprint = blueprint
	return blprntExecCtxBuilder
}

func (blprntExecCtxBuilder *builder) WithDiagnostics(diags diagnostics.Diagnosable) *builder {
	blprntExecCtxBuilder.diags = diags
	return blprntExecCtxBuilder
}

func (blprntExecCtxBuilder *builder) withSource(alias string, sourceConfig elements.ConfigurationSource) {
	sourceDiags := diagnostics.NewDiagnostics()
	defer func() { blprntExecCtxBuilder.diags.Append(sourceDiags.Entries()...) }()

	providerFactory := providers.ProviderFactory()
	sourceDiags.Append(diagnostics.Builder().
		Info().
		Summary(fmt.Sprintf("\nConfiguring source '%s' with '%s' provider...", alias, sourceConfig.Provider)).
		Build(),
	)

	provider, err := providerFactory.Provider(sourceConfig.Provider)
	if err == nil {
		providerConfigurer := blueprint.ProviderConfigurer()
		if err = providerConfigurer.Configure(provider, sourceConfig); err != nil {
			sourceDiags.Append(diagnostics.Builder().
				Error().
				Details(
					fmt.Sprintf("Couldn't configure source '%s' with '%s' provider, reason:\n%s",
						alias, sourceConfig.Provider, err.Error()),
				).
				Build(),
			)
			return
		}

		sourceDiags.Append(diagnostics.Builder().
			Info().
			Summary(fmt.Sprintf("Configured source '%s' with '%s' provider.", alias, sourceConfig.Provider)).
			Build(),
		)

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

		if !sourceDiags.HasErrors() {
			sourceDiags.Append(diagnostics.Builder().
				Info().
				Summary(fmt.Sprintf("\nReading from source '%s'...", alias)).
				Build(),
			)
		}

		providerData, providerDiags := provider.Provide(provisionDataReq)
		if !providerDiags.HasErrors() {
			blprntExecCtxBuilder.valRefStore.AddReference(alias, providerData)
			sourceDiags.Append(diagnostics.Builder().
				Info().
				Summary(fmt.Sprintf("Done reading from source '%s'.", alias)).
				Build(),
			)
			return
		}

		sourceDiags.Append(diagnostics.Builder().
			Error().
			Details(
				fmt.Sprintf("\nCouldn't read source '%s' using '%s' provider, reason:\n",
					alias, sourceConfig.Provider),
			).
			Build())
		sourceDiags.Append(providerDiags.Entries()...)
	} else {
		sourceDiags.Append(diagnostics.Builder().
			Error().
			Details(
				fmt.Sprintf("\nCouldn't configure source '%s' with provider '%s', reason:\n%s",
					alias, sourceConfig.Provider, err.Error()),
			).
			Build(),
		)
	}
}

func (blprntExecCtxBuilder *builder) withRule(rIdx int, r elements.Rule) {
	ruleIdx := rIdx + 1
	pathParser := pathparser.NewPathParser()
	valuePathSteps, err := pathParser.Parse(r.Value)
	if err != nil {
		blprntExecCtxBuilder.diags.
			Append(diagnostics.Builder().
				Error().
				Summary(fmt.Sprintf("\nRule %d - couldn't parse value path %s", ruleIdx, r.Value)).
				Details(err.Error()).
				Build(),
			)
		return
	}

	valuePath := path.NewPath(valuePathSteps)
	valRef, err := blprntExecCtxBuilder.valRefStore.GetAtPath(valuePath)
	if err != nil {
		blprntExecCtxBuilder.diags.
			Append(diagnostics.Builder().
				Error().
				Summary(fmt.Sprintf("\nRule %d - couldn't resolve value path %s, reason:", ruleIdx, r.Value)).
				Details(err.Error()).
				Build(),
			)
		return
	}

	var v interface{}
	if _, ok := valRef.Value.(*common.IterFnNodeValue); !ok {
		v, err = functions.ParseRawValue(valRef.Value)
		if err != nil {
			blprntExecCtxBuilder.diags.
				Append(diagnostics.Builder().
					Error().
					Details(fmt.Sprintf("\nRule %d - couldn't parse value at path %s, reason:\n%s", ruleIdx, r.Value, err.Error())).
					Build(),
				)
			return
		}
	} else {
		v = valRef.Value
	}

	ruleBldr := rule.NewRuleBuilder()
	ruleBldr.Value(v)
	predicateCondition := condition.FromString(r.Predicate)
	predicate, err := predicatefactory.Instance().Build(predicateCondition)
	if err != nil {
		blprntExecCtxBuilder.diags.
			Append(diagnostics.Builder().
				Error().
				Summary(fmt.Sprintf("\nRule %d - couldn't build predicate '%s', reason:", ruleIdx, r.Predicate)).
				Details(err.Error()).
				Build(),
			)
		return
	}

	ruleBldr.Predicate(predicate)
	ruleMeta := &elements.RuleMeta{
		Name:      r.Name,
		ValuePath: r.Value,
		Predicate: r.Predicate,
	}

	argMeta := elements.ArgumentMeta{Sensitive: r.Arguments.IsSensitive()}
	switch argVal := r.Arguments.(type) {
	case *elements.RawValue:
		argMeta.Value = argVal.Value
	case *elements.PathValue:
		valuePathSteps, err := pathParser.Parse(argVal.Path)
		if err != nil {
			blprntExecCtxBuilder.diags.
				Append(diagnostics.Builder().
					Error().
					Summary(fmt.Sprintf("Rule %d - Couldn't parse value path %s for argument %", ruleIdx, argVal.Path)).
					Details(err.Error()).
					Build(),
				)
			return
		}

		val, err := blprntExecCtxBuilder.valRefStore.GetAtPath(path.NewPath(valuePathSteps))
		if err != nil {
			blprntExecCtxBuilder.diags.
				Append(diagnostics.Builder().
					Error().
					Summary(fmt.Sprintf("Rule %d - Couldn't resolve value path %s for argumenт", ruleIdx, argVal.Path)).
					Details(err.Error()).
					Build(),
				)
			return
		}
		argMeta.Value = val.Value
		argMeta.Path = argVal.Path
	}

	ruleMeta.ArgumentsMeta = argMeta
	ruleBldr.Arguments(argMeta.Value)

	blprntExecCtxBuilder.rulesMeta[rIdx] = ruleMeta
	rule, err := ruleBldr.Build()
	if err != nil {
		blprntExecCtxBuilder.diags.
			Append(diagnostics.Builder().
				Error().
				Summary(fmt.Sprintf("\nRule %d:\n", ruleIdx)).
				Details(err.Error()).
				Build(),
			)
		return
	}
	blprntExecCtxBuilder.ruleSet[rIdx] = rule
}

func (blprntExecCtxBuilder *builder) Build() *BlueprintExecutionContext {
	blprntValidation := &validation.BlueprintValidation{}
	if blprntValidDiags := blprntValidation.Validate(blprntExecCtxBuilder.blueprint); blprntValidDiags.HasErrors() {
		blprntExecCtxBuilder.diags.Append(blprntValidDiags.Entries()...)
		return nil
	}

	availableCPUs := max(1, runtime.NumCPU()-1)
	cpus := runtime.GOMAXPROCS(availableCPUs)
	defer runtime.GOMAXPROCS(cpus)

	var wg sync.WaitGroup
	sourcesCount := len(blprntExecCtxBuilder.blueprint.Sources)
	wg.Add(sourcesCount)

	sourcesTasks := make(chan struct{}, min(availableCPUs, sourcesCount))
	defer close(sourcesTasks)
	for configSourceAlias, configSource := range blprntExecCtxBuilder.blueprint.Sources {
		sourcesTasks <- struct{}{}
		go func(configSourceAlias string, configSource elements.ConfigurationSource) {
			defer wg.Done()
			blprntExecCtxBuilder.withSource(configSourceAlias, configSource)
			<-sourcesTasks
		}(configSourceAlias, configSource)
	}

	wg.Wait()
	if blprntExecCtxBuilder.diags.HasErrors() {
		return nil
	}

	resolveRefsChan := make(chan struct{})
	refResolver := NewReferencesResolver(blprntExecCtxBuilder.valRefStore, availableCPUs, resolveRefsChan)
	refResolver.Resolve(blprntExecCtxBuilder.blueprint.References, blprntExecCtxBuilder.diags)
	<-resolveRefsChan

	if blprntExecCtxBuilder.diags.HasErrors() {
		return nil
	}

	ruleSetLen := len(blprntExecCtxBuilder.blueprint.Ruleset)
	blprntExecCtxBuilder.rulesMeta = make([]*elements.RuleMeta, ruleSetLen)
	blprntExecCtxBuilder.ruleSet = make([]*rule.Rule, ruleSetLen)

	wg.Add(ruleSetLen)
	rulesTasks := make(chan struct{}, min(availableCPUs, ruleSetLen))
	defer close(rulesTasks)

	for rIdx, r := range blprntExecCtxBuilder.blueprint.Ruleset {
		rulesTasks <- struct{}{}
		go func(rIdx int, r elements.Rule) {
			defer wg.Done()
			blprntExecCtxBuilder.withRule(rIdx, r)
			<-rulesTasks
		}(rIdx, r)
	}

	wg.Wait()
	if blprntExecCtxBuilder.diags.HasErrors() {
		return nil
	}

	blueprintExecutionCtx := NewBlueprintExecutionContext()
	blueprintExecutionCtx.rulesMeta = blprntExecCtxBuilder.rulesMeta
	blueprintExecutionCtx.ruleSet = ruleset.NewRuleSet(blprntExecCtxBuilder.ruleSet)
	return blueprintExecutionCtx
}
