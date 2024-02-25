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
	entries     Diags
	entriesMask byte
	rwLock      *sync.RWMutex
	stream      *streams.Stream[Diagnostic]
	doStream    bool
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
			}
		}
	}
	return strBldr.String()
}

func NewDiagnostics() *Diagnostics {
	d := &Diagnostics{
		entries:     make(Diags, 0, 100),
		entriesMask: 0,
		rwLock:      &sync.RWMutex{},
	}
	return d
}

func (d *Diagnostics) Append(diags ...Diagnostic) {
	for _, diag := range diags {
		if diag != nil {
			d.rwLock.Lock()
			d.entries = append(d.entries, diag)
			d.entriesMask |= 1 << diag.GetType()
			d.rwLock.Unlock()

			if d.doStream {
				if err := d.stream.Write(diag); err != nil {
					d.doStream = false
				}
			}

		}
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
	res := d.entriesMask&0b10000000 != 0
	d.rwLock.RUnlock()
	return res
}

func (d *Diagnostics) HasWarnings() bool {
	d.rwLock.RLock()
	res := d.entriesMask&0b01000000 != 0
	d.rwLock.RUnlock()
	return res
}

func (d *Diagnostics) Stream(diagsStream *streams.Stream[Diagnostic]) {
	if diagsStream != nil {
		d.rwLock.Lock()
		d.doStream = true
		d.stream = diagsStream
		d.rwLock.Unlock()
	}
}
