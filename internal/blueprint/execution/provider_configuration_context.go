package execution

import (
	"github.com/conformize/conformize/common/ds"
	"github.com/conformize/conformize/internal/blueprint/elements"
)

type ProviderConfigurationContext struct {
	Alias                    string
	Config                   *elements.ConfigurationSource
	ProvidersDependencyGraph *ds.DependencyGraph[string]
	ProvidersRegistry        *ConfiguredProvidersRegistry
	References               *map[string]string
}
