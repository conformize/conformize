// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package diagnostics

import (
	"strings"
	"time"
)

type diagnostic struct {
	Type      DiagnosticType
	Summary   string
	Details   string
	Timestamp time.Time
}

func (d *diagnostic) GetType() DiagnosticType {
	return d.Type
}

func (d *diagnostic) GetSummary() string {
	return d.Summary
}

func (d *diagnostic) GetDetails() string {
	return d.Details
}

func (d *diagnostic) GetTimestamp() time.Time {
	return d.Timestamp
}

func (d *diagnostic) String() string {
	var sb strings.Builder
	if len(d.Summary) > 0 {
		sb.WriteString(d.Summary)
	}

	if len(d.Details) > 0 {
		if sb.Len() > 0 {
			sb.WriteString("\n")
		}
		sb.WriteString(d.Details)
	}
	return sb.String()
}

type builder struct {
	diagnostic
}

func (d *builder) Type(diagnosticType DiagnosticType) *builder {
	d.diagnostic.Type = diagnosticType
	return d
}

func (d *builder) Error() *builder {
	d.diagnostic.Type = Error
	return d
}

func (d *builder) Warning() *builder {
	d.diagnostic.Type = Warning
	return d
}

func (d *builder) Info() *builder {
	d.diagnostic.Type = Info
	return d
}

func (d *builder) Summary(summary string) *builder {
	d.diagnostic.Summary = summary
	return d
}

func (d *builder) Details(details string) *builder {
	d.diagnostic.Details = details
	return d
}

func (d *builder) Build() Diagnostic {
	d.diagnostic.Timestamp = time.Now()
	return &d.diagnostic
}

func Builder() *builder {
	return &builder{
		diagnostic: diagnostic{},
	}
}
