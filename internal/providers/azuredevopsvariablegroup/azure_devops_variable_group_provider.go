// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package azuredevopsvariablegroup

import (
	"context"

	"github.com/conformize/conformize/common/diagnostics"
	"github.com/conformize/conformize/common/ds"
	sdk "github.com/conformize/conformize/internal/providers/api"
	"github.com/conformize/conformize/internal/providers/api/schema"
	"github.com/conformize/conformize/internal/providers/api/schema/attributes"
	"github.com/microsoft/azure-devops-go-api/azuredevops"
	"github.com/microsoft/azure-devops-go-api/azuredevops/taskagent"
)

type azureDevOpsVariableGroupClientConfig struct {
	OrganizationUrl string `cnfrmz:"organizationUrl"`
	Project         string `cnfrmz:"project"`
	GroupId         int    `cnfrmz:"groupId"`
	Token           string `cnfrmz:"token"`
}

type AzureDevOpsVariableGroupProvider struct {
	alias   string
	project string
	groupId int
	client  taskagent.Client
}

func (azureDevOpsVarGrpPrvdr *AzureDevOpsVariableGroupProvider) Configure(req *sdk.ConfigurationRequest) error {
	var clientConfig azureDevOpsVariableGroupClientConfig
	if err := req.Get(&clientConfig); err != nil {
		return err
	}

	connection := azuredevops.NewPatConnection(clientConfig.OrganizationUrl, clientConfig.Token)
	client, err := taskagent.NewClient(context.Background(), connection)
	if err != nil {
		return err
	}

	azureDevOpsVarGrpPrvdr.client = client
	azureDevOpsVarGrpPrvdr.project = clientConfig.Project
	azureDevOpsVarGrpPrvdr.groupId = clientConfig.GroupId
	return nil
}

func (azureDevOpsVarGrpPrvdr *AzureDevOpsVariableGroupProvider) ConfigurationSchema() *schema.Schema {
	return &schema.Schema{
		Description: "Configuration for the Azure DevOps Variable Group provider",
		Version:     1,
		Attributes: map[string]schema.Attributeable{
			"organizationUrl": &attributes.StringAttribute{},
			"project":         &attributes.StringAttribute{},
			"groupId":         &attributes.NumberAttribute{},
			"token":           &attributes.StringAttribute{},
		},
	}
}

func (azureDevOpsVarGrpPrvdr *AzureDevOpsVariableGroupProvider) ProvisionDataRequestSchema() *schema.Schema {
	return &schema.Schema{}
}

func (azureDevOpsVarGrpPrvdr *AzureDevOpsVariableGroupProvider) Provide(queryRequest *sdk.ProviderDataRequest) (*ds.Node[string, any], *diagnostics.Diagnostics) {
	diags := diagnostics.NewDiagnostics()

	varGrp, err := azureDevOpsVarGrpPrvdr.client.GetVariableGroup(context.Background(), taskagent.GetVariableGroupArgs{
		Project: &azureDevOpsVarGrpPrvdr.project,
		GroupId: &azureDevOpsVarGrpPrvdr.groupId,
	})

	if err != nil {
		diags.Append(diagnostics.Builder().Error().Summary(err.Error()).Build())
		return nil, diags
	}

	result := ds.NewNode[string, any]()
	if varGrp.Variables != nil {
		var node *ds.Node[string, any]
		for varName, varValue := range *varGrp.Variables {
			node = nil
			if nodes, found := result.GetChildren(varName); !found {
				node = result.AddChild(varName)
			} else {
				node = nodes.First()
			}

			node.Value = varValue
		}
	}
	return result, diags
}

func (azureDevOpsVarGrpPrvdr *AzureDevOpsVariableGroupProvider) Alias() string {
	return azureDevOpsVarGrpPrvdr.alias
}

func New(alias string) *AzureDevOpsVariableGroupProvider {
	return &AzureDevOpsVariableGroupProvider{alias: alias}
}
