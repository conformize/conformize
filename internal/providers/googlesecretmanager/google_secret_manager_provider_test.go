package googlesecretmanager

import (
	"testing"

	sdk "github.com/conformize/conformize/internal/providers/api"
)

func TestGoogleSecretManagerProviderConfigurationWithCredentialsFile(t *testing.T) {
	googleSecretManagerPrvdr := GoogleSecretManagerProvider{}
	cfgReq := sdk.NewConfigurationRequest(googleSecretManagerPrvdr.ConfigurationSchema())
	cfgReq.SetAtPath("project", "FAKE_PROJECT")
	cfgReq.SetAtPath("credentials.credentialsFile", "FAKE_CREDS_FILE")

	err := googleSecretManagerPrvdr.Configure(cfgReq)
	if err != nil {
		t.Errorf("failed to configure google secret manager provider, reason: %s", err.Error())
	}
}

func TestGoogleSecretManagerProvideSecretsByName(t *testing.T) {
	googleSecretManagerPrvdr := GoogleSecretManagerProvider{}
	cfgReq := sdk.NewConfigurationRequest(googleSecretManagerPrvdr.ConfigurationSchema())
	cfgReq.SetAtPath("project", "FAKE_PROJECT")
	cfgReq.SetAtPath("credentials.credentialsFile", "FAKE_CREDS_FILE")

	googleSecretManagerPrvdr.Configure(cfgReq)

	queryRequest := sdk.NewProviderDataRequest(googleSecretManagerPrvdr.ProvisionDataRequestSchema())
	queryRequest.SetAtPath("secretNames", []string{"apiConfig", "appConfig"})

	data, diags := googleSecretManagerPrvdr.Provide(queryRequest)
	if data == nil || diags == nil {
		t.Fail()
	}

	if diags.HasErrors() {
		t.Errorf("Failed to retrieve data, reason: %s", diags.Entries().String())
	}

	data.PrintTree()
}

func TestGoogleSecretManagerProvideSecretsByFilter(t *testing.T) {
	googleSecretManagerPrvdr := GoogleSecretManagerProvider{}
	cfgReq := sdk.NewConfigurationRequest(googleSecretManagerPrvdr.ConfigurationSchema())
	cfgReq.SetAtPath("project", "FAKE_PROJECT")
	cfgReq.SetAtPath("credentials.credentialsFile", "FAKE_CREDS_FILE")

	googleSecretManagerPrvdr.Configure(cfgReq)

	queryRequest := sdk.NewProviderDataRequest(googleSecretManagerPrvdr.ProvisionDataRequestSchema())
	queryRequest.SetAtPath("filter", "labels.env=stage")

	data, diags := googleSecretManagerPrvdr.Provide(queryRequest)
	if data == nil || diags == nil {
		t.Fail()
	}

	if diags.HasErrors() {
		t.Errorf("Failed to retrieve data, reason: %s", diags.Entries().String())
	}

	data.PrintTree()
}
