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

type commandRunner struct {
	cmdRegistry commands.CommandRegistrar
}

func NewCommandRunner() *commandRunner {
	commandRegistry := commandregistry.Instance()
	commandRegistry.Register(
		&commands.RunnableCommand{
			Command: &commands.Command{
				Expression:  "help",
				Description: "display usage instructions",
				Aliases:     []string{"-h", "-help", "--help"},
			},
			Handler: &commands.HelpCommandHandler{
				CommandRegistry: commandRegistry,
			},
		},
	)
	cmdRunner := &commandRunner{
		cmdRegistry: commandRegistry,
	}
	return cmdRunner
}

func (cmdRun *commandRunner) Run(args []string, diags diagnostics.Diagnosable) int {
	cmd, cmdArgs := cmdRun.resolveCommandAndArgs(args)

	if cmd == nil || !cmd.HasHandler() {
		streams.Output().Writeln(format.Formatter().Color(colors.Red).Dimmed().Format(fmt.Sprintf("Unrecognized command %s", strings.Join(args, " "))))
		return 1
	}

	cmd.Run(cmdArgs, diags)
	if diags.HasErrors() {
		return 1
	}
	return 0
}

func (cmdRun *commandRunner) resolveCommandAndArgs(args []string) (commands.CommandRunnable, []string) {
	if len(args) == 0 {
		args = []string{"conformize"}
	}

	argIdx := 0
	cmdNode := cmdRun.cmdRegistry.GetCommand(args[argIdx])
	if cmdNode == nil {
		return nil, args
	}

	cmd := cmdNode.Value
	argIdx++

	for argIdx < len(args) {
		arg := args[argIdx]
		if shouldShowHelp(arg) {
			helpNode := cmdRun.cmdRegistry.GetCommand("help")
			if helpNode != nil {
				return helpNode.Value, args[:argIdx]
			}
			return nil, args
		}

		subCmds, ok := cmdNode.GetChildren(arg)
		next := subCmds.First()
		if !ok || next == nil {
			break
		}

		cmdNode = next
		cmd = cmdNode.Value
		argIdx++
	}

	if !cmd.HasHandler() && argIdx == len(args) {
		helpNode := cmdRun.cmdRegistry.GetCommand("help")
		if helpNode != nil {
			return helpNode.Value, args[:argIdx]
		}
		return nil, args
	}

	return cmd, args[argIdx:]
}

func shouldShowHelp(arg string) bool {
	return cmdutil.IsHelpCommand(arg)
}
