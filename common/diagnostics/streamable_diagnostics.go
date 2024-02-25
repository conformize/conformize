// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package diagnostics

import (
	"strings"
	"sync"

	"github.com/conformize/conformize/common/ds"
	"github.com/conformize/conformize/common/streams"
)

type StreamableDiagnostics struct {
	diags      Diagnosable
	wg         *sync.WaitGroup
	stream     *ds.CircularBuffer[Diagnostic]
	signalChan chan struct{}
}

func (d StreamableDiagnostics) Print() string {
	strBldr := strings.Builder{}
	for _, diag := range d.Entries() {
		if diag != nil {
			if diag.GetSummary() != "" {
				strBldr.WriteString(diag.GetSummary())
				strBldr.WriteString("\n")
			}
			if diag.GetDetails() != "" {
				strBldr.WriteString(diag.GetDetails())
				strBldr.WriteString("\n\n")
			}
		}
	}
	return strBldr.String()
}

func NewStreamableDiagnostics(wg *sync.WaitGroup) *StreamableDiagnostics {
	d := &StreamableDiagnostics{
		diags:      NewDiagnostics(),
		stream:     ds.NewCircularBuffer[Diagnostic](100),
		signalChan: make(chan struct{}, 10),
		wg:         wg,
	}
	return d
}

func (d *StreamableDiagnostics) Append(diags ...Diagnostic) {
	for _, diag := range diags {
		if diag != nil {
			d.wg.Add(1)
			d.diags.Append(diag)
			d.stream.Write(diag)
			d.signalChan <- struct{}{}
		}
	}
}

func (d *StreamableDiagnostics) Entries() Diags {
	return d.diags.Entries()
}

func (d *StreamableDiagnostics) Errors() Diags {
	return d.diags.Errors()
}

func (d *StreamableDiagnostics) Warnings() Diags {
	return d.diags.Warnings()
}

func (d *StreamableDiagnostics) Infos() Diags {
	return d.diags.Infos()
}

func (d *StreamableDiagnostics) HasErrors() bool {
	return d.diags.HasErrors()
}

func (d *StreamableDiagnostics) HasWarnings() bool {
	return d.diags.HasWarnings()

}

func (d *StreamableDiagnostics) Stream(diagsStream *streams.Stream[Diagnostic]) {
	go func() {
		defer close(d.signalChan)
		for !diagsStream.Closed() {
			<-d.signalChan
			for isDrained := d.stream.IsEmpty(); !isDrained; isDrained = d.stream.IsEmpty() {
				diag := d.stream.Read()
				if err := diagsStream.Write(diag); err != nil {
					return
				}
			}
		}
	}()
}
