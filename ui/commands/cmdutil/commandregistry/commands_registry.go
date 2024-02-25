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
		return &commands.Command{
				Expression:  "conformize",
				Description: "Conformize enables validation of application service configurations against defined rules.",
				Hidden:      true,
				Handler:     &commands.AppCommandHandler{},
			},
			args,
			true
	}

	idx := 0
	var cmdExpr strings.Builder
	var found commands.CommandEntry
	for idx < argsLen {
		cmdExpr.WriteString(args[idx])
		if strings.HasPrefix(args[idx], "-") {
			return found, args[idx:], true
		}

		cmd, ok := reg.cmds[cmdExpr.String()]
		if !ok {
			return nil, args[idx:], false
		}
		found = cmd
		idx++
		cmdExpr.WriteString(" ")
	}

	if found != nil {
		return found, args[idx:], true
	}

	return nil, args[idx:], false
}

func (reg *commandRegistry) GetCommands() []commands.CommandEntry {
	cmds := make([]commands.CommandEntry, 0, len(reg.cmds))
	seen := make(map[string]struct{})
	for _, cmd := range reg.cmds {
		if _, ok := seen[cmd.GetExpression()]; ok {
			continue
		}
		cmds = append(cmds, cmd)
		seen[cmd.GetExpression()] = struct{}{}
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
