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

type Command struct {
	Expression  string
	Description string
	Aliases     []string
	Hidden      bool
}

func (c *Command) GetExpression() string {
	return c.Expression
}

func (c *Command) GetDescription() string {
	return c.Description
}

func (c *Command) IsHidden() bool {
	return c.Hidden
}

func (c *Command) GetFlags() *flag.FlagSet {
	return nil
}

func (c *Command) GetSubcommands() []CommandRunnable {
	return []CommandRunnable{}
}

func (c *Command) GetAliases() []string {
	if c.Aliases == nil {
		return []string{}
	}
	return c.Aliases
}
