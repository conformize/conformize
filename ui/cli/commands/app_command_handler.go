// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package commands

import (
	"github.com/conformize/conformize/common/diagnostics"
	"github.com/conformize/conformize/common/streams"
	"github.com/conformize/conformize/common/streams/colors"
	"github.com/conformize/conformize/common/streams/format"
	"github.com/conformize/conformize/resources"
)

type AppCommandHandler struct{}

func (h *AppCommandHandler) Handle(c Commandable, args []string, diags diagnostics.Diagnosable) {
	streams.Instance().
		Output().
		Writef(
			"%s\n%s\n\n%s\n",
			format.Formatter().Color(colors.Grey).Format(resources.ASCII_LOGO()),
			c.GetDescription(),
			"To see details on available commands and how to use them, please run:\n\nconformize help",
		)
}
