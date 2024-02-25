// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package pathparser

import (
	"testing"

	"github.com/conformize/conformize/common/path"
)

func TestPathParser(t *testing.T) {
	var testCases = []struct {
		pathStr  string
		expected path.Steps
	}{
		{
			pathStr:  "$appDev.'test'.attributes.'test'",
			expected: path.Steps{path.ObjectStep("appDev"), path.KeyStep("test"), path.AttributeStep("test")},
		},
		{
			pathStr:  "$appDev.'test'.length",
			expected: path.Steps{path.ObjectStep("appDev"), path.KeyStep("test"), path.PropertyStep("length")},
		},
	}

	for _, testCase := range testCases {
		pathParser := NewPathParser()
		var steps, err = pathParser.Parse(testCase.pathStr)
		if err != nil {
			t.Errorf("Failed to parse path, reason: %s", err.Error())
		}

		if len(steps) != len(testCase.expected) {
			t.Errorf("expected %d number of steps, got %d", len(testCase.expected), len(steps))
		}

		var expectedSteps = testCase.expected
		for step, hasNext := steps.Next(); hasNext; step, hasNext = steps.Next() {
			if expectedStep, _ := expectedSteps.Next(); !step.Equal(expectedStep) {
				t.Errorf("expected %v, got %v", expectedStep, step)
			}
		}
	}
}

func TestPathParserWithIndex(t *testing.T) {
	var testCases = []struct {
		pathStr  string
		expected path.Steps
	}{
		{
			pathStr:  "$appDev.'test'.0",
			expected: path.Steps{path.ObjectStep("appDev"), path.KeyStep("test"), path.IndexStep("0")},
		},
	}

	for _, testCase := range testCases {
		pathParser := NewPathParser()
		var steps, err = pathParser.Parse(testCase.pathStr)
		if err != nil {
			t.Errorf("Failed to parse path, reason: %s", err.Error())
		}
		var expectedSteps = testCase.expected
		for step, hasNext := steps.Next(); hasNext; step, hasNext = steps.Next() {
			if expectedStep, _ := expectedSteps.Next(); !step.Equal(expectedStep) {
				t.Errorf("expected %v, got %v", expectedSteps, steps)
			}
		}
	}
}
