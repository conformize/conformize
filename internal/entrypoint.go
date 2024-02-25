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

	diagsChan := make(chan diagnostics.Diagnostic, 100)
	defer close(diagsChan)

	var wg sync.WaitGroup
	diags := diagnostics.NewStreamableDiagnostics(&wg)

	diagsStream := streams.NewStream(diagsChan, 100)
	diags.Stream(diagsStream)
	defer diagsStream.Close()

	done := false
	go func() {
		for !done {
			for diag := range diagsChan {
				switch diag.GetType() {
				case diagnostics.Info, diagnostics.Warning:
					streams.Instance().Output().Writeln(diag.String())
				case diagnostics.Error:
					streams.Instance().Error().Writeln(diag.String())
				}
				wg.Done()
				<-time.After(20 * time.Millisecond)
			}
		}
	}()

	ui.Run(diags)
	wg.Wait()
	done = true
	if diags.HasErrors() {
		return 1
	}
	return 0
}
