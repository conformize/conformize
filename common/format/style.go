// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package format

type styleCode string

const (
	bold       styleCode = "1"
	underlined styleCode = "4"
	dimmed     styleCode = "2"
)

type Style string

var styleMappings = map[Style]styleCode{
	Bold:       bold,
	Underlined: underlined,
	Dimmed:     dimmed,
}

const (
	Bold       = Style("bold")
	Underlined = Style("underlined")
	Dimmed     = Style("dimmed")
)

func (s Style) String() string {
	return string(s)
}

func (s Style) Code() string {
	return string(styleMappings[s])
}
