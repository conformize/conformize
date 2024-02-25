// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package awsparameterstore

import (
	"testing"

	sdk "github.com/conformize/conformize/internal/providers/api"
)

func awsParameterStoreConfigurationRequest() *sdk.ConfigurationRequest {
	awsParameterStoreProvider := New("awsParameterStore")
	configurationRequest := sdk.NewConfigurationRequest(awsParameterStoreProvider.ConfigurationSchema())
	configurationRequest.SetAtPath("region", "us-east-1")
	configurationRequest.SetAtPath("endpointUrl", "http://localhost:4566")
	return configurationRequest
}

func awsParameterStoreQueryByNameRequestOptions() *sdk.ProviderDataRequest {
	awsParameterStoreProvider := New("awsParameterStore")
	queryRequest := sdk.NewProviderDataRequest(awsParameterStoreProvider.ProvisionDataRequestSchema())
	queryRequest.SetAtPath("path", "/app/config/db/host")
	return queryRequest
}

func awsParameterStoreQueryByPathRequestOptions() *sdk.ProviderDataRequest {
	awsParameterStoreProvider := New("awsParameterStore")
	queryRequest := sdk.NewProviderDataRequest(awsParameterStoreProvider.ProvisionDataRequestSchema())
	queryRequest.SetAtPath("path", "/app")
	queryRequest.SetAtPath("recursive", true)
	queryRequest.SetAtPath("maxResults", 1)
	return queryRequest
}

func TestAwsParameterStoreQueryByName(t *testing.T) {
	awsParameterStoreProvider := New("awsSecretsManager")
	configurationRequest := awsParameterStoreConfigurationRequest()
	if err := awsParameterStoreProvider.Configure(configurationRequest); err != nil {
		t.Fatalf("Failed to configure AWS Parameter Store provider: %v", err)
	}
	queryRequest := awsParameterStoreQueryByNameRequestOptions()
	if data, diags := awsParameterStoreProvider.Provide(queryRequest); diags.HasErrors() {
		t.Fatalf("Failed to query AWS Parameter Store: %s", diags.Errors())
	} else {
		data.PrintTree()
	}
}

func TestAwsParameterStoreQueryByPath(t *testing.T) {
	awsParameterStoreProvider := New("awsParameterStore")
	configurationRequest := awsParameterStoreConfigurationRequest()
	if err := awsParameterStoreProvider.Configure(configurationRequest); err != nil {
		t.Fatalf("Failed to configure AWS Parameter Store provider: %v", err)
	}
	queryRequest := awsParameterStoreQueryByPathRequestOptions()
	if data, diags := awsParameterStoreProvider.Provide(queryRequest); diags.HasErrors() {
		t.Fatalf("Failed to query AWS Parameter Store: %s", diags.Errors())
	} else {
		data.PrintTree()
	}
}
