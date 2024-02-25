// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package aggregate

import (
	"testing"

	"github.com/conformize/conformize/common/ds"
	api "github.com/conformize/conformize/internal/providers/api"
	"github.com/conformize/conformize/internal/valuereferencesstore"
)

func BenchmarkAggregateProvider_Provide(b *testing.B) {
	valueStore := valuereferencesstore.NewValueReferencesStore()
	provider := NewAggregateProvider("testAggregate", valueStore, api.NewConfiguredProvidersRegistry(), ds.NewDependencyGraph[string]())

	configReq := api.NewConfigurationRequest(provider.ConfigurationSchema())
	configReq.SetAtPath("paths", []string{"test1", "test2", "test3"})

	err := provider.Configure(configReq)
	if err != nil {
		b.Fatalf("Failed to configure provider: %v", err)
	}

	dataReq := api.NewProviderDataRequest(provider.ProvisionDataRequestSchema())
	dataReq.SetAtPath("sources", []string{"test1", "test2", "test3"})

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		result, diags := provider.Provide(dataReq)

		if result == nil || diags == nil {
			b.Fatal("Unexpected nil result")
		}
	}
}

func BenchmarkAggregateProvider_Provide_EmptyPaths(b *testing.B) {
	// Setup
	valueStore := valuereferencesstore.NewValueReferencesStore()
	provider := NewAggregateProvider("testAggregate", valueStore, api.NewConfiguredProvidersRegistry(), ds.NewDependencyGraph[string]())

	// Configure with empty paths
	configReq := api.NewConfigurationRequest(provider.ConfigurationSchema())
	configReq.SetAtPath("paths", []string{})

	err := provider.Configure(configReq)
	if err != nil {
		b.Fatalf("Failed to configure provider: %v", err)
	}

	// Create provider data request
	dataReq := api.NewProviderDataRequest(provider.ProvisionDataRequestSchema())
	dataReq.SetAtPath("paths", []string{})

	// Reset timer and run benchmark
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		result, diags := provider.Provide(dataReq)

		// Ensure we don't optimize away the call
		if result == nil || diags == nil {
			b.Fatal("Unexpected nil result")
		}
	}
}

func BenchmarkAggregateProvider_Configure(b *testing.B) {
	// Setup
	valueStore := valuereferencesstore.NewValueReferencesStore()
	provider := NewAggregateProvider("testAggregate", valueStore, api.NewConfiguredProvidersRegistry(), ds.NewDependencyGraph[string]())

	// Reset timer and run benchmark
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		configReq := api.NewConfigurationRequest(provider.ConfigurationSchema())
		configReq.SetAtPath("paths", []string{"test1", "test2", "test3"})

		err := provider.Configure(configReq)
		if err != nil {
			b.Fatalf("Failed to configure provider: %v", err)
		}
	}
}
