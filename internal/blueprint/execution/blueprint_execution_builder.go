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

	"github.com/conformize/conformize/common/diagnostics"
	"github.com/conformize/conformize/common/functions"
	"github.com/conformize/conformize/common/typed"
	"github.com/conformize/conformize/internal/blueprint"
	"github.com/conformize/conformize/internal/blueprint/elements"
	"github.com/conformize/conformize/internal/blueprint/validation"
	"github.com/conformize/conformize/internal/providers"
	api "github.com/conformize/conformize/internal/providers/api"
	"github.com/conformize/conformize/predicates/condition"
	"github.com/conformize/conformize/predicates/rule"
	"github.com/conformize/conformize/predicates/ruleset"
)

type builder struct {
	blueprint   *blueprint.Blueprint
	valRefStore *blueprint.ValueReferencesStore
	rulesMeta   []*elements.RuleMeta
	ruleSet     []*rule.Rule
	diags       diagnostics.Diagnosable
}

func NewBlueprintExecutionBuilder() *builder {
	return &builder{
		valRefStore: blueprint.NewValueReferencesStore(),
	}
}

func (blprntExecBuilder *builder) WithBlueprint(blueprint *blueprint.Blueprint) *builder {
	blprntExecBuilder.blueprint = blueprint
	return blprntExecBuilder
}

func (blprntExecBuilder *builder) WithDiagnostics(diags diagnostics.Diagnosable) *builder {
	blprntExecBuilder.diags = diags
	return blprntExecBuilder
}

func (blprntExecBuilder *builder) withSource(alias string, sourceConfig elements.ConfigurationSource, diags diagnostics.Diagnosable) *builder {
	sourceDiags := diagnostics.NewDiagnostics()
	defer func() { diags.Append(sourceDiags.Entries()...) }()

	providerFactory := providers.ProviderFactory()
	sourceDiags.Append(diagnostics.Builder().
		Info().
		Summary(fmt.Sprintf("\nConfiguring source '%s' with '%s' provider...", alias, sourceConfig.Provider)).
		Build(),
	)

	if provider, err := providerFactory.Provider(sourceConfig.Provider); err == nil {
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
			return blprntExecBuilder
		}
		sourceDiags.Append(diagnostics.Builder().
			Info().
			Summary(fmt.Sprintf("Configured source '%s' with '%s' provider.", alias, sourceConfig.Provider)).
			Build(),
		)

		var provisionDataReq *api.ProviderDataRequest
		if sourceConfig.QueryOptions != nil {
			var queryErr error
			if provisionDataReq, queryErr = providerConfigurer.Query(provider, sourceConfig.QueryOptions); queryErr != nil {
				sourceDiags.Append(diagnostics.Builder().
					Error().
					Details(
						fmt.Sprintf("Couldn't set query options for source '%s' with '%s' provider, reason:\n%s",
							alias, sourceConfig.Provider, queryErr.Error()),
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
			if providerData, providerDiags := provider.Provide(provisionDataReq); !providerDiags.HasErrors() {
				blprntExecBuilder.valRefStore.AddReference(alias, &elements.ValueReference{Node: providerData})
				sourceDiags.Append(diagnostics.Builder().
					Info().
					Summary(fmt.Sprintf("Done reading from source '%s'.", alias)).
					Build(),
				)
			} else {
				for _, providerDiag := range providerDiags.Entries() {
					sourceDiags.Append(diagnostics.Builder().
						Error().
						Summary(
							fmt.Sprintf("Couldn't read from source '%s', reason:", alias),
						).
						Details(fmt.Sprintf("%s\n%s", providerDiag.GetSummary(), providerDiag.GetDetails())).
						Build(),
					)
				}
			}
		}
	} else {
		sourceDiags.Append(diagnostics.Builder().
			Error().
			Details(
				fmt.Sprintf("Couldn't configure source '%s' with provider '%s', reason:\n%s",
					alias, sourceConfig.Provider, err.Error()),
			).
			Build(),
		)
	}
	return blprntExecBuilder
}

func (blprntExecBuilder *builder) withRule(rIdx int, r elements.Rule, diags diagnostics.Diagnosable) *builder {
	ruleIdx := rIdx + 1
	if valRef, err := blprntExecBuilder.valRefStore.GetAtPath(r.Value); err != nil {
		diags.Append(diagnostics.Builder().
			Error().
			Summary(fmt.Sprintf("Rule %d - couldn't resolve value path %s", ruleIdx, r.Value)).
			Details(err.Error()).
			Build(),
		)
	} else {
		val, err := functions.ParseRawValue(valRef.Value)
		if err != nil {
			diags.Append(diagnostics.Builder().
				Error().
				Details(fmt.Sprintf("Rule %d - couldn't parse value at path %s, reason:\n%s", ruleIdx, r.Value, err.Error())).
				Build(),
			)
			return blprntExecBuilder
		}
		ruleBldr := rule.NewRuleBuilder()
		ruleBldr.Value(val)
		predicate := condition.FromString(r.Predicate)
		ruleBldr.Predicate(predicate)

		ruleMeta := &elements.RuleMeta{
			ValuePath: r.Value,
			Predicate: r.Predicate,
		}

		for idx, arg := range r.Arguments {
			var argIdx = idx + 1
			var ruleArg typed.RawValue
			argMeta := elements.ArgumentMeta{Sensitive: arg.IsSensitive()}
			if argVal, ok := arg.(*elements.RawValue); ok {
				ruleArg = argVal.Value
			}

			if argPathVal, ok := arg.(*elements.PathValue); ok {
				argPath := argPathVal.GetValue().(string)
				val, err := blprntExecBuilder.valRefStore.GetAtPath(argPath)
				if err != nil {
					diags.Append(diagnostics.Builder().
						Error().
						Summary(
							fmt.Sprintf("Rule %d - Couldn't resolve value path %s for argument %d",
								ruleIdx, argPath, argIdx),
						).
						Details(err.Error()).
						Build(),
					)
					return blprntExecBuilder
				}
				ruleArg = val.Value
				argMeta.Path = argPath

			}

			argMeta.Value = ruleArg
			ruleMeta.ArgumentsMeta = append(ruleMeta.ArgumentsMeta, argMeta)
			if argVal, err := functions.ParseRawValue(ruleArg); err == nil {
				ruleBldr.Argument(argVal)
			} else {
				diags.Append(diagnostics.Builder().
					Error().
					Summary(
						fmt.Sprintf("Rule %d - Couldn't parse argument %d:\n", ruleIdx, argIdx),
					).
					Details(err.Error()).
					Build(),
				)
			}
		}
		blprntExecBuilder.rulesMeta[rIdx] = ruleMeta
		blprntExecBuilder.ruleSet[rIdx] = ruleBldr.Build()
	}
	return blprntExecBuilder
}

func (blprntExecBuilder *builder) withReference(refAlias, refPath string, diags diagnostics.Diagnosable) *builder {
	if ref, refErr := blprntExecBuilder.valRefStore.GetAtPath(refPath); refErr == nil {
		blprntExecBuilder.valRefStore.AddReference(refAlias, ref)
	} else {
		diags.Append(diagnostics.Builder().
			Error().
			Details(
				fmt.Sprintf("Couldn't resolve reference %s in path %s, reason:\n\n%s",
					refAlias, refPath, refErr.Error()),
			).
			Build(),
		)
	}
	return blprntExecBuilder
}

func (blprntExecBuilder *builder) Build() *BlueprintExecution {
	blprntValidation := &validation.BlueprintValidation{}
	blprntValidDiags := blprntValidation.Validate(blprntExecBuilder.blueprint)
	blprntExecBuilder.diags.Append(blprntValidDiags.Entries()...)
	if blprntValidDiags.HasErrors() {
		return nil
	}

	availableCPUs := max(1, runtime.NumCPU()-1)
	cpus := runtime.GOMAXPROCS(availableCPUs)
	defer runtime.GOMAXPROCS(cpus)

	var wg sync.WaitGroup
	sourcesCount := len(blprntExecBuilder.blueprint.Sources)
	wg.Add(sourcesCount)

	sourcesTasks := make(chan struct{}, min(availableCPUs, sourcesCount))
	defer close(sourcesTasks)
	for configSourceAlias, configSource := range blprntExecBuilder.blueprint.Sources {
		sourcesTasks <- struct{}{}
		go func(configSourceAlias string, configSource elements.ConfigurationSource) {
			defer wg.Done()
			blprntExecBuilder.withSource(configSourceAlias, configSource, blprntExecBuilder.diags)
			<-sourcesTasks
		}(configSourceAlias, configSource)
	}

	wg.Wait()
	if blprntExecBuilder.diags.HasErrors() {
		return nil
	}

	referencesCount := len(blprntExecBuilder.blueprint.References)
	wg.Add(referencesCount)

	referencesTasks := make(chan struct{}, min(availableCPUs, referencesCount))
	defer close(referencesTasks)
	for refAlias, refPath := range blprntExecBuilder.blueprint.References {
		referencesTasks <- struct{}{}
		go func(refAlias, refPath string) {
			defer wg.Done()
			blprntExecBuilder.withReference(refAlias, refPath, blprntExecBuilder.diags)
			<-referencesTasks
		}(refAlias, refPath)
	}

	wg.Wait()
	if blprntExecBuilder.diags.HasErrors() {
		return nil
	}

	ruleSetLen := len(blprntExecBuilder.blueprint.Ruleset)
	blprntExecBuilder.rulesMeta = make([]*elements.RuleMeta, ruleSetLen)
	blprntExecBuilder.ruleSet = make([]*rule.Rule, ruleSetLen)

	wg.Add(ruleSetLen)
	rulesTasks := make(chan struct{}, min(availableCPUs, ruleSetLen))
	defer close(rulesTasks)
	for rIdx, r := range blprntExecBuilder.blueprint.Ruleset {
		rulesTasks <- struct{}{}
		go func(rIdx int, r elements.Rule) {
			defer wg.Done()
			blprntExecBuilder.withRule(rIdx, r, blprntExecBuilder.diags)
			<-rulesTasks
		}(rIdx, r)
	}

	wg.Wait()
	if blprntExecBuilder.diags.HasErrors() {
		return nil
	}

	blueprintExecution := NewBlueprintExecution()
	blueprintExecution.rulesMeta = blprntExecBuilder.rulesMeta
	blueprintExecution.ruleSet = ruleset.NewRuleSet(blprntExecBuilder.ruleSet)
	return blueprintExecution
}
