// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package azuredevopsvariablegroup

import (
	"testing"

	sdk "github.com/conformize/conformize/internal/providers/api"
)

func TestAzureDevOpsVariableGroupProviderConfiguration(t *testing.T) {
	prvdr := &AzureDevOpsVariableGroupProvider{}
	cfgReq := sdk.NewConfigurationRequest(prvdr.ConfigurationSchema())
	cfgReq.SetAtPath("organizationUrl", "FAKE_ORG")
	cfgReq.SetAtPath("project", "FAKE_PROJECT")
	cfgReq.SetAtPath("groupId", 1)
	cfgReq.SetAtPath("token", "FAKE_TOKEN")

	err := prvdr.Configure(cfgReq)
	if err != nil {
		t.Errorf("failed to configure azure devops vairable group provider, reason: %s", err.Error())
	}
}

func TestAzureDevOpsVariableGroupProviderProvideGroupVariables(t *testing.T) {
	prvdr := &AzureDevOpsVariableGroupProvider{}
	cfgReq := sdk.NewConfigurationRequest(prvdr.ConfigurationSchema())
	cfgReq.SetAtPath("organizationUrl", "FAKE_ORG")
	cfgReq.SetAtPath("project", "FAKE_PROJECT")
	cfgReq.SetAtPath("groupId", 1)
	cfgReq.SetAtPath("token", "FAKE_TOKEN")

	if err := prvdr.Configure(cfgReq); err != nil {
		t.Errorf("failed to configure azure devops vairable group provider, reason: %s", err.Error())
	}

	data, diags := prvdr.Provide(nil)
	if data == nil || diags == nil {
		t.Fail()
	}

	if diags.HasErrors() {
		t.Errorf("Failed to retrieve data, reason: %s", diags.Entries().String())
	}

	data.PrintTree()
}
