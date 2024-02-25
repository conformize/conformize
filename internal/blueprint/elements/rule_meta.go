// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package elements

import (
	"fmt"
	"strings"

	"github.com/conformize/conformize/common/typed"
)

type ArgumentMeta struct {
	Value     typed.RawValue
	Sensitive bool
	Path      string
}

func (argMeta *ArgumentMeta) String() string {
	strBldr := strings.Builder{}
	strBldr.WriteString(strings.Repeat(" ", 2))

	var val typed.RawValue
	if argMeta.Value == nil {
		val = "<nil>"
	} else if argMeta.Sensitive {
		val = "<sensitive>"
	} else {
		val = argMeta.Value
	}

	strBldr.WriteString(fmt.Sprintf("value: %v", val))
	if argMeta.Path != "" {
		strBldr.WriteString(fmt.Sprintf(", path: %s", argMeta.Path))
	}
	return strBldr.String()
}

type RuleMeta struct {
	Provider      string
	Predicate     string
	ValuePath     string
	ArgumentsMeta []ArgumentMeta
}
