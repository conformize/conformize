// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package aggregate

import (
	"testing"
	"time"

	"github.com/conformize/conformize/common/ds"
	api "github.com/conformize/conformize/internal/providers/api"
	"github.com/conformize/conformize/internal/valuereferencesstore"
)

func TestAggregateProvider_Provide_BasicFunctionality(t *testing.T) {
	valueStore := valuereferencesstore.NewValueReferencesStore()
	provider := NewAggregateProvider("aggregateTestProvider", valueStore, api.NewConfiguredProvidersRegistry(), ds.NewDependencyGraph[string]())

	configReq := api.NewConfigurationRequest(provider.ConfigurationSchema())
	configReq.SetAtPath("paths", []string{})

	err := provider.Configure(configReq)
	if err != nil {
		t.Fatalf("Failed to configure provider: %v", err)
	}

	dataReq := api.NewProviderDataRequest(provider.ProvisionDataRequestSchema())
	dataReq.SetAtPath("paths", []string{})

	result, diags := provider.Provide(dataReq)
	if result == nil {
		t.Fatal("Expected non-nil result")
	}

	if diags == nil {
		t.Fatal("Expected non-nil diagnostics")
	}
}

func TestAggregateProvider_Provide_Timeout(t *testing.T) {
	valueStore := valuereferencesstore.NewValueReferencesStore()
	provider := NewAggregateProvider("aggregateTestProvider", valueStore, api.NewConfiguredProvidersRegistry(), ds.NewDependencyGraph[string]())

	configReq := api.NewConfigurationRequest(provider.ConfigurationSchema())
	configReq.SetAtPath("paths", []string{"nonexistent.path"})

	err := provider.Configure(configReq)
	if err != nil {
		t.Fatalf("Failed to configure provider: %v", err)
	}

	dataReq := api.NewProviderDataRequest(provider.ProvisionDataRequestSchema())
	dataReq.SetAtPath("paths", []string{"nonexistent.path"})

	start := time.Now()

	result, diags := provider.Provide(dataReq)

	elapsed := time.Since(start)
	if elapsed > 35*time.Second {
		t.Fatal("Provider took too long, possible goroutine leak - fix didn't work")
	}

	if result == nil {
		t.Fatal("Expected non-nil result")
	}
	if diags == nil {
		t.Fatal("Expected non-nil diagnostics")
	}

	t.Logf("Test completed in %v - goroutine management working correctly", elapsed)
}

func TestAggregateProvider_ConfigurationSchema(t *testing.T) {
	valueStore := valuereferencesstore.NewValueReferencesStore()
	provider := NewAggregateProvider("testAggregate", valueStore, api.NewConfiguredProvidersRegistry(), ds.NewDependencyGraph[string]())

	schema := provider.ConfigurationSchema()
	if schema == nil {
		t.Fatal("Expected non-nil configuration schema")
	}

	if schema.Description == "" {
		t.Error("Expected non-empty schema description")
	}

	if schema.Attributes == nil {
		t.Fatal("Expected non-nil schema attributes")
	}

	if _, exists := schema.Attributes["paths"]; !exists {
		t.Error("Expected 'paths' attribute in schema")
	}
}

func TestAggregateProvider_ProvisionDataRequestSchema(t *testing.T) {
	valueStore := valuereferencesstore.NewValueReferencesStore()
	provider := NewAggregateProvider("testAggregate", valueStore, api.NewConfiguredProvidersRegistry(), ds.NewDependencyGraph[string]())

	schema := provider.ProvisionDataRequestSchema()
	if schema == nil {
		t.Fatal("Expected non-nil provision data request schema")
	}

	if schema.Description == "" {
		t.Error("Expected non-empty schema description")
	}

	if schema.Attributes == nil {
		t.Fatal("Expected non-nil schema attributes")
	}

	if _, exists := schema.Attributes["paths"]; !exists {
		t.Error("Expected 'paths' attribute in schema")
	}
}

func TestAggregateProvider_Configure_ValidData(t *testing.T) {
	valueStore := valuereferencesstore.NewValueReferencesStore()
	provider := NewAggregateProvider("aggregateTestProvider", valueStore, api.NewConfiguredProvidersRegistry(), ds.NewDependencyGraph[string]())

	configReq := api.NewConfigurationRequest(provider.ConfigurationSchema())
	configReq.SetAtPath("paths", []string{"test.path"})

	err := provider.Configure(configReq)
	if err != nil {
		t.Fatalf("Expected no error for valid configuration data, got: %v", err)
	}
}
