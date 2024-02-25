// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package commandregistry

import (
	"flag"

	"github.com/conformize/conformize/ui/commands"
)

func blueprintCommmandFlags(name string) *flag.FlagSet {
	flags := flag.NewFlagSet(name, flag.ContinueOnError)
	flags.String("f", "", "specifies path to a blueprint file")
	return flags
}

func applyBlueprintCommandFlags() *flag.FlagSet {
	return blueprintCommmandFlags("blueprint apply")
}

func validateBlueprintCommandFlags() *flag.FlagSet {
	return blueprintCommmandFlags("blueprint validate")
}

func supportedCommands() []commands.CommandEntry {
	return []commands.CommandEntry{
		&commands.Command{
			Expression:  "blueprint",
			Description: "manage blueprints",
			Handler:     nil,
			Subcommands: []commands.CommandEntry{
				&commands.Command{
					Expression:  "validate",
					Description: "validate blueprint",
					Handler:     &commands.ValidateBlueprintCommandHandler{},
					Subcommands: nil,
					Flags:       applyBlueprintCommandFlags,
					Hidden:      false,
				},
				&commands.Command{
					Expression:  "apply",
					Description: "apply blueprint",
					Flags:       validateBlueprintCommandFlags,
					Handler:     &commands.ApplyBlueprintCommandHandler{},
				},
			},
			Hidden: false,
		},
		commands.BlueprintScaffoldCommand(),
		&commands.Command{
			Expression:  "version",
			Description: "output version information",
			Aliases:     []string{"--version", "-version"},
			Handler:     &commands.VersionCommandHandler{},
			Subcommands: nil,
			Flags:       nil,
			Hidden:      false,
		},
	}
}
