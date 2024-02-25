// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package ssm

import (
	"context"
	"strings"

	"github.com/conformize/conformize/common/diagnostics"
	"github.com/conformize/conformize/common/ds"
	"github.com/conformize/conformize/common/path"
	"github.com/conformize/conformize/common/typed"
	api "github.com/conformize/conformize/internal/providers/api"
	"github.com/conformize/conformize/internal/providers/api/schema"
	"github.com/conformize/conformize/internal/providers/api/schema/attributes"
	"github.com/conformize/conformize/serialization"
	"github.com/conformize/conformize/serialization/unmarshal/yaml"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
)

type AwsParameterStoreProvider struct {
	client *ssm.Client
}

type awsParameterStoreClientConfig struct {
	Region      string
	EndpointUrl string
}

type awsParameterStoreQueryOptions struct {
	Path             string
	Recursive        bool
	WithDecryption   bool
	MaxResults       int32
	ParameterFilters []string
}

func (awsParamStorePrvdr *AwsParameterStoreProvider) ConfigurationSchema() *schema.Schema {
	return &schema.Schema{
		Description: "Configuration for the AWS Parameter Store provider",
		Version:     1,
		Attributes: map[string]schema.Attributeable{
			"Region":      &attributes.StringAttribute{},
			"EndpointUrl": &attributes.StringAttribute{},
		},
	}
}

func (awsParamStorePrvdr *AwsParameterStoreProvider) ProvisionDataRequestSchema() *schema.Schema {
	return &schema.Schema{
		Description: "AWS Parameter Store resource request schema",
		Version:     1,
		Attributes: map[string]schema.Attributeable{
			"Path":             &attributes.StringAttribute{},
			"Recursive":        &attributes.BooleanAttribute{},
			"WithDecryption":   &attributes.BooleanAttribute{},
			"MaxResults":       &attributes.NumberAttribute{},
			"ParameterFilters": &attributes.ListAttribute{ElementsType: &typed.StringTyped{}},
		},
	}
}

func (awsParamStorePrvdr *AwsParameterStoreProvider) Configure(req *api.ConfigurationRequest) error {
	var clientConfig awsParameterStoreClientConfig
	if err := req.Get(&clientConfig); err != nil {
		return err
	}

	if cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(clientConfig.Region)); err != nil {
		return err
	} else {
		optsFunc := func(o *ssm.Options) {}
		if clientConfig.EndpointUrl != "" {
			optsFunc = func(o *ssm.Options) {
				o.BaseEndpoint = &clientConfig.EndpointUrl
			}
		}
		awsParamStorePrvdr.client = ssm.NewFromConfig(cfg, optsFunc)
		return nil
	}
}

func (awsParamStorePrvdr *AwsParameterStoreProvider) Provide(req *api.ProviderDataRequest) (*ds.Node[string, interface{}], *diagnostics.Diagnostics) {
	diags := diagnostics.NewDiagnostics()
	var queryOptions awsParameterStoreQueryOptions
	if err := req.Get(&queryOptions); err != nil {
		diags.Append(diagnostics.Builder().Error().Details(err.Error()).Build())
		return nil, diags
	}

	if queryOptions.Recursive {
		return getParametersByPath(awsParamStorePrvdr.client, &queryOptions)
	} else {
		return getParameterByName(awsParamStorePrvdr.client, &queryOptions)
	}
}

func getParameterByName(client *ssm.Client, queryOptions *awsParameterStoreQueryOptions) (*ds.Node[string, interface{}], *diagnostics.Diagnostics) {
	query := &ssm.GetParameterInput{
		Name:           &queryOptions.Path,
		WithDecryption: &queryOptions.WithDecryption,
	}

	diags := diagnostics.NewDiagnostics()
	if param, err := client.GetParameter(context.Background(), query); err != nil {
		diags.Append(diagnostics.Builder().Error().Details(err.Error()).Build())
		return nil, diags
	} else {
		result := ds.NewNode[string, interface{}]()
		valuePath, _ := path.NewFromStringWithSeparator(*param.Parameter.Name, '/')
		steps := valuePath.Steps()

		lastNodeRef := result
		for step, hasNext := steps.Next(); hasNext; step, hasNext = steps.Next() {
			stepName := step.String()
			lastNodeRef = lastNodeRef.AddChild(stepName)
		}

		if param.Parameter.Type == "StringList" {
			lastNodeRef.Value = strings.Split(*param.Parameter.Value, ",")
		} else {
			lastNodeRef.Value = *param.Parameter.Value
		}
		lastNodeRef.AddAttribute("DataType", *param.Parameter.DataType)
		lastNodeRef.AddAttribute("Type", param.Parameter.Type)
		lastNodeRef.AddAttribute("ARN", *param.Parameter.ARN)
		return result, nil
	}
}

func getParametersByPath(client *ssm.Client, queryOptions *awsParameterStoreQueryOptions) (*ds.Node[string, interface{}], *diagnostics.Diagnostics) {
	query := &ssm.GetParametersByPathInput{
		Path:           &queryOptions.Path,
		WithDecryption: &queryOptions.WithDecryption,
		MaxResults:     &queryOptions.MaxResults,
		Recursive:      &queryOptions.Recursive,
	}

	result := ds.NewNode[string, interface{}]()
	diags := diagnostics.NewDiagnostics()
	for {
		if params, err := client.GetParametersByPath(context.Background(), query); err == nil {
			for _, param := range params.Parameters {
				valuePath, _ := path.NewFromStringWithSeparator(*param.Name, '/')

				lastNodeRef := result
				steps := valuePath.Steps()
				for step, hasNext := steps.Next(); hasNext; step, hasNext = steps.Next() {
					stepName := step.String()
					if nodeRef, found := lastNodeRef.GetChild(stepName); found {
						lastNodeRef = nodeRef
					} else {
						lastNodeRef = lastNodeRef.AddChild(stepName)
					}
				}

				bufData := serialization.NewBufferedData([]byte(*param.Value))
				yamlUnmarshal := yaml.NewYamlUnmarshal(bufData)
				if node, err := yamlUnmarshal.Unmarshal(); err == nil {
					*lastNodeRef = *ds.MergeNodes(lastNodeRef, node)
				} else {
					if param.Type == "StringList" {
						lastNodeRef.Value = strings.Split(*param.Value, ",")
					} else {
						lastNodeRef.Value = *param.Value
					}
				}
				lastNodeRef.AddAttribute("DataType", *param.DataType)
				lastNodeRef.AddAttribute("Type", param.Type)
				lastNodeRef.AddAttribute("ARN", *param.ARN)
			}

			if params.NextToken == nil || len(*params.NextToken) == 0 {
				break
			} else {
				query.NextToken = params.NextToken
			}
		} else {
			diags.Append(diagnostics.Builder().Error().Details(err.Error()).Build())
			return nil, diags
		}
	}
	return result, nil
}

func New() api.ConfigurationProvider {
	return &AwsParameterStoreProvider{}
}
