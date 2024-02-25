package execution

import (
	"fmt"
	"sync"

	"github.com/conformize/conformize/common/diagnostics"
	"github.com/conformize/conformize/common/ds"
	"github.com/conformize/conformize/common/path"
	"github.com/conformize/conformize/common/pathparser"
)

type ReferencesResolver struct {
	dependecyGraph   *ds.DependencyGraph[string]
	referencesStore  *ValueReferencesStore
	signalChan       chan struct{}
	maxParallelTasks int
}

func (refResolver ReferencesResolver) Resolve(refs map[string]string, diags diagnostics.Diagnosable) {
	defer close(refResolver.signalChan)

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
		var wg sync.WaitGroup
		refsResolveDepOrder := refResolver.dependecyGraph.GetOrder()
		refResolveTasks := make(chan struct{}, min(refResolver.maxParallelTasks, len(refsResolveDepOrder)))
		defer close(refResolveTasks)

		resolvedRefs := make(map[string]chan struct{})
		for _, refAlias := range refsResolveDepOrder {
			resolveChan := make(chan struct{})
			resolvedRefs[refAlias] = resolveChan

			var depChan chan struct{}
			refPathSteps := refPaths[refAlias]
			_, resolved := refResolver.referencesStore.GetReference(refAlias)
			if resolved {
				close(resolveChan)
				continue
			}
			depRefRoot := refPathSteps[0]
			depChan = resolvedRefs[depRefRoot.String()]

			wg.Add(1)
			refResolveTasks <- struct{}{}
			go func(alias string, refSteps path.Steps, depChan <-chan struct{}, resolveChan chan struct{}) {
				defer wg.Done()
				defer close(resolveChan)
				defer func() { <-refResolveTasks }()

				<-depChan
				ref, err := refResolver.referencesStore.GetAtPath(path.NewPath(refSteps))
				if err == nil {
					refResolver.referencesStore.AddReference(alias, ref)
					return
				}

				refRawPath := refs[alias]
				diags.Append(diagnostics.Builder().
					Error().
					Details(
						fmt.Sprintf("\nCouldn't resolve reference '%s' in path %s, reason:\n%s",
							alias, refRawPath, err.Error()),
					).
					Build(),
				)
			}(refAlias, refPathSteps, depChan, resolveChan)
		}
		wg.Wait()
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

func NewReferencesResolver(valRefStore *ValueReferencesStore, maxParallelTasks int, signalChan chan struct{}) *ReferencesResolver {
	return &ReferencesResolver{
		dependecyGraph:   ds.NewDependencyGraph[string](),
		referencesStore:  valRefStore,
		maxParallelTasks: maxParallelTasks,
		signalChan:       signalChan,
	}
}
