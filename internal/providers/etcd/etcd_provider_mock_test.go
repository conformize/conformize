// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package etcd

import (
	"context"
	reflect "reflect"
	"testing"

	sdk "github.com/conformize/conformize/internal/providers/api"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type opOptionMatcher struct {
	expected clientv3.OpOption
}

func (opMatch *opOptionMatcher) Matches(x interface{}) bool {
	actual, ok := x.(clientv3.OpOption)
	if !ok {
		return false
	}

	expectedOp := clientv3.Op{}
	actualOp := clientv3.Op{}

	opMatch.expected(&expectedOp)
	actual(&actualOp)

	return reflect.DeepEqual(expectedOp, actualOp)
}

func (opMatch *opOptionMatcher) String() string {
	return "matches expected clientv3.OpOption"
}

func option(expected clientv3.OpOption) gomock.Matcher {
	return &opOptionMatcher{expected: expected}
}

func setupMockEtcdClient(t *testing.T) (*MockEtcdClient, *EtcdProvider) {
	mockCtrl := gomock.NewController(t)
	mockClient := NewMockEtcdClient(mockCtrl)

	clientFactory := func(config clientv3.Config) (EtcdClient, error) {
		return mockClient, nil
	}

	etcdClientFactory(clientFactory)

	provider := New()
	return mockClient, provider
}

func mockResponse(expectedValues map[string]string) (*clientv3.GetResponse, error) {
	res := &clientv3.GetResponse{}
	res.Kvs = make([]*mvccpb.KeyValue, 0)
	for k, v := range expectedValues {
		res.Kvs = append(res.Kvs, &mvccpb.KeyValue{Key: []byte(k), Value: []byte(v)})
	}
	return res, nil
}

func TestEtcdProvider_Configuration(t *testing.T) {
	_, etcdPrvdr := setupMockEtcdClient(t)
	cfgReq := sdk.NewConfigurationRequest(etcdPrvdr.ConfigurationSchema())
	cfgReq.SetAtPath("Endpoints", []string{"http://localhost:2379"})
	cfgReq.SetAtPath("ConnectionTimeout", 30)
	cfgReq.SetAtPath("Authentication.Basic.Username", "root")
	cfgReq.SetAtPath("Authentication.Basic.Password", "rootpw")

	err := etcdPrvdr.Configure(cfgReq)
	require.NoError(t, err)
}

func TestEtcdProvider_ProvideDataKeysList(t *testing.T) {
	mockClient, etcdPrvdr := setupMockEtcdClient(t)
	ctx := context.Background()

	expectedValues := map[string]string{
		"app/api/config/host": "localhost",
		"app/db/host":         "db.localhost",
		"app/db/schema":       "testdb",
		"app/db/user":         "dbUser",
		"app/db/password":     "dbPass",
	}

	for k, v := range expectedValues {
		mockClient.EXPECT().
			Get(ctx, k, gomock.Any()).
			Return(&clientv3.GetResponse{Kvs: []*mvccpb.KeyValue{{Key: []byte(k), Value: []byte(v)}}}, nil)
	}

	cfgReq := sdk.NewConfigurationRequest(etcdPrvdr.ConfigurationSchema())
	cfgReq.SetAtPath("Endpoints", []string{"http://localhost:2379"})
	cfgReq.SetAtPath("ConnectionTimeout", 5)
	cfgReq.SetAtPath("Authentication.Basic.Username", "root")
	cfgReq.SetAtPath("Authentication.Basic.Password", "rootpw")
	etcdPrvdr.Configure(cfgReq)

	provideDataReq := sdk.NewProviderDataRequest(etcdPrvdr.ProvisionDataRequestSchema())
	provideDataReq.SetAtPath("Keys", []string{
		"app/api/config/host",
		"app/db/host",
		"app/db/schema",
		"app/db/user",
		"app/db/password",
	})

	mockClient.EXPECT().Close().Times(1)
	data, diags := etcdPrvdr.Provide(provideDataReq)
	require.NotNil(t, diags)
	require.NotNil(t, data)

	for key, expectedValue := range expectedValues {
		childNode, ok := data.GetChild(key)
		if !ok {
			t.Error()
		}
		require.NotNil(t, childNode)
		require.Equal(t, expectedValue, childNode.Value)
	}
}

func TestEtcdProvider_ProvideDataPrefix(t *testing.T) {
	mockClient, etcdPrvdr := setupMockEtcdClient(t)
	ctx := context.Background()

	expectedValues := map[string]string{
		"app/api/config/host":       "localhost",
		"app/api/config/port":       "8080",
		"app/api/config/rate-limit": "100",
	}

	mockClient.EXPECT().
		Get(ctx, "app/api", option(clientv3.WithPrefix())).
		DoAndReturn(func(ctx context.Context, key string, opts ...clientv3.OpOption) (*clientv3.GetResponse, error) {
			return mockResponse(expectedValues)
		})
	mockClient.EXPECT().Close().Times(1)

	cfgReq := sdk.NewConfigurationRequest(etcdPrvdr.ConfigurationSchema())
	cfgReq.SetAtPath("Endpoints", []string{"http://localhost:2379"})
	cfgReq.SetAtPath("ConnectionTimeout", 30)
	cfgReq.SetAtPath("Authentication.Basic.Username", "root")
	cfgReq.SetAtPath("Authentication.Basic.Password", "rootpw")
	etcdPrvdr.Configure(cfgReq)

	provideDataReq := sdk.NewProviderDataRequest(etcdPrvdr.ProvisionDataRequestSchema())
	provideDataReq.SetAtPath("Prefix", "app/api")
	data, diags := etcdPrvdr.Provide(provideDataReq)
	require.NotNil(t, diags)
	require.NotNil(t, data)

	for key, expectedValue := range expectedValues {
		childNode, ok := data.GetChild(key)
		if !ok {
			t.Error()
		}
		require.NotNil(t, childNode)
		require.Equal(t, expectedValue, childNode.Value)
	}
}

func TestEtcdProvider_ProvideDataRange(t *testing.T) {
	mockClient, etcdPrvdr := setupMockEtcdClient(t)
	ctx := context.Background()

	expectedValues := map[string]string{
		"app/api/config/host": "localhost",
		"app/api/config/port": "8080",
	}

	startKey := "app/api/config/host"
	endKey := "app/api/config/tls"
	mockClient.EXPECT().
		Get(ctx, startKey, option(clientv3.WithRange(endKey))).
		DoAndReturn(func(ctx context.Context, key string, opts ...clientv3.OpOption) (*clientv3.GetResponse, error) {
			return mockResponse(expectedValues)
		})
	mockClient.EXPECT().Close().Times(1)

	cfgReq := sdk.NewConfigurationRequest(etcdPrvdr.ConfigurationSchema())
	cfgReq.SetAtPath("Endpoints", []string{"http://localhost:2379"})
	cfgReq.SetAtPath("ConnectionTimeout", 30)
	cfgReq.SetAtPath("Authentication.Basic.Username", "user")
	cfgReq.SetAtPath("Authentication.Basic.Password", "pass")
	etcdPrvdr.Configure(cfgReq)

	provideDataReq := sdk.NewProviderDataRequest(etcdPrvdr.ProvisionDataRequestSchema())
	provideDataReq.SetAtPath("Range.StartKey", "app/api/config/host")
	provideDataReq.SetAtPath("Range.EndKey", "app/api/config/tls")
	data, diags := etcdPrvdr.Provide(provideDataReq)
	require.NotNil(t, diags)
	require.NotNil(t, data)

	for key, expectedValue := range expectedValues {
		childNode, ok := data.GetChild(key)
		if !ok {
			t.Error()
		}
		require.NotNil(t, childNode)
		require.Equal(t, expectedValue, childNode.Value)
	}

}
