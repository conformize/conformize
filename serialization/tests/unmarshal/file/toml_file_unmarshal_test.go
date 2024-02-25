// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package tests

import (
	"fmt"
	"testing"
	"time"

	"github.com/conformize/conformize/common/ds"
	"github.com/conformize/conformize/serialization"
	toml "github.com/conformize/conformize/serialization/unmarshal/toml"
)

func testTomlFileUnmarshalling(filePath string) (*ds.Node[string, any], error) {
	var fileSource, _ = serialization.NewFileSource(filePath)
	var tomlFileUnmarshal = toml.TomlFilelUnmarshal{}
	return tomlFileUnmarshal.Unmarshal(fileSource)
}

func TestSimpleTomlFileUnmarshalling(t *testing.T) {
	startTime := time.Now()
	var content, err = testTomlFileUnmarshalling("../../mocks/app.toml")
	if err != nil {
		t.Fail()
	}
	duration := time.Since(startTime)
	ms := float64(duration) / float64(time.Millisecond)

	fmt.Printf("execution time: %.2f ms\n", ms)
	fmt.Println("unmarshalled content:")
	content.PrintTree()
}
