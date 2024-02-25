// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package secretsmanager

import (
	"context"
	"runtime"
	"sync"

	"github.com/conformize/conformize/common/diagnostics"
	"github.com/conformize/conformize/common/ds"
	"github.com/conformize/conformize/common/path"
	"github.com/conformize/conformize/common/typed"
	sdk "github.com/conformize/conformize/internal/providers/api"
	"github.com/conformize/conformize/internal/providers/api/schema"
	"github.com/conformize/conformize/internal/providers/api/schema/attributes"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager/types"
)

type awsSecretsManagerClientConfig struct {
	Region      string `cnfrmz:"region"`
	EndpointUrl string `cnfrmz:"endpointUrl"`
}

type awsSecretsManagerQueryOptions struct {
	SecretIdList []string              `cnfrmz:"secretIdList"`
	Filters      []map[string][]string `cnfrmz:"filters"`
	MaxResults   int32                 `cnfrmz:"maxResults"`
}

type AwsSecretsManagerProvider struct {
	alias  string
	client *secretsmanager.Client
}

func (awsSecretsManagerPrvdr *AwsSecretsManagerProvider) Alias() string {
	return awsSecretsManagerPrvdr.alias
}

func New(alias string) *AwsSecretsManagerProvider {
	return &AwsSecretsManagerProvider{alias: alias}
}

func (awsSecretsManagerPrvdr *AwsSecretsManagerProvider) ConfigurationSchema() *schema.Schema {
	return &schema.Schema{
		Description: "Configuration for the AWS Secrets Manager provider",
		Version:     1,
		Attributes: map[string]schema.Attributeable{
			"region":      &attributes.StringAttribute{},
			"endpointUrl": &attributes.StringAttribute{},
		},
	}
}

func (awsSecretsManagerPrvdr *AwsSecretsManagerProvider) ProvisionDataRequestSchema() *schema.Schema {
	return &schema.Schema{
		Description: "AWS Secrets Manager resource request schema",
		Version:     1,
		Attributes: map[string]schema.Attributeable{
			"secretIdList": &attributes.ListAttribute{ElementsType: &typed.StringTyped{}},
			"maxResults":   &attributes.NumberAttribute{},
			"filters": &attributes.ListAttribute{
				ElementsType: &typed.MapTyped{
					ElementsType: &typed.ListTyped{ElementsType: &typed.StringTyped{}},
				},
			},
		},
	}
}

func (awsSecretsManagerPrvdr *AwsSecretsManagerProvider) Configure(req *sdk.ConfigurationRequest) error {
	var clientConfig awsSecretsManagerClientConfig
	if err := req.Get(&clientConfig); err != nil {
		return err
	}

	if cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(clientConfig.Region)); err != nil {
		return err
	} else {
		optsFunc := func(o *secretsmanager.Options) {}
		if clientConfig.EndpointUrl != "" {
			optsFunc = func(o *secretsmanager.Options) {
				o.BaseEndpoint = &clientConfig.EndpointUrl
			}
		}
		awsSecretsManagerPrvdr.client = secretsmanager.NewFromConfig(cfg, optsFunc)
		return nil
	}
}

func (awsSecretsManagerPrvdr *AwsSecretsManagerProvider) Provide(req *sdk.ProviderDataRequest) (*ds.Node[string, any], *diagnostics.Diagnostics) {
	diags := diagnostics.NewDiagnostics()
	var queryOptions awsSecretsManagerQueryOptions
	if err := req.Get(&queryOptions); err != nil {
		diags.Append(diagnostics.Builder().Error().Details(err.Error()).Build())
		return nil, diags
	}

	serectIdsLen := len(queryOptions.SecretIdList)
	filtersLen := len(queryOptions.Filters)

	queryBySecredId := serectIdsLen > 0
	queryByFilters := filtersLen > 0
	if !queryBySecredId && !queryByFilters {
		diags.Append(diagnostics.Builder().Error().Details("either 'Filters' or 'SecretIdList' must be specified").Build())
		return nil, diags
	}

	if queryBySecredId && queryByFilters {
		diags.Append(diagnostics.Builder().Error().Details("only one of 'Filters' or 'SecretIdList' must be specified, but not both").Build())
		return nil, diags
	}

	var data *ds.Node[string, any]
	var err error

	batchQuery := queryByFilters || serectIdsLen > 1
	if !batchQuery {
		data, err = getSecretValue(awsSecretsManagerPrvdr.client, &queryOptions)
	} else {
		data, err = getSecretValuesBatch(awsSecretsManagerPrvdr.client, &queryOptions)
	}

	if err != nil {
		diags.Append(diagnostics.Builder().Error().Details(err.Error()).Build())
		return nil, diags
	}
	return data, diags
}

func getSecretValue(client *secretsmanager.Client, queryOptions *awsSecretsManagerQueryOptions) (*ds.Node[string, any], error) {
	query := &secretsmanager.GetSecretValueInput{
		SecretId: &queryOptions.SecretIdList[0],
	}

	if secretValue, err := client.GetSecretValue(context.Background(), query); err != nil {
		return nil, err
	} else {
		result := ds.NewNode[string, any]()
		valuePath, _ := path.NewFromStringWithSeparator(*secretValue.Name, '/')
		steps := valuePath.Steps()

		lastNodeRef := result
		for step, hasNext := steps.Next(); hasNext; step, hasNext = steps.Next() {
			stepName := step.String()
			lastNodeRef = lastNodeRef.AddChild(stepName)
		}
		lastNodeRef.Value = *secretValue.SecretString
		lastNodeRef.AddAttribute("ARN", *secretValue.ARN)
		return result, nil
	}
}

func getSecretValuesBatch(client *secretsmanager.Client, queryOptions *awsSecretsManagerQueryOptions) (*ds.Node[string, any], error) {
	query := &secretsmanager.BatchGetSecretValueInput{
		SecretIdList: queryOptions.SecretIdList,
		MaxResults:   &queryOptions.MaxResults,
		Filters:      prepareQueryFilters(queryOptions.Filters),
	}

	result := ds.NewNode[string, any]()

	var rwLock sync.RWMutex
	var wg sync.WaitGroup
	errChan := make(chan error, 1)
	defer close(errChan)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wg.Add(1)
	go func() {
		defer wg.Done()
		availableCPUs := max(1, runtime.NumCPU()-1)
		cpus := runtime.GOMAXPROCS(availableCPUs)
		tasks := make(chan struct{}, availableCPUs)
		defer close(tasks)
		defer runtime.GOMAXPROCS(cpus)

		for {
			secretValuesBatch, err := client.BatchGetSecretValue(ctx, query)
			if err != nil {
				errChan <- err
				return
			}

			for _, secretValue := range secretValuesBatch.SecretValues {
				wg.Add(1)
				tasks <- struct{}{}
				go func(secretValue types.SecretValueEntry) {
					defer wg.Done()
					valuePath, _ := path.NewFromStringWithSeparator(*secretValue.Name, '/')
					steps := valuePath.Steps()

					lastNodeRef := result
					for step, hasNext := steps.Next(); hasNext; step, hasNext = steps.Next() {
						stepName := step.String()

						rwLock.RLock()
						nodes, found := lastNodeRef.GetChildren(stepName)
						rwLock.RUnlock()

						if found {
							lastNodeRef = nodes.First()
							continue
						}
						rwLock.Lock()
						lastNodeRef = lastNodeRef.AddChild(stepName)
						rwLock.Unlock()
					}

					lastNodeRef.Value = *secretValue.SecretString
					lastNodeRef.AddAttribute("ARN", *secretValue.ARN)
					<-tasks
				}(secretValue)
			}

			if secretValuesBatch.NextToken == nil || len(*secretValuesBatch.NextToken) == 0 {
				break
			}
			query.NextToken = secretValuesBatch.NextToken
		}
	}()

	wg.Wait()
	select {
	case err := <-errChan:
		return nil, err
	default:
	}
	return result, nil
}

func prepareQueryFilters(filters []map[string][]string) []types.Filter {
	var queryFilters []types.Filter
	for _, filter := range filters {
		for key, values := range filter {
			queryFilters = append(queryFilters, types.Filter{
				Key:    types.FilterNameStringType(key),
				Values: values,
			})
		}
	}
	return queryFilters
}
