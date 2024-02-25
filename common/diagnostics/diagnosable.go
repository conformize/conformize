// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package diagnostics

type Diagnosable interface {
	Append(diags ...Diagnostic)
	Entries() Diags
	Errors() Diags
	Warnings() Diags
	Infos() Diags
	HasErrors() bool
	HasWarnings() bool
}
