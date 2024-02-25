// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package format

import (
	"testing"

	"github.com/conformize/conformize/common/format/colors"
)

func TestBold(t *testing.T) {
	builder := Formatter().Bold()
	expected := "\033[1mThis is bold text\033[0m"
	result := builder.Format("This is bold text")

	if result != expected {
		t.Errorf("Expected %q but got %q", expected, result)
	}
}

func TestUnderlined(t *testing.T) {
	builder := Formatter().Underlined()
	expected := "\033[4mThis is underlined text\033[0m"
	result := builder.Format("This is underlined text")

	if result != expected {
		t.Errorf("Expected %q but got %q", expected, result)
	}
}

func TestBoldAndUnderlined(t *testing.T) {
	builder := Formatter().Bold().Underlined()
	expected := "\033[1;4mThis is bold and underlined text\033[0m"
	result := builder.Format("This is bold and underlined text")

	if result != expected {
		t.Errorf("Expected %q but got %q", expected, result)
	}
}

func TestColoredText(t *testing.T) {
	builder := Formatter().Color(colors.Red)
	expected := "\033[38;5;9mThis is red text\033[0m"
	result := builder.Format("This is red text")

	if result != expected {
		t.Errorf("Expected %q but got %q", expected, result)
	}
}

func TestBoldAndColored(t *testing.T) {
	builder := Formatter().Bold().Color(colors.Red)
	expected := "\x1b[38;5;9;1mThis is bold and red text\x1b[0m"
	result := builder.Format("This is bold and red text")

	if result != expected {
		t.Errorf("Expected %q but got %q", expected, result)
	}
}

func TestMultipleStyles(t *testing.T) {
	builder := Formatter().Bold().Underlined().Color(colors.Yellow)
	expected := "\x1b[38;5;190;1;4mThis is bold, underlined, and yellow text\x1b[0m"
	result := builder.Format("This is bold, underlined, and yellow text")

	if result != expected {
		t.Errorf("Expected %q but got %q", expected, result)
	}
}
