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
)

type BlueprintCommandHandler struct{}

func (h *BlueprintCommandHandler) Handle(c CommandRunnable, args []string, diags diagnostics.Diagnosable) {
	streams.Instance().Output().Writeln("Execute blueprint tasks")
}
