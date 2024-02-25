// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package commandregistry

import (
	"strings"
	"sync"

	"github.com/conformize/conformize/ui/commands"
	"github.com/conformize/conformize/ui/commands/cmdutil"
)

type commandRegistry struct {
	cmds map[string]commands.CommandEntry
}

func newCommandRegistry() *commandRegistry {
	reg := &commandRegistry{
		cmds: make(map[string]commands.CommandEntry, 0),
	}

	for _, cmd := range supportedCommands() {
		reg.Register(cmd)
	}
	return reg
}

func (reg *commandRegistry) Register(cmd commands.CommandEntry) {
	reg.register([]string{}, cmd)
}

func (reg *commandRegistry) register(parentPath []string, cmd commands.CommandEntry) {
	names := append([]string{cmd.GetExpression()}, cmd.GetAliases()...)

	for _, name := range names {
		fullPath := append([]string{}, parentPath...)
		fullPath = append(fullPath, name)
		pathExpr := strings.Join(fullPath, " ")
		if _, exists := reg.cmds[pathExpr]; !exists {
			reg.cmds[pathExpr] = cmd
		}

		for _, subcmd := range cmd.GetSubcommands() {
			reg.register(fullPath, subcmd)
		}
	}
}

func (reg *commandRegistry) GetCommand(args []string) (commands.CommandEntry, []string, bool) {
	argsLen := len(args)
	if argsLen == 0 {
		return &commands.Command{Handler: &commands.AppCommandHandler{}}, args, true
	}

	if argsLen == 1 {
		cmd, ok := reg.cmds[args[0]]
		return cmd, []string{}, ok
	}

	idx := 0
	sep := ""

	var cmdExpr strings.Builder
	for idx < argsLen {
		if strings.HasPrefix(args[idx], "-") || cmdutil.IsHelpCommand(args[idx]) {
			break
		}
		cmdExpr.WriteString(sep)
		cmdExpr.WriteString(args[idx])
		sep = " "
		idx++
	}

	cmdExprStr := cmdExpr.String()
	found, ok := reg.cmds[cmdExprStr]
	return found, args[idx:], ok
}

func (reg *commandRegistry) GetCommands() []commands.CommandEntry {
	cmds := make([]commands.CommandEntry, 0)
	seen := make(map[string]struct{})
	
	for pathExpr, cmd := range reg.cmds {
		// Only include top-level commands (no spaces in path = no parent)
		if !strings.Contains(pathExpr, " ") {
			if _, ok := seen[cmd.GetExpression()]; ok {
				continue
			}
			cmds = append(cmds, cmd)
			seen[cmd.GetExpression()] = struct{}{}
		}
	}
	return cmds
}

var (
	instance *commandRegistry
	once     sync.Once
)

func Instance() *commandRegistry {
	once.Do(func() {
		instance = newCommandRegistry()
	})
	return instance
}
