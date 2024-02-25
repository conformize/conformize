// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package commands

import (
	"flag"
	"strings"
)

type BlueprintScaffoldCommand struct {
	*Command
	version          float64
	sourceAliases    []string
	providers        []string
	referenceAliases []string
	predicates       []string
	format           string
}

func (c *BlueprintScaffoldCommand) GetFlags() *flag.FlagSet {
	flags := flag.NewFlagSet(c.GetExpression(), flag.ContinueOnError)
	flags.Float64Var(&c.version, "version", 1,
		"specifies the schema version to be used for blueprint scaffold.\ndefault value of 1 will be used unless provided.",
	)
	flags.Func("source", "specifies the alias for a configuration source to be added to blueprint scaffold",
		func(v string) error {
			c.sourceAliases = append(c.sourceAliases, v)
			return nil
		})
	flags.Func("provider", "specifies the provider to be used to retrieve data from a configuration source",
		func(v string) error {
			c.providers = append(c.providers, v)
			return nil
		})

	flags.Func("refs", "specifies a comma-separated list of reference aliases to be defined in the blueprint scaffold", func(v string) error {
		c.referenceAliases = strings.Split(v, ",")
		return nil
	})

	flags.Func("predicates", "specifies a comma-separated list of predicates to ba added to the blueprint scaffold", func(v string) error {
		c.predicates = strings.Split(v, ",")
		return nil
	})

	flags.StringVar(&c.format, "format", "yaml",
		"specifies the output format for blueprint scaffold - JSON or YAML, e.g. -format yaml. YAML format will be used if not specified.",
	)
	return flags
}
