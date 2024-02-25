// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package etcd

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"time"

	"github.com/conformize/conformize/common/diagnostics"
	"github.com/conformize/conformize/common/ds"
	"github.com/conformize/conformize/common/typed"
	sdk "github.com/conformize/conformize/internal/providers/api"
	"github.com/conformize/conformize/internal/providers/api/schema"
	"github.com/conformize/conformize/internal/providers/api/schema/attributes"

	"go.etcd.io/etcd/api/v3/mvccpb"
	"go.etcd.io/etcd/client/pkg/v3/transport"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type EtcdClient interface {
	Get(ctx context.Context, key string, opts ...clientv3.OpOption) (*clientv3.GetResponse, error)
	Close() error
}

type EtcdClientFactory func(clientv3.Config) (EtcdClient, error)

type etcdClientConfig struct {
	Endpoints         []string        `cnfrmz:"endpoints"`
	ConnectionTimeout int16           `cnfrmz:"connectionTimeout"`
	Authentication    *etcdClientAuth `cnfrmz:"authentication"`
}

type etcdClientAuth struct {
	Basic *etcdBasicAuth `cnfrmz:"basic"`
	TLS   *etcdTlsAuth   `cnfrmz:"tls"`
}

type etcdBasicAuth struct {
	Username string `cnfrmz:"username"`
	Password string `cnfrmz:"password"`
}

type etcdTlsAuth struct {
	KeyFile       string `cnfrmz:"keyFile"`
	CertFile      string `cnfrmz:"certFile"`
	TrustedCAFile string `cnfrmz:"trustedCaFile"`
}

type queryOptions struct {
	Keys   []string        `cnfrmz:"keys"`
	Prefix string          `cnfrmz:"prefix"`
	Range  *keysRangeQuery `cnfrmz:"range"`
}

type keysRangeQuery struct {
	StartKey string `cnfrmz:"startKey"`
	EndKey   string `cnfrmz:"endKey"`
}

type EtcdProvider struct {
	alias  string
	client EtcdClient
}

type etcdClient struct {
	client *clientv3.Client
}

func (ec *etcdClient) Get(ctx context.Context, key string, opts ...clientv3.OpOption) (*clientv3.GetResponse, error) {
	return ec.client.Get(ctx, key, opts...)
}

func (ec *etcdClient) Close() error {
	return ec.client.Close()
}

const maxKVBatchSize = 10

var (
	clientFactory EtcdClientFactory
	once          sync.Once
)

func defaultClientFactory(config clientv3.Config) (EtcdClient, error) {
	c, err := clientv3.New(config)
	return &etcdClient{
		client: c,
	}, err
}

func etcdClientFactory(factory EtcdClientFactory) {
	once.Do(func() {
		clientFactory = factory
	})
}

func New(alias string) *EtcdProvider {
	once.Do(func() {
		clientFactory = defaultClientFactory
	})
	return &EtcdProvider{alias: alias}
}

func (etcdPrvdr *EtcdProvider) Alias() string {
	return etcdPrvdr.alias
}

func (etcdPrvdr *EtcdProvider) ConfigurationSchema() *schema.Schema {
	return &schema.Schema{
		Attributes: map[string]schema.Attributeable{
			"endpoints":         &attributes.ListAttribute{ElementsType: &typed.StringTyped{}},
			"connectionTimeout": &attributes.NumberAttribute{},
			"authentication": &attributes.ObjectAttribute{
				FieldsTypes: map[string]typed.Typeable{
					"basic": &typed.ObjectTyped{
						FieldsTypes: map[string]typed.Typeable{
							"username": &typed.StringTyped{},
							"password": &typed.StringTyped{},
						},
					},
					"tls": &typed.ObjectTyped{
						FieldsTypes: map[string]typed.Typeable{
							"keyFile":       &typed.StringTyped{},
							"certFile":      &typed.StringTyped{},
							"trustedCAFile": &typed.StringTyped{},
						},
					},
				},
			},
		},
	}
}

func (etcdPrvdr *EtcdProvider) ProvisionDataRequestSchema() *schema.Schema {
	return &schema.Schema{
		Attributes: map[string]schema.Attributeable{
			"keys":   &attributes.ListAttribute{ElementsType: &typed.StringTyped{}},
			"prefix": &attributes.StringAttribute{},
			"range": &attributes.ObjectAttribute{
				FieldsTypes: map[string]typed.Typeable{
					"startKey": &typed.StringTyped{},
					"endKey":   &typed.StringTyped{},
				},
			},
		},
	}
}

func (etcdPrvdr *EtcdProvider) Configure(req *sdk.ConfigurationRequest) error {
	var clientConfig etcdClientConfig
	var err error

	if err = req.Get(&clientConfig); err != nil {
		return err
	}

	config := clientv3.Config{
		Endpoints:   clientConfig.Endpoints,
		DialTimeout: time.Duration(clientConfig.ConnectionTimeout) * time.Second,
	}

	if authConfig := clientConfig.Authentication; authConfig != nil {
		if authConfig.Basic != nil {
			config.Username = authConfig.Basic.Username
			config.Password = authConfig.Basic.Password
		} else if authConfig.TLS != nil {
			config.TLS, err = transport.TLSInfo{
				KeyFile:        authConfig.TLS.KeyFile,
				CertFile:       authConfig.TLS.CertFile,
				TrustedCAFile:  authConfig.TLS.TrustedCAFile,
				ClientCertAuth: true,
			}.ClientConfig()
		}
	}

	if err == nil {
		etcdPrvdr.client, err = clientFactory(config)
	}
	return err
}

func (etcdPrvdr *EtcdProvider) Provide(queryRequest *sdk.ProviderDataRequest) (*ds.Node[string, any], *diagnostics.Diagnostics) {
	diags := diagnostics.NewDiagnostics()
	if etcdPrvdr.client == nil {
		diags.Append(diagnostics.Builder().Error().Details("etcd client is not configured").Build())
		return nil, diags
	}
	defer etcdPrvdr.client.Close()

	var queryOptions queryOptions
	if err := queryRequest.Get(&queryOptions); err != nil {
		diags.Append(diagnostics.Builder().Error().Details(err.Error()).Build())
		return nil, diags
	}

	kvs := make([]*mvccpb.KeyValue, 0)
	if len(queryOptions.Keys) > 0 {
		keysLen := len(queryOptions.Keys)
		batchesCount := (keysLen + maxKVBatchSize - 1) / maxKVBatchSize

		kvChan := make(chan []*mvccpb.KeyValue, batchesCount)
		errChan := make(chan error)
		doneChan := make(chan struct{})

		defer close(errChan)
		defer close(kvChan)
		defer close(doneChan)

		go func() {
			done := false
			for !done {
				select {
				case kvsBatch := <-kvChan:
					kvs = append(kvs, kvsBatch...)
				case err := <-errChan:
					diags.Append(diagnostics.Builder().Error().Details(err.Error()).Build())
				case <-doneChan:
					done = true
				}
			}
		}()
		var wg sync.WaitGroup
		wg.Add(batchesCount)

		availableCPUs := max(1, runtime.NumCPU()-1)
		cpus := runtime.GOMAXPROCS(availableCPUs)
		defer runtime.GOMAXPROCS(cpus)

		maxParallelTasksCount := min(availableCPUs, max(1, batchesCount))
		tasks := make(chan struct{}, maxParallelTasksCount)
		defer close(tasks)
		for i, offset := 0, 0; i <= batchesCount && offset < keysLen; i, offset = i+1, offset+maxKVBatchSize {
			upperBound := min(offset+maxKVBatchSize, keysLen)

			keys := queryOptions.Keys[offset:upperBound]
			tasks <- struct{}{}
			go func(keys []string) {
				defer wg.Done()
				for _, key := range keys {
					if resp, err := etcdPrvdr.client.Get(context.Background(), key); err == nil {
						if len(resp.Kvs) > 0 {
							kvChan <- resp.Kvs
						} else {
							errChan <- fmt.Errorf("key %s not found", key)
						}
					} else {
						errChan <- fmt.Errorf("failed to retrieve key %s, reason: %s", key, err.Error())
					}
				}
				<-tasks
			}(keys)
		}
		wg.Wait()
		doneChan <- struct{}{}
	}

	if queryOptions.Prefix != "" {
		prefixOpt := clientv3.WithPrefix()
		if resp, err := etcdPrvdr.client.Get(context.Background(), queryOptions.Prefix, prefixOpt); err == nil {
			kvs = append(kvs, resp.Kvs...)
		} else {
			diags.Append(diagnostics.Builder().Error().Details(err.Error()).Build())
		}
	}

	if queryOptions.Range != nil {
		rangeOpt := clientv3.WithRange(queryOptions.Range.EndKey)
		if resp, err := etcdPrvdr.client.Get(context.Background(), queryOptions.Range.StartKey, rangeOpt); err == nil {
			kvs = append(kvs, resp.Kvs...)
		} else {
			diags.Append(diagnostics.Builder().Error().Details(err.Error()).Build())
		}
	}

	result := ds.NewNode[string, any]()
	var node *ds.Node[string, any]
	for _, kv := range kvs {
		node = nil
		key := string(kv.Key)
		value := string(kv.Value)
		nodes, found := result.GetChildren(key)
		if !found {
			node = result.AddChild(key)
		} else {
			node = nodes.First()
		}
		node.Value = value
	}
	return result, diags
}
