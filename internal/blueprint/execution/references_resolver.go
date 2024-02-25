// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package execution

import (
	"fmt"

	"github.com/conformize/conformize/common/diagnostics"
	"github.com/conformize/conformize/common/ds"
	"github.com/conformize/conformize/common/path"
	"github.com/conformize/conformize/common/pathparser"
)

type ReferencesResolver struct {
	dependecyGraph  *ds.DependencyGraph[string]
	referencesStore *ValueReferencesStore
}

func (refResolver *ReferencesResolver) Resolve(refs map[string]string, diags diagnostics.Diagnosable) {

	refPaths := make(map[string]path.Steps)
	pathParser := pathparser.NewPathParser()
	for refAlias, refPath := range refs {
		pathSteps, _ := pathParser.Parse(refPath)
		root := pathSteps[0]
		refResolver.dependecyGraph.AddEdge(refAlias, root.String())
		refPaths[refAlias] = pathSteps
	}

	refResolver.detectCycles(diags)
	if !diags.HasErrors() {
		refsResolveDepOrder := refResolver.dependecyGraph.GetOrder()

		for _, refAlias := range refsResolveDepOrder {
			refPathSteps := refPaths[refAlias]
			_, resolved := refResolver.referencesStore.GetReference(refAlias)
			if resolved {
				continue
			}

			ref, err := refResolver.referencesStore.GetAtPath(path.NewPath(refPathSteps))
			if err == nil {
				refResolver.referencesStore.AddReference(refAlias, ref)
				continue
			}

			refRawPath := refs[refAlias]
			diags.Append(diagnostics.Builder().
				Error().
				Details(
					fmt.Sprintf("\nCouldn't resolve reference '%s' in path %s, reason:\n%s",
						refAlias, refRawPath, err.Error()),
				).
				Build(),
			)
		}
	}
}

func (refResolver *ReferencesResolver) detectCycles(diags diagnostics.Diagnosable) {
	refResolver.dependecyGraph.Run()
	if refResolver.dependecyGraph.HasCycles() {
		refCycles := refResolver.dependecyGraph.GetCycles()
		for _, cycle := range refCycles {
			ref := cycle[0]
			otherRef := cycle[1]
			diags.Append(diagnostics.Builder().
				Error().
				Summary(
					fmt.Sprintf("\ncyclic dependency detected between references '%s' and '%s'", ref, otherRef),
				).
				Build(),
			)
		}
	}
}

func NewReferencesResolver(valRefStore *ValueReferencesStore) *ReferencesResolver {
	return &ReferencesResolver{
		dependecyGraph:  ds.NewDependencyGraph[string](),
		referencesStore: valRefStore,
	}
}
