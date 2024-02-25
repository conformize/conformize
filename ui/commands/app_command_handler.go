// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package commands

import (
	"fmt"

	"github.com/conformize/conformize/common/diagnostics"
	"github.com/conformize/conformize/common/format"
	"github.com/conformize/conformize/common/format/colors"
	"github.com/conformize/conformize/common/streams"
	"github.com/conformize/conformize/resources"
)

type AppCommandHandler struct{}

func (h *AppCommandHandler) Handle(c CommandEntry, args []string, _ *diagnostics.Diagnostics) {
	streams.Output().Writef(
		"%s\n\n%s\n\n%s\n",
		format.Formatter().Color(colors.Grey).Format(resources.ASCII_LOGO()),
		format.Formatter().Bold().Format("Rule-based validation engine for your application configuration."),
		fmt.Sprintf("Run %s to see available commands.", format.Formatter().Bold().Color(colors.Grey).Format("conformize --help")),
	)
}
