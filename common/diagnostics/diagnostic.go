// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package diagnostics

import "time"

type DiagnosticType int

const (
	Info    DiagnosticType = iota
	Warning                = 6
	Error                  = 7
)

type Diagnostic interface {
	GetTimestamp() time.Time
	GetType() DiagnosticType
	GetSummary() string
	GetDetails() string
	String() string
}
