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
	"github.com/conformize/conformize/common/format"
	"github.com/conformize/conformize/common/format/colors"
	"github.com/conformize/conformize/common/util"
	"github.com/conformize/conformize/internal/blueprint"
	"github.com/conformize/conformize/internal/blueprint/execution"
	"github.com/conformize/conformize/internal/ui/options"
	"github.com/conformize/conformize/ui/commands/cmdutil"
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
		if blueprintFilePath, err = util.ResolveFilePath(path); err == nil {
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

	if len(options.Options().WorkDir) == 0 {
		workDir, err := util.ResolveFileBasePath(blueprintFilePath)
		if err == nil {
			if err = os.Chdir(workDir); err != nil {
				diags.Append(diagnostics.Builder().
					Error().
					Details(fmt.Sprintf("Failed to change working directory to %s, reason:\n\n%s", workDir, err.Error())).
					Build(),
				)
				return
			}

			msg := format.Formatter().Detail(format.Tool).Color(colors.Blue).Dimmed().Format(fmt.Sprintf("Working directory: %s\n", workDir))
			diags.Append(diagnostics.Builder().Info().Summary(msg))
		} else {
			diags.Append(diagnostics.Builder().
				Error().
				Details(fmt.Sprintf("Failed to get working directory from blueprint file %s, reason:\n\n%s", blueprintFilePath, err.Error())).
				Build(),
			)
			return
		}
	}

	blprntExecutor := execution.NewBlueprintExecutor()
	blprntExecutor.Execute(blprnt, diags)
	if !diags.HasErrors() {
		diags.Append(diagnostics.Builder().
			Info().
			Summary(format.Formatter().Bold().Detail(format.Ok).Format("Blueprint applied successfully.")).
			Build(),
		)
		return
	}

	diags.Append(
		diagnostics.Builder().
			Error().
			Summary(
				format.Formatter().
					Color(colors.Red).
					Bold().
					Detail(format.FailureWarning).
					Format(fmt.Sprintf("%d rule assertons failed.", len(diags.Errors()))),
			),
	)

	line := format.Formatter().
		Detail(format.Error).
		Color(colors.Red).
		Format("Couldn't apply blueprint.")

	diags.Append(diagnostics.Builder().
		Error().
		Summary(line).
		Build(),
	)
}
