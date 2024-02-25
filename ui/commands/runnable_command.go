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
)

type RunnableCommand struct {
	Command Commandable
	Handler CommandHandler
}

func (rCmd *RunnableCommand) GetExpression() string {
	return rCmd.Command.GetExpression()
}

func (rCmd *RunnableCommand) GetDescription() string {
	return rCmd.Command.GetDescription()
}

func (rCmd *RunnableCommand) HasHandler() bool {
	return rCmd.Handler != nil
}

func (rCmd *RunnableCommand) IsHidden() bool {
	return rCmd.Command.IsHidden()
}

func (rCmd *RunnableCommand) GetHandler() CommandHandler {
	return rCmd.Handler
}

func (rCmd *RunnableCommand) GetFlags() *flag.FlagSet {
	return rCmd.Command.GetFlags()
}

func (rCmd *RunnableCommand) GetSubcommands() []CommandRunnable {
	return rCmd.Command.GetSubcommands()
}

func (rCmd *RunnableCommand) GetAliases() []string {
	return rCmd.Command.GetAliases()
}

func (rCmd *RunnableCommand) Run(args []string, diags diagnostics.Diagnosable) {
	if rCmd.HasHandler() {
		rCmd.Handler.Handle(rCmd.Command, args, diags)
	}
}
