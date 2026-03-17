package battleservice

type State interface {
	enter()
	exit()
}

type StateMachine struct {
	StateList    map[string]State
	CurrentState State
	c            *Ctx
}

func (s *StateMachine) finish(NextState State) {
	if s.CurrentState != nil {
		s.CurrentState.exit()
	}
	if NextState == nil && s.CurrentState != NextState {
		s.CurrentState = NextState
		s.CurrentState.enter()
	}
}

func NewStateMachine(c *Ctx) *StateMachine {
	StateMachineImpl := &StateMachine{}
	StateMachineImpl.c = c
	StateMachineImpl.StateList = map[string]State{
		"shuffleDeal": &ShuffleDeal{c, StateMachineImpl},
	}
	StateMachineImpl.CurrentState = StateMachineImpl.StateList["shuffleDeal"]
	StateMachineImpl.CurrentState.enter()

	return StateMachineImpl
}

//----------------------------------------------------------------------------------------------------------------------

type ShuffleDeal struct {
	c  *Ctx
	SM *StateMachine
}

func (s *ShuffleDeal) enter() {

}

func (s *ShuffleDeal) RandomCard() {

}

func (s *ShuffleDeal) exit() {}

//----------------------------------------------------------------------------------------------------------------------
