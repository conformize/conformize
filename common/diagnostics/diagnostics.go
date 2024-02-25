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

	"github.com/conformize/conformize/common/streams"
)

type Diagnostics struct {
	entries      Diags
	hasEntries   byte
	rwLock       *sync.RWMutex
	streamSignal chan struct{}
	doStream     bool
}

type Diags []Diagnostic

func (d Diags) String() string {
	strBldr := strings.Builder{}
	for _, diag := range d {
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

func NewDiagnostics() *Diagnostics {
	d := &Diagnostics{
		streamSignal: make(chan struct{}, 10),
		entries:      make(Diags, 0, 100),
		hasEntries:   0,
		rwLock:       &sync.RWMutex{},
	}
	return d
}

func (d *Diagnostics) Append(diags ...Diagnostic) {
	d.rwLock.Lock()
	for _, diag := range diags {
		if diag != nil {
			d.entries = append(d.entries, diag)
			d.hasEntries |= 1 << diag.GetType()
		}
	}
	d.rwLock.Unlock()

	if d.doStream {
		d.streamSignal <- struct{}{}
	}
}

func (d *Diagnostics) Entries() Diags {
	d.rwLock.RLock()
	defer d.rwLock.RUnlock()
	return d.entries
}

func (d *Diagnostics) Errors() Diags {
	d.rwLock.RLock()
	defer d.rwLock.RUnlock()
	return d.entriesOfType(Error)
}

func (d *Diagnostics) Warnings() Diags {
	d.rwLock.RLock()
	defer d.rwLock.RUnlock()
	return d.entriesOfType(Warning)
}

func (d *Diagnostics) Infos() Diags {
	d.rwLock.RLock()
	defer d.rwLock.RUnlock()
	return d.entriesOfType(Info)
}

func (d *Diagnostics) entriesOfType(diagnosticType DiagnosticType) Diags {
	var diagnostics Diags
	for _, diag := range d.entries {
		if diag.GetType() == diagnosticType {
			diagnostics = append(diagnostics, diag)
		}
	}
	return diagnostics
}

func (d *Diagnostics) HasErrors() bool {
	d.rwLock.RLock()
	defer d.rwLock.RUnlock()
	return d.hasEntries&0b10000000 != 0
}

func (d *Diagnostics) HasWarnings() bool {
	d.rwLock.RLock()
	defer d.rwLock.RUnlock()
	return d.hasEntries&0b01000000 != 0
}

func (d *Diagnostics) Stream(diagsStream *streams.Stream[Diagnostic]) {
	d.doStream = true
	go func() {
		defer close(d.streamSignal)

		streamPos := 0
		for !diagsStream.Closed() {
			<-d.streamSignal

			d.rwLock.RLock()
			streamLen := len(d.entries) - streamPos
			for ; streamLen > 0; streamLen-- {
				err := diagsStream.Write(d.entries[streamPos])
				if err != nil {
					break
				}
				streamPos++
			}
			d.rwLock.RUnlock()
		}
	}()
}
