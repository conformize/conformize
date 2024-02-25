// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package execution

type BlueprintReferencesResolutionPhase struct {
	refs *map[string]string
}

func NewBlueprintReferencesResolutionPhase(refs *map[string]string) *BlueprintReferencesResolutionPhase {
	return &BlueprintReferencesResolutionPhase{
		refs: refs,
	}
}

func (phase *BlueprintReferencesResolutionPhase) Execute(blprntExecCtx *BlueprintExecutionContext) {
	refResolver := NewReferencesResolver(blprntExecCtx.valueReferencesStore)
	refResolver.Resolve(phase.refs, blprntExecCtx.diags)
}
