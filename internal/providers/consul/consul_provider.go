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

func New() *ConsulProvider {
	once.Do(func() {
		clientFactory = defaultClientFactory
	})
	return &ConsulProvider{}
}

func (consulPrvdr *ConsulProvider) ConfigurationSchema() *schema.Schema {
	return &schema.Schema{
		Description: "Configuration for the Consul provider",
		Version:     1,
		Attributes: map[string]schema.Attributeable{
			"Address":    &attributes.StringAttribute{},
			"Scheme":     &attributes.StringAttribute{},
			"PathPrefix": &attributes.StringAttribute{},
			"Datacenter": &attributes.StringAttribute{},
			"Token":      &attributes.StringAttribute{},
			"TlsEnabled": &attributes.BooleanAttribute{},
			"Tls": &attributes.ObjectAttribute{
				FieldsTypes: map[string]typed.Typeable{
					"Address":            &typed.StringTyped{},
					"CAFile":             &typed.StringTyped{},
					"CAPath":             &typed.StringTyped{},
					"CertFile":           &typed.StringTyped{},
					"KeyFile":            &typed.StringTyped{},
					"InsecureSkipVerify": &typed.BooleanTyped{},
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
			"Path":              &attributes.StringAttribute{},
			"Namespace":         &attributes.StringAttribute{},
			"Partition":         &attributes.StringAttribute{},
			"Datacenter":        &attributes.StringAttribute{},
			"Peer":              &attributes.StringAttribute{},
			"AllowStale":        &attributes.BooleanAttribute{},
			"RequireConsistent": &attributes.BooleanAttribute{},
			"UseCache":          &attributes.BooleanAttribute{},
			"MaxCacheAge":       &attributes.NumberAttribute{},
			"Token":             &attributes.StringAttribute{},
			"Recursive":         &attributes.BooleanAttribute{},
		},
	}
}

type consulQueryOptions struct {
	Path              string
	Namespace         string
	Partition         string
	Datacenter        string
	Peer              string
	AllowStale        bool
	RequireConsistent bool
	UseCache          bool
	MaxCacheAge       int64
	Token             string
	Recursive         bool
}

type consulClientConfig struct {
	Address    string
	Scheme     string
	PathPrefix string
	Datacenter string
	Token      string
	TlsEnabled bool
	Tls        api.TLSConfig
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

	if clientConfig.TlsEnabled {
		apiConfig.TLSConfig = clientConfig.Tls
	}

	var err error
	consulPrvdr.client, err = clientFactory(apiConfig)
	return err
}

func (consulPrvdr *ConsulProvider) Provide(req *sdk.ProviderDataRequest) (*ds.Node[string, interface{}], *diagnostics.Diagnostics) {
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

func consulQueryKey(client ConsulClient, key string, query *api.QueryOptions) (*ds.Node[string, interface{}], *diagnostics.Diagnostics) {
	diags := diagnostics.NewDiagnostics()
	pair, _, err := client.KV().Get(key, query)
	if err != nil {
		diags.Append(diagnostics.Builder().Error().Details(err.Error()).Build())
		return nil, diags
	}

	result := ds.NewNode[string, interface{}]()
	if pair != nil {
		nodeRef := result.AddChild(key)
		bufData := serialization.NewBufferedData(pair.Value)

		yamlUnmarshal := yaml.NewYamlUnmarshal(bufData)
		if node, err := yamlUnmarshal.Unmarshal(); err == nil {
			*nodeRef = *ds.MergeNodes(nodeRef, node)
		} else {
			nodeRef.Value = string(pair.Value)
		}
		nodeRef.AddAttribute("Flags", pair.Flags)
	}
	return result, diags
}

func consulRecursiveQuery(client ConsulClient, prefix string, query *api.QueryOptions) (*ds.Node[string, interface{}], *diagnostics.Diagnostics) {
	diags := diagnostics.NewDiagnostics()
	pairs, _, err := client.KV().List(prefix, query)
	if err != nil {
		diags.Append(diagnostics.Builder().Error().Details(err.Error()).Build())
		return nil, diags
	}

	result := ds.NewNode[string, interface{}]()
	for _, pair := range pairs {
		isFolderPath := strings.HasSuffix(pair.Key, "/")
		if !isFolderPath {
			valuePath, _ := path.NewFromStringWithSeparator(pair.Key, '/')

			lastNodeRef := result
			steps := valuePath.Steps()
			for step, hasNext := steps.Next(); hasNext; step, hasNext = steps.Next() {
				stepName := step.String()
				if nodeRef, found := lastNodeRef.GetChild(stepName); found {
					lastNodeRef = nodeRef
				} else {
					lastNodeRef = lastNodeRef.AddChild(stepName)
				}
			}

			bufData := serialization.NewBufferedData(pair.Value)
			yamlUnmarshal := yaml.NewYamlUnmarshal(bufData)
			if node, err := yamlUnmarshal.Unmarshal(); err == nil {
				*lastNodeRef = *ds.MergeNodes(lastNodeRef, node)
			} else {
				lastNodeRef.Value = string(pair.Value)
			}
		}
	}
	return result, diags
}
