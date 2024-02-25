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
	"github.com/conformize/conformize/internal/ui/options"
)

var detailSpacing = " "

type styleElement struct {
	precedence int
	style      StyleAttribute
	next       *styleElement
}

type formatter struct {
	styles *styleElement
	reset  bool
	plain  bool
}

func Formatter() *formatter {
	plain := options.Options().Ui.Plain

	return &formatter{
		styles: nil,
		reset:  true,
		plain:  plain,
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
	newStyle := &styleElement{
		precedence: precedence,
		style:      style,
		next:       nil,
	}

	if f.styles == nil || f.styles.precedence > precedence {
		newStyle.next = f.styles
		f.styles = newStyle
		return
	}

	head := f.styles
	for head.next != nil && head.next.precedence <= precedence {
		head = head.next
	}
	newStyle.next = head.next
	head.next = newStyle
}

func (f *formatter) Bold() *formatter {
	if f.plain {
		return f
	}
	f.enqueueStyle(3, Bold)
	return f
}

func (f *formatter) Underlined() *formatter {
	if f.plain {
		return f
	}
	f.enqueueStyle(3, Underlined)
	return f
}

func (f *formatter) Dimmed() *formatter {
	if f.plain {
		return f
	}
	f.enqueueStyle(3, Dimmed)
	return f
}

func (f *formatter) Color(color colors.Color) *formatter {
	if f.plain {
		return f
	}
	f.enqueueStyle(2, color)
	return f
}

func (f *formatter) Detail(detail Detail) *formatter {
	f.enqueueStyle(1, detail)
	return f
}

func (f *formatter) Format(in string) string {
	var bldr strings.Builder
	var styleCodes strings.Builder
	var detailPrefix string
	var sepNeeded bool

	for style := f.dequeueStyle(); style != nil; style = f.dequeueStyle() {
		switch s := style.(type) {
		case Detail:
			detailPrefix = s.Code() + detailSpacing
		case colors.Color:
			if sepNeeded {
				styleCodes.WriteByte(';')
			}
			styleCodes.WriteString(s.Code())
			sepNeeded = true
		case Style:
			if sepNeeded {
				styleCodes.WriteByte(';')
			}
			styleCodes.WriteString(s.Code())
			sepNeeded = true
		}
	}

	if styleCodes.Len() > 0 {
		bldr.WriteString("\033[")
		bldr.WriteString(styleCodes.String())
		bldr.WriteString("m")
	}

	bldr.WriteString(detailPrefix)
	bldr.WriteString(in)

	if styleCodes.Len() > 0 && f.reset {
		bldr.WriteString("\033[0m")
	}

	return bldr.String()
}
