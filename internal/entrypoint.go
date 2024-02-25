// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package internal

import (
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/conformize/conformize/common/diagnostics"
	"github.com/conformize/conformize/common/format"
	"github.com/conformize/conformize/common/format/colors"
	"github.com/conformize/conformize/common/streams"
	"github.com/conformize/conformize/ui"
)

type Entrypoint struct {
	Args []string
}

func (e *Entrypoint) Run() int {
	binName := filepath.Base(os.Args[0])
	args := os.Args[1:]

	appUi := &ui.UI{
		AppName: binName,
		Args:    args,
	}

	var wg sync.WaitGroup
	diagsChan := make(chan diagnostics.Diagnostic, 10)
	defer close(diagsChan)

	diags := diagnostics.NewDiagnostics()
	diagsStream := streams.NewStream(&wg, diagsChan, 10)
	defer diagsStream.Close()

	diags.Stream(diagsStream)

	go func() {
		const ACCUMULATED_DELAY_THRESHOLD = 10 * time.Millisecond

		outputDelay := 5 * time.Millisecond
		accumulatedDelay := time.Duration(0)

		var lastOutputTime time.Time
		beautySleep := true

		for diag := range diagsChan {
			if beautySleep && !lastOutputTime.IsZero() {
				elapsed := time.Since(lastOutputTime)
				delay := outputDelay - elapsed
				if delay > 0 {
					time.Sleep(delay)
					accumulatedDelay += delay
				}
				outputDelay -= 50 * time.Microsecond
			}

			switch diag.GetType() {
			case diagnostics.Info:
				streams.Output().Writeln(format.Formatter().Color(colors.Blue).Dimmed().Format(diag.String()))
			case diagnostics.Warning:
				streams.Output().Writeln(format.Formatter().Bold().Color(colors.Yellow).Format(diag.String()))
			case diagnostics.Error:
				streams.Error().Writeln(format.Formatter().Color(colors.Red).Format(diag.String()))
			}

			lastOutputTime = time.Now()
			wg.Done()

			beautySleep = accumulatedDelay < ACCUMULATED_DELAY_THRESHOLD && outputDelay >= 1*time.Millisecond
		}
	}()

	appUi.Run(diags)
	wg.Wait()
	if diags.HasErrors() {
		return 1
	}
	return 0
}
