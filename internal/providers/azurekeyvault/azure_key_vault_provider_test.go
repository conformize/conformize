// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package azurekeyvault

import (
	"testing"

	sdk "github.com/conformize/conformize/internal/providers/api"
)

func TestAzureKeyVaultConfiguration(t *testing.T) {
	azureKvPrvdr := AzureKeyVaultProvider{}
	cfgReq := sdk.NewConfigurationRequest(azureKvPrvdr.ConfigurationSchema())
	cfgReq.SetAtPath("vaultUrl", "FAKE_VAULT_URL")
	cfgReq.SetAtPath("credentials.provided", "managedIdentity")

	err := azureKvPrvdr.Configure(cfgReq)
	if err != nil {
		t.Errorf("failed to configure azure keyvault provider, reason: %s", err.Error())
	}
}

func TestAzureKeyVaultProvideSecret(t *testing.T) {
	azureKvPrvdr := AzureKeyVaultProvider{}
	cfgReq := sdk.NewConfigurationRequest(azureKvPrvdr.ConfigurationSchema())
	cfgReq.SetAtPath("vaultUrl", "FAKE_VAULT_URL")
	cfgReq.SetAtPath("credentials.clientSecret.tenantId", "FAKE_TENANT_ID")
	cfgReq.SetAtPath("credentials.clientSecret.clientId", "FAKE_CLIENT_ID")
	cfgReq.SetAtPath("credentials.clientSecret.secret", "FAKE_SECRET")

	if err := azureKvPrvdr.Configure(cfgReq); err != nil {
		t.Errorf("failed to configure provider, reason: %s", err)
	}

	queryRequest := sdk.NewProviderDataRequest(azureKvPrvdr.ProvisionDataRequestSchema())
	queryRequest.SetAtPath("secretNames", []string{"apiEndpoint", "AppConfig"})
	data, diags := azureKvPrvdr.Provide(queryRequest)
	if data == nil || diags == nil {
		t.Fail()
	}

	if diags.HasErrors() {
		t.Errorf("Failed to retrieve data, reason: %s", diags.Entries().String())
	}

	data.PrintTree()
}
