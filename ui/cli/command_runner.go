// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package cli

import (
	"fmt"
	"strings"

	"github.com/conformize/conformize/common/diagnostics"
	"github.com/conformize/conformize/common/format"
	"github.com/conformize/conformize/common/format/colors"
	"github.com/conformize/conformize/common/streams"
	"github.com/conformize/conformize/ui/commands"
	"github.com/conformize/conformize/ui/commands/cmdutil"
	"github.com/conformize/conformize/ui/commands/cmdutil/commandregistry"
)

type commandRunner struct{}

func NewCommandRunner() *commandRunner {
	commandRegistry := commandregistry.Instance()
	commandRegistry.Register(
		&commands.Command{
			Expression:  "help",
			Description: "display usage instructions",
			Aliases:     []string{"-h", "-help", "--help"},
			Handler: &commands.HelpCommandHandler{
				CommandRegistry: commandRegistry,
			},
		},
	)

	return &commandRunner{}
}

func (cmdRun *commandRunner) Run(args []string, diags *diagnostics.Diagnostics, doneChan chan struct{}) int {
	defer func() { doneChan <- struct{}{} }()

	argsLen := len(args)
	commandRegistry := commandregistry.Instance()

	showHelp := false
	var cmd commands.CommandEntry

	var cmdArgs []string = args
	if argsLen > 0 {
		if cmdutil.IsHelpCommand(args[0]) {
			cmdArgs = cmdArgs[1:]
			showHelp = true
		} else if cmdutil.IsHelpCommand(args[argsLen-1]) {
			cmdArgs = cmdArgs[:argsLen-1]
			showHelp = true
		}
	}

	if showHelp {
		cmd, _, _ = commandRegistry.GetCommand([]string{"help"})
	} else {
		var found bool
		cmd, cmdArgs, found = commandRegistry.GetCommand(cmdArgs)
		if !found {
			streams.Output().Writeln(
				format.Formatter().Color(colors.Red).
					Format(fmt.Sprintf("Unrecognized command %s", strings.Join(args, " "))),
			)
			return 1
		}

		if cmd.GetHandler() == nil && !cmdutil.IsHelpCommand(cmd.GetExpression()) {
			cmd, _, _ = commandRegistry.GetCommand([]string{"help"})
			cmdArgs = args
		}
	}

	cmd.Run(cmdArgs, diags)
	if !diags.HasErrors() {
		return 0
	}
	return 1
}
