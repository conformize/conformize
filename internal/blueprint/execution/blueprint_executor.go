// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package execution

import (
	"fmt"
	"strings"

	"github.com/conformize/conformize/common/diagnostics"
	"github.com/conformize/conformize/common/format"
	"github.com/conformize/conformize/common/format/colors"
	"github.com/conformize/conformize/internal/blueprint"
	"github.com/conformize/conformize/internal/blueprint/elements"
)

type BlueprintExecutor struct{}

func NewBlueprintExecutor() *BlueprintExecutor {
	return &BlueprintExecutor{}
}

func (blprntExec *BlueprintExecutor) Execute(blueprint *blueprint.Blueprint, diags diagnostics.Diagnosable) {
	blueprintExecutionCtx := BlueprintExecutionContext().
		WithBlueprint(blueprint).
		WithDiagnostics(diags)

	blueprintExecutionCtx.Execute()
}

const bulletIndent = "   "
const labelWidth = -14

var msgBldr strings.Builder = strings.Builder{}

func ruleViolationErrorMessage(ruleMeta *elements.RuleMeta, ruleIdx int) string {
	msgBldr.Reset()

	var ruleHeader string
	if len(ruleMeta.Name) > 0 {
		ruleHeader = fmt.Sprintf("Rule '%s':\n\n", ruleMeta.Name)
	} else {
		ruleHeader = fmt.Sprintf("Rule %d:\n\n", ruleIdx+1)
	}

	msgBldr.WriteString(format.Formatter().
		Detail(format.Failure).
		Color(colors.Red).
		Bold().
		Format(ruleHeader))

	writeLine := func(label, value string) {
		line := fmt.Sprintf("%-*s%s\n", labelWidth, label+":", value)
		msgBldr.WriteString(
			bulletIndent +
				format.Formatter().
					Detail(format.Bullet).
					Dimmed().
					Color(colors.Red).
					Format(line),
		)
	}

	writeLine("$value", ruleMeta.ValuePath)
	writeLine("Predicate", ruleMeta.Predicate)
	if ruleMeta.ArgumentsMeta.Value != nil {
		if ok, _ := ruleMeta.ArgumentsMeta.Value.([]interface{}); ok != nil {
			writeLine("Arguments", fmt.Sprintf("%v", ruleMeta.ArgumentsMeta.Value))
		} else {
			writeLine("Argument", ruleMeta.ArgumentsMeta.String())
		}
	}

	return msgBldr.String()
}
