// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package commands

import (
	"fmt"
	"strings"

	"github.com/conformize/conformize/common/diagnostics"
	"github.com/conformize/conformize/internal/blueprint"
	"github.com/conformize/conformize/internal/blueprint/scaffold"
)

type scaffoldMetaData struct {
	version          float64
	sourceAliases    []string
	providers        []string
	referenceAliases []string
	predicates       []string
	format           string
}

type BlueprintScaffoldCommandHandler struct{}

func (ch *BlueprintScaffoldCommandHandler) Handle(c CommandEntry, args []string, diags *diagnostics.Diagnostics) {
	blprntScaffoldBldr := scaffold.NewBuilder()
	cmdFlags := c.GetFlags()
	err := cmdFlags.Parse(args)
	if err != nil {
		diags.Append(diagnostics.Builder().Error().Summary(err.Error()).Build())
		return
	}

	cmdMeta := c.GetMeta()
	if cmdMeta == nil {
		diags.Append(diagnostics.Builder().
			Error().
			Summary("Command metadata is not set").
			Details("Command metadata is required to build the blueprint scaffold.").
			Build(),
		)
		return
	}

	var cmd scaffoldMetaData
	if err = cmdMeta.Get(&cmd); err != nil {
		diags.Append(diagnostics.Builder().
			Error().
			Summary("Failed to retrieve command metadata").
			Details(fmt.Sprintf("Error: %s", err.Error())).
			Build(),
		)
		return
	}

	blprntSchemaFormat := strings.ToLower(cmd.format)
	if blprntSchemaFormat != "json" && blprntSchemaFormat != "yaml" {
		diags.Append(diagnostics.Builder().
			Error().
			Summary(
				fmt.Sprintf(
					"Blueprint must be in JSON or YAML format. Format '%s' is not supported.", blprntSchemaFormat,
				),
			).
			Build(),
		)
		return
	}

	blprntScaffoldBldr.WithVersion(cmd.version)

	fileName := fmt.Sprintf("blueprint.cnfrm.%s", blprntSchemaFormat)
	for idx, sourceAlias := range cmd.sourceAliases {
		blprntScaffoldBldr.WithSource(sourceAlias, cmd.providers[idx])
	}

	for _, refAlias := range cmd.referenceAliases {
		blprntScaffoldBldr.WithReference(refAlias)
	}

	for _, predicate := range cmd.predicates {
		blprntScaffoldBldr.WithPredicate(predicate)
	}

	blprntScaffold, scaffoldDiags := blprntScaffoldBldr.Build()
	diags.Append(scaffoldDiags.Entries()...)

	blprntMarshaller := blueprint.NewBlueprintMarshaller(fileName, blprntScaffold)
	if err = blprntMarshaller.Marshal(); err != nil {
		diags.Append(diagnostics.Builder().
			Error().
			Summary("Couldn't create blueprint scaffold").
			Details(err.Error()).
			Build(),
		)
		return
	}

	diags.Append(diagnostics.Builder().
		Info().
		Summary(fmt.Sprintf("Blueprint scaffold created at ./%s", fileName)).
		Build(),
	)
}
