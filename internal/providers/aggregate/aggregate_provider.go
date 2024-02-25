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
	config *aggregateConfig
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

	valRefStore := valuereferencesstore.Instance()
	root := ds.NewNode[string, interface{}]()

	var sourceData *ds.Node[string, interface{}]
	var pathParser = pathparser.NewPathParser()

	var refPath *path.Path
	for _, valuePath := range paths {
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
		refPath = path.NewPath(pathSteps)
		sourceData, err = valRefStore.GetAtPath(refPath)
		if err != nil {
			diags.Append(diagnostics.Builder().
				Error().
				Details(fmt.Sprintf("Error resolving path '%s': %s", valuePath, err.Error())).
				Build(),
			)
			break
		}

		if sourceData == nil {
			diags.Append(diagnostics.Builder().
				Error().
				Details(fmt.Sprintf("Data at path '%s' is nil", valuePath)).
				Build(),
			)
			break
		}
		root = functions.MergeTrees(root, sourceData)
	}
	return root, diags
}
