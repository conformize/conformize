// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package blueprint

import (
	"testing"
	"time"

	"math/rand"

	"github.com/conformize/conformize/common/diagnostics"
	"github.com/conformize/conformize/ui/cli"
)

func BenchmarkApplyBlueprint(b *testing.B) {
	blprntMocks := []string{"./internal/blueprint/mocks/blueprint.cnfrm.yaml", "./internal/blueprint/mocks/blueprint.cnfrm.json"}
	for n := 0; n < b.N; n++ {
		cmdRunnner := cli.NewCommandRunner()
		diags := diagnostics.NewDiagnostics()

		rand.NewSource(time.Now().UnixNano())
		mockIdx := rand.Intn(2)
		cmdRunnner.Run([]string{"blueprint", "apply", "-f", blprntMocks[mockIdx]}, diags)
	}
}
