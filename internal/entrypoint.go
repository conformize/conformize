// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package internal

import (
	"bytes"
	"fmt"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/conformize/conformize/common/diagnostics"
	"github.com/conformize/conformize/common/format"
	"github.com/conformize/conformize/common/format/colors"
	"github.com/conformize/conformize/common/streams"
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

	diags := diagnostics.NewDiagnostics()

	diagsStream := streams.NewStream[diagnostics.Diagnostic](50)
	diags.Stream(diagsStream)

	appUi := &ui.UI{
		AppName: binName,
		Args:    parsedArgs,
	}

	doneSignal := sync.NewCond(&sync.Mutex{})
	doStream := atomic.Bool{}
	doStream.Store(true)

	done := atomic.Bool{}
	done.Store(false)

	go func() {
		const ACCUMULATED_DELAY_THRESHOLD = 20 * time.Millisecond
		interval := 4 * time.Millisecond
		minInterval := 1 * time.Millisecond
		intervalDecrement := 100 * time.Microsecond
		accumulatedDelay := time.Duration(0)

		beautySleep := options.Options().Ui.BeautySleep
		showTimestamps := options.Options().Ui.Timestamps
		const columnWidth = 14
		const padding = 2

		buf := bytes.Buffer{}
		buf.Grow(64)

		var diag diagnostics.Diagnostic

		var isEmpty bool
		var processed int = 0
		for {
			isEmpty = diagsStream.IsEmpty()
			if !isEmpty {
				diag, err = diagsStream.Read()
				if beautySleep {
					time.Sleep(interval)

					accumulatedDelay += interval
					if interval > minInterval {
						interval -= intervalDecrement
					}
				}

				if showTimestamps {
					buf.WriteString(
						fmt.Sprintf("%-*s%s\n",
							columnWidth,
							fmt.Sprintf("[%s]%s", diag.GetTimestamp().Format("15:04:05.000"), strings.Repeat(" ", padding)),
							diag.String(),
						),
					)
				} else {
					buf.WriteString(fmt.Sprintf("%s\n", diag.String()))
				}

				output := buf.String()
				switch diag.GetType() {
				case diagnostics.Info, diagnostics.Warning:
					streams.Output().Writeln(output)
				default:
					streams.Error().Writeln(output)
				}
				buf.Reset()

				beautySleep = beautySleep &&
					accumulatedDelay < ACCUMULATED_DELAY_THRESHOLD &&
					interval >= minInterval

				processed++
				continue
			}

			if done.Load() && processed == len(diags.Entries()) {
				doneSignal.L.Lock()
				doStream.Store(false)
				doneSignal.Signal()
				doneSignal.L.Unlock()
				return
			}
		}
	}()

	appUi.Run(diags, &done)
	doneSignal.L.Lock()
	for doStream.Load() {
		doneSignal.Wait()
	}
	doneSignal.L.Unlock()

	diagsStream.Close()
	if diags.HasErrors() {
		return 1
	}
	return 0
}
