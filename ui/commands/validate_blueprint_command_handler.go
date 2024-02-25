// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package commands

import (
	"fmt"

	"github.com/conformize/conformize/common/diagnostics"
	"github.com/conformize/conformize/common/util"
	"github.com/conformize/conformize/internal/blueprint"
	"github.com/conformize/conformize/internal/blueprint/validation"
	"github.com/conformize/conformize/ui/commands/cmdutil"
)

var defaultBlueprintFileLocations = []string{"./blueprint.cnfrm.json", "./blueprint.cnfrm.yaml"}

type ValidateBlueprintCommandHandler struct{}

func (h *ValidateBlueprintCommandHandler) Handle(c CommandEntry, args []string, diags *diagnostics.Diagnostics) {
	flags := c.GetFlags()
	flags.Parse(args)

	var paths []string
	defaultPath := false
	if defaultPath = !cmdutil.FlagIsSet(flags, "f"); defaultPath {
		paths = defaultBlueprintFileLocations
	} else {
		flag := flags.Lookup("f")
		flagVal := flag.Value.String()
		paths = []string{flagVal}
	}

	var blueprintFilePath string
	var filePath string
	var err error
	for _, path := range paths {
		filePath = path
		if blueprintFilePath, err = util.ResolveFilePath(path); err == nil {
			break
		}
	}

	if err != nil {
		var errMsg string
		if defaultPath {
			errMsg = "No blueprint file found in current directory."
		} else {
			errMsg = fmt.Sprintf("Blueprint file %s not found.", filePath)
		}

		diags.Append(diagnostics.Builder().
			Error().
			Summary(errMsg).
			Build(),
		)
		return
	}

	blprntValid := &validation.BlueprintValidation{}
	blprntUnmarshaller := &blueprint.BlueprintUnmarshaller{Path: blueprintFilePath}
	blprnt, err := blprntUnmarshaller.Unmarshal()
	if err != nil {
		diags.Append(diagnostics.Builder().
			Error().
			Details(fmt.Sprintf("Failed to read blueprint, reason:\n\n%s", err.Error())).
			Build(),
		)
		return
	}

	validateDiags := blprntValid.Validate(blprnt)
	if !validateDiags.HasErrors() {
		diags.Append(diagnostics.Builder().
			Info().
			Details("No errors found.").
			Build(),
		)
		return
	}
	diags.Append(validateDiags.Entries()...)
}
