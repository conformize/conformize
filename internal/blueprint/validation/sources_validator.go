// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package validation

import (
	"fmt"
	"strings"

	"github.com/conformize/conformize/common/diagnostics"
	"github.com/conformize/conformize/internal/blueprint"
)

type BlueprintSourcesValidator struct{}

func (blprntSrcsVld *BlueprintSourcesValidator) Validate(blueprint *blueprint.Blueprint) *diagnostics.Diagnostics {
	diags := diagnostics.NewDiagnostics()
	if len(blueprint.Sources) == 0 {
		diags.Append(diagnostics.Builder().
			Error().
			Summary("No sources specified in blueprint.\n").
			Build(),
		)
	} else {
		for alias, src := range blueprint.Sources {
			if src.Config != nil && src.ConfigFile != nil {
				diags.Append(diagnostics.Builder().
					Error().
					Summary(fmt.Sprintf("\nConfiguration for source '%s', provider '%s' is not valid, reason:\n", alias, src.Provider)).
					Details("Both inlined and externalized provider configurations found - please choose only one.").
					Build(),
				)
				return diags
			}

			if src.ConfigFile != nil && !strings.HasSuffix(*src.ConfigFile, ".yaml") && !strings.HasSuffix(*src.ConfigFile, ".json") {
				diags.Append(diagnostics.Builder().
					Error().
					Summary(fmt.Sprintf("\nConfiguration for source '%s', provider '%s' is not valid, reason:\n", alias, src.Provider)).
					Details("Provider configuration file must have .json or .yaml extension.\n").
					Build(),
				)
			}
		}
	}
	return diags
}
