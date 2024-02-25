package execution

import (
	"fmt"

	"github.com/conformize/conformize/common/diagnostics"
	"github.com/conformize/conformize/internal/blueprint/elements"
	"github.com/conformize/conformize/internal/providers"
)

type ProviderInitializationStep struct {
	alias  string
	config *elements.ConfigurationSource
}

func NewProviderInitializationStep(alias string, config *elements.ConfigurationSource) *ProviderInitializationStep {
	return &ProviderInitializationStep{
		alias:  alias,
		config: config,
	}
}

func (step *ProviderInitializationStep) Run(blprntExecCtx *BlueprintExecutionContext) {
	prvdrInitCtx := &providers.ProviderInitializationContext{
		ValueReferencesStore:       blprntExecCtx.valueReferencesStore,
		Alias:                      step.alias,
		ProvidersDependenciesGraph: blprntExecCtx.providersDependenciesGraph,
	}

	providerFactory := providers.ProviderFactory()
	provider, err := providerFactory.Provider(step.config.Provider, prvdrInitCtx)
	if err != nil {
		blprntExecCtx.diags.Append(
			diagnostics.Builder().Error().
				Summary("Failed to initialize provider").
				Details(fmt.Sprintf("Provider '%s' initialization failed, reason: %v", step.config.Provider, err)).
				Build(),
		)
		return
	}

	err = blprntExecCtx.providersRegistry.Register(step.alias, provider)
	if err != nil {
		blprntExecCtx.diags.Append(
			diagnostics.Builder().Error().
				Details(fmt.Sprintf("Provider '%s' registration failed, reason: %v", step.alias, err)).
				Build(),
		)
	}
}
