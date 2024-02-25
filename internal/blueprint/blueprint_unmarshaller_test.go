// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package blueprint

import (
	"testing"
)

func TestJSONBlueprintUnmarshalling(t *testing.T) {
	filePath := "./mocks/blueprint.cnfrm.json"
	blueprintUnmarshaller := BlueprintUnmarshaller{Path: filePath}
	if blueprint, err := blueprintUnmarshaller.Unmarshal(); err != nil {
		t.Errorf("Failed to unmarshal blueprint: %v", err.Error())
	} else {
		if blueprint == nil {
			t.Errorf("Blueprint is nil")
		}
	}
}

func TestYAMLBlueprintUnmarshalling(t *testing.T) {
	filePath := "./mocks/blueprint.cnfrm.yaml"
	blueprintUnmarshaller := BlueprintUnmarshaller{Path: filePath}
	if blueprint, err := blueprintUnmarshaller.Unmarshal(); err != nil {
		t.Errorf("Failed to unmarshal blueprint: %v", err.Error())
	} else {
		if blueprint == nil {
			t.Errorf("Blueprint is nil")
		}
	}
}
