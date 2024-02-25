package execution

type ProvidersInitializationPhase struct {
	steps []ExecutionStep
}

func (phase *ProvidersInitializationPhase) AddStep(step ExecutionStep) {
	phase.steps = append(phase.steps, step)
}

func (phase *ProvidersInitializationPhase) Execute(blprntExecCtx *BlueprintExecutionContext) {
	for _, step := range phase.steps {
		step.Run(blprntExecCtx)
	}
}

func NewProvidersInitializationPhase() *ProvidersInitializationPhase {
	return &ProvidersInitializationPhase{
		steps: make([]ExecutionStep, 0, 10),
	}
}
