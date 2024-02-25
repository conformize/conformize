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

func vaultClientConfigWithUserPassAuth() map[string]interface{} {
	return map[string]interface{}{
		"Address":   "http://localhost:8200",
		"MountPath": "kv/data/",
		"Auth": map[string]interface{}{
			"UserPass": map[string]interface{}{
				"Username": "admin",
				"Password": "admin",
			},
		},
	}
}

func TestVaultProviderConfiguration(t *testing.T) {
	vaultProvider := New()

	config := vaultClientConfigWithUserPassAuth()
	configurationRequest := sdk.NewConfigurationRequest(vaultProvider.ConfigurationSchema())
	configurationRequest.Set(&config)
	if err := vaultProvider.Configure(configurationRequest); err != nil {
		t.Errorf("failed to configure Vault provider: %v", err)
	}
}

func TestVaultProviderProvisionDataRequest(t *testing.T) {
	vaultProvider := New()

	configurationRequest := sdk.NewConfigurationRequest(vaultProvider.ConfigurationSchema())
	configurationRequest.SetAtPath("Address", "http://localhost:8200")
	configurationRequest.SetAtPath("MountPath", "kv/data/")
	configurationRequest.SetAtPath("Auth.UserPass.Username", "admin")
	configurationRequest.SetAtPath("Auth.UserPass.Password", "admin")
	vaultProvider.Configure(configurationRequest)

	provisionDataRequest := sdk.NewProviderDataRequest(vaultProvider.ProvisionDataRequestSchema())
	provisionDataRequest.SetAtPath("Paths", []string{"/app/api/endpoint", "/app/api/config"})

	if data, diags := vaultProvider.Provide(provisionDataRequest); diags.HasErrors() {
		t.Errorf("failed to query Vault: %v\n", diags.Errors().Print())
	} else {
		data.PrintTree()
	}
}
