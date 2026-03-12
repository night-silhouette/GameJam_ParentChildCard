package BattleService

type State interface {
	enter(ctx Ctx)
	exit(ctx Ctx)
}

type StateMachine struct {
	StateList    map[string]State
	CurrentState State
	InitState    State
}

func (s *StateMachine) finish(ctx Ctx, NextState State) {
	if s.CurrentState != nil {
		s.CurrentState.exit(ctx)
	}
	if NextState == nil && s.CurrentState != NextState {
		s.CurrentState = NextState
		s.CurrentState.enter(ctx)
	}
}

func NewStateMachine() *StateMachine {
	StateMachineImpl := &StateMachine{}
	StateMachineImpl.StateList = make(map[string]State)

	return StateMachineImpl
}
