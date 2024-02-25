// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package ui

import (
	"github.com/conformize/conformize/common/diagnostics"
	"github.com/conformize/conformize/internal/ui"
	"github.com/conformize/conformize/ui/cli"
)

type UI struct {
	AppName string
	Args    []string
	Options ui.Options
}

func (ui *UI) Run(diags diagnostics.Diagnosable) {
	cli := &cli.Cli{
		AppName: ui.AppName,
		Args:    ui.Args,
		Options: ui.Options,
	}
	cli.Run(diags)
}
