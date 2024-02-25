// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package azurekeyvault

import (
	"context"
	"fmt"
	"runtime"
	"sync"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azsecrets"
	"github.com/conformize/conformize/common/diagnostics"
	"github.com/conformize/conformize/common/ds"
	"github.com/conformize/conformize/common/typed"
	sdk "github.com/conformize/conformize/internal/providers/api"
	"github.com/conformize/conformize/internal/providers/api/schema"
	"github.com/conformize/conformize/internal/providers/api/schema/attributes"
)

const (
	cliCredentials  string = "cli"
	envCredentials  string = "environment"
	managedIdentity string = "managedIdentity"
)

type azureKeyVaultClient struct {
	client *azsecrets.Client
}

type clientSecretAuth struct {
	TenantId string `cnfrmz:"tenantId"`
	ClientId string `cnfrmz:"clientId"`
	Secret   string `cnfrmz:"secret"`
}

type pipelineAuth struct {
	TenantId            string `cnfrmz:"tenantId"`
	ClientId            string `cnfrmz:"clientId"`
	ServiceConnectionId string `cnfrmz:"serviceConnectionId"`
	Token               string `cnfrmz:"token"`
}

type keyVaultCredentials struct {
	ClientSecret *clientSecretAuth `cnfrmz:"clientSecret"`
	Provided     *string           `cnfrmz:"provided"`
	Pipeline     *pipelineAuth     `cnfrmz:"pipeline"`
}

type azureKeyVaultConfig struct {
	VaultUrl    string               `cnfrmz:"vaultUrl"`
	Credentials *keyVaultCredentials `cnfrmz:"credentials"`
}

type AzureKeyVaultProvider struct {
	alias  string
	client *azureKeyVaultClient
}

type keyVaultQueryOptions struct {
	SecretNames []string `cnfrmz:"secretNames"`
}

func (azureKVPrvdr *AzureKeyVaultProvider) ConfigurationSchema() *schema.Schema {
	return &schema.Schema{
		Description: "Configuration for the Azure Key Vault provider",
		Version:     1,
		Attributes: map[string]schema.Attributeable{
			"vaultUrl": &attributes.StringAttribute{},
			"credentials": &attributes.ObjectAttribute{
				FieldsTypes: map[string]typed.Typeable{
					"clientSecret": &typed.ObjectTyped{
						FieldsTypes: map[string]typed.Typeable{
							"tenantId": &typed.StringTyped{},
							"clientId": &typed.StringTyped{},
							"secret":   &typed.StringTyped{},
						},
					},
					"provided": &typed.StringTyped{},
					"pipeline": &typed.ObjectTyped{
						FieldsTypes: map[string]typed.Typeable{
							"tenantId":            &typed.StringTyped{},
							"clientId":            &typed.StringTyped{},
							"serviceConnectionId": &typed.StringTyped{},
							"token":               &typed.StringTyped{},
						},
					},
				},
			},
		},
	}
}

func (azurekvPrvdr *AzureKeyVaultProvider) ProvisionDataRequestSchema() *schema.Schema {
	return &schema.Schema{
		Attributes: map[string]schema.Attributeable{
			"secretNames": &attributes.ListAttribute{ElementsType: &typed.StringTyped{}},
		},
	}
}

func getCredentials(creds *keyVaultCredentials) (azcore.TokenCredential, error) {
	if creds.ClientSecret != nil {
		clientSecretCreds := creds.ClientSecret
		return azidentity.NewClientSecretCredential(clientSecretCreds.TenantId, clientSecretCreds.ClientId, clientSecretCreds.Secret, nil)
	}

	if creds.Pipeline != nil {
		pipelineCreds := creds.Pipeline
		return azidentity.NewAzurePipelinesCredential(pipelineCreds.TenantId, pipelineCreds.ClientId, pipelineCreds.ServiceConnectionId, pipelineCreds.Token, nil)
	}

	if provided := creds.Provided; provided != nil {
		switch *provided {
		case cliCredentials:
			return azidentity.NewAzureCLICredential(nil)
		case envCredentials:
			return azidentity.NewEnvironmentCredential(nil)
		case managedIdentity:
			return azidentity.NewManagedIdentityCredential(nil)
		default:
			return nil, fmt.Errorf("unsupported provided credentials type %s", *provided)
		}
	}

	return nil, fmt.Errorf("no credentials specified")
}

func (azureKVPrvdr *AzureKeyVaultProvider) Configure(req *sdk.ConfigurationRequest) error {
	var clientConfig azureKeyVaultConfig
	if err := req.Get(&clientConfig); err != nil {
		return err
	}

	creds, credsErr := getCredentials(clientConfig.Credentials)
	if credsErr != nil {
		return credsErr
	}

	azsecretsClient, clientErr := azsecrets.NewClient(clientConfig.VaultUrl, creds, nil)
	if clientErr != nil {
		return clientErr
	}

	azureKVPrvdr.client = &azureKeyVaultClient{client: azsecretsClient}
	return nil
}

const maxSecretBatchSize = 10

func (azureKVPrvdr *AzureKeyVaultProvider) Provide(queryRequest *sdk.ProviderDataRequest) (*ds.Node[string, any], *diagnostics.Diagnostics) {
	diags := diagnostics.NewDiagnostics()
	if azureKVPrvdr.client == nil {
		diags.Append(diagnostics.Builder().Error().Details("Azure Key Vault provider is not configured").Build())
		return nil, diags
	}

	var queryOptions keyVaultQueryOptions
	if err := queryRequest.Get(&queryOptions); err != nil {
		diags.Append(diagnostics.Builder().Error().Details(err.Error()).Build())
		return nil, diags
	}

	result := ds.NewNode[string, any]()
	if len(queryOptions.SecretNames) > 0 {
		secretsLen := len(queryOptions.SecretNames)
		batchesCount := (secretsLen + maxSecretBatchSize - 1) / maxSecretBatchSize

		secretChan := make(chan []azsecrets.Secret, batchesCount)
		errChan := make(chan error)
		doneChan := make(chan struct{})

		defer close(errChan)
		defer close(secretChan)
		defer close(doneChan)

		secrets := make([]azsecrets.Secret, 0)
		go func() {
			done := false
			for !done {
				select {
				case secretsBatch := <-secretChan:
					secrets = append(secrets, secretsBatch...)
				case err := <-errChan:
					diags.Append(diagnostics.Builder().Error().Details(err.Error()).Build())
				case <-doneChan:
					done = true
				}
			}
		}()

		var wg sync.WaitGroup
		wg.Add(batchesCount)

		maxParallelTasks := min(runtime.NumCPU()-1, batchesCount)
		tasks := make(chan struct{}, maxParallelTasks)
		defer close(tasks)

		for i, offset := 0, 0; i <= batchesCount && offset < secretsLen; i, offset = i+1, offset+maxSecretBatchSize {
			upperBound := min(offset+maxSecretBatchSize, secretsLen)
			secretNames := queryOptions.SecretNames[offset:upperBound]

			tasks <- struct{}{}
			go func(secretNames []string) {
				defer wg.Done()
				secretsBatch := make([]azsecrets.Secret, 0)
				for _, secretName := range secretNames {
					resp, err := azureKVPrvdr.client.client.GetSecret(context.Background(), secretName, "", nil)
					if err == nil {
						secretsBatch = append(secretsBatch, resp.Secret)
						continue
					}
					errChan <- fmt.Errorf("failed to retrieve secret %s: %w", secretName, err)
				}
				secretChan <- secretsBatch
				<-tasks
			}(secretNames)
		}

		wg.Wait()
		doneChan <- struct{}{}

		var node *ds.Node[string, any]
		for _, secret := range secrets {
			node = nil
			if nodes, found := result.GetChildren(secret.ID.Name()); !found {
				node = result.AddChild(secret.ID.Name())
			} else {
				node = nodes.First()
			}

			if secret.Value != nil {
				node.Value = *secret.Value
			}

			addAttributes(node, secret.Attributes)
			if secret.ContentType != nil {
				node.AddAttribute("contentType", *secret.ContentType)
			}

			if len(secret.Tags) > 0 {
				tags := make(map[string]string)
				for tagName, tagValue := range secret.Tags {
					tags[tagName] = *tagValue
				}
				node.AddAttribute("tags", tags)
			}
		}
	}

	return result, diags
}

func (azureKVPrvdr *AzureKeyVaultProvider) Alias() string {
	return azureKVPrvdr.alias
}

func New(alias string) *AzureKeyVaultProvider {
	return &AzureKeyVaultProvider{alias: alias}
}

func addAttributes(node *ds.Node[string, any], attrs *azsecrets.SecretAttributes) {
	if attrs != nil {
		if attrs.Enabled != nil {
			node.AddAttribute("enabled", *attrs.Enabled)
		}
		if attrs.Created != nil {
			node.AddAttribute("created", attrs.Created.String())
		}
		if attrs.Expires != nil {
			node.AddAttribute("expires", attrs.Expires.String())
		}
	}
}
