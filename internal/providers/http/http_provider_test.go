package http

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	sdk "github.com/conformize/conformize/internal/providers/api"
)

func TestHttpProviderConfiguration(t *testing.T) {
	httpPrvdr := &HttpProvider{}
	cfgReq := sdk.NewConfigurationRequest(httpPrvdr.ConfigurationSchema())
	cfgReq.SetAtPath("endpoint", "http://localhost")

	err := httpPrvdr.Configure(cfgReq)
	if err != nil {
		t.Fatalf("configuration failed: %v", err)
	}
}

func TestHttpProviderProvideWithMockServer(t *testing.T) {
	mockResponse, err := os.ReadFile("../../../mocks/app-dev.json")
	if err != nil {
		t.Fatalf("failed to read mock response file: %v", err)
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write(mockResponse)
	}))
	defer server.Close()

	httpPrvdr := &HttpProvider{}
	cfgReq := sdk.NewConfigurationRequest(httpPrvdr.ConfigurationSchema())
	cfgReq.SetAtPath("endpoint", server.URL)

	if err := httpPrvdr.Configure(cfgReq); err != nil {
		t.Fatalf("failed to configure http provider: %v", err)
	}

	queryReq := sdk.NewProviderDataRequest(httpPrvdr.ProvisionDataRequestSchema())
	queryReq.SetAtPath("method", "GET")
	data, diags := httpPrvdr.Provide(queryReq)
	if data == nil || diags == nil {
		t.Fatalf("expected non-nil data and diagnostics")
	}
	if diags.HasErrors() {
		t.Errorf("diagnostics has errors: %s", diags.Errors().String())
	}
}

func TestHttpProviderPostWithJsonStringBody(t *testing.T) {
	mockResponse, _ := os.ReadFile("../../../mocks/app-dev.json")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST request, got %s", r.Method)
		}
		body, _ := io.ReadAll(r.Body)
		defer r.Body.Close()
		if strings.Compare(string(body), "{\"request\":\"data\"}") != 0 {
			t.Errorf("expected request body to match JSON, got: %s", string(body))
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(mockResponse)
	}))
	defer server.Close()

	httpPrvdr := &HttpProvider{}
	cfgReq := sdk.NewConfigurationRequest(httpPrvdr.ConfigurationSchema())
	cfgReq.SetAtPath("endpoint", server.URL)

	if err := httpPrvdr.Configure(cfgReq); err != nil {
		t.Fatalf("failed to configure http provider: %v", err)
	}

	queryReq := sdk.NewProviderDataRequest(httpPrvdr.ProvisionDataRequestSchema())
	queryReq.SetAtPath("method", "POST")
	queryReq.SetAtPath("body", `{"request":"data"}`)

	data, diags := httpPrvdr.Provide(queryReq)
	if diags.HasErrors() {
		t.Errorf("diagnostics has errors: %s", diags.Errors().String())
	}
	if data == nil {
		t.Fatal("expected data node")
	}
}

func TestHttpProviderPostWithJsonStringBodyAndHeaders(t *testing.T) {
	mockResponse, _ := os.ReadFile("../../../mocks/app-dev.json")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Test-Header") != "123" {
			t.Errorf("missing or incorrect X-Test-Header")
		}
		if r.Method != http.MethodPost {
			t.Errorf("expected POST request, got %s", r.Method)
		}
		var m map[string]any
		if err := json.NewDecoder(r.Body).Decode(&m); err != nil || m["request"] != "mapdata" {
			t.Errorf("unexpected request body: %+v", m)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(mockResponse)
	}))
	defer server.Close()

	httpPrvdr := &HttpProvider{}
	cfgReq := sdk.NewConfigurationRequest(httpPrvdr.ConfigurationSchema())
	cfgReq.SetAtPath("endpoint", server.URL)

	if err := httpPrvdr.Configure(cfgReq); err != nil {
		t.Fatalf("failed to configure http provider: %v", err)
	}

	queryReq := sdk.NewProviderDataRequest(httpPrvdr.ProvisionDataRequestSchema())
	queryReq.SetAtPath("method", "POST")
	queryReq.SetAtPath("headers", map[string]string{"X-Test-Header": "123"})
	queryReq.SetAtPath("body", "{\"request\": \"mapdata\"}")
	data, diags := httpPrvdr.Provide(queryReq)
	if diags.HasErrors() {
		t.Errorf("diagnostics has errors: %s", diags.Errors().String())
	}
	if data == nil {
		t.Fatal("expected data node")
	}
}

func TestHttpProviderQueryParams(t *testing.T) {
	mockResponse := []byte(`{"query":"ok"}`)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("env") != "dev" {
			t.Errorf("missing or wrong query param 'env'")
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(mockResponse)
	}))
	defer server.Close()

	httpPrvdr := &HttpProvider{}
	cfgReq := sdk.NewConfigurationRequest(httpPrvdr.ConfigurationSchema())
	cfgReq.SetAtPath("endpoint", server.URL)

	if err := httpPrvdr.Configure(cfgReq); err != nil {
		t.Fatalf("config failed: %v", err)
	}

	queryReq := sdk.NewProviderDataRequest(httpPrvdr.ProvisionDataRequestSchema())
	queryReq.SetAtPath("method", "GET")
	queryReq.SetAtPath("queryParams", map[string]string{"env": "dev"})
	data, diags := httpPrvdr.Provide(queryReq)
	if diags.HasErrors() {
		t.Errorf("unexpected errors: %s", diags.Errors().String())
	}
	if data == nil {
		t.Fatal("expected data node")
	}
}

func TestHttpProviderServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "error", http.StatusInternalServerError)
	}))
	defer server.Close()

	httpPrvdr := &HttpProvider{}
	cfgReq := sdk.NewConfigurationRequest(httpPrvdr.ConfigurationSchema())
	cfgReq.SetAtPath("endpoint", server.URL)
	cfgReq.SetAtPath("method", "GET")

	if err := httpPrvdr.Configure(cfgReq); err != nil {
		t.Fatalf("config failed: %v", err)
	}

	data, diags := httpPrvdr.Provide(nil)
	if data != nil {
		t.Errorf("expected nil data on server error")
	}
	if !diags.HasErrors() {
		t.Errorf("expected error diagnostics from server 500")
	}
}

func TestHttpProviderWithAuthHeaders(t *testing.T) {
	mockResponse := []byte(`{"auth":"ok"}`)

	tests := []struct {
		name        string
		authConfig  map[string]any
		expectKey   string
		expectValue string
	}{
		{
			name: "basic auth",
			authConfig: map[string]any{
				"basic": map[string]any{
					"username": "testuser",
					"password": "testpass",
				},
			},
			expectKey:   "Authorization",
			expectValue: "Basic",
		},
		{
			name: "bearer auth",
			authConfig: map[string]any{
				"bearer": map[string]any{
					"token": "sometoken",
				},
			},
			expectKey:   "Authorization",
			expectValue: "Bearer sometoken",
		},
		{
			name: "api key",
			authConfig: map[string]any{
				"apiKey": map[string]any{
					"key":   "X-API-KEY",
					"value": "my-secret",
				},
			},
			expectKey:   "X-API-KEY",
			expectValue: "my-secret",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if got := r.Header.Get(tc.expectKey); !strings.Contains(got, tc.expectValue) {
					t.Errorf("expected header %s to contain %q, got %q", tc.expectKey, tc.expectValue, got)
				}
				w.Header().Set("Content-Type", "application/json")
				w.Write(mockResponse)
			}))
			defer server.Close()

			httpPrvdr := &HttpProvider{}
			cfgReq := sdk.NewConfigurationRequest(httpPrvdr.ConfigurationSchema())
			cfgReq.SetAtPath("endpoint", server.URL)
			cfgReq.SetAtPath("auth", tc.authConfig)

			if err := httpPrvdr.Configure(cfgReq); err != nil {
				t.Fatalf("configuration failed: %v", err)
			}

			queryReq := sdk.NewProviderDataRequest(httpPrvdr.ProvisionDataRequestSchema())
			queryReq.SetAtPath("method", "GET")
			data, diags := httpPrvdr.Provide(queryReq)
			if diags.HasErrors() {
				t.Errorf("unexpected diagnostics: %s", diags.Errors().String())
			}
			if data == nil {
				t.Fatal("expected non-nil data")
			}
		})
	}
}
