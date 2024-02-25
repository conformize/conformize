// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package providers

import (
	"fmt"
	"sync"

	"github.com/conformize/conformize/internal/providers/aggregate"
	sdk "github.com/conformize/conformize/internal/providers/api"
	"github.com/conformize/conformize/internal/providers/awsparameterstore"
	"github.com/conformize/conformize/internal/providers/azuredevopsvariablegroup"
	"github.com/conformize/conformize/internal/providers/azurekeyvault"
	"github.com/conformize/conformize/internal/providers/consul"
	environment "github.com/conformize/conformize/internal/providers/env"
	"github.com/conformize/conformize/internal/providers/etcd"
	"github.com/conformize/conformize/internal/providers/file"
	"github.com/conformize/conformize/internal/providers/googlesecretmanager"
	"github.com/conformize/conformize/internal/providers/http"
	"github.com/conformize/conformize/internal/providers/secretsmanager"
	"github.com/conformize/conformize/internal/providers/vault"
	"github.com/conformize/conformize/serialization/unmarshal/env"
	"github.com/conformize/conformize/serialization/unmarshal/hcl"
	"github.com/conformize/conformize/serialization/unmarshal/properties"
	"github.com/conformize/conformize/serialization/unmarshal/toml"
	"github.com/conformize/conformize/serialization/unmarshal/xml"
	"github.com/conformize/conformize/serialization/unmarshal/yaml"
)

type ProviderName string

const (
	XML                      ProviderName = "xml"
	TOML                     ProviderName = "toml"
	JSON                     ProviderName = "json"
	DotEnvFile               ProviderName = "dotenv"
	Properties               ProviderName = "properties"
	YAML                     ProviderName = "yaml"
	Consul                   ProviderName = "consul"
	AwsParameterStore        ProviderName = "aws_parameter_store"
	AwsSecretsManager        ProviderName = "aws_secrets_manager"
	Vault                    ProviderName = "vault"
	Etcd                     ProviderName = "etcd"
	AzureKeyVaultSecrets     ProviderName = "azure_keyvault_secrets"
	AzureDevOpsVariableGroup ProviderName = "azure_devops_variable_group"
	Env                      ProviderName = "env"
	GoogleSecretManager      ProviderName = "google_secret_manager"
	Http                     ProviderName = "http"
	HCL                      ProviderName = "hcl"
	Aggregate                ProviderName = "aggregate"
)

var supportedProviders = map[ProviderName]providerFactoryFn{
	XML:                      func() sdk.ConfigurationProvider { return file.NewFileProvider(&xml.XmlFileUnmarshal{}) },
	DotEnvFile:               func() sdk.ConfigurationProvider { return file.NewFileProvider(&env.EnvFileUnmarshal{}) },
	TOML:                     func() sdk.ConfigurationProvider { return file.NewFileProvider(&toml.TomlFilelUnmarshal{}) },
	JSON:                     func() sdk.ConfigurationProvider { return file.NewFileProvider(&yaml.YamlUnmarshal{}) },
	Properties:               func() sdk.ConfigurationProvider { return file.NewFileProvider(&properties.PropertiesFileUnmarshal{}) },
	YAML:                     func() sdk.ConfigurationProvider { return file.NewFileProvider(&yaml.YamlUnmarshal{}) },
	Consul:                   func() sdk.ConfigurationProvider { return consul.New() },
	AwsParameterStore:        func() sdk.ConfigurationProvider { return &awsparameterstore.AwsParameterStoreProvider{} },
	AwsSecretsManager:        func() sdk.ConfigurationProvider { return &secretsmanager.AwsSecretsManagerProvider{} },
	Vault:                    func() sdk.ConfigurationProvider { return &vault.VaultProvider{} },
	Etcd:                     func() sdk.ConfigurationProvider { return etcd.New() },
	AzureKeyVaultSecrets:     func() sdk.ConfigurationProvider { return &azurekeyvault.AzureKeyVaultProvider{} },
	AzureDevOpsVariableGroup: func() sdk.ConfigurationProvider { return &azuredevopsvariablegroup.AzureDevOpsVariableGroupProvider{} },
	Env:                      func() sdk.ConfigurationProvider { return &environment.EnvProvider{} },
	GoogleSecretManager:      func() sdk.ConfigurationProvider { return &googlesecretmanager.GoogleSecretManagerProvider{} },
	Http:                     func() sdk.ConfigurationProvider { return &http.HttpProvider{} },
	HCL:                      func() sdk.ConfigurationProvider { return file.NewFileProvider(&hcl.HclFileUnmarshal{}) },
	Aggregate:                func() sdk.ConfigurationProvider { return &aggregate.AggregateProvider{} },
}

func (pn ProviderName) build() (sdk.ConfigurationProvider, error) {
	if provider, ok := supportedProviders[pn]; ok {
		return provider(), nil
	}
	return nil, fmt.Errorf("provider %s not found", pn)
}

func (pn ProviderName) String() string {
	return string(pn)
}

type providerFactoryFn func() sdk.ConfigurationProvider

type providersFactory struct {
	providers map[ProviderName]providerFactoryFn
}

func (pf *providersFactory) Provider(name string) (sdk.ConfigurationProvider, error) {
	return ProviderName(name).build()
}

var instance *providersFactory
var once = sync.Once{}

func ProviderFactory() *providersFactory {
	once.Do(func() {
		instance = &providersFactory{providers: supportedProviders}
	})
	return instance
}
