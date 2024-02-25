// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package execution

import (
	"errors"
	"sync"

	sdk "github.com/conformize/conformize/internal/providers/api"
)

type ConfiguredConfigurationProvidersRegistry struct {
	providers map[string]sdk.ConfigurationProvider
}

var _configurationProvidersRegistry *ConfiguredConfigurationProvidersRegistry
var _configurationProvidersRegistryOnce sync.Once

func (r *ConfiguredConfigurationProvidersRegistry) Register(alias string, provider sdk.ConfigurationProvider) error {
	if provider == nil {
		return errors.New("provider cannot be nil")
	}

	if alias == "" {
		return errors.New("alias cannot be empty")
	}

	if _, exists := r.providers[alias]; exists {
		return errors.New("provider already registered")
	}

	r.providers[alias] = provider
	return nil
}

func (r *ConfiguredConfigurationProvidersRegistry) Get(name string) (sdk.ConfigurationProvider, bool) {
	provider, exists := r.providers[name]
	return provider, exists
}

func ConfiguredProvidersRegistry() *ConfiguredConfigurationProvidersRegistry {
	_configurationProvidersRegistryOnce.Do(func() {
		_configurationProvidersRegistry = &ConfiguredConfigurationProvidersRegistry{
			providers: make(map[string]sdk.ConfigurationProvider),
		}
	})
	return _configurationProvidersRegistry
}
