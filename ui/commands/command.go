// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package commands

import (
	"flag"

	"github.com/conformize/conformize/common/diagnostics"
	"github.com/conformize/conformize/internal/providers/api/schema"
)

type CommandHandler interface {
	Handle(cmd CommandEntry, args []string, diags *diagnostics.Diagnostics)
}

type CommandEntry interface {
	GetExpression() string
	GetDescription() string
	GetAliases() []string
	GetFlags() *flag.FlagSet
	GetHandler() CommandHandler
	GetSubcommands() []CommandEntry
	IsHidden() bool
	GetMeta() *schema.Data
	Run(args []string, diags *diagnostics.Diagnostics)
}

type Command struct {
	Expression  string
	Description string
	Aliases     []string
	Handler     CommandHandler
	Subcommands []CommandEntry
	Hidden      bool
	Flags       func() *flag.FlagSet
	Meta        *schema.Data
}

func (c *Command) Run(args []string, diags *diagnostics.Diagnostics) {
	if c.Handler != nil {
		c.Handler.Handle(c, args, diags)
	}
}

func (c *Command) GetExpression() string {
	return c.Expression
}

func (c *Command) GetDescription() string {
	return c.Description
}

func (c *Command) GetAliases() []string {
	return c.Aliases
}

func (c *Command) GetFlags() *flag.FlagSet {
	if c.Flags != nil {
		return c.Flags()
	}
	return flag.NewFlagSet(c.Expression, flag.ContinueOnError)
}

func (c *Command) GetSubcommands() []CommandEntry {
	return c.Subcommands
}

func (c *Command) GetHandler() CommandHandler {
	return c.Handler
}

func (c *Command) IsHidden() bool {
	return c.Hidden
}

func (c *Command) GetMeta() *schema.Data {
	if c.Meta == nil {
		c.Meta = &schema.Data{}
	}
	return c.Meta
}
