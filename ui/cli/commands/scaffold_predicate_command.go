// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package commands

import (
	"flag"
)

type ScaffoldPredicateCommand struct {
	*Command
}

func (c *ScaffoldPredicateCommand) GetFlags() *flag.FlagSet {
	flags := flag.NewFlagSet(c.GetExpression(), flag.ExitOnError)
	flags.String("predicate", "", "Specifies the predicate to be used for evaluation")
	flags.String("type", "", "Specifies the type of the value to be evaluated with a given predicate")
	flags.String("format", "yaml", "Specifies the format for the predicate scaffold - YAML or JSON")
	return flags
}
