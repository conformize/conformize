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
	"github.com/conformize/conformize/ui/cli/commands"
	"github.com/conformize/conformize/ui/cli/commands/cmdutil"
	"github.com/conformize/conformize/ui/cli/commands/commandregistry"
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
	cmdIdx := 0
	command := args[cmdIdx]
	showHelp := false

	lastArgIdx := len(args) - 1
	var cmd commands.CommandRunnable
	if cmdNode := cmdRun.cmdRegistry.GetCommand(command); cmdNode != nil {
		cmd = cmdNode.Value
		for argIdx := cmdIdx + 1; cmdNode != nil && len(args) > argIdx; {
			command = args[argIdx]
			subCmdNode, ok := cmdNode.GetChild(command)
			cmdNode = subCmdNode
			if ok {
				cmd = cmdNode.Value
				cmdIdx++
			} else {
				if showHelp = cmdutil.IsHelpCommand(command); showHelp {
					cmdIdx++
					break
				}
			}
			argIdx++
		}

		showHelp = showHelp ||
			(cmd != nil && !cmd.HasHandler()) && (lastArgIdx == cmdIdx)
	}

	if cmd == nil || (!cmd.HasHandler() && !(lastArgIdx-cmdIdx == 0)) {
		diags.Append(diagnostics.Builder().
			Error().
			Summary(fmt.Sprintf("Unrecognized command %s", strings.Join(args, " "))).
			Build(),
		)
		return 1
	}

	if showHelp {
		cmdNode := cmdRun.cmdRegistry.GetCommand("help")
		cmd = cmdNode.Value
		if lastArgIdx > 0 {
			args = args[:cmdIdx]
		}
	} else {
		args = args[cmdIdx+1:]
	}

	cmd.Run(args, diags)
	if diags.HasErrors() {
		return 1
	}
	return 0
}
