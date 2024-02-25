package providers

import (
	"github.com/conformize/conformize/common/ds"
	sdk "github.com/conformize/conformize/internal/providers/api"
	"github.com/conformize/conformize/internal/valuereferencesstore"
)

type ProviderInitializationContext struct {
	Alias                      string
	ValueReferencesStore       *valuereferencesstore.ValueReferencesStore
	ProvidersDependenciesGraph *ds.DependencyGraph[string]
	ProvidersRegistry          sdk.ProvidersRegistrar
}
