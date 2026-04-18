package battleservice

import (
	"context"
	"fmt"
	"math/rand"
	"pcc_card/Util"
	"pcc_card/application/entity/BattleData"
	"pcc_card/application/entity/Card/CardAbstract"
	"pcc_card/global"
	"pcc_card/presentation/handler/battlehandler/BattleDto"
	"sync"
	"time"

	"github.com/mitchellh/mapstructure"
)

type State interface {
	enter()
	exit()
	Init(id1 int, id2 int, c *Ctx, Nt *NotifyManager, SM *StateMachine, sub State)
	process(GoCtx context.Context)
	AddTaskCount()
	SetName(name string)
	GetName() string
	SpecialInit()
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

//#region StateMachine

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
			handleAction(s.Id2, action, s.Nt.ChanMap[s.Id2].ResponseChan) //!!!
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
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
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
		element.Init(id1, id2, c, Nt, StateMachineImpl, element)
	}
	StateMachineImpl.finish("ShuffleDeal")
	return StateMachineImpl
}

//#endregion
//#region StateTemplate

type StateTemplate struct {
	name      string
	Id1       int
	Id2       int
	c         *Ctx
	Nt        *NotifyManager
	SM        *StateMachine
	TaskCount int
}

func (s *StateTemplate) Init(id1 int, id2 int, c *Ctx, Nt *NotifyManager, SM *StateMachine, sub State) {
	s.Id1 = id1
	s.Id2 = id2
	s.c = c
	s.Nt = Nt
	s.SM = SM
	s.TaskCount = 0
	sub.SpecialInit()
}
func (s *StateTemplate) SpecialInit() {}
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

//#endregion
//#region State:ShuffleDeal

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
	go s.SM.finish("Judge")
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

//#endregion
//#region State:SelectCharacterCard

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

//#endregion
//#region State:SelectSkillCard

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
						if s.TaskMap[s.Id1] && s.TaskMap[s.Id2] {
							s.SM.SendActionById(s.Id1, BattleDto.NewAction(BattleDto.DeployCard, BattleDto.Finish, "技能牌全部选择完毕"))
						}
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

//#endregion
//region State:Judge

type Judge struct {
	StateTemplate
	Mutex     sync.Mutex
	TaskMap   map[int]int
	ChanStop  chan struct{}
	ChanCrash chan struct{}
}

type JudgeData struct {
	JudgeData int `json:"judge_data" mapstructure:"judge_data"`
}

func (J *Judge) SpecialInit() {
	J.TaskMap = make(map[int]int)
}

func JudgeWin(Jd1 int, Jd2 int) int { //输出Jd1是否win
	if Jd1 == Jd2 {
		return 0
	}
	if (Jd1+1)%3 == Jd2 {
		return 1
	}
	return -1
}

type JudgeEndTime struct {
	JudgeEndTime int64 `json:"judge_end_time" mapstructure:"judge_end_time"`
}

func NewJudgeEndTime() JudgeEndTime {
	result := JudgeEndTime{}
	result.JudgeEndTime = Util.SendTime(global.JudgeWaitTime)
	return result
}

func (J *Judge) EndJudge() {
	J.Mutex.Lock()
	defer J.Mutex.Unlock()
	for Key, value := range J.TaskMap {
		if value == 3 {
			J.TaskMap[Key] = Util.RandomRange(0, 2)
		}
	}

	J.SM.SendActionById(J.Id1, BattleDto.NewAction(BattleDto.Judge, BattleDto.Finish, NewJudgeRes(J.TaskMap[J.Id1], J.TaskMap[J.Id2], JudgeWin(J.TaskMap[J.Id1], J.TaskMap[J.Id2]))))
	J.SM.SendActionById(J.Id2, BattleDto.NewAction(BattleDto.Judge, BattleDto.Finish, NewJudgeRes(J.TaskMap[J.Id2], J.TaskMap[J.Id1], JudgeWin(J.TaskMap[J.Id2], J.TaskMap[J.Id1]))))

}

func (J *Judge) enter() {
	J.TaskMap[J.Id1] = 3 //设为一个不可能值作为检查是否返回了
	J.TaskMap[J.Id2] = 3

	chanStop, chanCrash := Util.CreateTimer(time.Second*global.JudgeWaitTime, J.EndJudge)
	J.ChanCrash = chanCrash
	J.ChanStop = chanStop

	J.SM.SendActionById(J.Id1, BattleDto.NewAction(BattleDto.Judge, BattleDto.Query, NewJudgeEndTime()))
	J.SM.SendActionById(J.Id2, BattleDto.NewAction(BattleDto.Judge, BattleDto.Query, NewJudgeEndTime()))
}
func (J *Judge) exit() {
	J.TaskMap[J.Id1] = 3
	J.TaskMap[J.Id1] = 3
	J.ChanStop = nil
	J.ChanCrash = nil
}

type JudgeRes struct {
	Self     int `json:"self" mapstructure:"self"`
	Opponent int `json:"opponent" mapstructure:"opponent"`
	IsWin    int `json:"is_win" mapstructure:"is_win"`
}

func NewJudgeRes(self int, opponent int, IsWin int) *JudgeRes {
	J := &JudgeRes{}
	J.Self = self
	J.Opponent = opponent
	J.IsWin = IsWin
	return J
}

func (J *Judge) process(GoCtx context.Context) {
	handleAction := func(id int, action BattleDto.Action, ResponseChan chan<- BattleDto.Action) {
		if action.ActionCode == BattleDto.Judge && action.Predicates == BattleDto.Result {
			J.Mutex.Lock()
			var data JudgeData
			err := mapstructure.Decode(action.ActionData, &data)
			if err != nil {
				fmt.Println(err)
				J.SM.SendActionById(id, BattleDto.NewErrAction(global.ResponseInvalidReqParams))
				return
			}
			Jd := data.JudgeData
			if !(0 <= Jd && Jd <= 2) {
				J.SM.SendActionById(id, BattleDto.NewErrAction(global.ResponseInvalidReqParams))
			}
			J.TaskMap[id] = Jd
			J.SM.SendActionById(id, BattleDto.NewAction(BattleDto.Judge, BattleDto.Succeed, "")) //单方选好了，存储进去了

			flag := true
			for _, value := range J.TaskMap {
				if value == 3 {
					flag = false
				}
			}
			J.Mutex.Unlock()
			if flag { //双方都已经选好了
				J.ChanStop <- struct{}{}
			}

		}

	}
	J.SM.AcceptAction(GoCtx, handleAction)
}

//endregion

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
