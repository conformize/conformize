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

func TestBlueprintExecutorExecutesBlueprintSuccessfully(t *testing.T) {
	blueprintPath := "../mocks/blueprint.cnfrm.yaml"
	BlueprintUnmarshaller := &blueprint.BlueprintUnmarshaller{Path: blueprintPath}

	if blueprint, err := BlueprintUnmarshaller.Unmarshal(); err != nil {
		t.Errorf("Blueprint Unmarshalling failed: %v", err)
	} else {
		blueprintExecutor := BlueprintExecutor{}
		diags := diagnostics.NewDiagnostics()

		blueprintExecutor.Execute(blueprint, diags)
		if diags.HasErrors() {
			t.Errorf("Blueprint Execution failed, rason: %s", diags.Errors())
		}
	}
}
