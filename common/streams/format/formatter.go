// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package format

import (
	"strings"

	"github.com/conformize/conformize/common/streams/colors"
)

type styleEntry struct {
	weight int
	style  StyleAttribute
	next   *styleEntry
}

type formatter struct {
	styles *styleEntry
	reset  bool
}

func Formatter() *formatter {
	return &formatter{
		styles: nil,
		reset:  false,
	}
}

func (f *formatter) popStyle() StyleAttribute {
	if f.styles == nil {
		return nil
	}
	var style = f.styles
	f.styles = f.styles.next
	return style.style
}

func (f *formatter) pushStyle(weight int, style StyleAttribute) {
	var newStyle = &styleEntry{
		weight: weight,
		style:  style,
		next:   nil,
	}

	if f.styles == nil || f.styles.weight < weight {
		newStyle.next = f.styles
		f.styles = newStyle
		return
	}

	var head = f.styles
	for head.next != nil && head.next.weight >= weight {
		head = head.next
	}
	newStyle.next = head.next
	head.next = newStyle
}

func (f *formatter) Bold() *formatter {
	f.pushStyle(1, Bold)
	return f
}

func (f *formatter) Underlined() *formatter {
	f.pushStyle(1, Underlined)
	return f
}

func (f *formatter) Dimmed() *formatter {
	f.pushStyle(1, Dimmed)
	return f
}

func (f *formatter) Color(color colors.Color) *formatter {
	f.pushStyle(0, color)
	return f
}

func (f *formatter) Format(in string) string {
	var formatBldr strings.Builder

	formatBldr.WriteString("\033[")

	prependSeparator := false
	for style := f.popStyle(); style != nil; style = f.popStyle() {
		if prependSeparator {
			formatBldr.WriteString(";")
		} else {
			prependSeparator = true
		}
		formatBldr.WriteString(style.Code())
	}

	formatBldr.WriteString("m")
	formatBldr.WriteString(in)
	formatBldr.WriteString("\033[0m")
	return formatBldr.String()
}
