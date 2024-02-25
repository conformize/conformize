package env

import (
	"os"
	"strings"

	"github.com/conformize/conformize/common/diagnostics"
	"github.com/conformize/conformize/common/ds"
	"github.com/conformize/conformize/common/typed"
	sdk "github.com/conformize/conformize/internal/providers/api"
	"github.com/conformize/conformize/internal/providers/api/schema"
	"github.com/conformize/conformize/internal/providers/api/schema/attributes"
)

type EnvProvider struct {
	alias string
}

type queryOptions struct {
	Vars []string `cnfrmz:"vars"`
}

func New(alias string) *EnvProvider {
	return &EnvProvider{alias: alias}
}

func (envPrvdr *EnvProvider) Alias() string {
	return envPrvdr.alias
}

func (envPrvdr *EnvProvider) ConfigurationSchema() *schema.Schema {
	return &schema.Schema{
		Description: "Configuration for the Environment variables provider",
		Version:     1,
		Attributes:  map[string]schema.Attributeable{},
	}
}

func (envPrvdr *EnvProvider) Configure(req *sdk.ConfigurationRequest) error {
	return nil
}

func (EnvProvider *EnvProvider) ProvisionDataRequestSchema() *schema.Schema {
	return &schema.Schema{
		Description: "Configuration for the Environment variables provider",
		Version:     1,
		Attributes: map[string]schema.Attributeable{
			"vars": &attributes.ListAttribute{ElementsType: &typed.StringTyped{}},
		},
	}
}

func (envPrvdr *EnvProvider) Provide(queryRequest *sdk.ProviderDataRequest) (*ds.Node[string, any], *diagnostics.Diagnostics) {
	diags := diagnostics.NewDiagnostics()
	var queryOpts queryOptions
	if queryRequest != nil {
		err := queryRequest.Get(&queryOpts)
		if err != nil {
			diags.Append(diagnostics.Builder().Error().Details(err.Error()).Build())
			return nil, diags
		}
	}

	result := ds.NewNode[string, any]()

	doFilter := len(queryOpts.Vars) > 0
	var lookupVars map[string]struct{}
	if doFilter {
		lookupVars = map[string]struct{}{}
		for _, varName := range queryOpts.Vars {
			lookupVars[varName] = struct{}{}
		}
	}

	envVars := os.Environ()
	for _, envVarLine := range envVars {
		envVar := strings.SplitN(envVarLine, "=", 2)
		key := envVar[0]
		if doFilter {
			if _, ok := lookupVars[key]; !ok {
				continue
			}
		}

		envVar = envVar[1:]
		node := result.AddChild(key)
		if len(envVar) > 0 {
			node.Value = envVar[0]
		}
	}

	return result, diags
}
