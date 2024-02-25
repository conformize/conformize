// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/conformize/conformize/common/diagnostics"
	"github.com/conformize/conformize/common/format"
	"github.com/conformize/conformize/common/streams"
	"github.com/conformize/conformize/internal/ui/options"
	"github.com/conformize/conformize/ui"
)

type Entrypoint struct {
	Args []string
}

func (e *Entrypoint) Run() int {
	binName := filepath.Base(os.Args[0])
	args := os.Args[1:]

	var wg sync.WaitGroup
	diagsChan := make(chan diagnostics.Diagnostic, 10)
	defer close(diagsChan)

	diags := diagnostics.NewDiagnostics()
	diagsStream := streams.NewStream(&wg, diagsChan, 10)
	defer diagsStream.Close()

	var err error
	var parsedArgs []string
	if parsedArgs, err = options.ParseOptions(args); err != nil {
		diags.Append(diagnostics.Builder().Error().Summary(err.Error()))
	}

	appUi := &ui.UI{
		AppName: binName,
		Args:    parsedArgs,
	}

	diags.Stream(diagsStream)

	go func() {
		const ACCUMULATED_DELAY_THRESHOLD = 30 * time.Millisecond

		outputDelay := 5 * time.Millisecond
		accumulatedDelay := time.Duration(0)
		var lastOutputTime time.Time

		beautySleepEnabled := options.Options().Ui.BeautySleep
		beautySleep := beautySleepEnabled
		showTimestamps := options.Options().Ui.Timestamps
		const columnWidth = 14
		const columnGap = 2 // add gap between timestamp and content
		var now time.Time

		for diag := range diagsChan {
			now = time.Now()
			if beautySleep {
				elapsed := time.Since(lastOutputTime)
				delay := outputDelay - elapsed
				if delay > 0 {
					time.Sleep(delay)
					accumulatedDelay += delay
				}
				outputDelay -= 100 * time.Microsecond
			}

			timestamp := format.Formatter().Dimmed().Format(fmt.Sprintf("[%s]", now.Format("15:04:05.000")))
			block := diag.String()

			if showTimestamps {
				lines := strings.Split(block, "\n")
				pad := strings.Repeat(" ", columnWidth+columnGap)
				for i, line := range lines {
					prefix := pad
					if i == 0 {
						prefix = fmt.Sprintf("\n%-*s%s", columnWidth, timestamp, strings.Repeat(" ", columnGap))
					}
					out := prefix + line
					switch diag.GetType() {
					case diagnostics.Info, diagnostics.Warning:
						streams.Output().Writeln(out)
					default:
						streams.Error().Writeln(out)
					}
				}
			} else {
				switch diag.GetType() {
				case diagnostics.Info, diagnostics.Warning:
					streams.Output().Writeln(block)
				default:
					streams.Error().Writeln(block)
				}
			}

			lastOutputTime = time.Now()
			wg.Done()

			beautySleep = beautySleepEnabled && accumulatedDelay < ACCUMULATED_DELAY_THRESHOLD && outputDelay >= 1*time.Millisecond
		}
	}()

	appUi.Run(diags)
	wg.Wait()
	if diags.HasErrors() {
		return 1
	}
	return 0
}
