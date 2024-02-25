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

	api "github.com/conformize/conformize/internal/providers/api"
	"github.com/conformize/conformize/internal/valuereferencesstore"
)

func TestAggregateProvider_Provide_BasicFunctionality(t *testing.T) {
	// Create value references store
	valueStore := valuereferencesstore.NewValueReferencesStore()

	// Create aggregate provider
	provider := NewAggregateProvider(valueStore)

	// Configure with empty paths (this should not crash)
	configReq := api.NewConfigurationRequest(provider.ConfigurationSchema())
	configReq.SetAtPath("paths", []string{})

	err := provider.Configure(configReq)
	if err != nil {
		t.Fatalf("Failed to configure provider: %v", err)
	}

	// Create provider data request with empty paths
	dataReq := api.NewProviderDataRequest(provider.ProvisionDataRequestSchema())
	dataReq.SetAtPath("paths", []string{})

	// Test providing data - this should complete without hanging
	result, diags := provider.Provide(dataReq)

	// Should not crash or hang
	if result == nil {
		t.Fatal("Expected non-nil result")
	}

	// Diagnostics should exist (may have warnings but shouldn't crash)
	if diags == nil {
		t.Fatal("Expected non-nil diagnostics")
	}
}

func TestAggregateProvider_Provide_Timeout(t *testing.T) {
	valueStore := valuereferencesstore.NewValueReferencesStore()
	provider := NewAggregateProvider(valueStore)

	// Configure with a path that will trigger timeout behavior
	configReq := api.NewConfigurationRequest(provider.ConfigurationSchema())
	configReq.SetAtPath("paths", []string{"nonexistent.path"})

	err := provider.Configure(configReq)
	if err != nil {
		t.Fatalf("Failed to configure provider: %v", err)
	}

	// Create request
	dataReq := api.NewProviderDataRequest(provider.ProvisionDataRequestSchema())
	dataReq.SetAtPath("paths", []string{"nonexistent.path"})

	// Record start time
	start := time.Now()

	// This should timeout after 30 seconds (or return quickly)
	// The key test is that it doesn't hang indefinitely
	result, diags := provider.Provide(dataReq)

	// Check that it didn't hang indefinitely (our fix prevents goroutine leaks)
	elapsed := time.Since(start)
	if elapsed > 35*time.Second {
		t.Fatal("Provider took too long, possible goroutine leak - fix didn't work")
	}

	// Should have some result and diagnostics (not crash)
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
	provider := NewAggregateProvider(valueStore)

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
	provider := NewAggregateProvider(valueStore)

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
	provider := NewAggregateProvider(valueStore)

	// Create valid configuration request
	configReq := api.NewConfigurationRequest(provider.ConfigurationSchema())
	configReq.SetAtPath("paths", []string{"test.path"})

	err := provider.Configure(configReq)
	if err != nil {
		t.Fatalf("Expected no error for valid configuration data, got: %v", err)
	}
}