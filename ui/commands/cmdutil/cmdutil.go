// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package cmdutil

import "flag"

func FlagIsSet(flags *flag.FlagSet, name string) bool {
	isSet := false
	flags.Visit(func(f *flag.Flag) {
		if f.Name == name {
			isSet = true
		}
	})
	return isSet
}

func IsHelpCommand(cmd string) bool {
	helpCmds := map[string]struct{}{
		"help":   {},
		"-help":  {},
		"--help": {},
		"-h":     {},
	}
	_, ok := helpCmds[cmd]
	return ok
}

func IsVersionCommand(cmd string) bool {
	versionCmds := map[string]struct{}{
		"version":   {},
		"-version":  {},
		"--version": {},
	}
	_, ok := versionCmds[cmd]
	return ok
}
