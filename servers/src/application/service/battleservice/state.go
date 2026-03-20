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
	Id1          int
	Id2          int
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

func NewStateMachine(c *Ctx, id1 int, id2 int) *StateMachine {
	StateMachineImpl := &StateMachine{}
	StateMachineImpl.c = c //ctx的注入
	StateMachineImpl.Id1 = id1
	StateMachineImpl.Id2 = id2
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
	for {
		OK := s.RandomCard()
		if OK {
			break
		}
	}

}

func (s *ShuffleDeal) RandomCard() bool {
	cList := s.SM.CardListCopy
	rand.Shuffle(len(*cList), func(i, j int) {
		(*cList)[i], (*cList)[j] = (*cList)[j], (*cList)[i]
	})

	numA := global.InitCardNum
	numB := global.InitCardNum
	i := 0
	CardInHandA := make([]CardAbstract.Card, 0, numA)
	s.c.PlayerDataMap[s.SM.Id1].CardInHand = &CardInHandA
	CardInHandB := make([]CardAbstract.Card, 0, numB)
	s.c.PlayerDataMap[s.SM.Id1].CardInHand = &CardInHandB
	CharacterNumA := 0
	CharacterNumB := 0

	for ; i < len(*cList); i++ {
		if (*cList)[i].GetInfo()["is_parent"] == true {
			CardInHandA = append(CardInHandA, (*cList)[i])
			if _, ok := (*cList)[i].(CardAbstract.Character); ok {
				CharacterNumA++
			}
			numA -= 1
			if numA == 0 {
				break
			}
		}
	}
	for ; i < len(*cList); i++ {
		if (*cList)[i].GetInfo()["is_parent"] == true {
			CardInHandB = append(CardInHandB, (*cList)[i])
			if _, ok := (*cList)[i].(CardAbstract.Character); ok {
				CharacterNumB++
			}
			numB -= 1
			if numB == 0 {
				break
			}
		}
	}
	if CharacterNumA <= 3 || CharacterNumB <= 3 {
		return false
	}
	return true
}

func (s *ShuffleDeal) exit() {}

//----------------------------------------------------------------------------------------------------------------------
