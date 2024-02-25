package providers

import (
	"github.com/conformize/conformize/common/ds"
	"github.com/conformize/conformize/internal/valuereferencesstore"
)

type ProviderInitializationContext struct {
	Alias                      string
	ValueReferencesStore       *valuereferencesstore.ValueReferencesStore
	ProvidersDependenciesGraph *ds.DependencyGraph[string]
}
