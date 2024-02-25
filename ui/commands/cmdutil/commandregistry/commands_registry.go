// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package commandregistry

import (
	"sync"

	"github.com/conformize/conformize/common/ds"
	"github.com/conformize/conformize/ui/commands"
)

type commandRegistry struct {
	cmdTree *ds.Node[string, commands.CommandRunnable]
}

func newCommandRegistry() *commandRegistry {
	cmdR := &commandRegistry{
		cmdTree: ds.NewNode[string, commands.CommandRunnable](),
	}
	cmdR.registerSupportedCommands()
	return cmdR
}

func (cmdR *commandRegistry) Register(cmd commands.CommandRunnable) {
	cmdNode := cmdR.cmdTree.AddChild(cmd.GetExpression())
	cmdNode.Value = cmd
	for _, alias := range cmd.GetAliases() {
		cmdAlias := cmdR.cmdTree.AddChild(alias)
		cmdAlias.Value = cmd
	}
	cmdR.registerSubcommands(cmdNode)
}

func (cmdR *commandRegistry) GetCommands() []commands.CommandRunnable {
	cmds := make([]commands.CommandRunnable, 0)
	processedComamnds := make(map[commands.CommandRunnable]struct{})
	for _, cmdNodes := range cmdR.cmdTree.Children() {
		cmd := cmdNodes.First().Value
		if _, exists := processedComamnds[cmd]; !exists {
			if cmd.IsHidden() {
				continue
			}
			cmds = append(cmds, cmd)
			processedComamnds[cmd] = struct{}{}
		}
	}
	return cmds
}

func (cmdR *commandRegistry) GetCommand(cmd string) *ds.Node[string, commands.CommandRunnable] {
	cmdNode, ok := cmdR.cmdTree.GetChildren(cmd)
	if !ok {
		return nil
	}
	return cmdNode.First()
}

func (cmdR *commandRegistry) registerSupportedCommands() {
	for _, c := range supportedCommands() {
		cmdR.Register(c)
	}
}

func (cmdR *commandRegistry) registerSubcommands(cmdNode *ds.Node[string, commands.CommandRunnable]) {
	cmd := cmdNode.Value
	for _, subCmd := range cmd.GetSubcommands() {
		subCmdNode := cmdNode.AddChild(subCmd.GetExpression())
		subCmdNode.Value = subCmd

		for _, alias := range subCmd.GetAliases() {
			aliasNode := cmdNode.AddChild(alias)
			aliasNode.Value = subCmd
		}
		cmdR.registerSubcommands(subCmdNode)
	}
}

var (
	instance *commandRegistry
	once     = sync.Once{}
)

func Instance() *commandRegistry {
	once.Do(func() {
		instance = newCommandRegistry()
	})
	return instance
}
