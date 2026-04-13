package battleservice

import (
	"context"
	"math/rand"
	"pcc_card/Util"
	"pcc_card/application/entity/Card/CardAbstract"
	"pcc_card/global"
	"pcc_card/presentation/handler/battlehandler/BattleDto"
	"sync"
	"time"
)

type State interface {
	enter()
	exit()
	Init(id1 int, id2 int, c *Ctx, Nt *NotifyManager, SM *StateMachine)
	process(GoCtx context.Context)
	AddTaskCount()
}

type StateMachine struct {
	Mutex         sync.RWMutex
	ParentNodeCtx context.Context

	Id1          int
	Id2          int
	StateList    map[string]State
	CurrentState State
	c            *Ctx
	Nt           *NotifyManager
	CardListCopy *[]CardAbstract.Card
	cancelFunc   context.CancelFunc
}

func (s *StateMachine) AcceptAction(goCtx context.Context, handleAction func(id int, action BattleDto.Action, ResponseChan chan<- BattleDto.Action)) {
	for {
		select {
		case <-goCtx.Done():
			return
		case action := <-s.Nt.ChanMap[s.Id1].AcceptChan:
			handleAction(s.Id1, action, s.Nt.ChanMap[s.Id1].ResponseChan)
			s.SharedProcess(s.Id1, action, s.Nt.ChanMap[s.Id1].ResponseChan)
		case action := <-s.Nt.ChanMap[s.Id2].AcceptChan:
			handleAction(s.Id2, action, s.Nt.ChanMap[s.Id2].ResponseChan)
			s.SharedProcess(s.Id2, action, s.Nt.ChanMap[s.Id2].ResponseChan)
		}
	}
}

func (s *StateMachine) SharedProcess(id int, action BattleDto.Action, ResponseChan chan<- BattleDto.Action) {
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

func (s *StateMachine) SendActionById(id int, action BattleDto.Action) {
	s.Nt.ChanMap[id].ResponseChan <- action
}

func (s *StateMachine) finish(NextState string) {
	NextStateObj, _ := s.StateList[NextState]

	if s.CurrentState == NextStateObj {
		s.CurrentState.AddTaskCount()
		return
	}
	if s.cancelFunc != nil {
		s.cancelFunc()
	}
	if s.CurrentState != nil {
		s.CurrentState.exit()
	}
	if NextState != "" {
		s.CurrentState = NextStateObj
		s.CurrentState.enter()
		var GoCtx context.Context
		GoCtx, s.cancelFunc = context.WithCancel(s.ParentNodeCtx)
		go s.CurrentState.process(GoCtx)

	}
}

func NewStateMachine(c *Ctx, id1 int, id2 int, Nt *NotifyManager, ParentNodeCtx context.Context) *StateMachine {

	StateMachineImpl := &StateMachine{}
	StateMachineImpl.ParentNodeCtx = ParentNodeCtx
	StateMachineImpl.c = c //ctx的注入
	StateMachineImpl.Id1 = id1
	StateMachineImpl.Id2 = id2
	StateMachineImpl.Nt = Nt //Nt的注入
	StateMachineImpl.CardListCopy = c.CardPool

	StateMachineImpl.RegisterState()
	for _, element := range StateMachineImpl.StateList {
		element.Init(id1, id2, c, Nt, StateMachineImpl)
	}
	StateMachineImpl.finish("ShuffleDeal")
	return StateMachineImpl
}

func (s *StateMachine) RegisterState() {
	s.StateList = map[string]State{
		"ShuffleDeal":         &ShuffleDeal{},
		"SelectCharacterCard": &SelectCharacterCard{},
		"SelectSkillCard":     &SelectSkillCard{},
		"Judge":               &Judge{},
		"Combat":              &Combat{},
		"SkillCardCalc":       &SkillCardCalc{},
	}
}

//----------------------------------------------------------------------------------------------------------------------

type StateTemplate struct {
	Id1       int
	Id2       int
	c         *Ctx
	Nt        *NotifyManager
	SM        *StateMachine
	TaskCount int
}

func (s *StateTemplate) Init(id1 int, id2 int, c *Ctx, Nt *NotifyManager, SM *StateMachine) {
	s.Id1 = id1
	s.Id2 = id2
	s.c = c
	s.Nt = Nt
	s.SM = SM
	s.TaskCount = 0
}
func (s *StateTemplate) exit() {
	s.TaskCount = 0
}
func (s *StateTemplate) AddTaskCount() {
	s.SM.Mutex.Lock()
	defer s.SM.Mutex.Unlock()
	s.TaskCount++
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
	s.SM.finish("SelectCharacterCard")
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

func (s *ShuffleDeal) exit() {
	s.StateTemplate.exit()
}

//---------------------------------------SelectCharacterCard-------------------------------------------------------------------------------

type SelectCharacterCard struct {
	IsFirst bool
	StateTemplate
}

func (s *SelectCharacterCard) Init(id1 int, id2 int, c *Ctx, Nt *NotifyManager, SM *StateMachine) { //复写init，初始化isfrist
	s.Id1 = id1
	s.Id2 = id2
	s.c = c
	s.Nt = Nt
	s.SM = SM
	s.IsFirst = true
}

func (s *SelectCharacterCard) enter() {
	var waitTime time.Duration
	waitTime = global.SelectCharacterTime * time.Second
	if s.IsFirst {
		s.IsFirst = false
		waitTime = 25
	}
	act := BattleDto.NewAction(BattleDto.SelectCharacterCard, BattleDto.Query, Util.SendTime(waitTime))
	s.SM.SendActionById(s.Id1, act)
	s.SM.SendActionById(s.Id1, act)

}

func (s *SelectCharacterCard) exit() {
	s.StateTemplate.exit()
}

func (s *SelectCharacterCard) process(GoCtx context.Context) {

	handleAction := func(id int, action BattleDto.Action, ResponseChan chan<- BattleDto.Action) {
		if action.ActionCode == BattleDto.SelectCharacterCard && action.Predicates == BattleDto.Result {

			//todo 上牌

			s.SM.Mutex.Lock()
			s.TaskCount--
			if s.TaskCount == 0 {
				s.SM.finish("SelectSkillCard")
			}
		}
	}
	s.SM.AcceptAction(GoCtx, handleAction)

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
