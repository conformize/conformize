// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package elements

import (
	"fmt"

	"github.com/conformize/conformize/common/diagnostics"
	"github.com/conformize/conformize/common/typed"
)

type ArgumentMeta struct {
	Value     typed.RawValue
	Sensitive bool
	Path      string
}

func (argMeta *ArgumentMeta) String() string {
	if argMeta.Value == nil {
		return "<nil>"
	}

	if argMeta.Sensitive {
		return "<sensitive>"
	}

	return fmt.Sprintf("%v", argMeta.Value)
}

type RuleMeta struct {
	Name          string
	Provider      string
	Predicate     string
	ValuePath     string
	ArgumentsMeta *ArgumentMeta
	Diagnostics   *diagnostics.Diagnostics
}
