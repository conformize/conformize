package aggregate

import (
	"github.com/conformize/conformize/common/diagnostics"
	"github.com/conformize/conformize/common/ds"
	"github.com/conformize/conformize/common/typed"
	api "github.com/conformize/conformize/internal/providers/api"
	"github.com/conformize/conformize/internal/providers/api/schema"
	"github.com/conformize/conformize/internal/providers/api/schema/attributes"
	"github.com/conformize/conformize/internal/valuereferencesstore"
)

type aggregateConfig struct {
	Paths []string `cnfrmz:"paths"`
}

type aggregateQueryOptions struct {
	Paths []string `cnfrmz:"paths"`
}

type AggregateProvider struct {
	config               *aggregateConfig
	valueReferencesStore *valuereferencesstore.ValueReferencesStore
}

func NewAggregateProvider(valRefStore *valuereferencesstore.ValueReferencesStore) *AggregateProvider {
	return &AggregateProvider{
		valueReferencesStore: valRefStore,
	}
}

func (ap *AggregateProvider) ConfigurationSchema() *schema.Schema {
	return &schema.Schema{
		Description: "Configuration for the Aggregate provider",
		Version:     1,
		Attributes: map[string]schema.Attributeable{
			"paths": &attributes.ListAttribute{ElementsType: &typed.StringTyped{}},
		},
	}
}

func (ap *AggregateProvider) ProvisionDataRequestSchema() *schema.Schema {
	return &schema.Schema{
		Description: "Aggregate resource request schema",
		Version:     1,
		Attributes: map[string]schema.Attributeable{
			"paths": &attributes.ListAttribute{ElementsType: &typed.StringTyped{}},
		},
	}
}

func (ap *AggregateProvider) Configure(req *api.ConfigurationRequest) error {
	var config aggregateConfig
	if err := req.Get(&config); err != nil {
		return err
	}
	ap.config = &config
	return nil
}

func (ap *AggregateProvider) Provide(req *api.ProviderDataRequest) (*ds.Node[string, interface{}], *diagnostics.Diagnostics) {
	var queryOptions aggregateQueryOptions
	diags := diagnostics.NewDiagnostics()
	if err := req.Get(&queryOptions); err != nil {
		diags.Append(diagnostics.Builder().Error().Details(err.Error()).Build())
		return nil, diags
	}

	root := ds.NewNode[string, interface{}]()
	return root, diags
}
