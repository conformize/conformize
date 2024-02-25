package scaffold

import (
	"fmt"

	"github.com/conformize/conformize/common/diagnostics"
	"github.com/conformize/conformize/internal/blueprint"
	"github.com/conformize/conformize/internal/blueprint/elements"
	"github.com/conformize/conformize/predicates/condition"
)

type BlueprintScaffoldBuilder struct {
	version    float64
	sources    map[string]elements.ConfigurationSource
	references map[string]string
	ruleset    []elements.Rule
	diags      *diagnostics.Diagnostics
}

func NewBuilder() *BlueprintScaffoldBuilder {
	return &BlueprintScaffoldBuilder{
		sources:    make(map[string]elements.ConfigurationSource),
		references: make(map[string]string),
		ruleset:    make([]elements.Rule, 0),
		diags:      diagnostics.NewDiagnostics(),
	}
}

func (blprntScaffoldBldr *BlueprintScaffoldBuilder) WithVersion(version float64) *BlueprintScaffoldBuilder {
	blprntScaffoldBldr.version = version
	return blprntScaffoldBldr
}

func (blprntScaffoldBldr *BlueprintScaffoldBuilder) WithSource(alias string, provider string) *BlueprintScaffoldBuilder {
	configSrc := elements.ConfigurationSource{Provider: provider}
	blprntScaffoldBldr.sources[alias] = configSrc
	return blprntScaffoldBldr
}

func (blprntScaffoldBldr *BlueprintScaffoldBuilder) WithReference(ref string) *BlueprintScaffoldBuilder {
	blprntScaffoldBldr.references[ref] = ""
	return blprntScaffoldBldr
}

func (blprntScaffoldBldr *BlueprintScaffoldBuilder) WithPredicate(predicateCondition string) *BlueprintScaffoldBuilder {
	cond := condition.FromString(predicateCondition)
	if cond == condition.UNKNOWN {
		blprntScaffoldBldr.diags.Append(diagnostics.Builder().
			Warning().
			Summary(fmt.Sprintf(
				"Unknown predicate '%s' is ignored and won't be added to the blueprint scaffold\n", predicateCondition,
			)).
			Build(),
		)
	}

	// predicate, _ := predicatefactory.Instance().Build(cond)
	// if argsPrd, ok := predicate.(predicates.Predicate); ok {
	//	predicateArgsCount := argsPrd.ArgumentsCount()
	//
	//	rule := elements.Rule{Predicate: predicateCondition}
	//	if predicateArgsCount > 0 {
	//		rule.Arguments = &elements.RawValue{Value: "<argument>"}
	//	}
	//	blprntScaffoldBldr.ruleset = append(blprntScaffoldBldr.ruleset, rule)
	// }

	return blprntScaffoldBldr
}

func (blprntScaffoldBldr *BlueprintScaffoldBuilder) Build() (*blueprint.Blueprint, *diagnostics.Diagnostics) {
	return &blueprint.Blueprint{
		Version:    blprntScaffoldBldr.version,
		Sources:    blprntScaffoldBldr.sources,
		References: blprntScaffoldBldr.references,
		Ruleset:    blprntScaffoldBldr.ruleset,
	}, blprntScaffoldBldr.diags
}
