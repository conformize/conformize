// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package commands

type BlueprintCommand struct {
	*Command
}

func (c *BlueprintCommand) GetSubcommands() []CommandRunnable {
	return []CommandRunnable{
		&RunnableCommand{
			Command: &BlueprintTaskCommand{
				Command: &Command{
					Expression:  "validate",
					Description: "validate blueprint",
				},
			},
			Handler: &ValidateBlueprintCommandHandler{},
		},
		&RunnableCommand{
			Command: &BlueprintTaskCommand{
				Command: &Command{
					Expression:  "apply",
					Description: "apply blueprint",
				},
			},
			Handler: &ApplyBlueprintCommandHandler{},
		},
	}
}
