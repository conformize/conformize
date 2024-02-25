package providers

import (
	"github.com/conformize/conformize/internal/valuereferencesstore"
)

type ProviderInitializationContext struct {
	ValueReferencesStore *valuereferencesstore.ValueReferencesStore
}
