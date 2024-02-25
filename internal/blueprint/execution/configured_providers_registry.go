// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package execution

import (
	"errors"

	sdk "github.com/conformize/conformize/internal/providers/api"
)

type ConfiguredProvidersRegistry struct {
	providers map[string]sdk.ConfigurationProvider
}

func (r *ConfiguredProvidersRegistry) Register(alias string, provider sdk.ConfigurationProvider) error {
	if provider == nil {
		return errors.New("provider cannot be nil")
	}

	if alias == "" {
		return errors.New("alias cannot be empty")
	}

	r.providers[alias] = provider
	return nil
}

func (r *ConfiguredProvidersRegistry) Get(name string) (sdk.ConfigurationProvider, bool) {
	provider, exists := r.providers[name]
	return provider, exists
}

func NewConfiguredProvidersRegistry() *ConfiguredProvidersRegistry {
	return &ConfiguredProvidersRegistry{
		providers: make(map[string]sdk.ConfigurationProvider, 10),
	}
}
