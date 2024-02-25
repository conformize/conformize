// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package benchmarks

import (
	"testing"
	"time"

	"math/rand"

	"github.com/conformize/conformize/common/diagnostics"
	"github.com/conformize/conformize/internal"
	"github.com/conformize/conformize/ui/cli"
)

func BenchmarkApplyBlueprint(b *testing.B) {
	blprntMocks := []string{"../internal/blueprint/mocks/blueprint.cnfrm.yaml", "../internal/blueprint/mocks/blueprint.cnfrm.json"}
	for b.Loop() {
		cmdRunnner := cli.NewCommandRunner()
		diags := diagnostics.NewDiagnostics()
		rand.NewSource(time.Now().UnixNano())
		mockIdx := rand.Intn(2)
		cmdRunnner.Run([]string{"blueprint", "apply", "-f", blprntMocks[mockIdx]}, diags)
	}
}

func BenchmarkApplyBlueprintFull(b *testing.B) {
	blprnt := "../internal/blueprint/mocks/blueprint.cnfrm.yaml"
	args := []string{"conformize", "blueprint", "apply", "-f", blprnt, "--no-beauty-sleep"}
	entrypoint := internal.Entrypoint{Args: args}
	for b.Loop() {
		entrypoint.Run()
	}
}
