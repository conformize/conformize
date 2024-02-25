// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package commands

import (
	"github.com/conformize/conformize/common/diagnostics"
)

type predicateMeta struct {
	Value     interface{}   `json:"$value" yaml:"$value"`
	Predicate string        `json:"predicate" yaml:"predicate"`
	Arguments []interface{} `json:"arguments" yaml:"arguments"`
}

type ScaffoldPredicateCommandHandler struct{}

func (h *ScaffoldPredicateCommandHandler) Handle(c Commandable, args []string, diags diagnostics.Diagnosable) {

}
