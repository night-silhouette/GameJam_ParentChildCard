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
		"CardCalc":            &CardCalc{},
	}
	for key, element := range s.StateList {
		element.SetName(key)
	}
}
func (s *StateMachine) SharedProcess(id int, action BattleDto.Action, ResponseChan chan<- BattleDto.Action) bool {
	if action.ActionCode == BattleDto.GetSelfCardInHard && action.Predicates == BattleDto.Query { //获取自己手牌
		s.Mutex.Lock()
		res := s.c.GetCardInHard(id)
		s.Mutex.Unlock()

		ResponseChan <- BattleDto.NewAction(BattleDto.GetSelfCardInHard, BattleDto.Result, res.Self)

		return true
	}
	if action.ActionCode == BattleDto.GetOpponentCardInHard && action.Predicates == BattleDto.Query { //获取对方手牌
		s.Mutex.Lock()
		res := s.c.GetCardInHard(id)
		s.Mutex.Unlock()

		ResponseChan <- BattleDto.NewAction(BattleDto.GetOpponentCardInHard, BattleDto.Result, res.Opponent)

		return true
	}
	if action.ActionCode == BattleDto.OverBattle && action.Predicates == BattleDto.Notify { //结束战斗
		ResponseChan <- BattleDto.NewAction(BattleDto.OverBattle, BattleDto.Notify, "ok")
		return true
	}
	if action.ActionCode == BattleDto.GetBtCardInfo && action.Predicates == BattleDto.Query { //获取战斗卡信息
		s.Mutex.Lock()
		res := s.c.GetBtCardInfo(id)
		s.Mutex.Unlock()

		ResponseChan <- BattleDto.NewAction(BattleDto.GetBtCardInfo, BattleDto.Result, res)
		return true
	}
	return false
}

//#region StateMachine

type StateMachine struct {
	StateChangeMtx sync.Mutex
	Mutex          sync.RWMutex
	ParentNodeCtx  context.Context

	Id1          int
	Id2          int
	StateList    map[string]State
	CurrentState State
	StateStack   []State
	c            *Ctx
	Nt           *NotifyManager
	CardListCopy *[]CardAbstract.Card
	cancelFunc   context.CancelFunc

	//stateData
	Winner         int
	Loser          int
	CombatDataChan chan BattleData.CombatDto
	CombatTime     time.Duration
}

type StateWaitTime struct {
	StateWaitTime int64 `json:"state_wait_time" mapstructure:"state_wait_time"`
}

func NewStateWaitTime(time time.Duration) StateWaitTime {
	result := StateWaitTime{}
	result.StateWaitTime = Util.SendTime(time)
	return result
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

func (s *StateMachine) AcceptAction(goCtx context.Context, handleAction func(id int, action BattleDto.Action, ResponseChan chan<- BattleDto.Action) bool) {
	for {
		select {
		case <-goCtx.Done():

			return
		case action := <-s.Nt.ChanMap[s.Id1].AcceptChan:
			if handleAction(s.Id1, action, s.Nt.ChanMap[s.Id1].ResponseChan) {
				continue
			}

			if s.SharedProcess(s.Id1, action, s.Nt.ChanMap[s.Id1].ResponseChan) {
				continue
			}
			s.SendActionById(s.Id1, BattleDto.NewErrAction(global.BattleInvalidTiming))
		case action := <-s.Nt.ChanMap[s.Id2].AcceptChan:
			if handleAction(s.Id2, action, s.Nt.ChanMap[s.Id2].ResponseChan) {
				continue
			}
			if s.SharedProcess(s.Id2, action, s.Nt.ChanMap[s.Id2].ResponseChan) {
				continue
			}
			s.SendActionById(s.Id2, BattleDto.NewErrAction(global.BattleInvalidTiming))
		}
	}
}

func (s *StateMachine) SendActionById(id int, action BattleDto.Action) {
	s.Nt.ChanMap[id].ResponseChan <- action
}

func (s *StateMachine) finish(NextState string) {
	s.StateChangeMtx.Lock()
	defer s.StateChangeMtx.Unlock()
	NextStateObj, _ := s.StateList[NextState]

	if s.cancelFunc != nil {
		s.cancelFunc()
	}
	if s.CurrentState != nil {
		s.CurrentState.exit()
		fmt.Print(s.CurrentState.GetName() + "->")
	}

	if NextState != "" {
		s.CurrentState = NextStateObj
		s.CurrentState.enter()

		var GoCtx context.Context
		GoCtx, s.cancelFunc = context.WithCancel(s.ParentNodeCtx) //stateMachine死掉，你也得死
		go s.CurrentState.process(GoCtx)
		fmt.Print(s.CurrentState.GetName() + "\n")

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
	StateMachineImpl.CombatDataChan = make(chan BattleData.CombatDto, 1)
	StateMachineImpl.CombatTime = global.CombatWaitTime * time.Second

	StateMachineImpl.RegisterState()
	for _, element := range StateMachineImpl.StateList {
		element.Init(id1, id2, c, Nt, StateMachineImpl, element)
	}
	go func() { //游戏结束，发通知
		select {
		case <-ParentNodeCtx.Done():
			StateMachineImpl.SendActionById(StateMachineImpl.Id1, BattleDto.NewAction(BattleDto.OverBattle, BattleDto.Notify, ""))
			StateMachineImpl.SendActionById(StateMachineImpl.Id2, BattleDto.NewAction(BattleDto.OverBattle, BattleDto.Notify, ""))
		}
	}()

	StateMachineImpl.finish("ShuffleDeal")
	return StateMachineImpl
}

//#endregion
//#region StateTemplate

type StateTemplate struct {
	name string
	Id1  int
	Id2  int
	c    *Ctx
	Nt   *NotifyManager
	SM   *StateMachine
}

func (s *StateTemplate) Init(id1 int, id2 int, c *Ctx, Nt *NotifyManager, SM *StateMachine, sub State) {
	s.Id1 = id1
	s.Id2 = id2
	s.c = c
	s.Nt = Nt
	s.SM = SM
	sub.SpecialInit()
}
func (s *StateTemplate) SpecialInit() {}
func (s *StateTemplate) exit() {
}

func (s *StateTemplate) process(GoCtx context.Context) {

	handleAction := func(id int, action BattleDto.Action, ResponseChan chan<- BattleDto.Action) bool {
		return false
	}
	s.SM.AcceptAction(GoCtx, handleAction)
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
	s.SM.SendActionById(s.Id1, BattleDto.NewAction(BattleDto.MatchSuccess, BattleDto.Notify, NewStateWaitTime(global.BattleWaitTime*time.Second))) //通知匹配成功
	s.SM.SendActionById(s.Id2, BattleDto.NewAction(BattleDto.MatchSuccess, BattleDto.Notify, NewStateWaitTime(global.BattleWaitTime*time.Second)))

	Util.CreateTimer(global.BattleWaitTime*time.Second, func() { //准备时间过后，正式开始战斗
		s.SM.Mutex.Lock()
		if !s.c.CheckCard(s.Id1) {
			s.c.RandomSelectCard(s.Id1)
			s.SM.SendActionById(s.Id1, BattleDto.NewAction(BattleDto.DeployCard, BattleDto.Succeed, "自动选择"))
		}
		if !s.c.CheckCard(s.Id2) {
			s.c.RandomSelectCard(s.Id2)
			s.SM.SendActionById(s.Id2, BattleDto.NewAction(BattleDto.DeployCard, BattleDto.Succeed, "自动选择"))
		}

		s.SM.SendActionById(s.Id1, BattleDto.NewAction(BattleDto.StartBattle, BattleDto.Notify, ""))
		s.SM.SendActionById(s.Id2, BattleDto.NewAction(BattleDto.StartBattle, BattleDto.Notify, ""))
		go s.SM.finish("SelectSkillCard")
		s.SM.Mutex.Unlock()
	}) //定时开始战斗
}

func (s *ShuffleDeal) process(GoCtx context.Context) {
	handleAction := func(id int, action BattleDto.Action, ResponseChan chan<- BattleDto.Action) bool { //监听上战斗牌

		if action.ActionCode == BattleDto.DeployCard && action.Predicates == BattleDto.Result {
			var data BattleData.SelectCard
			err := mapstructure.Decode(action.ActionData, &data)
			if err != nil {
				fmt.Println(err)
				s.SM.SendActionById(id, BattleDto.NewErrAction(global.ResponseInvalidReqParams))
				return true
			}
			if data.Where != BattleData.SkillCard {

				cardTempId := data.CardTempId

				if card, ok := s.c.PlayerDataMap[id].CardInHand[cardTempId]; ok { //手牌里有不有
					if _, ok := card.(CardAbstract.SkillCard); !ok {
						s.c.SetCardBt(id, card)
						s.SM.SendActionById(id, BattleDto.NewAction(BattleDto.DeployCard, BattleDto.Succeed, "选择成功"))

					} else {
						s.SM.SendActionById(id, BattleDto.NewErrAction(global.BattleCardCategoryError))
						return true
					}
				} else {
					s.SM.SendActionById(id, BattleDto.NewErrAction(global.BattleCardNotFound))
					return true
				}

			} else {
				s.SM.SendActionById(id, BattleDto.NewErrAction(global.BattleInvalidTiming))
			}
			return false

		}
		return false
	}
	s.SM.AcceptAction(GoCtx, handleAction)
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

}

//#endregion
//#region State:SelectCharacterCard

type SelectCharacterCard struct {
	StateTemplate
	ChanCrash chan struct{}
	ChanStop  chan struct{}
}

func (s *SelectCharacterCard) SpecialInit() {
	s.ChanCrash = make(chan struct{})
	s.ChanStop = make(chan struct{})
}

func (s *SelectCharacterCard) enter() {

}

func (s *SelectCharacterCard) exit() {
	s.StateTemplate.exit()
}

func (s *SelectCharacterCard) process(GoCtx context.Context) {

	handleAction := func(id int, action BattleDto.Action, ResponseChan chan<- BattleDto.Action) bool {

		if action.ActionCode == BattleDto.DeployCard && action.Predicates == BattleDto.Result {
			var data BattleData.SelectCard //解析actionData
			err := mapstructure.Decode(action.ActionData, &data)
			if err != nil {
				fmt.Println(err)
				s.SM.SendActionById(id, BattleDto.NewErrAction(global.ResponseInvalidReqParams))
				return true
			}
			if data.Where == BattleData.SkillCard {

			}
		}
		return false
	}
	s.SM.AcceptAction(GoCtx, handleAction)

}

//#endregion
//#region State:SelectSkillCard

type SelectSkillCard struct {
	Mutex sync.RWMutex
	StateTemplate
	TaskMap   map[int]bool
	ChanCrash chan struct{}
	ChanStop  chan struct{}
}

func (s *SelectSkillCard) SpecialInit() {
	s.TaskMap = make(map[int]bool)
	s.TaskMap[s.Id1] = false
	s.TaskMap[s.Id2] = false
}

func (s *SelectSkillCard) SelectEnd() {
	//s.Mutex.Lock()
	//defer s.Mutex.Unlock() //先不用上锁，毕竟没有race操作
	s.SM.SendActionById(s.Id1, BattleDto.NewAction(BattleDto.DeployCard, BattleDto.Finish, "技能牌全部选择完毕"))
	s.SM.SendActionById(s.Id2, BattleDto.NewAction(BattleDto.DeployCard, BattleDto.Finish, "技能牌全部选择完毕"))
	go s.SM.finish("Judge")

}

func (s *SelectSkillCard) enter() {

	chanStop, chanCrash := Util.CreateTimer(time.Second*global.SelectSkillCardTime, s.SelectEnd)
	s.ChanCrash = chanCrash
	s.ChanStop = chanStop

	s.SM.SendActionById(s.Id1, BattleDto.NewAction(BattleDto.DeployCard, BattleDto.Query, map[string]any{"state_wait_time": Util.SendTime(time.Second * global.SelectSkillCardTime), "where": BattleData.SkillCard}))
	s.SM.SendActionById(s.Id2, BattleDto.NewAction(BattleDto.DeployCard, BattleDto.Query, map[string]any{"state_wait_time": Util.SendTime(time.Second * global.SelectSkillCardTime), "where": BattleData.SkillCard}))

}
func (s *SelectSkillCard) exit() {
	s.StateTemplate.exit()
	s.TaskMap[s.Id1] = false
	s.TaskMap[s.Id2] = false
	s.ChanCrash = nil
	s.ChanStop = nil
}
func (s *SelectSkillCard) process(GoCtx context.Context) {
	handleAction := func(id int, action BattleDto.Action, ResponseChan chan<- BattleDto.Action) bool {
		if s.TaskMap[id] {
			s.SM.SendActionById(s.Id1, BattleDto.NewErrAction(global.ResponseRepeatRequest))
			return true
		}

		if action.ActionCode == BattleDto.DeployCard && action.Predicates == BattleDto.Result {

			var data BattleData.SelectCard
			err := mapstructure.Decode(action.ActionData, &data)
			if err != nil {
				fmt.Println(err)
				s.SM.SendActionById(id, BattleDto.NewErrAction(global.ResponseInvalidReqParams))
				return true
			}
			if data.Where == BattleData.SkillCard {

				cardTempId := data.CardTempId
				if cardTempId == -1 {
					s.TaskMap[id] = true
					s.SM.SendActionById(id, BattleDto.NewAction(BattleDto.DeployCard, BattleDto.Succeed, "技能牌选择成功"))
					if s.TaskMap[s.Id1] && s.TaskMap[s.Id2] { //都上牌了
						s.ChanStop <- struct{}{}
					}
					return true
				}

				if card, ok := s.c.PlayerDataMap[id].CardInHand[cardTempId]; ok { //手牌里有不有
					if _, ok := card.(CardAbstract.SkillCard); ok { //上的是不是skillcard
						delete(s.c.PlayerDataMap[id].CardInHand, cardTempId)
						s.c.SetSkillCardBT(id, card)
						s.SM.SendActionById(id, BattleDto.NewAction(BattleDto.DeployCard, BattleDto.Succeed, "技能牌选择成功"))
						s.TaskMap[id] = true

						if s.TaskMap[s.Id1] && s.TaskMap[s.Id2] { //都上牌了

							s.ChanStop <- struct{}{}
						}
						return true
					} else {
						s.SM.SendActionById(id, BattleDto.NewErrAction(global.BattleCardCategoryError))
						return true
					}
				} else {
					s.SM.SendActionById(id, BattleDto.NewErrAction(global.BattleCardNotFound))
					return true
				}

			} else {
				s.SM.SendActionById(id, BattleDto.NewErrAction(global.BattleCardCategoryError))
			}
			return true
		}
		return false
	}
	s.SM.AcceptAction(GoCtx, handleAction)

}

//#endregion
//region State:Judge

type Judge struct {
	StateTemplate
	Mutex             sync.Mutex
	TaskMap           map[int]int
	ChanStop          chan struct{} //这东西是不用初始化的
	ChanCrash         chan struct{}
	IsTie             bool
	WaitAnimationPlay bool
}

type JudgeData struct {
	JudgeData int `json:"judge_data" mapstructure:"judge_data"`
}

func (J *Judge) SpecialInit() {
	J.TaskMap = make(map[int]int)
	J.IsTie = false
	J.WaitAnimationPlay = false
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

	if JudgeWin(J.TaskMap[J.Id1], J.TaskMap[J.Id2]) == 0 {
		J.IsTie = true
	} else {
		J.SM.Winner = J.Id1
		J.SM.Loser = J.Id2
		if JudgeWin(J.TaskMap[J.Id1], J.TaskMap[J.Id2]) == -1 {
			J.SM.Winner = J.Id2
			J.SM.Loser = J.Id1
		}
	}

	J.WaitAnimationPlay = true

}

func (J *Judge) enter() {
	J.TaskMap[J.Id1] = 3 //设为一个不可能值作为检查是否返回了
	J.TaskMap[J.Id2] = 3

	chanStop, chanCrash := Util.CreateTimer(time.Second*global.JudgeWaitTime, J.EndJudge)
	J.ChanCrash = chanCrash
	J.ChanStop = chanStop

	J.SM.SendActionById(J.Id1, BattleDto.NewAction(BattleDto.Judge, BattleDto.Query, NewStateWaitTime(global.JudgeWaitTime)))
	J.SM.SendActionById(J.Id2, BattleDto.NewAction(BattleDto.Judge, BattleDto.Query, NewStateWaitTime(global.JudgeWaitTime)))
}
func (J *Judge) exit() {
	J.TaskMap[J.Id1] = 3
	J.TaskMap[J.Id1] = 3
	J.ChanStop = nil
	J.ChanCrash = nil
	J.IsTie = false
	J.WaitAnimationPlay = false
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
	handleAction := func(id int, action BattleDto.Action, ResponseChan chan<- BattleDto.Action) bool {

		J.Mutex.Lock()

		if J.WaitAnimationPlay && action.ActionCode == BattleDto.AnimationPlayEnd && action.Predicates == BattleDto.Notify {
			if !J.IsTie {
				go J.SM.finish("Combat")
			} else {
				go J.SM.finish("Judge")
			}
			J.Mutex.Unlock()
			return true
		}
		J.Mutex.Unlock()

		if action.ActionCode == BattleDto.Judge && action.Predicates == BattleDto.Result {
			J.Mutex.Lock()
			var data JudgeData
			err := mapstructure.Decode(action.ActionData, &data)
			if err != nil {
				fmt.Println(err)
				J.SM.SendActionById(id, BattleDto.NewErrAction(global.ResponseInvalidReqParams))
				J.Mutex.Unlock()
				return true
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
			return true
		}
		return false

	}
	J.SM.AcceptAction(GoCtx, handleAction)
}

//endregion
//region State:Combat

type Combat struct {
	StateTemplate
	ChanCrash chan struct{}
	ChanStop  chan struct{}
	TimeStart time.Time
}

func (c *Combat) enter() {
	c.TimeStart = time.Now()
	NeedWait := c.SM.CombatTime
	c.SM.SendActionById(c.SM.Winner, BattleDto.NewAction(BattleDto.Combat, BattleDto.Query, NewStateWaitTime(NeedWait)))
	c.SM.SendActionById(c.SM.Loser, BattleDto.NewAction(BattleDto.Combat, BattleDto.Notify, NewStateWaitTime(NeedWait)))
	c.ChanStop, c.ChanCrash = Util.CreateTimer(NeedWait, c.CombatEnd)
}
func (c *Combat) CombatEnd() {
	println("CombatEnd 跳转到结算技能牌")
}

func (c *Combat) exit() {
	TimeEnd := time.Now()
	c.SM.CombatTime -= TimeEnd.Sub(c.TimeStart)
	if c.SM.CombatTime < 0 {
		c.SM.CombatTime = 0
	}
}
func (c *Combat) process(GoCtx context.Context) {
	handleAction := func(id int, action BattleDto.Action, ResponseChan chan<- BattleDto.Action) bool {
		if id == c.SM.Loser {
			c.SM.SendActionById(id, BattleDto.NewErrAction(global.BattleNotInYourRound))
			return true
		}

		if action.ActionCode == BattleDto.Combat && action.Predicates == BattleDto.Result {
			c.SM.Mutex.Lock()
			var data BattleData.CombatDto
			err := mapstructure.Decode(action.ActionData, &data)
			if err != nil {
				fmt.Println(err)
				c.SM.SendActionById(id, BattleDto.NewErrAction(global.ResponseInvalidReqParams))
				return true
			}
			if data.SelfWhere == BattleData.SkillCard {
				c.SM.SendActionById(id, BattleDto.NewErrAction(global.BattleCardCategoryError))
				return true
			}
			if !c.c.CheckCardByWhere(id, data.SelfWhere) {
				c.SM.SendActionById(id, BattleDto.NewErrAction(global.BattleCardNotFound))
				return true
			}
			c.SM.CombatDataChan <- data
			c.SM.Mutex.Unlock()
			go c.SM.finish("CardCalc")
			return true

		}
		return false
	}
	c.SM.AcceptAction(GoCtx, handleAction)
}

//endregion
//region State:SkillCardCalc

type SkillCardCalc struct {
	StateTemplate
}

func (s *SkillCardCalc) enter() {
	if s.c.PlayerDataMap[s.Id1].SkillCardBT != nil {
		s.c.PlayerDataMap[s.Id1].SkillCardBT.(CardAbstract.SkillCard).PlayMagic() //触发法术，然后，在法术这个函数里面，用和ctx的协议，把通知前端的action传出来
	}
	if s.c.PlayerDataMap[s.Id2].SkillCardBT != nil {
		s.c.PlayerDataMap[s.Id2].SkillCardBT.(CardAbstract.SkillCard).PlayMagic()
	}

	go s.SM.finish("SelectSkillCard")
}
func (s *SkillCardCalc) exit() { //这里其实主要可以初始化下一个循环的参数
	s.SM.CombatTime = global.CombatWaitTime * time.Second //这很重要，每次循环重新初始化一下战斗倒计时

}
func (s *SkillCardCalc) process(GoCtx context.Context) {}

//endregion

//region State:CardCalc

type CardCalc struct {
	StateTemplate
}

func (s *CardCalc) enter() {
	data := <-s.SM.CombatDataChan
	if data.Behavior == BattleData.Attack {
		s.c.GetCardBt(s.SM.Winner, data.SelfWhere).(CardAbstract.Character).Attack()
	} else if data.Behavior == BattleData.Skill {
		s.c.GetCardBt(s.SM.Winner, data.SelfWhere).(CardAbstract.Character).Skill()
	}

}
func (s *CardCalc) exit()                         {}
func (s *CardCalc) process(GoCtx context.Context) {}
func (s *CardCalc) SpecialInit()                  {}

//endregion
