// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package secretsmanager

import (
	"testing"

	sdk "github.com/conformize/conformize/internal/providers/api"
)

func awsSecretsManagerConfigurationRequest() *sdk.ConfigurationRequest {
	awsSecretsManagerProvider := New("awsSecretsManager")
	configurationRequest := sdk.NewConfigurationRequest(awsSecretsManagerProvider.ConfigurationSchema())
	configurationRequest.SetAtPath("region", "us-east-1")
	configurationRequest.SetAtPath("endpointUrl", "http://localhost:4566")
	return configurationRequest
}

func awsSecretsManagerGetSecretValueRequestOptions() *sdk.ProviderDataRequest {
	awsSecretsManagerProvider := New("awsSecretsManager")
	queryRequest := sdk.NewProviderDataRequest(awsSecretsManagerProvider.ProvisionDataRequestSchema())
	queryRequest.SetAtPath("secretIdList", []string{"app/config/db/username"})
	return queryRequest
}

func awsSecretsManagerSecretValuesBatchRequestOptions() *sdk.ProviderDataRequest {
	queryFilters := []map[string][]string{
		{"description": {"Database"}},
		{"description": {"Database"}},
	}

	awsSecretsManagerProvider := New("awsSecretsManager")
	queryRequest := sdk.NewProviderDataRequest(awsSecretsManagerProvider.ProvisionDataRequestSchema())
	queryRequest.SetAtPath("maxResults", 10)
	queryRequest.SetAtPath("filters", queryFilters)
	return queryRequest
}

func TestAwsSecretsManagerGetSecretValue(t *testing.T) {
	awsSecretsManagerProvider := New("awsSecretsManager")
	configurationRequest := awsSecretsManagerConfigurationRequest()
	if err := awsSecretsManagerProvider.Configure(configurationRequest); err != nil {
		t.Fatalf("Failed to configure AWS Secrets Manager provider: %v", err)
	}
	queryRequest := awsSecretsManagerGetSecretValueRequestOptions()
	if data, diags := awsSecretsManagerProvider.Provide(queryRequest); diags.HasErrors() {
		t.Fatalf("Failed to query AWS Secrets Manager: %s", diags.Errors())
	} else {
		data.PrintTree()
	}
}

func TestAwsSecretsManagerGetSecretValuesBatch(t *testing.T) {
	awsSecretsManagerProvider := New("awsSecretsManager")
	configurationRequest := awsSecretsManagerConfigurationRequest()
	if err := awsSecretsManagerProvider.Configure(configurationRequest); err != nil {
		t.Fatalf("Failed to configure AWS Secrets Manager provider: %v", err)
	}
	queryRequest := awsSecretsManagerSecretValuesBatchRequestOptions()
	if data, diags := awsSecretsManagerProvider.Provide(queryRequest); diags.HasErrors() {
		t.Fatalf("Failed to query AWS Secrets Manager: %v", diags.Errors())
	} else {
		data.PrintTree()
	}
}
