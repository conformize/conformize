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
	XML: func(initCtx *ProviderInitializationContext) sdk.ConfigurationProvider {
		return file.NewFileProvider(initCtx.Alias, &xml.XmlFileUnmarshal{})
	},
	DotEnvFile: func(initCtx *ProviderInitializationContext) sdk.ConfigurationProvider {
		return file.NewFileProvider(initCtx.Alias, &env.EnvFileUnmarshal{})
	},
	TOML: func(initCtx *ProviderInitializationContext) sdk.ConfigurationProvider {
		return file.NewFileProvider(initCtx.Alias, &toml.TomlFilelUnmarshal{})
	},
	JSON: func(initCtx *ProviderInitializationContext) sdk.ConfigurationProvider {
		return file.NewFileProvider(initCtx.Alias, &yaml.YamlUnmarshal{})
	},
	Properties: func(initCtx *ProviderInitializationContext) sdk.ConfigurationProvider {
		return file.NewFileProvider(initCtx.Alias, &properties.PropertiesFileUnmarshal{})
	},
	YAML: func(initCtx *ProviderInitializationContext) sdk.ConfigurationProvider {
		return file.NewFileProvider(initCtx.Alias, &yaml.YamlUnmarshal{})
	},
	Consul: func(initCtx *ProviderInitializationContext) sdk.ConfigurationProvider {
		return consul.New(initCtx.Alias)
	},
	AwsParameterStore: func(initCtx *ProviderInitializationContext) sdk.ConfigurationProvider {
		return awsparameterstore.New(initCtx.Alias)
	},
	AwsSecretsManager: func(initCtx *ProviderInitializationContext) sdk.ConfigurationProvider {
		return secretsmanager.New(initCtx.Alias)
	},
	Vault: func(initCtx *ProviderInitializationContext) sdk.ConfigurationProvider {
		return vault.New(initCtx.Alias)
	},
	Etcd: func(initCtx *ProviderInitializationContext) sdk.ConfigurationProvider { return etcd.New(initCtx.Alias) },
	AzureKeyVaultSecrets: func(initCtx *ProviderInitializationContext) sdk.ConfigurationProvider {
		return azurekeyvault.New(initCtx.Alias)
	},
	AzureDevOpsVariableGroup: func(initCtx *ProviderInitializationContext) sdk.ConfigurationProvider {
		return azuredevopsvariablegroup.New(initCtx.Alias)
	},
	Env: func(initCtx *ProviderInitializationContext) sdk.ConfigurationProvider {
		return environment.New(initCtx.Alias)
	},
	GoogleSecretManager: func(initCtx *ProviderInitializationContext) sdk.ConfigurationProvider {
		return googlesecretmanager.New(initCtx.Alias)
	},
	Http: func(initCtx *ProviderInitializationContext) sdk.ConfigurationProvider { return http.New(initCtx.Alias) },
	HCL: func(initCtx *ProviderInitializationContext) sdk.ConfigurationProvider {
		return file.NewFileProvider(initCtx.Alias, &hcl.HclFileUnmarshal{})
	},
	Aggregate: func(initCtx *ProviderInitializationContext) sdk.ConfigurationProvider {
		return aggregate.NewAggregateProvider(initCtx.Alias,
			initCtx.ValueReferencesStore, initCtx.ProvidersRegistry, initCtx.ProvidersDependenciesGraph,
		)
	},
}

func (pn ProviderName) build(ctx *ProviderInitializationContext) (sdk.ConfigurationProvider, error) {
	if provider, ok := supportedProviders[pn]; ok {
		return provider(ctx), nil
	}
	return nil, fmt.Errorf("provider %s not found", pn)
}

func (pn ProviderName) String() string {
	return string(pn)
}

type providerFactoryFn func(ctx *ProviderInitializationContext) sdk.ConfigurationProvider

type providersFactory struct {
	providers map[ProviderName]providerFactoryFn
}

func (pf *providersFactory) Provider(name string, ctx *ProviderInitializationContext) (sdk.ConfigurationProvider, error) {
	return ProviderName(name).build(ctx)
}

var instance *providersFactory
var once = sync.Once{}

func ProviderFactory() *providersFactory {
	once.Do(func() {
		instance = &providersFactory{providers: supportedProviders}
	})
	return instance
}
