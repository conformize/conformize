// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package commands

import (
	"runtime"

	"github.com/conformize/conformize/common/diagnostics"
	"github.com/conformize/conformize/common/format"
	"github.com/conformize/conformize/common/format/colors"
	"github.com/conformize/conformize/common/streams"
	"github.com/conformize/conformize/resources"
)

type VersionCommandHandler struct {
}

func (h *VersionCommandHandler) Handle(c CommandEntry, args []string, _ *diagnostics.Diagnostics) {
	streams.Output().Writef(
		"%s\nConformize v%s\nrunning on %s %s\n",
		format.Formatter().Color(colors.Grey).Format(resources.ASCII_LOGO()),
		resources.VersionString(),
		runtime.GOOS,
		runtime.GOARCH,
	)
}
