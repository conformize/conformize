// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package commands

import (
	"flag"
	"fmt"
	"strings"

	"github.com/conformize/conformize/common/diagnostics"
	"github.com/conformize/conformize/common/typed"
	"github.com/conformize/conformize/internal/providers/api/schema"
	"github.com/conformize/conformize/internal/providers/api/schema/attributes"
)

type blueprintScaffoldCommand struct {
	Handler          CommandHandler
	Meta             *schema.Data
	Expression       string
	Description      string
	Subcommands      []CommandEntry
	Aliases          []string
	version          float64
	sourceAliases    []string
	providers        []string
	referenceAliases []string
	predicates       []string
	format           string
	Hidden           bool
}

func BlueprintScaffoldCommand() CommandEntry {
	return &Command{
		Expression: "scaffold",
		Subcommands: []CommandEntry{
			&blueprintScaffoldCommand{
				Expression:  "blueprint",
				Description: "create a blueprint scaffold",
				Handler:     &BlueprintScaffoldCommandHandler{},
				Meta: schema.NewData(&schema.Schema{
					Attributes: map[string]schema.Attributeable{
						"version":          &attributes.NumberAttribute{},
						"sourceAliases":    &attributes.ListAttribute{ElementsType: &typed.StringTyped{}},
						"providers":        &attributes.ListAttribute{ElementsType: &typed.StringTyped{}},
						"referenceAliases": &attributes.ListAttribute{ElementsType: &typed.StringTyped{}},
						"predicates":       &attributes.ListAttribute{ElementsType: &typed.StringTyped{}},
						"format":           &attributes.StringAttribute{},
					},
				}),
				version:          1,
				sourceAliases:    make([]string, 0),
				providers:        make([]string, 0),
				referenceAliases: make([]string, 0),
				predicates:       make([]string, 0),
				format:           "yaml",
			},
		},
	}
}

func (c *blueprintScaffoldCommand) GetExpression() string {
	return c.Expression
}

func (c *blueprintScaffoldCommand) GetDescription() string {
	return c.Description
}

func (c *blueprintScaffoldCommand) GetAliases() []string {
	return c.Aliases
}

func (c *blueprintScaffoldCommand) GetHandler() CommandHandler {
	return c.Handler
}

func (c *blueprintScaffoldCommand) GetSubcommands() []CommandEntry {
	return c.Subcommands
}

func (c *blueprintScaffoldCommand) IsHidden() bool {
	return c.Hidden
}

func (c *blueprintScaffoldCommand) GetMeta() *schema.Data {
	return c.Meta
}

func (c *blueprintScaffoldCommand) Run(args []string, diags *diagnostics.Diagnostics) {
	if c.Handler != nil {
		c.Handler.Handle(c, args, diags)
	}
}

func (c *blueprintScaffoldCommand) GetFlags() *flag.FlagSet {
	flags := flag.NewFlagSet(c.GetExpression(), flag.ContinueOnError)
	flags.Float64Var(&c.version, "version", 1,
		"specifies the schema version to be used for blueprint scaffold.\ndefault value of 1 will be used unless provided.",
	)
	flags.Func("source", "specifies the alias for a configuration source to be added to blueprint scaffold",
		func(v string) error {
			var sources []string
			aliases, err := c.Meta.GetAtPath("sourceAliases")
			if err == nil {
				if err = aliases.As(&sources); err != nil {
					return err
				}
				sources = append(sources, v)
			}
			return c.Meta.SetAtPath("sourceAliases", sources)
		})
	flags.Func("provider", "specifies the provider to be used to retrieve data from a configuration source",
		func(v string) error {
			var providers []string
			providerAliases, err := c.Meta.GetAtPath("providers")
			if err == nil {
				if err = providerAliases.As(&providers); err != nil {
					return err
				}
				providers = append(providers, v)
			}
			return c.Meta.SetAtPath("providers", providers)
		})

	flags.Func("refs", "specifies a comma-separated list of reference aliases to be defined in the blueprint scaffold", func(v string) error {
		c.referenceAliases = strings.Split(v, ",")
		return c.Meta.SetAtPath("referenceAliases", c.referenceAliases)
	})

	flags.Func("predicates", "specifies a comma-separated list of predicates to ba added to the blueprint scaffold", func(v string) error {
		c.predicates = strings.Split(v, ",")
		return c.Meta.SetAtPath("predicates", c.predicates)
	})

	flags.Func("format", "specifies the output format for blueprint scaffold - JSON or YAML, e.g. -format yaml. YAML format will be used if not specified.", func(v string) error {
		if v != "json" && v != "yaml" {
			return fmt.Errorf("invalid format '%s' specified. Supported formats are 'json' and 'yaml'", v)
		}
		c.format = strings.ToLower(v)
		return c.Meta.SetAtPath("format", c.format)
	})
	return flags
}
