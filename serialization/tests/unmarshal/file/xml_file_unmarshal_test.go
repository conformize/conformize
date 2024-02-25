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
	xml "github.com/conformize/conformize/serialization/unmarshal/xml"
)

func testXmlFileUnmarshalling(filePath string) (*ds.Node[string, any], error) {
	var fileSource, _ = serialization.NewFileSource(filePath)
	var xmlFileUnmarshal = xml.XmlFileUnmarshal{}
	return xmlFileUnmarshal.Unmarshal(fileSource)
}

func TestSimpleXmlFileUnmarshalling(t *testing.T) {
	startTime := time.Now()
	var content, err = testXmlFileUnmarshalling("../../mocks/app.xml")
	if err != nil {
		t.Fail()
	}
	duration := time.Since(startTime)
	ms := float64(duration) / float64(time.Millisecond)

	fmt.Printf("execution time: %.2f ms\n", ms)
	fmt.Println("unmarshalled content:")
	content.PrintTree()
}

func TestComplexXmlFileUnmarshalling(t *testing.T) {
	startTime := time.Now()
	var content, err = testXmlFileUnmarshalling("../../mocks/logback.xml")
	if err != nil {
		t.Fail()
	}
	duration := time.Since(startTime)
	ms := float64(duration) / float64(time.Millisecond)

	fmt.Printf("execution time: %.2f ms\n", ms)
	fmt.Println("unmarshalled content:")
	content.PrintTree()
}
