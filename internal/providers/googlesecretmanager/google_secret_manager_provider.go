package googlesecretmanager

import (
	"context"
	"fmt"
	"regexp"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
	"github.com/conformize/conformize/common/diagnostics"
	"github.com/conformize/conformize/common/ds"
	"github.com/conformize/conformize/common/typed"
	sdk "github.com/conformize/conformize/internal/providers/api"
	"github.com/conformize/conformize/internal/providers/api/schema"
	"github.com/conformize/conformize/internal/providers/api/schema/attributes"
	"github.com/conformize/conformize/internal/providers/googlesecretmanager/builder"
	"google.golang.org/api/iterator"
)

const defaultPageSize = 10
const secretNameRegExpr string = `^projects/[^/]+/secrets/([^/]+)$`
const defaultScope string = "https://www.googleapis.com/auth/cloud-platform"

var secretNameRegExp *regexp.Regexp = regexp.MustCompile(secretNameRegExpr)

type googleSecretMangerCredentials struct {
	UseADC          bool             `cnfrmz:"useADC"`
	CredentialsFile *credentialsFile `cnfrmz:"credentialsFile"`
	Token           *string          `cnfrmz:"token"`
}

type credentialsFile struct {
	Path  string `cnfrmz:"path"`
	Scope string `cnfrmz:"scope"`
}

type googleSecretManagerClientConfig struct {
	Project     string                        `cnfrmz:"project"`
	Credentials googleSecretMangerCredentials `cnfrmz:"credentials"`
}

type googleSecretManagerQueryOptions struct {
	SecretNames []string `cnfrmz:"secretNames"`
	Filter      string   `cnfrmz:"filter"`
	PageSize    int      `cnfrmz:"pageSize"`
}

type GoogleSecretManagerProvider struct {
	alias     string
	projectId string
	client    *secretmanager.Client
}

func New(alias string) *GoogleSecretManagerProvider {
	return &GoogleSecretManagerProvider{alias: alias}
}

func (googleSecretMngrPrvd *GoogleSecretManagerProvider) Alias() string {
	return googleSecretMngrPrvd.alias
}

func (googleSecretMngrPrvd *GoogleSecretManagerProvider) ConfigurationSchema() *schema.Schema {
	return &schema.Schema{
		Description: "Configuration for the Google Secret Manager Provider",
		Version:     1,
		Attributes: map[string]schema.Attributeable{
			"project": &attributes.StringAttribute{},
			"credentials": &attributes.ObjectAttribute{
				FieldsTypes: map[string]typed.Typeable{
					"useADC": &typed.BooleanTyped{},
					"credentialsFile": &typed.ObjectTyped{
						FieldsTypes: map[string]typed.Typeable{
							"path":  &typed.StringTyped{},
							"scope": &typed.StringTyped{},
						},
					},
					"token": &typed.StringTyped{},
				},
			},
		},
	}
}

func (googleSecretMngrPrvd *GoogleSecretManagerProvider) ProvisionDataRequestSchema() *schema.Schema {
	return &schema.Schema{
		Description: "Google Secret Manager secrets request schema",
		Version:     1,
		Attributes: map[string]schema.Attributeable{
			"secretNames": &attributes.ListAttribute{ElementsType: &typed.StringTyped{}},
			"filter":      &attributes.StringAttribute{},
			"pageSize":    &attributes.NumberAttribute{},
		},
	}
}

func (googleSecretMngrPvdr *GoogleSecretManagerProvider) Configure(req *sdk.ConfigurationRequest) error {
	var clientConfig googleSecretManagerClientConfig
	if err := req.Get(&clientConfig); err != nil {
		return err
	}

	googleSecretMngrBldr := builder.New()
	if clientConfig.Credentials.UseADC {
		googleSecretMngrBldr.WithADC()
	}

	if clientConfig.Credentials.CredentialsFile != nil {
		scope := clientConfig.Credentials.CredentialsFile.Scope
		if len(scope) == 0 {
			scope = defaultScope
		}
		googleSecretMngrBldr.WithCredentialsFile(clientConfig.Credentials.CredentialsFile.Path, scope)
	}

	if clientConfig.Credentials.Token != nil {
		googleSecretMngrBldr.WithToken(*clientConfig.Credentials.Token)
	}

	var diags *diagnostics.Diagnostics
	googleSecretMngrPvdr.client, diags = googleSecretMngrBldr.Build()
	if diags.HasErrors() {
		return fmt.Errorf("couldn't configure google secret manager provider, reason:\n%s", diags.Entries().String())
	}
	googleSecretMngrPvdr.projectId = clientConfig.Project
	return nil
}

func (googleSecretMngrPvdr *GoogleSecretManagerProvider) Provide(req *sdk.ProviderDataRequest) (*ds.Node[string, any], *diagnostics.Diagnostics) {
	diags := diagnostics.NewDiagnostics()
	if googleSecretMngrPvdr.client == nil {
		diags.Append(diagnostics.Builder().Error().Details("Google Secret Manager provider is not configured").Build())
		return nil, diags
	}
	defer googleSecretMngrPvdr.client.Close()

	var queryOptions googleSecretManagerQueryOptions
	if err := req.Get(&queryOptions); err != nil {
		diags.Append(diagnostics.Builder().Error().Details(err.Error()).Build())
		return nil, diags
	}

	if len(queryOptions.SecretNames)+len(queryOptions.Filter) == 0 {
		diags.Append(diagnostics.Builder().
			Error().
			Summary("\nInvalid query, reason:\n").
			Details("No search criteria specified - please provide either 'secretNames' or 'filter'.").
			Build(),
		)
		return nil, diags
	}

	if len(queryOptions.SecretNames) > 0 && len(queryOptions.Filter) > 0 {
		diags.Append(diagnostics.Builder().
			Error().
			Summary("\nInvalid query, reason:\n").
			Details("Both 'secretNames' and 'filter' provided - please choose only one.").
			Build(),
		)
		return nil, diags
	}

	if len(queryOptions.SecretNames) > 0 {
		secrets, secretsDiags := getSecrets(googleSecretMngrPvdr.client, googleSecretMngrPvdr.projectId, queryOptions.SecretNames)
		if secretsDiags.HasErrors() {
			return nil, secretsDiags
		}
		return getSecretValues(googleSecretMngrPvdr.client, secrets)
	}

	if len(queryOptions.Filter) > 0 {
		var pageSize = defaultPageSize
		if queryOptions.PageSize > 0 {
			pageSize = queryOptions.PageSize
		}
		return queryByFilter(googleSecretMngrPvdr.client, googleSecretMngrPvdr.projectId, queryOptions.Filter, int32(pageSize))
	}
	return nil, diags
}

func getSecretValues(client *secretmanager.Client, secrets []*secretmanagerpb.Secret) (*ds.Node[string, any], *diagnostics.Diagnostics) {
	result := ds.NewNode[string, any]()
	diags := diagnostics.NewDiagnostics()

	var node *ds.Node[string, any]
	for _, secret := range secrets {
		secretVersionUrl := fmt.Sprintf("%s/versions/latest", secret.Name)
		secretVersionReq := &secretmanagerpb.AccessSecretVersionRequest{Name: secretVersionUrl}

		secretVersionRes, err := client.AccessSecretVersion(context.Background(), secretVersionReq)
		if err != nil {
			diags.Append(diagnostics.Builder().Error().Summary(err.Error()))
			return nil, diags
		}

		matches := secretNameRegExp.FindStringSubmatch(secret.Name)
		secretName := matches[1]

		secretVal := string(secretVersionRes.Payload.Data)
		nodes, found := result.GetChildren(secretName)
		if !found {
			node = result.AddChild(secretName)
		} else {
			node = nodes.First()
		}
		node.Value = secretVal

		node.AddAttribute("annotations", secret.Annotations)
		node.AddAttribute("labels", secret.Labels)
		node.AddAttribute("etag", secret.Etag)
	}
	return result, diags
}

func getSecrets(client *secretmanager.Client, projectId string, secretNames []string) ([]*secretmanagerpb.Secret, *diagnostics.Diagnostics) {
	secrets := make([]*secretmanagerpb.Secret, 0)
	diags := diagnostics.NewDiagnostics()
	for _, secretName := range secretNames {
		secretUrl := fmt.Sprintf("projects/%s/secrets/%s", projectId, secretName)
		secretReq := &secretmanagerpb.GetSecretRequest{Name: secretUrl}
		secretRes, err := client.GetSecret(context.Background(), secretReq)
		if err != nil {
			diags.Append(diagnostics.Builder().Error().Summary(err.Error()))
			return nil, diags
		}
		secrets = append(secrets, secretRes)
	}
	return secrets, diags
}

func queryByFilter(client *secretmanager.Client, project string, filter string, pageSize int32) (*ds.Node[string, any], *diagnostics.Diagnostics) {
	diags := diagnostics.NewDiagnostics()
	listSecretReq := &secretmanagerpb.ListSecretsRequest{
		Parent:   fmt.Sprintf("projects/%s", project),
		Filter:   filter,
		PageSize: pageSize,
	}

	secretsIt := client.ListSecrets(context.Background(), listSecretReq)
	var secrets = make([]*secretmanagerpb.Secret, 0)
	for {
		secretResp, err := secretsIt.Next()
		if err == iterator.Done {
			break
		}

		if err != nil {
			diags.Append(diagnostics.Builder().Error().Summary(err.Error()).Build())
			return nil, diags
		}
		secrets = append(secrets, secretResp)
	}
	return getSecretValues(client, secrets)
}
