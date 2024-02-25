package execution

import (
	"github.com/conformize/conformize/common/ds"
	"github.com/conformize/conformize/internal/blueprint/elements"
	sdk "github.com/conformize/conformize/internal/providers/api"
)

type ProviderConfigurationContext struct {
	Alias                    string
	Config                   *elements.ConfigurationSource
	ProvidersDependencyGraph *ds.DependencyGraph[string]
	ProvidersRegistry        sdk.ProvidersRegistrar
	References               *map[string]string
}
