// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package blueprint

import (
	"testing"

	"github.com/conformize/conformize/internal/providers"
)

func TestProviderConfigurerWithConfigurationBlock(t *testing.T) {
	blprntUnmarshaller := BlueprintUnmarshaller{Path: "./mocks/blueprint.cnfrm.json"}
	if blprnt, err := blprntUnmarshaller.Unmarshal(); err != nil {
		t.Errorf("Failed to unmarshal blueprint: %s", err.Error())
	} else {
		if blprnt == nil {
			t.Errorf("Blueprint is nil")
		}
		providerConfigurer := ProviderConfigurer()
		src := blprnt.Sources["devEnv"]
		prvdrFactory := providers.ProviderFactory()
		if provider, err := prvdrFactory.Provider(src.Provider); err != nil {
			t.Errorf("Failed to get provider: %s", err.Error())
		} else {
			if err := providerConfigurer.Configure(provider, &src); err != nil {
				t.Errorf("Failed to configure provider: %s", err.Error())
			}
		}
	}
}
