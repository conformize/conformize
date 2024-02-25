// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package commandregistry

import (
	"github.com/conformize/conformize/ui/cli/commands"
)

func supportedCommands() []commands.CommandRunnable {
	return []commands.CommandRunnable{
		&commands.RunnableCommand{
			Command: &commands.RunnableCommand{
				Command: &commands.Command{
					Expression:  "conformize",
					Description: "Conformize is a tool for validating configuration values accross multiple sources",
					Hidden:      true,
				},
			},
			Handler: &commands.AppCommandHandler{},
		},
		&commands.RunnableCommand{
			Command: &commands.Command{
				Expression:  "version",
				Description: "displays the version in use",
				Aliases:     []string{"-version", "--version"},
			},
			Handler: &commands.VersionCommandHandler{},
		},
		&commands.RunnableCommand{
			Command: &commands.BlueprintCommand{
				Command: &commands.Command{
					Expression:  "blueprint",
					Description: "executes blueprint tasks",
				},
			},
		},
		&commands.RunnableCommand{
			Command: &commands.ScaffoldCommand{
				Command: &commands.Command{
					Expression:  "scaffold",
					Description: "creates a blueprint or predicte scaffold",
				},
			},
		},
	}
}
