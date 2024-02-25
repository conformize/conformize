// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package consul

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/conformize/conformize/common/ds"
	sdk "github.com/conformize/conformize/internal/providers/api"
	"github.com/golang/mock/gomock"
	"github.com/hashicorp/consul/api"
)

func getConfigurationRequest(provider sdk.ConfigurationProvider) *sdk.ConfigurationRequest {
	configurationRequest := sdk.NewConfigurationRequest(provider.ConfigurationSchema())
	configurationRequest.SetAtPath("address", "localhost:8500")
	configurationRequest.SetAtPath("scheme", "http")
	configurationRequest.SetAtPath("datacenter", "dc1")
	return configurationRequest
}

func getProvider(factory ConsulClientFactory) sdk.ConfigurationProvider {
	consulClientFactory(factory)
	return &ConsulProvider{}
}

func TestConsulProviderConfiguration(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConsulClient := NewMockConsulClient(ctrl)
	mockConsulKV := NewMockKV(ctrl)
	mockConsulClient.EXPECT().KV().Return(mockConsulKV).AnyTimes()

	mockClientFactory := func(config *api.Config) (ConsulClient, error) {
		return mockConsulClient, nil
	}

	provider := getProvider(mockClientFactory)
	configurationRequest := sdk.NewConfigurationRequest(provider.ConfigurationSchema())

	err := provider.Configure(configurationRequest)
	if err != nil {
		t.Error(err)
	}
}

func TestMockConsulProviderProvideRequestSucceeds(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConsulClient := NewMockConsulClient(ctrl)
	mockConsulKV := NewMockKV(ctrl)

	mockConsulClient.EXPECT().KV().Return(mockConsulKV).AnyTimes()

	mockClientFactory := func(config *api.Config) (ConsulClient, error) {
		return mockConsulClient, nil
	}

	provider := getProvider(mockClientFactory)
	configurationRequest := getConfigurationRequest(provider)
	if err := provider.Configure(configurationRequest); err != nil {
		t.Error(err)
		return
	}

	resourceRequest := sdk.NewProviderDataRequest(provider.ProvisionDataRequestSchema())
	resourceRequest.SetAtPath("path", "/")
	resourceRequest.SetAtPath("recursive", true)

	mockConsulKV.
		EXPECT().
		List("/", gomock.Any()).
		DoAndReturn(func(prefix string, opts *api.QueryOptions) (api.KVPairs, *api.QueryMeta, error) {
			return api.KVPairs{
				&api.KVPair{Key: "app/api/host", Value: []byte("localhost")},
				&api.KVPair{Key: "app/api/key", Value: []byte("Test")},
			}, nil, nil
		}).Times(1)

	data, diags := provider.Provide(resourceRequest)
	if diags.HasErrors() {
		t.Error()
	}
	expected := ds.NewNode[string, any]()
	appNode := expected.AddChild("app")
	apiNode := appNode.AddChild("api")
	hostNode := apiNode.AddChild("host")
	hostNode.Value = "localhost"
	keyNode := apiNode.AddChild("key")
	keyNode.Value = "Test"

	fmt.Printf("Expected:\n")
	expected.PrintTree()

	fmt.Printf("\nGot:\n")
	data.PrintTree()

	if !reflect.DeepEqual(expected, data) {
		t.Error()
	}
}
