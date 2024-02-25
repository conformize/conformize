// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package blueprint

import (
	"reflect"
	"sync"

	"github.com/conformize/conformize/common/functions"
	"github.com/conformize/conformize/common/util"
	"github.com/conformize/conformize/internal/blueprint/elements"
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
	if strVal, ok := v.Interface().(string); ok {
		if str, err := functions.InterpolateEnvVars(strVal); err != nil {
			return err
		} else {
			m.SetMapIndex(k, reflect.ValueOf(str))
		}
	}
	return nil
}

type providerConfigurer struct{}

func (prvdrCnfgr *providerConfigurer) Configure(provider api.ConfigurationProvider, sourceConfig elements.ConfigurationSource) error {
	providerConf := sourceConfig.Config
	if providerConf == nil && sourceConfig.ConfigFile != nil {
		configFilePath, err := functions.InterpolateEnvVars(*sourceConfig.ConfigFile)
		if err != nil {
			return err
		}

		absConfigFilePath, err := util.ResolveFilePath(configFilePath)
		if err != nil {
			return err
		}

		fileSrc, err := serialization.NewFileSource(absConfigFilePath)
		if err != nil {
			return err
		}

		fileContent, err := fileSrc.Read()
		if err != nil {
			return err
		}

		err = yaml.Unmarshal(fileContent, &providerConf)
		if err != nil {
			return err
		}
	}

	if err := reflectwalk.Walk(providerConf, &mapValueWalker{}); err != nil {
		return err
	}

	providerConfigSchema := provider.ConfigurationSchema()
	providerConfReq := api.NewConfigurationRequest(providerConfigSchema)
	if err := providerConfReq.Set(providerConf); err != nil {
		return err
	}
	return provider.Configure(providerConfReq)
}

func (prvdrCnfgr *providerConfigurer) Query(provider api.ConfigurationProvider, queryOptions map[string]interface{}) (*api.ProviderDataRequest, error) {
	var provisionDataReq *api.ProviderDataRequest
	provisionDataReqSchema := provider.ProvisionDataRequestSchema()
	provisionDataReq = api.NewProviderDataRequest(provisionDataReqSchema)

	if err := reflectwalk.Walk(&queryOptions, &mapValueWalker{}); err != nil {
		return nil, err
	}
	if err := provisionDataReq.Set(queryOptions); err != nil {
		return nil, err
	}
	return provisionDataReq, nil
}

var instance *providerConfigurer
var once sync.Once

func ProviderConfigurer() *providerConfigurer {
	once.Do(func() {
		instance = &providerConfigurer{}
	})
	return instance
}
