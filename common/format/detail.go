// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package format

type detailCode string

type Detail int

const (
	Info Detail = iota
	Ok
	Error
	FailureWarning
	Warning
	Tool
	Bullet
	Item
	Box
	Pencil
	Failure
)

type detail struct {
	name  string
	code  detailCode
	plain detailCode
}

var details = []detail{
	{"info", "ℹ️", "[info]"},
	{"ok", "✅", "[ok]"},
	{"error", "❌", "[error]"},
	{"failure_warning", "❗", "[failure]"},
	{"warning", "⚠️", "[warning]"},
	{"tool", "🔧", "==>"},
	{"bullet", "•", "-"},
	{"item", "✓", "✓"},
	{"box", "📦", "--"},
	{"pencil", "✏️", "==>"},
	{"failure", "✗", "✗"},
}

var fallbacks = map[Detail]detailCode{
	Item:    "v",
	Failure: "x",
}

func (d Detail) String() string {
	if int(d) < 0 || int(d) >= len(details) {
		return "unknown"
	}
	entry := details[d]
	return entry.name
}

func (d Detail) Code() string {
	if int(d) < 0 || int(d) >= len(details) {
		return ""
	}
	entry := details[d]
	return string(entry.code)
}

func (d Detail) Plain() string {
	if int(d) < 0 || int(d) >= len(details) {
		return ""
	}
	entry := details[d]
	return string(entry.plain)
}

func (d Detail) Fallback() string {
	if fallback, ok := fallbacks[d]; ok {
		return string(fallback)
	}
	return ""
}
