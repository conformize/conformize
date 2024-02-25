// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package cli

import (
	"fmt"
	"os"

	"github.com/conformize/conformize/common/diagnostics"
	"github.com/conformize/conformize/internal/ui/options"
)

type Cli struct {
	AppName string
	Args    []string
	Options options.GlobalOptions
}

func (cli *Cli) Run(diags diagnostics.Diagnosable) {
	args := cli.prepareArguments(cli.Args)

	commandRunner := NewCommandRunner()

	var err error
	var parsedArgs []string
	if parsedArgs, err = options.ParseOptions(args); err != nil {
		diags.Append(diagnostics.Builder().Error().Summary(err.Error()))
		return
	}

	workDir := options.Options().WorkDir
	if (len(workDir) > 0) && workDir != "." {
		if err = os.Chdir(workDir); err != nil {
			diags.Append(diagnostics.Builder().Error().Summary(fmt.Sprintf("Error changing working directory: %s", err.Error())))
			return
		}
		diags.Append(diagnostics.Builder().Info().Summary(fmt.Sprintf("Working directory: %s", workDir)))
	}

	commandRunner.Run(parsedArgs, diags)
}

func (cli *Cli) prepareArguments(args []string) []string {
	if len(args) == 0 {
		return []string{cli.AppName}
	}
	return args
}
