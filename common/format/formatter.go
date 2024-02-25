// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package format

import (
	"strings"

	"github.com/conformize/conformize/common/format/colors"
)

type styleElement struct {
	precedence int
	style      StyleAttribute
	next       *styleElement
}

type formatter struct {
	styles *styleElement
	reset  bool
}

func Formatter() *formatter {
	return &formatter{
		styles: nil,
		reset:  false,
	}
}

func (f *formatter) dequeueStyle() StyleAttribute {
	if f.styles == nil {
		return nil
	}
	var style = f.styles
	f.styles = f.styles.next
	return style.style
}

func (f *formatter) enqueueStyle(precedence int, style StyleAttribute) {
	var newStyle = &styleElement{
		precedence: precedence,
		style:      style,
		next:       nil,
	}

	if f.styles == nil || f.styles.precedence < precedence {
		newStyle.next = f.styles
		f.styles = newStyle
		return
	}

	var head = f.styles
	for head.next != nil && head.next.precedence >= precedence {
		head = head.next
	}
	newStyle.next = head.next
	head.next = newStyle
}

func (f *formatter) Bold() *formatter {
	f.enqueueStyle(1, Bold)
	return f
}

func (f *formatter) Underlined() *formatter {
	f.enqueueStyle(4, Underlined)
	return f
}

func (f *formatter) Dimmed() *formatter {
	f.enqueueStyle(2, Dimmed)
	return f
}

func (f *formatter) Color(color colors.Color) *formatter {
	f.enqueueStyle(0, color)
	return f
}

func (f *formatter) Format(in string) string {
	var formatBldr strings.Builder

	formatBldr.WriteString("\033[")

	for style, prependSeparator := f.dequeueStyle(), false; style != nil; style = f.dequeueStyle() {
		if prependSeparator {
			formatBldr.WriteString(";")
		}
		formatBldr.WriteString(style.Code())
		prependSeparator = true
	}

	formatBldr.WriteString("m")
	formatBldr.WriteString(in)
	formatBldr.WriteString("\033[0m")
	return formatBldr.String()
}
