// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package commandregistry

import (
	"github.com/conformize/conformize/ui/commands"
)

func supportedCommands() []commands.CommandRunnable {
	return []commands.CommandRunnable{
		&commands.RunnableCommand{
			Command: &commands.RunnableCommand{
				Command: &commands.Command{
					Expression:  "conformize",
					Description: "Conformize is a tool that enables the validation of application service configurations against a set of defined rules",
					Hidden:      true,
				},
			},
			Handler: &commands.AppCommandHandler{},
		},
		&commands.RunnableCommand{
			Command: &commands.Command{
				Expression:  "version",
				Description: "display the version in use",
				Aliases:     []string{"-version", "--version"},
			},
			Handler: &commands.VersionCommandHandler{},
		},
		&commands.RunnableCommand{
			Command: &commands.BlueprintCommand{
				Command: &commands.Command{
					Expression:  "blueprint",
					Description: "execute a blueprint task",
				},
			},
		},
		&commands.RunnableCommand{
			Command: &commands.ScaffoldCommand{
				Command: &commands.Command{
					Expression:  "scaffold",
					Description: "execute a scaffold task",
				},
			},
		},
	}
}
