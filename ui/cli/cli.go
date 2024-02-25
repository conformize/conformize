// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package cli

import (
	"github.com/conformize/conformize/common/diagnostics"
	"github.com/conformize/conformize/internal/ui/options"
)

type Cli struct {
	AppName string
	Args    []string
	Options options.GlobalOptions
}

func (cli *Cli) Run(diags *diagnostics.Diagnostics) int {
	commandRunner := NewCommandRunner()
	return commandRunner.Run(cli.Args, diags)
}
