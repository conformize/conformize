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

	streamCounter := atomic.Int32{}
	streamCounter.Store(0)
	diagsStream := streams.NewStream[diagnostics.Diagnostic](50)
	diags.Stream(diagsStream)

	appUi := &ui.UI{
		AppName: binName,
		Args:    parsedArgs,
	}

	doneChan := make(chan struct{})
	defer close(doneChan)

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()

		const ACCUMULATED_DELAY_THRESHOLD = 20 * time.Millisecond
		interval := 4 * time.Millisecond
		minInterval := 1 * time.Millisecond
		intervalDecrement := 100 * time.Microsecond
		accumulatedDelay := time.Duration(0)

		beautySleep := options.Options().Ui.BeautySleep
		showTimestamps := options.Options().Ui.Timestamps
		const columnWidth = 14
		const columnGap = 2

		buf := bytes.Buffer{}
		buf.Grow(1024)
		outputFormatter := format.Formatter()

		var diag diagnostics.Diagnostic

		var doCloseStream = false
	outputLoop:
		for !diagsStream.Closed() {
			select {
			case <-doneChan:
				doCloseStream = true
			default:
				isEmpty := diagsStream.IsEmpty()
				if isEmpty && doCloseStream {
					diagsStream.Close()
					break outputLoop
				}

				for !isEmpty {
					diag, err = diagsStream.Read()

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
								timestamp := outputFormatter.Format(fmt.Sprintf("[%s]", time.Now().Format("15:04:05.000")))
								prefix = fmt.Sprintf("%-*s%s", columnWidth, timestamp, strings.Repeat(" ", columnGap))
							}
							buf.WriteString("\n" + prefix + line)
						}
					} else {
						buf.WriteString("\n" + diag.String())
					}

					output := buf.String()
					switch diag.GetType() {
					case diagnostics.Info, diagnostics.Warning:
						streams.Output().Writeln(output)
					default:
						streams.Error().Writeln(output)
					}
					buf.Reset()
					buf.Grow(1024)

					beautySleep = beautySleep &&
						accumulatedDelay < ACCUMULATED_DELAY_THRESHOLD &&
						interval >= minInterval

					isEmpty = diagsStream.IsEmpty()
				}

			}
		}
	}()

	go appUi.Run(diags, doneChan)
	wg.Wait()
	if diags.HasErrors() {
		return 1
	}
	return 0
}
