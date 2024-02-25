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

	sdk "github.com/conformize/conformize/internal/providers/api"
	"github.com/conformize/conformize/internal/providers/awsparameterstore"
	"github.com/conformize/conformize/internal/providers/azuredevopsvariablegroup"
	"github.com/conformize/conformize/internal/providers/azurekeyvault"
	"github.com/conformize/conformize/internal/providers/consul"
	environment "github.com/conformize/conformize/internal/providers/env"
	"github.com/conformize/conformize/internal/providers/etcd"
	"github.com/conformize/conformize/internal/providers/file/env"
	"github.com/conformize/conformize/internal/providers/file/json"
	"github.com/conformize/conformize/internal/providers/file/properties"
	"github.com/conformize/conformize/internal/providers/file/toml"
	"github.com/conformize/conformize/internal/providers/file/xml"
	"github.com/conformize/conformize/internal/providers/file/yaml"
	"github.com/conformize/conformize/internal/providers/googlesecretmanager"
	"github.com/conformize/conformize/internal/providers/secretsmanager"
	"github.com/conformize/conformize/internal/providers/vault"
)

type ProviderName string

const (
	Xml                      ProviderName = "xml"
	Toml                     ProviderName = "toml"
	Json                     ProviderName = "json"
	DotEnvFile               ProviderName = "dotenv"
	Properties               ProviderName = "properties"
	Yaml                     ProviderName = "yaml"
	Consul                   ProviderName = "consul"
	AwsParameterStore        ProviderName = "aws_parameter_store"
	AwsSecretsManager        ProviderName = "aws_secrets_manager"
	Vault                    ProviderName = "vault"
	Etcd                     ProviderName = "etcd"
	AzureKeyVaultSecrets     ProviderName = "azure_keyvault_secrets"
	AzureDevOpsVariableGroup ProviderName = "azure_devops_variable_group"
	Env                      ProviderName = "env"
	GoogleSecretManager      ProviderName = "google_secret_manager"
)

var supportedProviders = map[ProviderName]providerFactoryFn{
	Xml:                      func() sdk.ConfigurationProvider { return &xml.XmlFileProvider{} },
	DotEnvFile:               func() sdk.ConfigurationProvider { return &env.DotEnvFileProvider{} },
	Toml:                     func() sdk.ConfigurationProvider { return &toml.TomlFileProvider{} },
	Json:                     func() sdk.ConfigurationProvider { return &json.JsonFileProvider{} },
	Properties:               func() sdk.ConfigurationProvider { return &properties.PropertiesFileProvider{} },
	Yaml:                     func() sdk.ConfigurationProvider { return &yaml.YamlFileProvider{} },
	Consul:                   func() sdk.ConfigurationProvider { return consul.New() },
	AwsParameterStore:        func() sdk.ConfigurationProvider { return &awsparameterstore.AwsParameterStoreProvider{} },
	AwsSecretsManager:        func() sdk.ConfigurationProvider { return &secretsmanager.AwsSecretsManagerProvider{} },
	Vault:                    func() sdk.ConfigurationProvider { return &vault.VaultProvider{} },
	Etcd:                     func() sdk.ConfigurationProvider { return etcd.New() },
	AzureKeyVaultSecrets:     func() sdk.ConfigurationProvider { return &azurekeyvault.AzureKeyVaultProvider{} },
	AzureDevOpsVariableGroup: func() sdk.ConfigurationProvider { return &azuredevopsvariablegroup.AzureDevOpsVariableGroupProvider{} },
	Env:                      func() sdk.ConfigurationProvider { return &environment.EnvProvider{} },
	GoogleSecretManager:      func() sdk.ConfigurationProvider { return &googlesecretmanager.GoogleSecretManagerProvider{} },
}

func (pn ProviderName) new() (sdk.ConfigurationProvider, error) {
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
	return ProviderName(name).new()
}

var instance *providersFactory
var once = sync.Once{}

func ProviderFactory() *providersFactory {
	once.Do(func() {
		instance = &providersFactory{providers: supportedProviders}
	})
	return instance
}
