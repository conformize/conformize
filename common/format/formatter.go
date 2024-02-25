// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package format

import (
	"strings"

	"github.com/conformize/conformize/common/ds"
	"github.com/conformize/conformize/common/format/colors"
	"github.com/conformize/conformize/internal/ui/options"
)

var detailSpacing = "  "

type styledEntry struct {
	precedence int
	style      StyleAttribute
}

type formatter struct {
	styles *ds.PriorityQueue[StyleAttribute]
	reset  bool
	plain  bool
	bldr   strings.Builder
}

func Formatter() *formatter {
	return &formatter{
		reset:  true,
		plain:  options.Options().Ui.Plain,
		styles: ds.NewPriorityQueue(ds.MaxPriorityComparator[StyleAttribute]),
	}
}

func (f *formatter) addStyle(precedence int, style StyleAttribute) {
	f.styles.Push(precedence, style)
}

func (f *formatter) Bold() *formatter {
	if !f.plain {
		f.addStyle(3, Bold)
	}
	return f
}

func (f *formatter) Underlined() *formatter {
	if !f.plain {
		f.addStyle(3, Underlined)
	}
	return f
}

func (f *formatter) Dimmed() *formatter {
	if !f.plain {
		f.addStyle(3, Dimmed)
	}
	return f
}

func (f *formatter) Color(color colors.Color) *formatter {
	if !f.plain {
		f.addStyle(2, color)
	}
	return f
}

func (f *formatter) Detail(detail Detail) *formatter {
	f.addStyle(1, detail)
	return f
}

func (f *formatter) Format(in string) string {
	var sepNeeded bool
	var detailPrefix string

	for !f.styles.IsEmpty() {
		entry, ok := f.styles.Pop()
		if !ok {
			continue
		}
		switch s := entry.(type) {
		case Detail:
			if !f.plain {
				detailPrefix = s.Code() + detailSpacing
				continue
			}

			detailPrefix = s.Plain() + detailSpacing
		case colors.Color, Style:
			if !sepNeeded {
				f.bldr.WriteString("\033[")
			} else {
				f.bldr.WriteByte(';')
			}
			f.bldr.WriteString(s.Code())
			sepNeeded = true
		}
	}

	if sepNeeded {
		f.bldr.WriteByte('m')
	}

	f.bldr.WriteString(detailPrefix)
	f.bldr.WriteString(in)

	if sepNeeded && f.reset {
		f.bldr.WriteString("\033[0m")
	}

	formatted := f.bldr.String()
	f.bldr.Reset()
	f.bldr.Grow(256)
	return formatted
}
