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

func (cmdRun *commandRunner) Run(args []string, diags *diagnostics.Diagnostics) int {
	commandRegistry := commandregistry.Instance()
	cmd, cmdArgs := parseCommand(args, commandRegistry)

	cmd.Run(cmdArgs, diags)
	if diags.HasErrors() {
		return 1
	}
	return 0
}

func parseCommand(args []string, registry commands.CommandRegistrar) (commands.CommandEntry, []string) {
	cmdArgs, showHelp := extractHelpFlag(args)

	if showHelp {
		cmd, _, _ := registry.GetCommand([]string{"help"})
		return cmd, cmdArgs
	}

	return resolveCommand(cmdArgs, args, registry)
}

func extractHelpFlag(args []string) ([]string, bool) {
	argsLen := len(args)
	if argsLen == 0 {
		return args, false
	}

	if cmdutil.IsHelpCommand(args[0]) {
		return args[1:], true
	}

	if cmdutil.IsHelpCommand(args[argsLen-1]) {
		return args[:argsLen-1], true
	}

	return args, false
}

func resolveCommand(cmdArgs, originalArgs []string, registry commands.CommandRegistrar) (commands.CommandEntry, []string) {
	cmd, parsedArgs, found := registry.GetCommand(cmdArgs)

	if !found {
		showUnrecognizedCommandError(originalArgs)
		return nil, originalArgs
	}

	if cmd.GetHandler() == nil && !cmdutil.IsHelpCommand(cmd.GetExpression()) {
		helpCmd, _, _ := registry.GetCommand([]string{"help"})
		return helpCmd, originalArgs
	}

	return cmd, parsedArgs
}

func showUnrecognizedCommandError(args []string) {
	streams.Output().Writeln(
		format.Formatter().Color(colors.Red).
			Format(fmt.Sprintf("Unrecognized command %s", strings.Join(args, " "))),
	)
}
