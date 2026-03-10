package BattleService

import "pcc_card/application/entity"

type State interface {
	enter(ctx entity.BattleCtx)
	exit(ctx entity.BattleCtx)
}

type StateMachine struct {
	StateList    map[string]State
	CurrentState State
	InitState    State
}

func (s *StateMachine) finish(ctx entity.BattleCtx, NextState State) {
	if s.CurrentState != nil {
		s.CurrentState.exit(ctx)
	}
	if NextState == nil && s.CurrentState != NextState {
		s.CurrentState = NextState
		s.CurrentState.enter(ctx)
	}
}

func NewStateMachine(StateList map[string]State) *StateMachine {
	StateMachineImpl := &StateMachine{}
	StateMachineImpl.StateList = StateList
	return StateMachineImpl
}
