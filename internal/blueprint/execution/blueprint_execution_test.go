// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package execution

import (
	"testing"

	"github.com/conformize/conformize/common/diagnostics"
	"github.com/conformize/conformize/internal/blueprint"
)

func TestBlueprintExecutionBuilder(t *testing.T) {
	blueprintPath := "./mocks/blueprint.cnfrm.yaml"
	blueprintUnmarshaller := blueprint.BlueprintUnmarshaller{Path: blueprintPath}
	if blueprint, err := blueprintUnmarshaller.Unmarshal(); err == nil {
		blprntExectBldr := NewBlueprintExecutionBuilder()
		diags := diagnostics.NewDiagnostics()
		blprntExectBldr.WithBlueprint(blueprint)
		blprntExectBldr.Build()
		if diags.HasErrors() {
			t.Errorf("Failed to build blueprint execution, reason: %v", diags.Errors().Print())
		}
	}
}
