// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package commands

type ScaffoldCommand struct {
	*Command
}

func (c *ScaffoldCommand) GetSubcommands() []CommandRunnable {
	return []CommandRunnable{
		&RunnableCommand{
			Command: &ScaffoldPredicateCommand{
				Command: &Command{
					Expression:  "predicate",
					Description: "Outputs a scaffold for specified predicate",
				},
			},
			Handler: &ScaffoldPredicateCommandHandler{},
		},
	}
}
