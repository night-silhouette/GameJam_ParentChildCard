package battleservice

import (
	"math/rand"
	"pcc_card/application/entity/Card/CardAbstract"
	"pcc_card/global"
)

type State interface {
	enter()
	exit()
}

type StateMachine struct {
	StateList    map[string]State
	CurrentState State
	c            *Ctx
	CardListCopy *[]CardAbstract.Card
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
	StateMachineImpl.CardListCopy = CardListImpl.Copy()

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
	s.RandomCard()
}

func (s *ShuffleDeal) RandomCard() {
	cList := s.SM.CardListCopy
	rand.Shuffle(len(*cList), func(i, j int) {
		(*cList)[i], (*cList)[j] = (*cList)[j], (*cList)[i]
	})

	numA := global.InitCardNum
	numB := global.InitCardNum
	i := 0
	CardInHandA := make([]CardAbstract.Card, 0, numA)
	s.c.DataA.CardInHand = &CardInHandA
	CardInHandB := make([]CardAbstract.Card, 0, numB)
	s.c.DataB.CardInHand = &CardInHandB
	for ; i < len(*cList); i++ {
		if (*cList)[i].GetInfo()["is_parent"] == true {
			CardInHandA = append(CardInHandA, (*cList)[i])
			numA -= 1
			if numA == 0 {
				break
			}
		}
	}
	for ; i < len(*cList); i++ {
		if (*cList)[i].GetInfo()["is_parent"] == true {
			CardInHandB = append(CardInHandB, (*cList)[i])
			numB -= 1
			if numB == 0 {
				break
			}
		}
	}
}

func (s *ShuffleDeal) exit() {}

//----------------------------------------------------------------------------------------------------------------------
