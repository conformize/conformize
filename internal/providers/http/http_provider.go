// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package http

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/conformize/conformize/common/diagnostics"
	"github.com/conformize/conformize/common/ds"
	"github.com/conformize/conformize/common/typed"
	sdk "github.com/conformize/conformize/internal/providers/api"
	"github.com/conformize/conformize/internal/providers/api/schema"
	"github.com/conformize/conformize/internal/providers/api/schema/attributes"
	"github.com/conformize/conformize/serialization"
	"github.com/conformize/conformize/serialization/unmarshal/yaml"
)

type httpClientBasicAuth struct {
	Username string `cnfrmz:"username"`
	Password string `cnfrmz:"password"`
}

type httpClientApiKeyAuth struct {
	Key   string `cnfrmz:"key"`
	Value string `cnfrmz:"value"`
}

type httpClientAuth struct {
	Basic  *httpClientBasicAuth  `cnfrmz:"basic"`
	Bearer *string               `cnfrmz:"bearer"`
	ApiKey *httpClientApiKeyAuth `cnfrmz:"apiKey"`
}

type httpClientTls struct {
	SkipVerify bool     `cnfrmz:"skipVerify"`
	CACert     string   `cnfrmz:"caCert"`
	ClientCert string   `cnfrmz:"clientCert"`
	ClientKey  string   `cnfrmz:"clientKey"`
	Ciphers    []string `cnfrmz:"ciphers"`
}

type httpClientOptions struct {
	FollowRedirects bool           `cnfrmz:"followRedirects"`
	Timeout         int            `cnfrmz:"timeout"`
	TLS             *httpClientTls `cnfrmz:"tls"`
}

type httpClientConfig struct {
	Endpoint string             `cnfrmz:"endpoint"`
	Auth     *httpClientAuth    `cnfrmz:"auth"`
	Options  *httpClientOptions `cnfrmz:"options"`
}

type httpClientRequest struct {
	Headers     map[string]string `cnfrmz:"headers"`
	Method      string            `cnfrmz:"method"`
	Body        string            `cnfrmz:"body"`
	QueryParams map[string]string `cnfrmz:"queryParams"`
}

type HttpProvider struct {
	endpoint   string
	client     *http.Client
	clientAuth *httpClientAuth
}

func (httpProvider *HttpProvider) Provide(req *sdk.ProviderDataRequest) (*ds.Node[string, interface{}], *diagnostics.Diagnostics) {
	diag := diagnostics.NewDiagnostics()

	url, err := url.Parse(httpProvider.endpoint)
	if err != nil {
		diag.Append(diagnostics.Builder().Error().Details(fmt.Sprintf("invalid endpoint URL: %s", err.Error())).Build())
		return nil, diag
	}

	var request httpClientRequest
	if req == nil {
		diag.Append(diagnostics.Builder().Error().Details("request is nil").Build())
		return nil, diag
	}

	if err := req.Get(&request); err != nil {
		diag.Append(diagnostics.Builder().Error().Details(fmt.Sprintf("failed to get request data: %s", err.Error())).Build())
		return nil, diag
	}

	query := url.Query()
	for k, v := range request.QueryParams {
		query.Set(k, v)
	}
	url.RawQuery = query.Encode()

	var body io.Reader
	if request.Body != "" {
		body = bytes.NewReader([]byte(request.Body))
	}

	reqHttp, err := http.NewRequest(request.Method, url.String(), body)
	if err != nil {
		diag.Append(diagnostics.Builder().Error().Details(fmt.Sprintf("couldn't create HTTP request: %s", err)).Build())
		return nil, diag
	}

	for k, v := range request.Headers {
		reqHttp.Header.Set(k, v)
	}

	setUpAuth(reqHttp, httpProvider.clientAuth)

	resp, err := httpProvider.client.Do(reqHttp)
	if resp.StatusCode > 299 {
		diag.Append(diagnostics.Builder().Error().Details(fmt.Sprintf("HTTP request failed with status code: %d", resp.StatusCode)).Build())
		return nil, diag
	}

	if err != nil {
		diag.Append(diagnostics.Builder().Error().Details(fmt.Sprintf("HTTP request failed: %s", err.Error())).Build())
		return nil, diag
	}
	defer resp.Body.Close()

	jsonUnmarshaller := yaml.YamlUnmarshal{}
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, resp.Body); err != nil {
		diag.Append(diagnostics.Builder().Error().Details(fmt.Sprintf("couldn't read response body: %s", err.Error())).Build())
		return nil, diag
	}

	if data, err := jsonUnmarshaller.Unmarshal(serialization.NewBufferedData(buf.Bytes())); err != nil {
		diag.Append(diagnostics.Builder().Error().Details(fmt.Sprintf("couldn't decode response: %s", err.Error())).Build())
		return nil, diag
	} else {
		return data, diag
	}
}

func (httpProvider *HttpProvider) ProvisionDataRequestSchema() *schema.Schema {
	return &schema.Schema{
		Description: "Http resource request schema",
		Version:     1,
		Attributes: map[string]schema.Attributeable{
			"headers":     &attributes.MapAttribute{ElementsType: &typed.StringTyped{}},
			"method":      &attributes.StringAttribute{},
			"body":        &attributes.StringAttribute{},
			"queryParams": &attributes.MapAttribute{ElementsType: &typed.StringTyped{}},
		},
	}
}

func (httpProvider *HttpProvider) ConfigurationSchema() *schema.Schema {
	return &schema.Schema{
		Description: "File provider schema",
		Version:     1,
		Attributes: map[string]schema.Attributeable{
			"endpoint": &attributes.StringAttribute{},
			"auth": &attributes.ObjectAttribute{
				FieldsTypes: map[string]typed.Typeable{
					"basic": &typed.ObjectTyped{
						FieldsTypes: map[string]typed.Typeable{
							"username": &typed.StringTyped{},
							"password": &typed.StringTyped{},
						},
					},
					"bearer": &typed.StringTyped{},
					"apiKey": &typed.ObjectTyped{
						FieldsTypes: map[string]typed.Typeable{
							"key":   &typed.StringTyped{},
							"value": &typed.StringTyped{},
						},
					},
				},
			},
			"options": &attributes.ObjectAttribute{
				FieldsTypes: map[string]typed.Typeable{
					"followRedirects": &typed.BooleanTyped{},
					"timeout":         &typed.NumberTyped{},
					"tls": &typed.ObjectTyped{
						FieldsTypes: map[string]typed.Typeable{
							"skipVerify": &typed.BooleanTyped{},
							"caCert":     &typed.StringTyped{},
							"clientCert": &typed.StringTyped{},
							"clientKey":  &typed.StringTyped{},
							"ciphers":    &typed.ListTyped{ElementsType: &typed.StringTyped{}},
						},
					},
				},
			},
		},
	}
}

func (httpProvider *HttpProvider) Configure(req *sdk.ConfigurationRequest) error {
	var clientConfig httpClientConfig
	if err := req.Get(&clientConfig); err != nil {
		return err
	}
	httpProvider.endpoint = clientConfig.Endpoint
	httpProvider.clientAuth = clientConfig.Auth

	transport := &http.Transport{}
	if clientConfig.Options != nil && clientConfig.Options.TLS != nil {
		tlsOptions := clientConfig.Options.TLS
		tlsConfig := &tls.Config{
			InsecureSkipVerify: tlsOptions.SkipVerify,
		}

		if tlsOptions.CACert != "" {
			caCert, err := os.ReadFile(tlsOptions.CACert)
			if err != nil {
				return fmt.Errorf("failed to read CA cert: %w", err)
			}
			caPool := x509.NewCertPool()
			caPool.AppendCertsFromPEM(caCert)
			tlsConfig.RootCAs = caPool
		}

		if tlsOptions.ClientCert != "" && tlsOptions.ClientKey != "" {
			cert, err := tls.LoadX509KeyPair(tlsOptions.ClientCert, tlsOptions.ClientKey)
			if err != nil {
				return fmt.Errorf("failed to load client cert/key: %w", err)
			}
			tlsConfig.Certificates = []tls.Certificate{cert}
		}

		transport.TLSClientConfig = tlsConfig
	}

	client := &http.Client{
		Transport: transport,
	}

	if clientConfig.Options != nil && !clientConfig.Options.FollowRedirects {
		client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}
	}

	if clientConfig.Options != nil && clientConfig.Options.Timeout > 0 {
		client.Timeout = time.Duration(clientConfig.Options.Timeout) * time.Second
	}

	httpProvider.client = client
	return nil
}

func setUpAuth(req *http.Request, auth *httpClientAuth) {
	if auth == nil {
		return
	}

	if auth.Basic != nil {
		req.SetBasicAuth(auth.Basic.Username, auth.Basic.Password)
		return
	}

	if auth.Bearer != nil {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", *auth.Bearer))
		return
	}

	if auth.ApiKey != nil {
		req.Header.Set(auth.ApiKey.Key, auth.ApiKey.Value)
		return
	}

	return
}
