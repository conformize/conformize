// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package format

import "github.com/conformize/conformize/internal/ui/options"

type detailCode string

const (
	info_code             detailCode = "ℹ️"
	info_ascii            detailCode = "[info]"
	ok_code               detailCode = "✅"
	ok_ascii              detailCode = "[ok]"
	error_code            detailCode = "❌"
	error_ascii           detailCode = "[error]"
	failure_warning_code  detailCode = "❗"
	failure_warning_ascii detailCode = "[failure]"
	warning_code          detailCode = "⚠️"
	warning_ascii         detailCode = "[warning]"
	tool_code             detailCode = "🔧"
	tool_ascii            detailCode = "==>"
	line_item_code        detailCode = "✓"
	line_item_ascii       detailCode = "->"
	bullet_code           detailCode = "•"
	bullet_ascii          detailCode = "-"
	box_code              detailCode = "📦"
	box_ascii             detailCode = "--"
	test_tube             detailCode = "🧪"
	failure                          = "✗"
)

type Detail string

const (
	Info           Detail = "info"
	Ok             Detail = "ok"
	Error          Detail = "error"
	FailureWarning Detail = "failure_warnign"
	Warning        Detail = "warning"
	Tool           Detail = "tool"
	Bullet         Detail = "bullet"
	Item           Detail = "item"
	Box            Detail = "box"
	TestTube       Detail = "pencil"
	Failure        Detail = "failure"
)

type detailPair struct {
	Code  detailCode
	Ascii detailCode
}

var detailsMappings = map[Detail]detailPair{
	Info:           {Code: info_code, Ascii: info_ascii},
	Ok:             {Code: ok_code, Ascii: ok_ascii},
	Error:          {Code: error_code, Ascii: error_ascii},
	FailureWarning: {Code: failure_warning_code, Ascii: failure_warning_ascii},
	Warning:        {Code: warning_code, Ascii: warning_ascii},
	Tool:           {Code: tool_code, Ascii: tool_ascii},
	Bullet:         {Code: bullet_code, Ascii: bullet_ascii},
	Item:           {Code: line_item_code, Ascii: line_item_code},
	Box:            {Code: box_code, Ascii: box_ascii},
	TestTube:       {Code: test_tube, Ascii: tool_ascii},
	Failure:        {Code: failure, Ascii: failure},
}

func (d Detail) String() string {
	return string(d)
}

func (d Detail) Code() string {
	if pair, ok := detailsMappings[d]; ok {
		if options.Options().Ui.Plain {
			return string(pair.Ascii)
		}
		return string(pair.Code)
	}
	return ""
}
