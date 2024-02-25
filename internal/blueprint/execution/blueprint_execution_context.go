// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package execution

import (
	"github.com/conformize/conformize/common/diagnostics"
	"github.com/conformize/conformize/common/ds"
	sdk "github.com/conformize/conformize/internal/providers/api"
	"github.com/conformize/conformize/internal/valuereferencesstore"
)

type BlueprintExecutionContext struct {
	diags                      *diagnostics.Diagnostics
	providersRegistry          sdk.ProvidersRegistrar
	providersDependenciesGraph *ds.DependencyGraph[string]
	valueReferencesStore       *valuereferencesstore.ValueReferencesStore
}
