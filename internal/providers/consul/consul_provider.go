// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package consul

import (
	"strings"
	"sync"
	"time"

	"github.com/conformize/conformize/common/diagnostics"
	"github.com/conformize/conformize/common/ds"
	"github.com/conformize/conformize/common/path"
	"github.com/conformize/conformize/common/typed"
	sdk "github.com/conformize/conformize/internal/providers/api"
	"github.com/conformize/conformize/internal/providers/api/schema"
	"github.com/conformize/conformize/internal/providers/api/schema/attributes"
	"github.com/conformize/conformize/serialization"
	"github.com/conformize/conformize/serialization/unmarshal/yaml"

	"github.com/hashicorp/consul/api"
)

type ConsulClientFactory func(*api.Config) (ConsulClient, error)

type ConsulClient interface {
	KV() KV
}

type KV interface {
	Get(key string, q *api.QueryOptions) (*api.KVPair, *api.QueryMeta, error)
	List(prefix string, q *api.QueryOptions) (api.KVPairs, *api.QueryMeta, error)
}

type consulClient struct {
	client *api.Client
}

func (cw *consulClient) KV() KV {
	return cw.client.KV()
}

type ConsulProvider struct {
	alias  string
	client ConsulClient
}

var (
	clientFactory ConsulClientFactory
	once          sync.Once
)

func defaultClientFactory(config *api.Config) (ConsulClient, error) {
	c, _ := api.NewClient(config)
	return &consulClient{
		client: c,
	}, nil
}

func consulClientFactory(factory ConsulClientFactory) {
	once.Do(func() {
		clientFactory = factory
	})
}

func New(alias string) *ConsulProvider {
	once.Do(func() {
		clientFactory = defaultClientFactory
	})
	return &ConsulProvider{alias: alias}
}

func (consulPrvdr *ConsulProvider) Alias() string {
	return consulPrvdr.alias
}

func (consulPrvdr *ConsulProvider) ConfigurationSchema() *schema.Schema {
	return &schema.Schema{
		Description: "Configuration for the Consul provider",
		Version:     1,
		Attributes: map[string]schema.Attributeable{
			"address":    &attributes.StringAttribute{},
			"scheme":     &attributes.StringAttribute{},
			"datacenter": &attributes.StringAttribute{},
			"token":      &attributes.StringAttribute{},
			"tlsEnabled": &attributes.BooleanAttribute{},
			"tls": &attributes.ObjectAttribute{
				FieldsTypes: map[string]typed.Typeable{
					"address":            &typed.StringTyped{},
					"caFile":             &typed.StringTyped{},
					"caPath":             &typed.StringTyped{},
					"certFile":           &typed.StringTyped{},
					"keyFile":            &typed.StringTyped{},
					"insecureSkipVerify": &typed.BooleanTyped{},
				},
			},
		},
	}
}

func (consulPrvdr *ConsulProvider) ProvisionDataRequestSchema() *schema.Schema {
	return &schema.Schema{
		Description: "Consul resource request schema",
		Version:     1,
		Attributes: map[string]schema.Attributeable{
			"path":              &attributes.StringAttribute{},
			"namespace":         &attributes.StringAttribute{},
			"partition":         &attributes.StringAttribute{},
			"datacenter":        &attributes.StringAttribute{},
			"peer":              &attributes.StringAttribute{},
			"allowStale":        &attributes.BooleanAttribute{},
			"requireConsistent": &attributes.BooleanAttribute{},
			"useCache":          &attributes.BooleanAttribute{},
			"maxCacheAge":       &attributes.NumberAttribute{},
			"token":             &attributes.StringAttribute{},
			"recursive":         &attributes.BooleanAttribute{},
		},
	}
}

type consulQueryOptions struct {
	Path              string `cnfrmz:"path"`
	Namespace         string `cnfrmz:"namespace"`
	Partition         string `cnfrmz:"partition"`
	Datacenter        string `cnfrmz:"datacenter"`
	Peer              string `cnfrmz:"peer"`
	AllowStale        bool   `cnfrmz:"allowStale"`
	RequireConsistent bool   `cnfrmz:"requireConsistent"`
	UseCache          bool   `cnfrmz:"useCache"`
	MaxCacheAge       int64  `cnfrmz:"maxCacheAge"`
	Token             string `cnfrmz:"token"`
	Recursive         bool   `cnfrmz:"recursive"`
}

type consulClientConfig struct {
	Address    string    `cnfrmz:"address"`
	Scheme     string    `cnfrmz:"scheme"`
	Datacenter string    `cnfrmz:"datacenter"`
	Token      string    `cnfrmz:"token"`
	TLSEnabled bool      `cnfrmz:"tlsEnabled"`
	TLS        tlsConfig `cnfrmz:"tls"`
}

type tlsConfig struct {
	Address string `cnfrmz:"address"`

	CAFile string `cnfrmz:"caFile"`

	CAPath string `cnfrmz:"caPath"`

	CertFile string `cnfrmz:"certFile"`

	KeyFile string `cnfrmz:"keyFile"`

	InsecureSkipVerify bool `cnfrmz:"insecureSkipVerify"`
}

func (consulPrvdr *ConsulProvider) Configure(req *sdk.ConfigurationRequest) error {
	var clientConfig consulClientConfig
	if err := req.Get(&clientConfig); err != nil {
		return err
	}

	var apiConfig = api.DefaultConfig()
	apiConfig.Address = clientConfig.Address
	apiConfig.Scheme = clientConfig.Scheme
	apiConfig.Datacenter = clientConfig.Datacenter
	apiConfig.Token = clientConfig.Token

	if clientConfig.TLSEnabled {
		apiConfig.TLSConfig = api.TLSConfig{
			Address:            clientConfig.TLS.Address,
			CAFile:             clientConfig.TLS.CAFile,
			CAPath:             clientConfig.TLS.CAPath,
			CertFile:           clientConfig.TLS.CertFile,
			KeyFile:            clientConfig.TLS.KeyFile,
			InsecureSkipVerify: clientConfig.TLS.InsecureSkipVerify,
		}
	}

	var err error
	consulPrvdr.client, err = clientFactory(apiConfig)
	return err
}

func (consulPrvdr *ConsulProvider) Provide(req *sdk.ProviderDataRequest) (*ds.Node[string, any], *diagnostics.Diagnostics) {
	var queryOptions consulQueryOptions
	req.Get(&queryOptions)

	var query = &api.QueryOptions{
		Namespace:         queryOptions.Namespace,
		Partition:         queryOptions.Partition,
		Datacenter:        queryOptions.Datacenter,
		Peer:              queryOptions.Peer,
		AllowStale:        queryOptions.AllowStale,
		RequireConsistent: queryOptions.RequireConsistent,
		UseCache:          queryOptions.UseCache,
		MaxAge:            time.Duration(queryOptions.MaxCacheAge),
		Token:             queryOptions.Token,
	}

	if queryOptions.Recursive {
		return consulRecursiveQuery(consulPrvdr.client, queryOptions.Path, query)
	}
	return consulQueryKey(consulPrvdr.client, queryOptions.Path, query)
}

func consulQueryKey(client ConsulClient, key string, query *api.QueryOptions) (*ds.Node[string, any], *diagnostics.Diagnostics) {
	diags := diagnostics.NewDiagnostics()
	pair, _, err := client.KV().Get(key, query)
	if err != nil {
		diags.Append(diagnostics.Builder().Error().Details(err.Error()).Build())
		return nil, diags
	}

	result := ds.NewNode[string, any]()
	if pair != nil {
		nodeRef := result.AddChild(key)
		bufData := serialization.NewBufferedData(pair.Value)

		yamlUnmarshal := yaml.YamlUnmarshal{}
		if node, err := yamlUnmarshal.Unmarshal(bufData); err == nil {
			*nodeRef = *ds.MergeNodes(nodeRef, node)
		} else {
			nodeRef.Value = string(pair.Value)
		}
		nodeRef.AddAttribute("Flags", pair.Flags)
	}
	return result, diags
}

func consulRecursiveQuery(client ConsulClient, prefix string, query *api.QueryOptions) (*ds.Node[string, any], *diagnostics.Diagnostics) {
	diags := diagnostics.NewDiagnostics()
	pairs, _, err := client.KV().List(prefix, query)
	if err != nil {
		diags.Append(diagnostics.Builder().Error().Details(err.Error()).Build())
		return nil, diags
	}

	result := ds.NewNode[string, any]()
	for _, pair := range pairs {
		isFolderPath := strings.HasSuffix(pair.Key, "/")
		if !isFolderPath {
			valuePath, _ := path.NewFromStringWithSeparator(pair.Key, '/')

			lastNodeRef := result
			steps := valuePath.Steps()
			for step, hasNext := steps.Next(); hasNext; step, hasNext = steps.Next() {
				stepName := step.String()
				if nodes, found := lastNodeRef.GetChildren(stepName); found {
					lastNodeRef = nodes.First()
				} else {
					lastNodeRef = lastNodeRef.AddChild(stepName)
				}
			}

			bufData := serialization.NewBufferedData(pair.Value)
			yamlUnmarshal := yaml.YamlUnmarshal{}
			if node, err := yamlUnmarshal.Unmarshal(bufData); err == nil {
				*lastNodeRef = *ds.MergeNodes(lastNodeRef, node)
			} else {
				lastNodeRef.Value = string(pair.Value)
			}
		}
	}
	return result, diags
}
