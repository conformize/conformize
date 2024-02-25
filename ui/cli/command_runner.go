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

func (cmdRun *commandRunner) Run(args []string, diags diagnostics.Diagnosable) int {
	var ok bool
	var cmd commands.CommandEntry

	var cmdArgs = args
	var argsLen = len(args)
	cmdReg := commandregistry.Instance()
	if argsLen > 0 && cmdutil.IsHelpCommand(args[0]) {
		cmd, _, ok = cmdReg.GetCommand([]string{"help"})
		cmdArgs = cmdArgs[1:]
	} else {
		cmd, cmdArgs, ok = cmdReg.GetCommand(args)
		if !ok && len(cmdArgs) > 0 && cmdutil.IsHelpCommand(args[argsLen-1]) {
			cmd, _, ok = cmdReg.GetCommand([]string{"help"})
			cmdArgs = cmdArgs[:argsLen-1]
		}
	}

	if !ok || cmd == nil {
		streams.Output().Writeln(format.Formatter().Color(colors.Red).Dimmed().Format(fmt.Sprintf("Unrecognized command %s", strings.Join(args, " "))))
		return 1
	}

	if cmd != nil && cmd.GetHandler() == nil && !cmdutil.IsHelpCommand(cmd.GetExpression()) {
		cmd, _, ok = cmdReg.GetCommand([]string{"help"})
		cmdArgs = args
	}

	cmd.Run(cmdArgs, diags)
	if diags.HasErrors() {
		return 1
	}
	return 0
}
