// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package format

type styleCode string

const (
	boldCode       styleCode = "1"
	underlinedCode styleCode = "4"
	dimmedCode     styleCode = "2"
)

type Style int

const (
	Bold Style = iota
	Underlined
	Dimmed
)

var styleNames = []string{
	"bold",
	"underlined",
	"dimmed",
}

var styleCodes = []styleCode{
	boldCode,
	underlinedCode,
	dimmedCode,
}

func (s Style) String() string {
	if int(s) < 0 || int(s) >= len(styleNames) {
		return "unknown"
	}
	return styleNames[s]
}

func (s Style) Code() string {
	if int(s) < 0 || int(s) >= len(styleCodes) {
		return ""
	}
	return string(styleCodes[s])
}
