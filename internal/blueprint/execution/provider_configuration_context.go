package execution

import (
	"github.com/conformize/conformize/common/ds"
	"github.com/conformize/conformize/internal/blueprint/elements"
	"github.com/conformize/conformize/internal/valuereferencesstore"
)

type ProviderConfigurationContext struct {
	Alias                      string
	Config                     *elements.ConfigurationSource
	ProvidersRegistry          *ConfiguredProvidersRegistry
	ValueReferencesStore       *valuereferencesstore.ValueReferencesStore
	ProvidersDependenciesGraph *ds.DependencyGraph[string]
}
