// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package blueprint

import (
	"fmt"
	"strings"

	"github.com/conformize/conformize/serialization"

	"gopkg.in/yaml.v2"
)

type blueprintFormat int

const (
	JSONBlueprint blueprintFormat = iota
	YAMLBlueprint
)

type BlueprintUnmarshaller struct {
	Path string
}

func (b *BlueprintUnmarshaller) Unmarshal() (*Blueprint, error) {
	if !strings.HasSuffix(b.Path, ".cnfrm.json") && !strings.HasSuffix(b.Path, ".cnfrm.yaml") {
		return nil, fmt.Errorf("Blueprint file must have .cnfrm.json or .cnfrm.yaml extension")
	}

	path := b.Path
	var err error
	fileSrc, err := serialization.NewFileSource(path)
	if err != nil {
		return nil, fmt.Errorf("failed to create file source for blueprint: %w", err)
	}

	var content []byte
	content, err = fileSrc.Read()
	if err != nil {
		return nil, err
	}

	var blueprint Blueprint
	if err = yaml.Unmarshal(content, &blueprint); err != nil {
		return nil, err
	}
	return &blueprint, nil
}
