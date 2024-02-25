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
	"github.com/conformize/conformize/common/format"
	"github.com/conformize/conformize/common/format/colors"
	"github.com/conformize/conformize/common/util"
	"github.com/conformize/conformize/internal/blueprint"
	"github.com/conformize/conformize/internal/blueprint/execution"
	"github.com/conformize/conformize/ui/commands/cmdutil"
)

type ApplyBlueprintCommandHandler struct{}

func (h *ApplyBlueprintCommandHandler) Handle(c CommandEntry, args []string, diags *diagnostics.Diagnostics) {
	workDir := util.GetWorkDir()
	cwd := workDir
	defer util.SetWorkDir(cwd)

	flags := c.GetFlags()
	flags.Parse(args)

	var paths []string
	defaultPath := false
	var blueprintFilePath string
	if defaultPath = !cmdutil.FlagIsSet(flags, "f"); defaultPath {
		paths = defaultBlueprintFileLocations
	} else {
		flag := flags.Lookup("f")
		blueprintFilePath = flag.Value.String()
		paths = []string{blueprintFilePath}
	}

	var filePath string
	var err error

	formatter := format.Formatter()
	for _, path := range paths {
		if filePath, err = util.ResolveFileRelativePath(workDir, path); err == nil {
			blueprintFilePath = filePath
			break
		}
	}

	if err != nil {
		var errMsg string
		if defaultPath {
			errMsg = "No blueprint file found in current directory."
		} else {
			errMsg = fmt.Sprintf("Blueprint file %s not found.", blueprintFilePath)
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

	workDir, err = util.ResolveFileBasePath(blueprintFilePath)
	if err != nil {
		diags.Append(diagnostics.Builder().
			Error().
			Details(fmt.Sprintf("Failed to resolve working directory for blueprint file %s, reason:\n\n%s", blueprintFilePath, err.Error())).
			Build(),
		)
		return
	}

	util.SetWorkDir(workDir)
	msg := formatter.Detail(format.Tool).Color(colors.Blue).Format(fmt.Sprintf("Working directory: %s", workDir))
	diags.Append(diagnostics.Builder().Info().Summary(msg).Build())

	blprntExecutor := execution.BlueprintExecutor{}
	blprntExecutor.Execute(blprnt, diags)
	if !diags.HasErrors() {
		diags.Append(diagnostics.Builder().
			Info().
			Summary(formatter.Bold().Detail(format.Ok).Format("Blueprint applied successfully.")).
			Build(),
		)
		return
	}

	diags.Append(diagnostics.Builder().
		Error().
		Summary(formatter.
			Detail(format.Error).
			Color(colors.Red).
			Bold().
			Format("Couldn't apply blueprint.")).
		Build(),
	)
}
