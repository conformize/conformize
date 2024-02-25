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
	"github.com/conformize/conformize/common/streams"
	"github.com/conformize/conformize/common/streams/colors"
	"github.com/conformize/conformize/common/streams/format"
	"github.com/conformize/conformize/ui"
)

type Entrypoint struct {
	Args []string
}

func (e *Entrypoint) Run() int {
	binName := filepath.Base(os.Args[0])
	args := os.Args[1:]
	ui := &ui.UI{
		AppName: binName,
		Args:    args,
	}

	var wg sync.WaitGroup
	diagsChan := make(chan diagnostics.Diagnostic, 10)
	diagsStream := streams.NewStream(&wg, diagsChan, 10)

	diags := diagnostics.NewDiagnostics()
	diags.Stream(diagsStream)
	defer diagsStream.Close()

	go func() {
		for diag := range diagsChan {
			switch diag.GetType() {
			case diagnostics.Info:
				streams.Instance().Output().Writeln(format.Formatter().Color(colors.Blue).Dimmed().Format(diag.String()))
			case diagnostics.Warning:
				streams.Instance().Output().Writeln(format.Formatter().Bold().Color(colors.Yellow).Format(diag.String()))
			case diagnostics.Error:
				streams.Instance().Error().Writeln(format.Formatter().Color(colors.Red).Format(diag.String()))
			}
			wg.Done()
			<-time.After(10 * time.Millisecond)
		}
	}()

	ui.Run(diags)
	wg.Wait()
	close(diagsChan)
	if diags.HasErrors() {
		return 1
	}
	return 0
}
