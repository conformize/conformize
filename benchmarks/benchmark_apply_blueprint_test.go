// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package blueprint

import (
	"testing"

	"github.com/conformize/conformize/common/diagnostics"
	"github.com/conformize/conformize/ui/cli"
)

func BenchmarkApplyBlueprint(b *testing.B) {
	for n := 0; n < b.N; n++ {
		cmdRunnner := cli.NewCommandRunner()
		diags := diagnostics.NewDiagnostics()
		cmdRunnner.Run([]string{"blueprint", "apply", "-f", "./internal/blueprint/mocks/blueprint_aws_source.cnfrm.yaml"}, diags)
	}
}
