// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package execution

import (
	"fmt"
	"reflect"
	"sync"

	"github.com/conformize/conformize/common/functions"
	api "github.com/conformize/conformize/internal/providers/api"
	"github.com/conformize/conformize/serialization"

	"github.com/mitchellh/reflectwalk"
	"gopkg.in/yaml.v2"
)

type mapValueWalker struct{}

func (mw *mapValueWalker) Map(m reflect.Value) error {
	return nil
}

func (mw *mapValueWalker) MapElem(m, k, v reflect.Value) error {
	strVal, ok := reflect.TypeAssert[string](v)
	var err error
	if ok {
		if strVal, err = functions.InterpolateEnvVars(strVal); err != nil {
			return err
		}
		m.SetMapIndex(k, reflect.ValueOf(strVal))
	}
	return nil
}

type providerConfigurer struct{}

func (prvdrCnfgr *providerConfigurer) Configure(prvdrConfigCtx *ProviderConfigurationContext) error {
	providerConf := prvdrConfigCtx.Config
	if providerConf.Config == nil && prvdrConfigCtx.Config.ConfigFile == nil {
		return fmt.Errorf("provider configuration for '%s' is empty", prvdrConfigCtx.Alias)
	}

	if providerConf.Config == nil && providerConf.ConfigFile != nil {
		configFilePath, err := functions.InterpolateEnvVars(*providerConf.ConfigFile)
		if err != nil {
			return err
		}

		fileSrc, err := serialization.NewFileSource(configFilePath)
		if err != nil {
			return err
		}

		fileContent, err := fileSrc.Read()
		if err != nil {
			return err
		}

		err = yaml.Unmarshal(fileContent, &providerConf.Config)
		if err != nil {
			return err
		}
	}

	if err := reflectwalk.Walk(providerConf.Config, &mapValueWalker{}); err != nil {
		return err
	}

	provider, found := prvdrConfigCtx.ProvidersRegistry.Get(prvdrConfigCtx.Alias)
	if !found {
		return fmt.Errorf("provider '%s' not found in registry", prvdrConfigCtx.Alias)
	}

	var err error
	providerConfigSchema := provider.ConfigurationSchema()
	providerConfReq := api.NewConfigurationRequest(providerConfigSchema)
	if err := providerConfReq.Set(providerConf.Config); err != nil {
		return err
	}

	err = provider.Configure(providerConfReq)
	if err != nil {
		return err
	}

	return nil
}

var instance *providerConfigurer
var once sync.Once

func ProviderConfigurer() *providerConfigurer {
	once.Do(func() {
		instance = &providerConfigurer{}
	})
	return instance
}
