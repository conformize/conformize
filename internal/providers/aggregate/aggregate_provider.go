package aggregate

import (
	"fmt"

	"github.com/conformize/conformize/common/diagnostics"
	"github.com/conformize/conformize/common/ds"
	"github.com/conformize/conformize/common/functions"
	"github.com/conformize/conformize/common/path"
	"github.com/conformize/conformize/common/pathparser"
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

	var paths []string
	if len(queryOptions.Paths) == 0 {
		paths = ap.config.Paths
	} else {
		paths = queryOptions.Paths
	}

	root := ds.NewNode[string, interface{}]()

	var aggregates = make([]*ds.Node[string, interface{}], len(paths))
	var pathParser = pathparser.NewPathParser()

	for idx, valuePath := range paths {
		pathSteps, err := pathParser.Parse(valuePath)
		if err != nil {
			diags.Append(
				diagnostics.Builder().
					Error().
					Details(err.Error()).
					Build(),
			)
			break
		}

		if len(pathSteps) == 0 {
			diags.Append(
				diagnostics.Builder().
					Error().
					Details(fmt.Sprintf("path '%s' is empty", valuePath)).
					Build(),
			)
			continue
		}

		refPath := path.NewPath(pathSteps)
		ap.valueReferencesStore.SubscribeRef(pathSteps[0].String(), func() {
			data, err := ap.valueReferencesStore.GetAtPath(refPath)
			if err != nil {
				diags.Append(diagnostics.Builder().
					Error().
					Details(fmt.Sprintf("error resolving path '%s': %s", valuePath, err.Error())).
					Build(),
				)
				return
			}

			if data == nil {
				diags.Append(diagnostics.Builder().
					Error().
					Details(fmt.Sprintf("data at path '%s' is nil", valuePath)).
					Build(),
				)
				return
			}
			aggregates[idx] = data
		})
	}

	for _, aggregate := range aggregates {
		root = functions.MergeTrees(root, aggregate)
	}
	return root, diags
}
