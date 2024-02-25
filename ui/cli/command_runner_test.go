// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package cli

import (
	"reflect"
	"testing"

	"github.com/conformize/conformize/common/diagnostics"
)

func TestNewCommandRunner(t *testing.T) {
	tests := []struct {
		name string
		want *commandRunner
	}{
		{
			name: "Successfully initializes Command Runner",
			want: &commandRunner{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewCommandRunner(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewCommandRunner() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_commandRunner_Run(t *testing.T) {
	type args struct {
		args  []string
		diags *diagnostics.Diagnostics
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Test version command execution is successful and displays  version info",
			args: args{
				args:  []string{"version"},
				diags: diagnostics.NewDiagnostics(),
			},
		},
		{
			name: "Test blueprint command execution is successful and displays usage instructions",
			args: args{
				args:  []string{"blueprint"},
				diags: diagnostics.NewDiagnostics(),
			},
		},
		{
			name: "Test blueprint help command execution is successful and displays usage instructions",
			args: args{
				args:  []string{"blueprint", "help"},
				diags: diagnostics.NewDiagnostics(),
			},
		},
		{
			name: "Test blueprint validate command execution fails without arguments",
			args: args{
				args:  []string{"blueprint", "validate"},
				diags: diagnostics.NewDiagnostics(),
			},
			wantErr: true,
		},
		{
			name: "Test blueprint validate command execution is successful with correct arguments",
			args: args{
				args:  []string{"blueprint", "validate", "-f", "../../internal/blueprint/mocks/blueprint.cnfrm.yaml"},
				diags: diagnostics.NewDiagnostics(),
			},
		},
		{
			name: "Test blueprint validate help command execution is successful and displays usage instruction",
			args: args{
				args:  []string{"blueprint", "validate", "help"},
				diags: diagnostics.NewDiagnostics(),
			},
		},
		{
			name: "Test blueprint apply command execution fails without arguments",
			args: args{
				args:  []string{"blueprint", "apply"},
				diags: diagnostics.NewDiagnostics(),
			},
			wantErr: true,
		},
		{
			name: "Test blueprint apply command execution is successful with correct arguments",
			args: args{
				args:  []string{"blueprint", "apply", "-f", "../../internal/blueprint/mocks/blueprint.cnfrm.yaml"},
				diags: diagnostics.NewDiagnostics(),
			},
		},
		{
			name: "Test blueprint apply help command execution is successful and displays usage instruction",
			args: args{
				args:  []string{"blueprint", "apply", "help"},
				diags: diagnostics.NewDiagnostics(),
			},
		},
		{
			name: "Test command execution fails for unrecognized command",
			args: args{
				args:  []string{"test"},
				diags: diagnostics.NewDiagnostics(),
			},
			wantErr: true,
		},
		{
			name: "Test command execution fails for command with unrecognized subcommand",
			args: args{
				args:  []string{"blueprint", "blueprint", "test"},
				diags: diagnostics.NewDiagnostics(),
			},
			wantErr: true,
		},
	}
	cmdRun := NewCommandRunner()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			res := cmdRun.Run(tt.args.args, tt.args.diags)
			if (tt.wantErr && !tt.args.diags.HasErrors() && res != 1) ||
				(!tt.wantErr && tt.args.diags.HasErrors() && res != 0) {
				t.Errorf("CommandRunner.Run() error = %v, wantErr %v", tt.args.diags.Errors(), tt.wantErr)
				return
			}
		})
	}
}
