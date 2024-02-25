// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package mocks

import (
	"fmt"
	"reflect"

	"github.com/conformize/conformize/common/ds"
	"github.com/conformize/conformize/common/functions"
)

func MergeConfigTreesWithOverride() bool {
	mergedConfigTree := functions.MergeConfigTrees(MockedConfigTrees()...)
	return reflect.DeepEqual(mergedConfigTree, MockedExpectedTree())
}

func MockedExpectedTree() *ds.Node[string, interface{}] {
	expectedTree := ds.NewNode[string, interface{}]()
	expectedAppNode := expectedTree.AddChild("app")
	expectedApiNode := expectedTree.AddChild("api")
	expectedApiRateLimitNode := expectedApiNode.AddChild("rate-limit")
	expectedApiRateLimitNode.Value = 60
	expectedAppVersionNode := expectedAppNode.AddChild("version")
	expectedAppVersionNode.Value = "1.1"
	expectedAppNameNode := expectedAppNode.AddChild("name")
	expectedAppNameNode.Value = "CoolestApp"
	expectedDbNode := expectedAppNode.AddChild("database")
	expectedDbType := expectedDbNode.AddChild("type")
	expectedDbType.Value = "sql"
	expectedDbPortNode := expectedDbNode.AddChild("port")
	expectedDbPortNode.Value = 3307
	expectedDbHostNode := expectedDbNode.AddChild("host")
	expectedDbHostNode.Value = "localhost"

	fmt.Printf("\nExpected tree:\n\n")
	expectedTree.PrintTree()

	return expectedTree
}

func MockedConfigTrees() []*ds.Node[string, interface{}] {
	initialConfigTree := ds.NewNode[string, interface{}]()
	appNode := initialConfigTree.AddChild("app")
	versionNode := appNode.AddChild("version")
	versionNode.Value = "1.0"
	appNameNode := appNode.AddChild("name")
	appNameNode.Value = "MyApp"
	dbNode := appNode.AddChild("database")
	dbType := dbNode.AddChild("type")
	dbType.Value = "sql"
	dbPortNode := dbNode.AddChild("port")
	dbPortNode.Value = 3306

	fmt.Printf("\nInitial config tree:\n\n")
	initialConfigTree.PrintTree()

	updatedConfigTree := ds.NewNode[string, interface{}]()
	updatedAppNode := updatedConfigTree.AddChild("app")
	updatedAppVersionNode := updatedAppNode.AddChild("version")
	updatedAppVersionNode.Value = "1.1"
	updatedAppNameNode := updatedAppNode.AddChild("name")
	updatedAppNameNode.Value = "CoolestApp"
	updatedDbNode := updatedAppNode.AddChild("database")
	updatedDbPortNode := updatedDbNode.AddChild("port")
	updatedDbPortNode.Value = 3307
	updatedDbHostNode := updatedDbNode.AddChild("host")
	updatedDbHostNode.Value = "localhost"

	fmt.Printf("\nOverridden config tree:\n\n")
	updatedConfigTree.PrintTree()

	apiConfigTree := ds.NewNode[string, interface{}]()
	apiNode := apiConfigTree.AddChild("api")
	apiRateLimittNode := apiNode.AddChild("rate-limit")
	apiRateLimittNode.Value = 100

	fmt.Printf("\nAPI config tree:\n\n")
	apiConfigTree.PrintTree()

	apiOverrideTree := ds.NewNode[string, interface{}]()
	apiOverrideNode := apiOverrideTree.AddChild("api")
	apiRateLimitOverrideNode := apiOverrideNode.AddChild("rate-limit")
	apiRateLimitOverrideNode.Value = 60

	fmt.Printf("\nOverridden API config tree:\n\n")
	apiOverrideTree.PrintTree()

	return []*ds.Node[string, interface{}]{initialConfigTree, apiConfigTree, updatedConfigTree, apiOverrideTree}
}
