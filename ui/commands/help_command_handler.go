// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package commands

import (
	"flag"
	"fmt"
	"strings"

	"github.com/conformize/conformize/common/diagnostics"
	"github.com/conformize/conformize/common/streams"
	"github.com/conformize/conformize/internal/ui/options"
	"github.com/conformize/conformize/resources"
)

const commandNameColumnWidth = 20

type HelpCommandHandler struct {
	CommandRegistry CommandRegistrar
}

func (h *HelpCommandHandler) Handle(c Commandable, args []string, diags diagnostics.Diagnosable) {
	helpBldr := strings.Builder{}
	if len(args) == 0 {
		availableCmdsBldr := strings.Builder{}
		for _, cmd := range h.CommandRegistry.GetCommands() {
			availableCmdsBldr.WriteString(fmt.Sprintf("%-*s\t%s\n",
				commandNameColumnWidth, cmd.GetExpression(), cmd.GetDescription()),
			)
		}
		helpBldr.WriteString(fmt.Sprintf(resources.HelpText(), availableCmdsBldr.String(), options.Help()))
		streams.Output().Writef("%s\n", helpBldr.String())
		return
	}

	cmdNode := h.CommandRegistry.GetCommand(args[0])
	if cmdNode == nil {
		diags.Append(diagnostics.Builder().
			Error().
			Summary(fmt.Sprintf("Unrecognized command %s", strings.Join(args, " "))).
			Build(),
		)
		return
	}
	cmd := cmdNode.Value

	cmdStr := cmd.GetExpression()
	for _, arg := range args[1:] {
		cmdNodes, ок := cmdNode.GetChildren(arg)
		if !ок || cmdNodes.First() == nil {
			diags.Append(diagnostics.Builder().
				Error().
				Summary(fmt.Sprintf("Unrecognized command %s", strings.Join(args, " "))).
				Build(),
			)
			return
		}
		cmdNode = cmdNodes.First()
		cmd = cmdNode.Value
		cmdStr = fmt.Sprintf("%s %s", cmdStr, arg)
	}
	helpBldr.WriteString(fmt.Sprintf("description: %s\n\n", cmd.GetDescription()))
	cmdUsageStr := fmt.Sprintf("usage: conformize [global options] %s [subcommand] [arguments]\n", cmdStr)
	helpBldr.WriteString(cmdUsageStr)

	helpBldr.WriteString("\navailable subcommands:\n\n")
	for _, subCmd := range cmd.GetSubcommands() {
		helpBldr.WriteString(fmt.Sprintf("%-*s\t%s\n",
			commandNameColumnWidth, subCmd.GetExpression(), subCmd.GetDescription()),
		)
	}

	helpBldr.WriteString(fmt.Sprintf("%-*s\t%s\n",
		commandNameColumnWidth, "help", "display usage instructions"),
	)

	if cmdFlags := cmd.GetFlags(); cmdFlags != nil {
		helpBldr.WriteString("\navailable arguments:\n\n")
		cmdFlags.VisitAll(func(flag *flag.Flag) {
			helpBldr.WriteString(fmt.Sprintf("%-*s\t%s\n",
				commandNameColumnWidth, fmt.Sprintf("-%s", flag.Name), flag.Usage),
			)
		})
	}
	streams.Output().Writef("%s\n", helpBldr.String())
}
