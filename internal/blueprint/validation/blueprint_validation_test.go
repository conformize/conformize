// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package validation

import (
	"testing"

	"github.com/conformize/conformize/internal/blueprint"
)

func TestBlueprintValidationSucceedsWithValidSchema(t *testing.T) {
	blueprintUnmarshaller := blueprint.BlueprintUnmarshaller{Path: "../mocks/blueprint.cnfrm.yaml"}
	blueprint, _ := blueprintUnmarshaller.Unmarshal()
	BlueprintValidation := BlueprintValidation{}
	if diags := BlueprintValidation.Validate(blueprint); diags.HasErrors() {
		t.Errorf("Validation failed")
	}
}

func TestBlueprintValidationFailsWithInvalidSchema(t *testing.T) {
	blueprintUnmarshaller := blueprint.BlueprintUnmarshaller{Path: "../mocks/blueprint.missing.version.cnfrm.json"}
	blueprint, _ := blueprintUnmarshaller.Unmarshal()
	BlueprintValidation := BlueprintValidation{}
	if diags := BlueprintValidation.Validate(blueprint); !diags.HasErrors() {
		t.Errorf("Validation succeeded when it should have failed due to missing version")
	}
}

func TestBlueprintValidation(t *testing.T) {
	blueprintPath := "../mocks/blueprint.cnfrm.yaml"
	blueprintUnmarshaller := blueprint.BlueprintUnmarshaller{Path: blueprintPath}
	blueprint, _ := blueprintUnmarshaller.Unmarshal()
	blueprintValidation := &BlueprintValidation{}
	if diags := blueprintValidation.Validate(blueprint); diags.HasErrors() {
		t.Errorf("Failed to validate blueprint")
	}
}
