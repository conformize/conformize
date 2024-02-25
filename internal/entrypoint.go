// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package internal

import (
	"fmt"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/conformize/conformize/common/diagnostics"
	"github.com/conformize/conformize/common/format"
	"github.com/conformize/conformize/common/format/colors"
	"github.com/conformize/conformize/common/streams"
	"github.com/conformize/conformize/common/util"
	"github.com/conformize/conformize/internal/ui/options"
	"github.com/conformize/conformize/ui"
)

type Entrypoint struct {
	Args []string
}

func (e Entrypoint) Run() int {
	binName := filepath.Base(e.Args[0])
	args := e.Args[1:]

	var err error
	var parsedArgs []string

	if parsedArgs, err = options.ParseOptions(args); err != nil {
		streams.Error().Writeln(format.Formatter().Color(colors.Red).Format(fmt.Sprintf("Error parsing options: %s", err.Error())))
		return 1
	}

	var wg sync.WaitGroup
	diagsChan := make(chan diagnostics.Diagnostic, 100)
	defer close(diagsChan)

	diags := diagnostics.NewDiagnostics()
	diagsStream := streams.NewStream(&wg, diagsChan, 10)
	diags.Stream(diagsStream)
	defer diagsStream.Close()

	appUi := &ui.UI{
		AppName: binName,
		Args:    parsedArgs,
	}

	workDir := options.Options().WorkDir
	if len(workDir) > 0 {
		util.SetWorkDir(workDir)
	} else {
		workDir = util.GetWorkDir()
	}

	go func() {
		const ACCUMULATED_DELAY_THRESHOLD = 20 * time.Millisecond
		interval := 4 * time.Millisecond
		minInterval := 1 * time.Millisecond
		intervalDecrement := 100 * time.Microsecond
		accumulatedDelay := time.Duration(0)

		beautySleep := options.Options().Ui.BeautySleep
		showTimestamps := options.Options().Ui.Timestamps
		const columnWidth = 14
		const columnGap = 2

		bldr := strings.Builder{}
		outputFormatter := format.Formatter()
		for diag := range diagsChan {
			if beautySleep {
				time.Sleep(interval)
				accumulatedDelay += interval
				if interval > minInterval {
					interval -= intervalDecrement
				}
			}

			if showTimestamps {
				lines := strings.Split(diag.String(), "\n")
				pad := strings.Repeat(" ", columnWidth+columnGap)
				for i, line := range lines {
					prefix := pad
					if i == 0 {
						timestamp := outputFormatter.Dimmed().Format(fmt.Sprintf("[%s]", time.Now().Format("15:04:05.000")))
						prefix = fmt.Sprintf("%-*s%s", columnWidth, timestamp, strings.Repeat(" ", columnGap))
					}
					bldr.WriteString("\n" + prefix + line)
				}
			} else {
				bldr.WriteString("\n" + diag.String())
			}

			output := bldr.String()
			switch diag.GetType() {
			case diagnostics.Info, diagnostics.Warning:
				streams.Output().Writeln(output)
			default:
				streams.Error().Writeln(output)
			}
			bldr.Reset()
			wg.Done()

			beautySleep = beautySleep &&
				accumulatedDelay < ACCUMULATED_DELAY_THRESHOLD &&
				interval >= minInterval
		}
	}()

	appUi.Run(diags)
	wg.Wait()
	if diags.HasErrors() {
		return 1
	}
	return 0
}
