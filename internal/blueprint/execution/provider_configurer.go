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
	"github.com/conformize/conformize/common/pathparser"
	"github.com/conformize/conformize/common/typed"
	"github.com/conformize/conformize/common/util"
	"github.com/conformize/conformize/internal/providers/aggregate"
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

		configFilePath, err = util.ResolveFilePath(configFilePath)

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

	if _, ok := provider.(*aggregate.AggregateProvider); !ok {
		prvdrConfigCtx.ProvidersDependenciesGraph.AddEdge(prvdrConfigCtx.Alias, "")
		return nil
	}

	var aggrPathsVal typed.Valuable
	aggrPathsVal, err = providerConfReq.GetAtPath("paths")
	if err != nil {
		return err
	}

	var aggrPaths []string
	err = aggrPathsVal.As(&aggrPaths)
	if err != nil {
		return err
	}

	for _, pathStr := range aggrPaths {
		pathParser := pathparser.NewPathParser()
		valPathSteps, err := pathParser.Parse(pathStr)
		if err != nil {
			return err
		}

		if len(valPathSteps) == 0 {
			return nil
		}

		prvdrConfigCtx.ProvidersDependenciesGraph.AddEdge(prvdrConfigCtx.Alias, valPathSteps[0].String())
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
