package battleservice

import (
	"context"
	"fmt"
	"math/rand"
	"pcc_card/application/entity/BattleData"
	"pcc_card/application/entity/Card/CardAbstract"
	"pcc_card/global"
	"pcc_card/presentation/handler/battlehandler/BattleDto"
	"sync"

	"github.com/mitchellh/mapstructure"
)

type State interface {
	enter()
	exit()
	Init(id1 int, id2 int, c *Ctx, Nt *NotifyManager, SM *StateMachine)
	process(GoCtx context.Context)
	AddTaskCount()
	SetName(name string)
	GetName() string
}

type StateMachine struct {
	Mutex         sync.RWMutex
	ParentNodeCtx context.Context

	Id1          int
	Id2          int
	StateList    map[string]State
	CurrentState State
	StateStack   []State
	c            *Ctx
	Nt           *NotifyManager
	CardListCopy *[]CardAbstract.Card
	cancelFunc   context.CancelFunc
}

func (s *StateMachine) StatePush(CurrentState string, NewState string) {
	temp := s.StateList[CurrentState]
	s.StateStack = append(s.StateStack, temp) //把现在的state压入栈
	s.finish(NewState)                        //切换到新的state
}
func (s *StateMachine) StatePop() { //切换到，上一次压栈的状态
	if len(s.StateStack) == 0 {
		return
	}
	lastIndex := len(s.StateStack) - 1
	pop := s.StateStack[lastIndex]
	s.finish(pop.GetName())
	s.StateStack[lastIndex] = nil
	s.StateStack = s.StateStack[:lastIndex]
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
		//s.CurrentState.AddTaskCount()
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
	c.StateMachine = StateMachineImpl
	StateMachineImpl.ParentNodeCtx = ParentNodeCtx
	StateMachineImpl.c = c //ctx的注入
	StateMachineImpl.Id1 = id1
	StateMachineImpl.Id2 = id2
	StateMachineImpl.Nt = Nt //Nt的注入
	StateMachineImpl.CardListCopy = c.CardPool
	StateMachineImpl.StateStack = make([]State, 0)

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
	for key, element := range s.StateList {
		element.SetName(key)
	}
}

//----------------------------------------------------------------------------------------------------------------------

type StateTemplate struct {
	name      string
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
	s.SpecialInit()
}
func (s *StateTemplate) SpecialInit() {

}
func (s *StateTemplate) exit() {
	s.TaskCount = 0
}

func (s *StateTemplate) process(GoCtx context.Context) {

	handleAction := func(id int, action BattleDto.Action, ResponseChan chan<- BattleDto.Action) {
	}
	s.SM.AcceptAction(GoCtx, handleAction)
}

func (s *StateTemplate) AddTaskCount() {
	s.SM.Mutex.Lock()
	defer s.SM.Mutex.Unlock()
	s.TaskCount++
}

func (s *StateTemplate) SetName(name string) {
	s.name = name
}

func (s *StateTemplate) GetName() string {
	return s.name
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
	s.SM.finish("SelectSkillCard")
}

func (s *ShuffleDeal) process(GoCtx context.Context) {
	//空的逻辑，不继承templete
}

func (s *ShuffleDeal) RandomCard() bool {
	cList := s.SM.CardListCopy
	for _, card := range *cList {
		card.SetBtCtx(s.c)
		card.SetTempId(s.c.entityCounter)
		s.c.entityCounter++
	}

	rand.Shuffle(len(*cList), func(i, j int) {
		(*cList)[i], (*cList)[j] = (*cList)[j], (*cList)[i]
	})

	numA := global.InitCardNum
	numB := global.InitCardNum
	i := 0
	CardInHandA := make(map[int]CardAbstract.Card)
	s.c.PlayerDataMap[s.SM.Id1].CardInHand = CardInHandA
	CardInHandB := make(map[int]CardAbstract.Card)
	s.c.PlayerDataMap[s.SM.Id2].CardInHand = CardInHandB
	CharacterNumA := 0
	CharacterNumB := 0

	for ; i < len(*cList); i++ {
		if (*cList)[i].GetInfo()["is_parent"] == true { //id1
			(*cList)[i].SetOwnerId(s.Id1)
			CardInHandA[(*cList)[i].GetTempId()] = (*cList)[i]
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
			(*cList)[i].SetOwnerId(s.Id2)
			CardInHandB[(*cList)[i].GetTempId()] = (*cList)[i]
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
	s.SM.SendActionById(s.Id1, BattleDto.NewAction(BattleDto.StartBattle, BattleDto.Notify, ""))
	s.SM.SendActionById(s.Id2, BattleDto.NewAction(BattleDto.StartBattle, BattleDto.Notify, ""))
}

//---------------------------------------SelectCharacterCard-------------------------------------------------------------------------------

type SelectCharacterCard struct {
	IsFirst bool
	StateTemplate
}

func (s *SelectCharacterCard) SpecialInit() {
	s.IsFirst = false
}

func (s *SelectCharacterCard) enter() {
	//var waitTime time.Duration
	//waitTime = global.SelectCharacterTime * time.Second
	//if s.IsFirst {
	//	s.IsFirst = false
	//	waitTime = 25
	//}
	//act := BattleDto.NewAction(BattleDto.SelectCharacterCard, BattleDto.Query, Util.SendTime(waitTime))
	//s.SM.SendActionById(s.Id1, act)
	//s.SM.SendActionById(s.Id1, act)

}

func (s *SelectCharacterCard) exit() {
	s.StateTemplate.exit()
}

func (s *SelectCharacterCard) process(GoCtx context.Context) {

	handleAction := func(id int, action BattleDto.Action, ResponseChan chan<- BattleDto.Action) {
		//if action.ActionCode == BattleDto.SelectCharacterCard && action.Predicates == BattleDto.Result {
		//
		//	//todo 上牌
		//
		//	s.SM.Mutex.Lock()
		//	s.TaskCount--
		//	if s.TaskCount <= 0 {
		//		s.SM.finish("SelectSkillCard")
		//	}
		//}
	}
	s.SM.AcceptAction(GoCtx, handleAction)

}

//---------------------------------------SelectSkillCard-------------------------------------------------------------------------------

type SelectSkillCard struct {
	StateTemplate
	TaskMap map[int]bool
}

func (s *SelectSkillCard) SpecialInit() {
	s.TaskMap = make(map[int]bool)
	s.TaskMap[s.Id1] = false
	s.TaskMap[s.Id2] = false
}

func (s *SelectSkillCard) enter() {
	s.SM.SendActionById(s.Id1, BattleDto.NewAction(BattleDto.DeployCard, BattleDto.Query, BattleData.SelectCard{Where: BattleData.SkillCard}))
	s.SM.SendActionById(s.Id2, BattleDto.NewAction(BattleDto.DeployCard, BattleDto.Query, BattleData.SelectCard{Where: BattleData.SkillCard}))

}
func (s *SelectSkillCard) exit() {
	s.StateTemplate.exit()
	s.TaskMap[s.Id1] = false
	s.TaskMap[s.Id2] = false
}
func (s *SelectSkillCard) process(GoCtx context.Context) {
	handleAction := func(id int, action BattleDto.Action, ResponseChan chan<- BattleDto.Action) {
		if s.TaskMap[id] {
			s.SM.SendActionById(s.Id1, BattleDto.NewErrAction(global.ResponseRepeatRequest))
			return
		}

		if action.ActionCode == BattleDto.DeployCard && action.Predicates == BattleDto.Result {

			var data BattleData.SelectCard
			err := mapstructure.Decode(action.ActionData, &data)
			if err != nil {
				fmt.Println(err)
				s.SM.SendActionById(id, BattleDto.NewErrAction(global.ResponseInvalidReqParams))
				return
			}
			if data.Where == BattleData.SkillCard {

				cardTempId := data.CardTempId
				if card, ok := s.c.PlayerDataMap[id].CardInHand[cardTempId]; ok { //手牌里有不有
					if _, ok := card.(CardAbstract.SkillCard); ok { //上的是不是skillcard
						delete(s.c.PlayerDataMap[id].CardInHand, cardTempId)
						s.c.SetSkillCardBT(id, card)
						s.SM.SendActionById(id, BattleDto.NewAction(BattleDto.DeployCard, BattleDto.Succeed, "技能牌选择成功"))
						s.TaskMap[id] = true
					} else {
						s.SM.SendActionById(id, BattleDto.NewErrAction(global.BattleCardCategoryError))
						return
					}
				} else {
					s.SM.SendActionById(id, BattleDto.NewErrAction(global.BattleCardNotFound))
					return
				}

			}

		}
	}
	s.SM.AcceptAction(GoCtx, handleAction)

}

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
