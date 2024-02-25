// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package common

import (
	"fmt"
	"testing"

	"github.com/conformize/conformize/common/ds"
	"github.com/conformize/conformize/common/path"
	"github.com/conformize/conformize/common/pathparser"
)

func TestValuePathEvaluator(t *testing.T) {
	rootNode := ds.NewNode[string, any]()
	nodeRef := rootNode
	for _, step := range []string{"app", "env", "dev", "region", "us-east-1"} {
		nodeRef = nodeRef.AddChild(step)
	}
	expectedValue := "enabled"
	nodeRef.Value = expectedValue

	var path, _ = path.NewFromString("app.env.dev.region.us-east-1")
	var exprPathEvaluator = &ExpressionPathEvaluator{}
	if node, err := exprPathEvaluator.Evaluate(rootNode, path); err != nil {
		t.Errorf("Error evaluating path: %s", err)
	} else {
		if nodeValue, ok := node.Value.(string); !ok || nodeValue != expectedValue {
			t.Errorf("Expected value: %s, but got: %T", expectedValue, node.Value)
		}
	}
}

func TestValuePathEvaluatorWithIndex(t *testing.T) {
	rootNode := ds.NewNode[string, any]()
	nodeRef := rootNode
	for _, step := range []string{"app", "env", "dev"} {
		nodeRef = nodeRef.AddChild(step)
	}

	for _, region := range []string{"us-east-1", "us-west-1"} {
		regionNode := nodeRef.AddChild("regions")
		regionNode.Value = region
	}

	expectedValue := "us-west-1"

	pathParser := pathparser.NewPathParser()
	var steps, _ = pathParser.Parse("$app.'env'.'dev'.'regions'.1")
	var path = path.NewPath(steps)

	var exprPathEvaluator = &ExpressionPathEvaluator{}
	if node, err := exprPathEvaluator.Evaluate(rootNode, path); err != nil {
		t.Errorf("Error walking path: %s", err)
	} else {
		if nodeValue, ok := node.Value.(string); !ok || nodeValue != expectedValue {
			t.Errorf("Expected value: %s, but got: %s", expectedValue, node.Value)
		}
	}
}

func TestValuePathEvaluatorWithIndexAndKey(t *testing.T) {
	rootNode := ds.NewNode[string, any]()
	nodeRef := rootNode
	for _, step := range []string{"app", "env", "dev"} {
		nodeRef = nodeRef.AddChild(step)
	}

	for idx, region := range []string{"us-east-1", "us-west-1"} {
		regionNode := nodeRef.AddChild("regions")
		regionNode.Value = region
		zoneNode := regionNode.AddChild("zone")
		zoneNode.Value = fmt.Sprintf("zone-%d", idx+1)
	}

	expectedValue := "zone-1"

	pathParser := pathparser.NewPathParser()
	var steps, _ = pathParser.Parse("$app.'env'.'dev'.'regions'.0.'zone'")
	var path = path.NewPath(steps)

	var exprPathEvaluator = &ExpressionPathEvaluator{}
	if node, err := exprPathEvaluator.Evaluate(rootNode, path); err != nil {
		t.Errorf("Error walking path: %s", err)
	} else {
		if nodeValue, ok := node.Value.(string); !ok || nodeValue != expectedValue {
			t.Errorf("Expected value: %s, but got: %s", expectedValue, node.Value)
		}
	}
}

func TestValuePathEvaluatorWithIndexAndAttribute(t *testing.T) {
	rootNode := ds.NewNode[string, any]()
	nodeRef := rootNode
	for _, step := range []string{"app", "env", "dev"} {
		nodeRef = nodeRef.AddChild(step)
	}

	for idx, region := range []string{"us-east-1", "us-west-1"} {
		regionNode := nodeRef.AddChild("regions")
		regionNode.Value = region
		zoneNode := regionNode.AddChild("zone")
		zoneNode.Value = fmt.Sprintf("zone-%d", idx+1)
		zoneNode.AddAttribute("meta", fmt.Sprintf("meta-%d", idx+1))
	}

	expectedValue := "meta-2"

	pathParser := pathparser.NewPathParser()
	var steps, _ = pathParser.Parse("$app.'env'.'dev'.'regions'.1.'zone'.attributes.'meta'")
	var path = path.NewPath(steps)

	var exprPathEvaluator = &ExpressionPathEvaluator{}
	if node, err := exprPathEvaluator.Evaluate(rootNode, path); err != nil {
		t.Errorf("Error walking path: %s", err)
	} else {
		if nodeValue, ok := node.Value.(string); !ok || nodeValue != expectedValue {
			t.Errorf("Expected value: %s, but got: %s", expectedValue, node.Value)
		}
	}
}

func TestValuePathEvaluatorReturnsErrorWithIndexOutOfRange(t *testing.T) {
	rootNode := ds.NewNode[string, any]()
	nodeRef := rootNode
	for _, step := range []string{"app", "env", "dev"} {
		nodeRef = nodeRef.AddChild(step)
	}

	for _, region := range []string{"us-east-1", "us-west-1"} {
		regionNode := nodeRef.AddChild("regions")
		regionNode.Value = region
	}

	pathParser := pathparser.NewPathParser()
	var steps, _ = pathParser.Parse("$app.'env'.'dev'.'regions'.3")
	var path = path.NewPath(steps)

	var exprPathEvaluator = &ExpressionPathEvaluator{}
	if _, err := exprPathEvaluator.Evaluate(rootNode, path); err == nil {
		t.Fail()
	}
}
