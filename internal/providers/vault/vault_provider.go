// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package vault

import (
	"context"
	"fmt"
	"runtime"
	"sync"

	"github.com/conformize/conformize/common/diagnostics"
	"github.com/conformize/conformize/common/ds"
	"github.com/conformize/conformize/common/path"
	"github.com/conformize/conformize/common/typed"
	sdk "github.com/conformize/conformize/internal/providers/api"
	"github.com/conformize/conformize/internal/providers/api/schema"
	"github.com/conformize/conformize/internal/providers/api/schema/attributes"

	"github.com/hashicorp/vault/api"
	appRoleAuth "github.com/hashicorp/vault/api/auth/approle"
	kubernetesAuth "github.com/hashicorp/vault/api/auth/kubernetes"
	userPassAuth "github.com/hashicorp/vault/api/auth/userpass"
)

type VaultProvider struct {
	alias            string
	secretsMountPath string
	client           *api.Client
}

type vaultClientConfig struct {
	Address   string                 `cnfrmz:"address"`
	MountPath string                 `cnfrmz:"mountPath"`
	Auth      *vaultClientAuthConfig `cnfrmz:"authentication"`
}

type vaultClientAuthConfig struct {
	AppRole    *vaultAppRoleAuthConfig          `cnfrmz:"appRole"`
	Github     *tokenAuthConfig                 `cnfrmz:"github"`
	Kubernetes *vaultKubernetesAuthConfig       `cnfrmz:"kubernetes"`
	UserPass   *vaultUsernamePasswordAuthConfig `cnfrmz:"userPass"`
	Token      *tokenAuthConfig                 `cnfrmz:"token"`
}

type vaultAppRoleAuthConfig struct {
	Path     string `cnfrmz:"path"`
	RoleId   string `cnfrmz:"roleId"`
	SecretId string `cnfrmz:"secretId"`
}

type tokenAuthConfig struct {
	Path  string `cnfrmz:"path"`
	Token string `cnfrmz:"token"`
}

type vaultKubernetesAuthConfig struct {
	Path  string `cnfrmz:"path"`
	Role  string `cnfrmz:"role"`
	Token string `cnfrmz:"token"`
}

type vaultUsernamePasswordAuthConfig struct {
	Path     string `cnfrmz:"path"`
	Username string `cnfrmz:"username"`
	Password string `cnfrmz:"password"`
}

type queryOptions struct {
	Paths []string `cnfrmz:"paths"`
}

const maxSecretsBatchSize = 10

func (vaultPrvdr *VaultProvider) ConfigurationSchema() *schema.Schema {
	return &schema.Schema{
		Description: "Configuration for the Vault provider",
		Version:     1,
		Attributes: map[string]schema.Attributeable{
			"address":   &attributes.StringAttribute{},
			"mountPath": &attributes.StringAttribute{},
			"auth": &attributes.ObjectAttribute{
				FieldsTypes: map[string]typed.Typeable{
					"appRole": &typed.ObjectTyped{
						FieldsTypes: map[string]typed.Typeable{
							"path":     &typed.StringTyped{},
							"roleId":   &typed.StringTyped{},
							"secretId": &typed.StringTyped{},
						},
					},
					"github": &typed.ObjectTyped{
						FieldsTypes: map[string]typed.Typeable{
							"path":  &typed.StringTyped{},
							"token": &typed.StringTyped{},
						},
					},
					"kubernetes": &typed.ObjectTyped{
						FieldsTypes: map[string]typed.Typeable{
							"path": &typed.StringTyped{},
							"role": &typed.StringTyped{},
						},
					},
					"userPass": &typed.ObjectTyped{
						FieldsTypes: map[string]typed.Typeable{
							"path":     &typed.StringTyped{},
							"username": &typed.StringTyped{},
							"password": &typed.StringTyped{},
						},
					},
					"token": &typed.ObjectTyped{
						FieldsTypes: map[string]typed.Typeable{
							"path":  &typed.StringTyped{},
							"token": &typed.StringTyped{},
						},
					},
				},
			},
		},
	}
}

func (vaultPrvdr *VaultProvider) ProvisionDataRequestSchema() *schema.Schema {
	return &schema.Schema{
		Version:     1,
		Description: "Vault resource request schema",
		Attributes: map[string]schema.Attributeable{
			"paths": &attributes.ListAttribute{ElementsType: &typed.StringTyped{}},
		},
	}
}

func (vaultPrvdr *VaultProvider) Configure(req *sdk.ConfigurationRequest) error {
	var clientConfig vaultClientConfig
	if err := req.Get(&clientConfig); err != nil {
		return err
	}
	vaultPrvdr.secretsMountPath = clientConfig.MountPath

	config := api.DefaultConfig()
	config.Address = clientConfig.Address
	if client, err := api.NewClient(config); err == nil {
		vaultPrvdr.client = client
	} else {
		return err
	}

	if clientConfig.Auth != nil {
		return vaultPrvdr.authenticate(clientConfig.Auth)
	}
	return nil
}

type secret struct {
	Path *string
	Data map[string]any
}

func (vaultPrvdr *VaultProvider) Provide(queryRequest *sdk.ProviderDataRequest) (*ds.Node[string, any], *diagnostics.Diagnostics) {
	diags := diagnostics.NewDiagnostics()
	var queryOptions queryOptions
	if err := queryRequest.Get(&queryOptions); err != nil {
		diags.Append(diagnostics.Builder().Error().Details(err.Error()).Build())
		return nil, diags
	}

	pathsLen := len(queryOptions.Paths)
	batchesCount := (pathsLen + maxSecretsBatchSize - 1) / maxSecretsBatchSize

	errChan := make(chan error)
	valsChan := make(chan []*secret, batchesCount)
	doneChan := make(chan struct{})

	defer close(errChan)
	defer close(valsChan)
	defer close(doneChan)

	var result = ds.NewNode[string, any]()
	go func() {
		done := false
		for !done {
			select {
			case err := <-errChan:
				diags.Append(diagnostics.Builder().Error().Details(err.Error()).Build())
			case vals := <-valsChan:
				for _, val := range vals {
					valuePath, _ := path.NewFromStringWithSeparator(*val.Path, '/')
					steps := valuePath.Steps()

					lastNodeRef := result
					for step, hasNext := steps.Next(); hasNext; step, hasNext = steps.Next() {
						stepName := step.String()
						if nodes, found := lastNodeRef.GetChildren(stepName); found {
							lastNodeRef = nodes.First()
						} else {
							lastNodeRef = lastNodeRef.AddChild(stepName)
						}
					}

					if data, ok := val.Data["data"]; ok {
						secretData := data.(map[string]any)
						for key, value := range secretData {
							if key != "metadata" {
								childNode := lastNodeRef.AddChild(key)
								childNode.Value = value
							}
						}
					}

					if metadata, ok := val.Data["metadata"]; ok {
						if metadataMap, ok := metadata.(map[string]any); ok {
							if customMetadata, ok := metadataMap["custom_metadata"]; ok {
								if meta, ok := customMetadata.(map[string]any); ok {
									for key, value := range meta {
										lastNodeRef.AddAttribute(key, value)
									}
								}
							}
						}
					}
				}
			case <-doneChan:
				done = true
			}
		}
	}()
	var wg sync.WaitGroup
	wg.Add(batchesCount)

	availableCPUs := max(1, runtime.NumCPU()-1)
	cpus := runtime.GOMAXPROCS(availableCPUs)
	defer runtime.GOMAXPROCS(cpus)

	maxParallelTasksCount := min(availableCPUs, max(1, batchesCount))
	tasks := make(chan struct{}, maxParallelTasksCount)
	defer close(tasks)

	for i, offset := 0, 0; i <= batchesCount && offset < pathsLen; i, offset = i+1, offset+maxSecretsBatchSize {
		upperBound := min(offset+maxSecretsBatchSize, pathsLen)

		paths := queryOptions.Paths[offset:upperBound]
		tasks <- struct{}{}
		go func(paths []string) {
			defer wg.Done()

			vals := make([]*secret, 0)
			for _, queryPath := range paths {
				var secretPath = vaultPrvdr.secretsMountPath + queryPath
				if secretVal, err := vaultPrvdr.client.Logical().ReadWithContext(context.Background(), secretPath); err != nil {
					errChan <- err
				} else if secretVal != nil {
					vals = append(vals, &secret{Path: &queryPath, Data: secretVal.Data})
				} else {
					errChan <- fmt.Errorf("secret at path %s not found", secretPath)
				}
			}
			valsChan <- vals
			<-tasks
		}(paths)
	}
	wg.Wait()
	doneChan <- struct{}{}
	return result, diags
}

func (vaultPrvdr *VaultProvider) authenticate(authConfig *vaultClientAuthConfig) error {
	var authMethod api.AuthMethod
	var authError error
	if authConfig.AppRole != nil {
		var opts []appRoleAuth.LoginOption
		if authConfig.AppRole.Path != "" {
			opts = append(opts, appRoleAuth.WithMountPath(authConfig.AppRole.Path))
		}

		authMethod, authError = appRoleAuth.NewAppRoleAuth(
			authConfig.AppRole.Path,
			&appRoleAuth.SecretID{FromString: authConfig.AppRole.SecretId},
			opts...,
		)
	}

	if authConfig.UserPass != nil {
		var opts []userPassAuth.LoginOption
		if authConfig.UserPass.Path != "" {
			opts = append(opts, userPassAuth.WithMountPath(authConfig.UserPass.Path))
		}

		authMethod, authError = userPassAuth.NewUserpassAuth(
			authConfig.UserPass.Username,
			&userPassAuth.Password{FromString: authConfig.UserPass.Password},
			opts...,
		)
	}

	if authConfig.Kubernetes != nil {
		var opts []kubernetesAuth.LoginOption
		if authConfig.Kubernetes.Path != "" {
			opts = append(opts, kubernetesAuth.WithMountPath(authConfig.Kubernetes.Path))
		}

		if authConfig.Kubernetes.Token != "" {
			opts = append(opts, kubernetesAuth.WithServiceAccountToken(authConfig.Kubernetes.Token))
		}

		authMethod, authError = kubernetesAuth.NewKubernetesAuth(
			authConfig.Kubernetes.Role,
			opts...,
		)
	}

	if authError == nil {
		if authMethod != nil {
			_, err := vaultPrvdr.client.Auth().Login(context.Background(), authMethod)
			authError = err
		} else {
			var loginPath string
			var loginData = map[string]any{}
			if authConfig.Github != nil {
				loginPath = "auth/github/login"

				if authConfig.Github.Path != "" {
					loginPath = authConfig.Github.Path
				}
				loginData["token"] = authConfig.Github.Token
			}

			if authConfig.Token != nil {
				loginPath = "auth/token/login"

				if authConfig.Token.Path != "" {
					loginPath = authConfig.Token.Path
				}
				loginData["token"] = authConfig.Github.Token
			}

			if len(loginData) == 0 {
				authError = fmt.Errorf("no authentication method provided")
			} else {
				if auth, err := vaultPrvdr.client.Logical().Write(loginPath, loginData); err != nil {
					authError = err
				} else if auth == nil {
					authError = fmt.Errorf("authentication failure")
				}
			}
		}
	}
	return authError
}

func (vaultPrvdr *VaultProvider) Alias() string {
	return vaultPrvdr.alias
}

func New(alias string) *VaultProvider {
	return &VaultProvider{alias: alias}
}
