// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package vault

import (
	"testing"

	sdk "github.com/conformize/conformize/internal/providers/api"
)

func vaultClientConfigWithUserPassAuth() map[string]any {
	return map[string]any{
		"address":   "http://localhost:8200",
		"mountPath": "kv/data/",
		"auth": map[string]any{
			"userPass": map[string]any{
				"username": "admin",
				"password": "admin",
			},
		},
	}
}

func TestVaultProviderConfiguration(t *testing.T) {
	vaultProvider := New("vaultProvider")

	config := vaultClientConfigWithUserPassAuth()
	configurationRequest := sdk.NewConfigurationRequest(vaultProvider.ConfigurationSchema())
	configurationRequest.Set(&config)
	if err := vaultProvider.Configure(configurationRequest); err != nil {
		t.Errorf("failed to configure Vault provider: %v", err)
	}
}

func TestVaultProviderProvisionDataRequest(t *testing.T) {
	vaultProvider := New("vaultProvider")

	configurationRequest := sdk.NewConfigurationRequest(vaultProvider.ConfigurationSchema())
	configurationRequest.SetAtPath("address", "http://localhost:8200")
	configurationRequest.SetAtPath("mountPath", "kv/data/")
	configurationRequest.SetAtPath("auth.userPass.username", "admin")
	configurationRequest.SetAtPath("auth.userPass.password", "admin")
	vaultProvider.Configure(configurationRequest)

	provisionDataRequest := sdk.NewProviderDataRequest(vaultProvider.ProvisionDataRequestSchema())
	provisionDataRequest.SetAtPath("paths", []string{"/app/api/endpoint", "/app/api/config"})

	if data, diags := vaultProvider.Provide(provisionDataRequest); diags.HasErrors() {
		t.Errorf("failed to query Vault: %s\n", diags.Errors())
	} else {
		data.PrintTree()
	}
}
