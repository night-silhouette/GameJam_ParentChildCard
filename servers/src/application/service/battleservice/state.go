package battleservice

import (
	"context"
	"math/rand"
	"pcc_card/application/entity/Card/CardAbstract"
	"pcc_card/global"
	"pcc_card/presentation/handler/battlehandler/BattleDto"
)

type State interface {
	enter()
	exit()
	Init(id1 int, id2 int, c *Ctx, Nt *NotifyManager, SM *StateMachine)
	process(GoCtx context.Context)
}

type StateMachine struct {
	Id1          int
	Id2          int
	StateList    map[string]State
	CurrentState State
	c            *Ctx
	Nt           *NotifyManager
	CardListCopy *[]CardAbstract.Card
	cancelFunc   context.CancelFunc
}

func (s *StateMachine) process(GoCtx context.Context) {

	handleAction := func(id int, action BattleDto.Action, AcceptChan <-chan BattleDto.Action, ResponseChan chan<- BattleDto.Action) {
		if action.ActionCode == BattleDto.GetSelfCardInHard && action.Predicates == BattleDto.Query { //获取自己手牌
			res := s.c.GetCardInHard(id)
			ResponseChan <- BattleDto.NewAction(BattleDto.GetSelfCardInHard, BattleDto.Result, res.Self)
		}
		if action.ActionCode == BattleDto.GetOpponentCardInHard && action.Predicates == BattleDto.Query { //获取对方手牌
			res := s.c.GetCardInHard(id)
			ResponseChan <- BattleDto.NewAction(BattleDto.GetOpponentCardInHard, BattleDto.Result, res.Opponent)
		}
		if action.ActionCode == BattleDto.OverBattle && action.Predicates == BattleDto.Notify { //结束战斗
			ResponseChan <- BattleDto.NewAction(BattleDto.OverBattle, BattleDto.Notify, "ok")
		}

	}

	for {
		select {
		case <-GoCtx.Done():

			return
		case action := <-s.Nt.ChanMap[s.Id1].AcceptChan:
			handleAction(s.Id1, action, s.Nt.ChanMap[s.Id1].AcceptChan, s.Nt.ChanMap[s.Id1].ResponseChan)
		case action := <-s.Nt.ChanMap[s.Id2].AcceptChan:
			handleAction(s.Id2, action, s.Nt.ChanMap[s.Id2].AcceptChan, s.Nt.ChanMap[s.Id2].ResponseChan)
		}
	}
}

func (s *StateMachine) finish(NextState string) {
	NextStateObj, _ := s.StateList[NextState]
	if s.cancelFunc != nil {
		s.cancelFunc()
	}
	if s.CurrentState != nil {
		s.CurrentState.exit()
	}
	if NextState != "" && s.CurrentState != NextStateObj {
		s.CurrentState = NextStateObj
		s.CurrentState.enter()
		var GoCtx context.Context
		GoCtx, s.cancelFunc = context.WithCancel(context.Background())
		go s.CurrentState.process(GoCtx)

	}
}

func NewStateMachine(c *Ctx, id1 int, id2 int, Nt *NotifyManager) (*StateMachine, context.CancelFunc) {
	StateMachineImpl := &StateMachine{}
	StateMachineImpl.c = c //ctx的注入
	StateMachineImpl.Id1 = id1
	StateMachineImpl.Id2 = id2
	StateMachineImpl.Nt = Nt //Nt的注入
	StateMachineImpl.CardListCopy = CardListImpl.Copy()

	StateMachineImpl.RegisterState()
	for _, element := range StateMachineImpl.StateList {
		element.Init(id1, id2, c, Nt, StateMachineImpl)
	}
	StateMachineImpl.finish("shuffleDeal") //这个状态的enter结束之后，才会开启state machine的process
	GoCtx, cancelFunc := context.WithCancel(context.Background())
	go StateMachineImpl.process(GoCtx)
	return StateMachineImpl, cancelFunc
}

func (s *StateMachine) RegisterState() {
	s.StateList = map[string]State{
		"shuffleDeal":         &ShuffleDeal{},
		"SelectCharacterCard": &SelectCharacterCard{},
		"SelectSkillCard":     &SelectSkillCard{},
		"Judge":               &Judge{},
		"Combat":              &Combat{},
		"SkillCardCalc":       &SkillCardCalc{},
	}
}

//----------------------------------------------------------------------------------------------------------------------

type StateTemplate struct {
	Id1 int
	Id2 int
	c   *Ctx
	Nt  *NotifyManager
	SM  *StateMachine
}

func (s *StateTemplate) Init(id1 int, id2 int, c *Ctx, Nt *NotifyManager, SM *StateMachine) {
	s.Id1 = id1
	s.Id2 = id2
	s.c = c
	s.Nt = Nt
	s.SM = SM
}

// -------------------------------------ShuffleDeal---------------------------------------------------------------------------------------

type ShuffleDeal struct {
	StateTemplate
}

func (s *ShuffleDeal) enter() {
	for {
		OK := s.RandomCard()
		if OK {
			break
		}
	}
}

func (s *ShuffleDeal) process(GoCtx context.Context) {

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
	s.c.PlayerDataMap[s.SM.Id2].CardInHand = &CardInHandB
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

//---------------------------------------SelectCharacterCard-------------------------------------------------------------------------------

type SelectCharacterCard struct {
	StateTemplate
}

func (i *SelectCharacterCard) enter() {
}

func (i *SelectCharacterCard) exit() {
}

func (i *SelectCharacterCard) process(GoCtx context.Context) {
}

//---------------------------------------SelectSkillCard-------------------------------------------------------------------------------

type SelectSkillCard struct {
	StateTemplate
}

func (s *SelectSkillCard) enter()                        {}
func (s *SelectSkillCard) exit()                         {}
func (s *SelectSkillCard) process(GoCtx context.Context) {}

//---------------------------------------Judge-------------------------------------------------------------------------------

type Judge struct {
	StateTemplate
}

func (J *Judge) enter()                        {}
func (J *Judge) exit()                         {}
func (J *Judge) process(GoCtx context.Context) {}

//---------------------------------------Combat-------------------------------------------------------------------------------

type Combat struct {
	StateTemplate
}

func (c *Combat) enter()                        {}
func (c *Combat) exit()                         {}
func (c *Combat) process(GoCtx context.Context) {}

//---------------------------------------SkillCardCalc-------------------------------------------------------------------------------

type SkillCardCalc struct {
	StateTemplate
}

func (s *SkillCardCalc) enter()                        {}
func (s *SkillCardCalc) exit()                         {}
func (s *SkillCardCalc) process(GoCtx context.Context) {}
