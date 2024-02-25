// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package diagnostics

import "testing"

func TestDiagnosticsWithErrorsOnly(t *testing.T) {
	diags := NewDiagnostics()
	diags.Append(&diagnostic{
		Type:    Error,
		Summary: "Test error",
		Details: "This is a test error",
	})
	diags.Append(&diagnostic{
		Type:    Info,
		Summary: "Test info",
		Details: "This is an info message",
	})
	if diags.HasWarnings() {
		t.Errorf("Expected HasErrors to be true and hasWarnings to be false, but got HasErrors: %t and HasWarnings: %t", diags.HasErrors(), diags.HasWarnings())
	}
}

func TestDiagnosticsWithWarningsOnly(t *testing.T) {
	diags := NewDiagnostics()
	diags.Append(&diagnostic{
		Type:    Info,
		Summary: "Test info",
		Details: "This is an info message",
	})
	diags.Append(&diagnostic{
		Type:    Warning,
		Summary: "Test warning",
		Details: "This is a test warning",
	})
	if diags.HasErrors() {
		t.Errorf("Expected HasErrors to be false and hasWarnings to be true, but got HasErrors: %t and HasWarnings: %t", diags.HasErrors(), diags.HasWarnings())
	}
}

func TestDiagnosticsWithBothErrorsAndWarnings(t *testing.T) {
	diags := NewDiagnostics()
	diags.Append(&diagnostic{
		Type:    Info,
		Summary: "Test info",
		Details: "This is an info message",
	})
	diags.Append(&diagnostic{
		Type:    Error,
		Summary: "Test error",
		Details: "This is a test error",
	})
	diags.Append(&diagnostic{
		Type:    Warning,
		Summary: "Test warning",
		Details: "This is a test warning",
	})
	diags.Append(&diagnostic{
		Type:    Error,
		Summary: "Another test error",
		Details: "This is another test error",
	})
	if !diags.HasErrors() && !diags.HasWarnings() {
		t.Errorf("Expected HasErrors to be true and hasWarnings to be true, but got HasErrors: %t and HasWarnings: %t", diags.HasErrors(), diags.HasWarnings())
	}
}
