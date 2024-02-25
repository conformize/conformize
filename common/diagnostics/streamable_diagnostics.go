// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package diagnostics

import (
	"sync"

	"github.com/conformize/conformize/common/ds"
	"github.com/conformize/conformize/common/streams"
)

type StreamableDiagnostics struct {
	diags     Diagnosable
	wg        *sync.WaitGroup
	streamBuf *ds.CircularBuffer[Diagnostic]
	signal    chan struct{}
}

func NewStreamableDiagnostics(wg *sync.WaitGroup) *StreamableDiagnostics {
	d := &StreamableDiagnostics{
		diags:     NewDiagnostics(),
		streamBuf: ds.NewCircularBuffer[Diagnostic](10),
		signal:    make(chan struct{}),
		wg:        wg,
	}
	return d
}

func (streamDiags *StreamableDiagnostics) Append(diags ...Diagnostic) {
	for _, diag := range diags {
		if diag != nil {
			streamDiags.diags.Append(diag)
			streamDiags.streamBuf.Write(diag)
			streamDiags.wg.Add(1)
			streamDiags.signal <- struct{}{}
		}
	}
}

func (streamDiags *StreamableDiagnostics) Entries() Diags {
	return streamDiags.diags.Entries()
}

func (streamDiags *StreamableDiagnostics) Errors() Diags {
	return streamDiags.diags.Errors()
}

func (streamDiags *StreamableDiagnostics) Warnings() Diags {
	return streamDiags.diags.Warnings()
}

func (streamDiags *StreamableDiagnostics) Infos() Diags {
	return streamDiags.diags.Infos()
}

func (streamDiags *StreamableDiagnostics) HasErrors() bool {
	return streamDiags.diags.HasErrors()
}

func (streamDiags *StreamableDiagnostics) HasWarnings() bool {
	return streamDiags.diags.HasWarnings()

}

func (streamDiags *StreamableDiagnostics) Stream(diagsStream *streams.Stream[Diagnostic]) {
	go func() {
		defer close(streamDiags.signal)
		for !diagsStream.Closed() {
			<-streamDiags.signal
			for isDrained := streamDiags.streamBuf.IsEmpty(); !isDrained; isDrained = streamDiags.streamBuf.IsEmpty() {
				diag := streamDiags.streamBuf.Read()
				err := diagsStream.Write(diag)
				if err != nil {
					return
				}
			}
		}
	}()
}
