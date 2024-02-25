// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package commands

import (
	"fmt"
	"os"

	"github.com/conformize/conformize/common/diagnostics"
	"github.com/conformize/conformize/common/util"
	"github.com/conformize/conformize/internal/blueprint"
	"github.com/conformize/conformize/internal/blueprint/execution"
	"github.com/conformize/conformize/ui/cli/commands/cmdutil"
)

type ApplyBlueprintCommandHandler struct{}

func (h *ApplyBlueprintCommandHandler) Handle(c Commandable, args []string, diags diagnostics.Diagnosable) {
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
		if blueprintFilePath, err = util.FilePath(path); err == nil {
			if _, err = os.Stat(blueprintFilePath); err == nil {
				break
			}
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

	blprntUnmarshaller := &blueprint.BlueprintUnmarshaller{Path: blueprintFilePath}
	blprnt, err := blprntUnmarshaller.Unmarshal()
	if err != nil {
		diags.Append(diagnostics.Builder().
			Error().
			Details(fmt.Sprintf("Failed to read blueprint file %s, reason:\n\n%s", blueprintFilePath, err.Error())).
			Build(),
		)
		return
	}

	blprntExecutor := execution.NewBlueprintExecutor()
	blprntExecutor.Execute(blprnt, diags)
	if !diags.HasErrors() {
		diags.Append(diagnostics.Builder().
			Info().
			Details("\nBlueprint has been applied successfully.").
			Build(),
		)
		return
	}

	diags.Append(diagnostics.Builder().
		Error().
		Summary("\nCouldn't apply blueprint.").
		Build(),
	)
}
