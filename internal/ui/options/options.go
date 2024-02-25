// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package options

import (
	"fmt"
	"strings"

	"github.com/conformize/conformize/common/typed"
	"github.com/conformize/conformize/common/util"
	"github.com/conformize/conformize/internal/providers/api/schema"
	"github.com/conformize/conformize/internal/providers/api/schema/attributes"
	"github.com/conformize/conformize/serialization/unmarshal/functions"
	"github.com/conformize/conformize/ui/commands/cmdutil"
)

type options struct {
	data *schema.Data
}

type optionSetter = func(data *schema.Data, optVal interface{}) error

type GlobalOptions struct {
	Ui          *UiOptions          `cnfrmz:"ui"`
	Diagnostics *DiagnosticsOptions `cnfrmz:"diagnostics"`
	WorkDir     string              `cnfrmz:"workdir"`
	Verbose     bool                `cnfrmz:"verbose"`
}

var optionsMap = map[string]optionSetter{
	"colorize": func(data *schema.Data, optVal interface{}) error {
		return data.SetAtPath("ui.colorize", optVal)
	},
	"include-diagnostics": func(data *schema.Data, optVal interface{}) error {
		incldueDiagsOpt, ok := optVal.(string)
		if !ok {
			return fmt.Errorf("invalid type for --include-diagnostics option: %T", optVal)
		}

		return data.SetAtPath("diagnostics.include", strings.Split(incldueDiagsOpt, ","))
	},
	"workdir": func(data *schema.Data, optVal interface{}) error {
		workDir, ok := optVal.(string)
		if !ok {
			return fmt.Errorf("invalid type for --workdir option: %T", optVal)
		}

		resolvedWorkDir, err := util.ResolveFilePath(workDir)
		if err != nil {
			return fmt.Errorf("couldn't resolve working directory: %v", err)
		}
		return data.SetAtPath("workdir", resolvedWorkDir)
	},
	"verbose": func(data *schema.Data, optVal interface{}) error {
		return data.SetAtPath("verbose", optVal)
	},
}

var optionsSchema = schema.Schema{
	Description: "Conformize Options schema",
	Version:     1,
	Attributes: map[string]schema.Attributeable{
		"ui": &attributes.ObjectAttribute{
			Description: "UI options",
			FieldsTypes: map[string]typed.Typeable{
				"colorize": &typed.BooleanTyped{},
			},
		},
		"diagnostics": &attributes.ObjectAttribute{
			Description: "Diagnostics options",
			FieldsTypes: map[string]typed.Typeable{
				"include": &typed.ListTyped{ElementsType: &typed.StringTyped{}},
			},
		},
		"workdir": &attributes.StringAttribute{},
		"verbose": &attributes.BooleanAttribute{},
	},
}

var globalOptions *GlobalOptions

func ParseOptions(args []string) ([]string, error) {
	if globalOptions != nil {
		return nil, fmt.Errorf("options already loaded")
	}

	var newArgs []string
	opts, newArgs, err := parseOptions(args)
	if err != nil {
		return nil, err
	}

	globalOptions = opts
	return newArgs, nil
}

func Options() *GlobalOptions {
	if globalOptions == nil {
		panic("options not loaded")
	}
	return globalOptions
}

func (opts *options) SetOption(optName string, optVal interface{}) error {
	optSetter, ok := optionsMap[optName]
	if !ok {
		return fmt.Errorf("unknown option: %s", optName)
	}

	return optSetter(opts.data, optVal)
}

func parseOptions(args []string) (*GlobalOptions, []string, error) {
	var err error

	newArgs := make([]string, 0)
	opts := &options{data: schema.NewData(&optionsSchema)}
	for _, arg := range args {
		if !strings.HasPrefix(arg, "--") || cmdutil.IsHelpCommand(arg) || cmdutil.IsVersionCommand(arg) {
			newArgs = append(newArgs, arg)
			continue
		}

		parts := strings.SplitN(arg[2:], "=", 2)
		optName := parts[0]

		optVal := ""
		if len(parts) > 1 {
			optVal = parts[1]
		}

		var val interface{}
		val, err = functions.DecodeStringValue(optVal)
		if err != nil {
			val = optVal
		}

		if err = opts.SetOption(optName, val); err != nil {
			return nil, nil, err
		}
	}

	var globalOpts GlobalOptions
	if err := opts.data.Get(&globalOpts); err != nil {
		return nil, nil, err
	}

	return &globalOpts, newArgs, nil
}

func Help() string {
	helpBldr := strings.Builder{}
	helpBldr.WriteString(fmt.Sprintf("%-45s\t%-s\n", "--colorize=< true|false >", "enable colorized output"))
	helpBldr.WriteString(fmt.Sprintf("%-45s\t%-s\n", "--include-diagnostics=< info,warn,error >", "include only specified diagnostics in output"))
	helpBldr.WriteString(fmt.Sprintf("%-45s\t%-s\n", "--workdir=< DIR >", "set the working directory"))
	helpBldr.WriteString(fmt.Sprintf("%-45s\t%-s\n", "--verbose=< true|false >", "enable verbose output"))

	return helpBldr.String()
}
