package env

import (
	"testing"

	sdk "github.com/conformize/conformize/internal/providers/api"
)

func TestEnvProviderConfiguration(t *testing.T) {
	envPrvdr := &EnvProvider{}
	cfgReq := sdk.NewConfigurationRequest(envPrvdr.ConfigurationSchema())
	err := envPrvdr.Configure(cfgReq)
	if err != nil {
		t.Fail()
	}
}

func TestEnvProviderProvideEnvironmentVariables(t *testing.T) {
	envPrvdr := &EnvProvider{}
	data, diags := envPrvdr.Provide(nil)
	if data == nil || diags == nil {
		t.Fail()
	}

	if diags.HasErrors() {
		t.Errorf("Failed to retrieve data, reason: %s", diags.Entries().String())
	}

	if _, ok := data.GetChildren("PATH"); !ok {
		t.Fail()
	}
}

func TestEnvProviderProvideEnvironmentVariablesWithLookup(t *testing.T) {
	envPrvdr := &EnvProvider{}

	vars := []string{"USER", "PATH"}
	queryReq := sdk.NewProviderDataRequest(envPrvdr.ProvisionDataRequestSchema())
	queryReq.SetAtPath("vars", vars)
	data, diags := envPrvdr.Provide(nil)
	if data == nil || diags == nil {
		t.Fail()
	}

	if diags.HasErrors() {
		t.Errorf("Failed to retrieve data, reason: %s", diags.Entries().String())
	}

	for _, envVar := range vars {
		if _, ok := data.GetChildren(envVar); !ok {
			t.Fail()
		}
	}
}
