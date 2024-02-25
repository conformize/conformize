// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package resources

import _ "embed"

//go:embed ASCII_LOGO
var ascii_logo string

//go:embed VERSION
var version_string string

//go:embed GENERAL_HELP_TMPL
var general_help_tmpl string

func ASCII_LOGO() string {
	return ascii_logo
}

func VersionString() string {
	return version_string
}

func HelpText() string {
	return general_help_tmpl
}
