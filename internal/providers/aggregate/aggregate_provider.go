package aggregate

import (
	"fmt"

	"github.com/conformize/conformize/common/diagnostics"
	"github.com/conformize/conformize/common/ds"
	"github.com/conformize/conformize/common/functions"
	"github.com/conformize/conformize/common/typed"
	api "github.com/conformize/conformize/internal/providers/api"
	"github.com/conformize/conformize/internal/providers/api/schema"
	"github.com/conformize/conformize/internal/providers/api/schema/attributes"
	"github.com/conformize/conformize/internal/valuereferencesstore"
)

type aggregateConfig struct {
	sources []string `cnfrmz:"sources"`
}

type AggregateProvider struct {
	alias                      string
	config                     *aggregateConfig
	valueReferencesStore       *valuereferencesstore.ValueReferencesStore
	providersRegistrar         api.ProvidersRegistrar
	providersDependenciesGraph *ds.DependencyGraph[string]
}

func NewAggregateProvider(alias string, valRefStore *valuereferencesstore.ValueReferencesStore,
	providersRegistrar api.ProvidersRegistrar, providersDependenciesGraph *ds.DependencyGraph[string]) *AggregateProvider {

	return &AggregateProvider{
		alias:                      alias,
		valueReferencesStore:       valRefStore,
		providersRegistrar:         providersRegistrar,
		providersDependenciesGraph: providersDependenciesGraph,
	}
}

func (ap *AggregateProvider) Alias() string {
	return ap.alias
}

func (ap *AggregateProvider) ConfigurationSchema() *schema.Schema {
	return &schema.Schema{
		Description: "Configuration for the Aggregate provider",
		Version:     1,
		Attributes: map[string]schema.Attributeable{
			"sources": &attributes.ListAttribute{ElementsType: &typed.StringTyped{}},
		},
	}
}

func (ap *AggregateProvider) ProvisionDataRequestSchema() *schema.Schema {
	return &schema.Schema{}
}

func (ap *AggregateProvider) Configure(req *api.ConfigurationRequest) error {
	var config aggregateConfig
	if err := req.Get(&config); err != nil {
		return err
	}
	ap.config = &config

	for _, sourceAlias := range ap.config.sources {
		depPrvd, exists := ap.providersRegistrar.Get(sourceAlias)
		if !exists {
			return fmt.Errorf("source provider with alias '%s' not found", sourceAlias)
		}

		if _, ok := depPrvd.(*AggregateProvider); ok {
			ap.providersDependenciesGraph.AddEdge(ap.alias, sourceAlias)
		}
	}

	return nil
}

func (ap *AggregateProvider) Provide(req *api.ProviderDataRequest) (*ds.Node[string, any], *diagnostics.Diagnostics) {
	diags := diagnostics.NewDiagnostics()
	root := ds.NewNode[string, any]()

	for _, sourceAlias := range ap.config.sources {
		sourceRef, exists := ap.valueReferencesStore.GetReference(sourceAlias)
		if !exists {
			diags.Append(diagnostics.Builder().Error().
				Summary("Failed to get source for aggregate provider").
				Details(fmt.Sprintf("source with alias '%s' does not exist", sourceAlias)).
				Build(),
			)
			continue
		}

		root = functions.MergeTrees(root, sourceRef)
	}

	if diags.HasErrors() {
		return nil, diags
	}

	return root, diags
}
