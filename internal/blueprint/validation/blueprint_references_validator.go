// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package validation

import (
	"fmt"

	"github.com/conformize/conformize/common/diagnostics"
	"github.com/conformize/conformize/common/pathparser"
	"github.com/conformize/conformize/internal/blueprint"
)

type BlueprintReferencesValidator struct{}

func (blprntRefValidator *BlueprintReferencesValidator) Validate(blueprint *blueprint.Blueprint) *diagnostics.Diagnostics {
	diags := diagnostics.NewDiagnostics()
	pathParser := pathparser.NewPathParser()

	for refAlias, refPath := range blueprint.References {
		if _, sourceRefOk := blueprint.Sources[refAlias]; sourceRefOk {
			if _, aliasRefOk := blueprint.References[refAlias]; aliasRefOk {
				diags.Append(diagnostics.Builder().
					Error().
					Summary(fmt.Sprintf("\nReference '%s' is already defined", refAlias)).
					Build(),
				)
				continue
			}
		}

		if len(refPath) == 0 {
			diags.Append(diagnostics.Builder().
				Error().
				Summary(fmt.Sprintf("\nPath not specified for reference '%s'", refAlias)).
				Build(),
			)
			continue
		}

		path, refPathErr := pathParser.Parse(refPath)
		if refPathErr != nil {
			diags.Append(diagnostics.Builder().
				Error().
				Summary(fmt.Sprintf("\nPath %s for reference '%s' is not valid, reason:\n%s", refPath, refAlias, refPathErr.Error())).
				Build(),
			)
			continue
		}

		root := path[0].String()
		if _, sourceRefOk := blueprint.Sources[root]; !sourceRefOk {
			if _, aliasRefOk := blueprint.References[root]; !aliasRefOk {
				diags.Append(diagnostics.Builder().
					Error().
					Summary(fmt.Sprintf("Couldn't resolve root '%s' in path %s for reference '%s'", root, refAlias, refPath)).
					Build(),
				)
			}
		}
	}
	return diags
}
