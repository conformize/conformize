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

type BlueprintTaskCommand struct {
	*Command
}

func (c *BlueprintTaskCommand) GetFlags() *flag.FlagSet {
	flags := flag.NewFlagSet(c.GetExpression(), flag.ContinueOnError)
	flags.String("f", "", "specifies path to a blueprint file")
	return flags
}
