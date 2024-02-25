// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package schema

import (
	"github.com/conformize/conformize/common/typed"
)

type Schema struct {
	Description     string
	Version         int64
	Deprecated      bool
	DeprecationHint string
	Attributes      map[string]Attributeable
}

func (s *Schema) GetDescription() string {
	return s.Description
}

func (s *Schema) GetVersion() int64 {
	return s.Version
}

func (s *Schema) IsDeprecated() bool {
	return s.Deprecated
}

func (s *Schema) GetDeprecationHint() string {
	return s.DeprecationHint
}

func (s *Schema) GetAttributes() map[string]Attributeable {
	return s.Attributes
}

func (s *Schema) Type() typed.Typeable {
	attributeTypes := map[string]typed.Typeable{}
	for attributeName, attributeValue := range s.Attributes {
		attributeTypes[attributeName] = attributeValue.Type()
	}
	return &typed.ObjectTyped{FieldsTypes: attributeTypes}
}
