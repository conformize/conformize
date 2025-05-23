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
	"github.com/conformize/conformize/common/ds"
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
	cmdIdx := 0
	argIdx := 0
	command := args[cmdIdx]
	showHelp := false

	argsLen := len(args)
	var cmd commands.CommandRunnable
	if cmdNode := cmdRun.cmdRegistry.GetCommand(command); cmdNode != nil {
		cmd = cmdNode.Value
		ok := false
		var subCmdNodes ds.NodeList[string, commands.CommandRunnable]
		for argIdx = cmdIdx + 1; cmdNode != nil && argsLen > argIdx; {
			command = args[argIdx]
			subCmdNodes, ok = cmdNode.GetChildren(command)
			cmdNode = subCmdNodes.First()
			if ok {
				cmd = cmdNode.Value
				cmdIdx++
				argIdx++
				continue
			}

			if showHelp = cmdutil.IsHelpCommand(command); showHelp {
				cmdIdx++
				argIdx++
				break
			}
		}

		showHelp = showHelp ||
			(cmd != nil && !cmd.HasHandler()) && (argsLen-1 == cmdIdx)
	}

	if showHelp {
		cmdNode := cmdRun.cmdRegistry.GetCommand("help")
		cmd = cmdNode.Value
		if len(args) > 1 {
			args = args[:cmdIdx]
		}
	} else {
		args = args[argIdx:]
	}

	if cmd == nil || (!cmd.HasHandler()) {
		diags.Append(diagnostics.Builder().
			Error().
			Summary(fmt.Sprintf("Unrecognized command %s", strings.Join(args, " "))).
			Build(),
		)
		return 1
	}

	cmd.Run(args, diags)
	if diags.HasErrors() {
		return 1
	}
	return 0
}
