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
	"github.com/conformize/conformize/common/format"
	"github.com/conformize/conformize/common/format/colors"
	"github.com/conformize/conformize/common/streams"
	"github.com/conformize/conformize/internal/ui/options"
	"github.com/conformize/conformize/resources"
)

const commandNameColumnWidth = 20

type HelpCommandHandler struct {
	CommandRegistry CommandRegistrar
}

func (h *HelpCommandHandler) Handle(c CommandEntry, args []string, diags *diagnostics.Diagnostics) {
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

	cmd, _, _ := h.CommandRegistry.GetCommand(args)
	if cmd == nil {
		streams.Error().
			Writeln(
				format.Formatter().
					Color(colors.Red).Format(fmt.Sprintf("Unrecognized command %s", strings.Join(args, " "))),
			)
		return
	}

	var usageStr string
	var subCommands = cmd.GetSubcommands()
	if len(subCommands) > 0 {
		usageStr = fmt.Sprintf("usage: conformize [global options] %s [subcommand] [arguments]\n", cmd.GetExpression())
	} else {
		usageStr = fmt.Sprintf("usage: conformize [global options] %s [arguments]\n", strings.Join(args, " "))
	}

	helpBldr.WriteString(fmt.Sprintf("description: %s\n\n", cmd.GetDescription()))
	helpBldr.WriteString(usageStr)

	helpBldr.WriteString("\navailable subcommands:\n\n")
	for _, subCmd := range subCommands {
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
