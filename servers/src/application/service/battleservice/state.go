package battleservice

import (
	"context"
	"fmt"
	"math/rand"
	"pcc_card/Util"
	"pcc_card/application/entity/BattleData"
	"pcc_card/application/entity/Card/CardAbstract"
	"pcc_card/application/entity/protocol"
	"pcc_card/global"
	"pcc_card/presentation/handler/battlehandler/BattleDto"
	"sync"
	"sync/atomic"

	"time"

	"github.com/mitchellh/mapstructure"
)

type State interface {
	enter()
	exit()
	main()
	Init(id1 int, id2 int, c *Ctx, Nt *NotifyManager, SM *StateMachine, sub State)
	process(GoCtx context.Context)
	SetName(name string)
	GetName() string
	SpecialInit()
}

func (s *StateMachine) RegisterState() {
	s.StateList = map[string]State{
		"ShuffleDeal":     &ShuffleDeal{},
		"SelectSkillCard": &SelectSkillCard{},
		"Judge":           &Judge{},
		"Combat":          &Combat{},
		"CardCalc":        &CardCalc{},
		"SelectWeather":   &SelectWeather{},
		"ActiveChildCard": &ActiveChildCard{},
	}
	for key, element := range s.StateList {
		element.SetName(key)
	}
}

// 减少现在的天气的持续时间,如果刚好减到0,那么就改为宁静,
func (s *StateMachine) ReduceWeatherLasting() {
	res := s.WeatherLasting - 1
	if res == 0 {
		s.c.Weather.Store(int64(protocol.Ningjing))
	}
	if res < 0 {
		res = 0
	}

	s.WeatherLasting = res
}

func (s *StateMachine) SharedProcess(id int, action BattleDto.Action, ResponseChan chan<- BattleDto.Action) bool {
	if action.ActionCode == BattleDto.GetUserId && action.Predicates == BattleDto.Query {
		res := make(map[string]int)
		res["self"] = id
		res["opponent"] = s.c.GetOpponentId(id)
		ResponseChan <- BattleDto.NewAction(BattleDto.GetWeather, BattleDto.Result, res)
	}
	if action.ActionCode == BattleDto.GetRoundNum && action.Predicates == BattleDto.Query {
		ResponseChan <- BattleDto.NewAction(BattleDto.GetRoundNum, BattleDto.Result, s.RoundNum.Load())
	}

	if action.ActionCode == BattleDto.ReConnect && action.Predicates == BattleDto.Query {
		//CurStateName := s.CurrentState.GetName()
		//这个是软重连,如果请求这个,就先给他一个dto,一个字段是CurStateName,还有一个是,告诉他用户有不有响应过,(直接在statemachine里面,搞两个原子bool就好了,idaflag)
		//然后,再把对应state的query给他

		//还有硬重连的问题

	}

	if action.ActionCode == BattleDto.Ping && action.Predicates == BattleDto.Query {
		ResponseChan <- BattleDto.NewAction(BattleDto.Ping, BattleDto.Result, "pong")
		return true
	}

	if action.ActionCode == BattleDto.GetWeather && action.Predicates == BattleDto.Query {
		res := protocol.Weather(s.c.Weather.Load())
		ResponseChan <- BattleDto.NewAction(BattleDto.GetWeather, BattleDto.Result, res)
		fmt.Println("前端获取天气为:", res)
		return true
	}
	if action.ActionCode == BattleDto.GetEnergy && action.Predicates == BattleDto.Query {
		res := make(map[string]int)
		res["self"] = s.c.PlayerDataMap[id].GetEnergy()
		res["opponent"] = s.c.PlayerDataMap[s.c.GetOpponentId(id)].GetEnergy()
		ResponseChan <- BattleDto.NewAction(BattleDto.GetEnergy, BattleDto.Result, res)
		return true
	}
	if action.ActionCode == BattleDto.GetChildCardList && action.Predicates == BattleDto.Query {
		res := s.c.GetChildCardDto()
		ResponseChan <- BattleDto.NewAction(BattleDto.GetChildCardList, BattleDto.Result, res)

		for _, t := range res {
			fmt.Println("childState", t)
		}

		return true
	}
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
	if action.ActionCode == BattleDto.Debug {
		s.Mutex.Lock()
		if action.ActionData == "currentState" {
			s.SendActionById(id, BattleDto.NewAction(BattleDto.Debug, BattleDto.Result, s.CurrentState.GetName()))
			s.Mutex.Unlock()
			return true
		}
		s.Mutex.Unlock()
	}
	if action.ActionCode == BattleDto.GetDisCard && action.Predicates == BattleDto.Query { //查看弃牌堆
		res := s.c.GetDisCardDto()
		s.SendActionById(id, BattleDto.NewAction(BattleDto.GetDisCard, BattleDto.Result, res))
		return true
	}

	return false
}

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
	RoundNum     atomic.Int32

	GoldMoreUserId int

	//stateData
	Winner         int
	Loser          int
	CombatDataChan chan map[string][]BattleData.CombatDto

	LoseMarkMap map[int]int //UserId:num

	WeatherLasting int         //天气持续时间
	WinnerIsAction atomic.Bool //赢的人是否行动
}

func NewStateMachine(c *Ctx, id1 int, id2 int, Nt *NotifyManager, ParentNodeCtx context.Context, GoldMoreUserId int) *StateMachine {

	StateMachineImpl := &StateMachine{}
	c.StateMachine = StateMachineImpl
	StateMachineImpl.ParentNodeCtx = ParentNodeCtx
	StateMachineImpl.c = c //ctx的注入
	StateMachineImpl.Id1 = id1
	StateMachineImpl.Id2 = id2
	StateMachineImpl.Nt = Nt //Nt的注入

	StateMachineImpl.StateStack = make([]State, 0)
	StateMachineImpl.CombatDataChan = make(chan map[string][]BattleData.CombatDto, 1)
	StateMachineImpl.LoseMarkMap = make(map[int]int)

	StateMachineImpl.LoseMarkMap[StateMachineImpl.Id1] = 0
	StateMachineImpl.LoseMarkMap[StateMachineImpl.Id2] = 0

	StateMachineImpl.GoldMoreUserId = GoldMoreUserId
	StateMachineImpl.RoundNum.Store(int32(0))
	StateMachineImpl.RegisterState()
	StateMachineImpl.WeatherLasting = 0
	StateMachineImpl.WinnerIsAction.Store(false)
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

//#region StateMachine

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

// DataDecode 返回值是绑定是否成功,data是数据的指针
func (s *StateMachine) DataDecode(action BattleDto.Action, data any, id int) bool {
	config := &mapstructure.DecoderConfig{
		// 开启弱类型转换，处理 float64 -> int 以及 []interface{} -> []int
		WeaklyTypedInput: true,
		// 目标结构体指针
		Result: data,
		// 显式指定使用 mapstructure 标签
		TagName: "mapstructure",
	}

	// 2. 初始化解码器
	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		fmt.Printf("Fatal: failed to create decoder: %v\n", err)
		return false
	}

	// 3. 执行解码
	err = decoder.Decode(action.ActionData)
	if err != nil {
		fmt.Printf("Decode Error: %v | ActionData: %+v\n", err, action.ActionData)
		s.SendActionById(id, BattleDto.NewErrAction(global.ResponseInvalidReqParams))
		return false
	}
	return true
}

func (s *StateMachine) AcceptAction(goCtx context.Context, handleAction func(id int, action BattleDto.Action, ResponseChan chan<- BattleDto.Action) bool) {
	for {
		select {
		case <-goCtx.Done():

			return
		case action := <-s.Nt.ChanMap[s.Id1].AcceptChan:
			InterruptListenFunc := s.c.InterruptListenFunc.Load().(func(id int, action BattleDto.Action, ResponseChan chan<- BattleDto.Action) bool)

			if InterruptListenFunc(s.Id1, action, s.Nt.ChanMap[s.Id1].ResponseChan) {
				continue
			}
			if handleAction(s.Id1, action, s.Nt.ChanMap[s.Id1].ResponseChan) {
				continue
			}
			if s.SharedProcess(s.Id1, action, s.Nt.ChanMap[s.Id1].ResponseChan) {
				continue
			}
			fmt.Println("战斗时机不对的原始字段", action)
			s.SendActionById(s.Id1, BattleDto.NewErrAction(global.BattleInvalidTiming))
		case action := <-s.Nt.ChanMap[s.Id2].AcceptChan:
			InterruptListenFunc := s.c.InterruptListenFunc.Load().(func(id int, action BattleDto.Action, ResponseChan chan<- BattleDto.Action) bool)

			if InterruptListenFunc(s.Id2, action, s.Nt.ChanMap[s.Id2].ResponseChan) {
				continue
			}
			if handleAction(s.Id2, action, s.Nt.ChanMap[s.Id2].ResponseChan) {
				continue
			}
			if s.SharedProcess(s.Id2, action, s.Nt.ChanMap[s.Id2].ResponseChan) {
				continue
			}
			fmt.Println("战斗时机不对的原始字段", action)
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
		fmt.Print(s.CurrentState.GetName() + "\n")
		s.CurrentState.enter()

		var GoCtx context.Context
		GoCtx, s.cancelFunc = context.WithCancel(s.ParentNodeCtx) //stateMachine死掉，你也得死
		go s.CurrentState.process(GoCtx)                          //监听

		s.CurrentState.main()
	}
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

func (s *StateTemplate) main() {

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
		go s.SM.finish("ActiveChildCard")
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
			fmt.Println("初始上牌阶段", data)
			if data.Where != BattleData.SkillCard {
				cardTempId := data.CardTempId
				if s.c.CheckCardByWhere(id, data.Where) { //判定这个上牌的位置是不是有牌了
					s.SM.SendActionById(id, BattleDto.NewErrAction(global.BattleHasCard))
					return true
				}
				if card, ok := s.c.PlayerDataMap[id].CardInHand[cardTempId]; ok { //手牌里有不有
					if _, ok := card.(CardAbstract.SkillCard); !ok {
						s.c.SetCardBt(id, card)
						s.SM.SendActionById(id, BattleDto.NewAction(BattleDto.DeployCard, BattleDto.Succeed, "选择成功"))
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
				s.SM.SendActionById(id, BattleDto.NewErrAction(global.BattleInvalidTiming))
				return true
			}

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
	//if CharacterNumA <= 1 || CharacterNumB <= 1 {
	//	return false
	//}
	return true
}

func (s *ShuffleDeal) exit() {
	s.StateTemplate.exit()
}

//#endregion
//#region State:ActiveChildCard

type ActiveChildCard struct {
	sync.Mutex
	StateTemplate
	TaskMap   map[int][]int
	DoneMap   map[int]bool
	ChanStop  chan struct{}
	ChanCrash chan struct{}
	Completed atomic.Bool
}

type ActiveChildCardDto struct {
	TempIdList []int `json:"temp_id_list" mapstructure:"temp_id_list"`
}

func (a *ActiveChildCard) SpecialInit() {
	a.TaskMap = make(map[int][]int)
	a.TaskMap[a.Id1] = make([]int, 0)
	a.TaskMap[a.Id2] = make([]int, 0)
	a.DoneMap = map[int]bool{a.Id1: false, a.Id2: false}
	a.Completed.Store(false)
}

func (a *ActiveChildCard) enter() {
	queryMap := map[string]any{
		"state_wait_time": Util.SendTime(global.ActiveChildCardTime * time.Second),
		"child_list":      a.SM.c.GetChildCardDto(),
	}
	a.SM.SendActionById(a.Id1, BattleDto.NewAction(BattleDto.ActiveChildCard, BattleDto.Query, queryMap))
	a.SM.SendActionById(a.Id2, BattleDto.NewAction(BattleDto.ActiveChildCard, BattleDto.Query, queryMap))
	a.ChanStop, a.ChanCrash = Util.CreateTimer(global.ActiveChildCardTime*time.Second, a.SelectEnd)
}

func (a *ActiveChildCard) exit() {
	a.StateTemplate.exit()
	a.TaskMap[a.Id1] = make([]int, 0)
	a.TaskMap[a.Id2] = make([]int, 0)
	a.DoneMap = map[int]bool{a.Id1: false, a.Id2: false}

	a.Completed.Store(false)
}

func (a *ActiveChildCard) process(GoCtx context.Context) {
	handleAction := func(id int, action BattleDto.Action, ResponseChan chan<- BattleDto.Action) bool {
		if action.ActionCode == BattleDto.ActiveChildCard && action.Predicates == BattleDto.Result {
			a.Lock()
			defer a.Unlock()
			if a.Completed.Load() {
				return true
			}
			if a.DoneMap[id] {
				a.SM.SendActionById(id, BattleDto.NewErrAction(global.ResponseRepeatRequest))
				return true
			}

			var data ActiveChildCardDto
			if !a.SM.DataDecode(action, &data, id) {
				return true
			}

			uniqueIds := uniqueIntSlice(data.TempIdList)
			if len(uniqueIds) > 5 {
				a.SM.SendActionById(id, BattleDto.NewErrAction(global.ResponseInvalidReqParams))
				return true
			}

			validIds := a.validChildTempIds()
			if !Util.VerifyIncludes(validIds, uniqueIds) {
				a.SM.SendActionById(id, BattleDto.NewErrAction(global.ResponseInvalidReqParams))
				return true
			}

			a.TaskMap[id] = uniqueIds
			a.DoneMap[id] = true
			a.SM.SendActionById(id, BattleDto.NewAction(BattleDto.ActiveChildCard, BattleDto.Succeed, "选择已接收"))
			if a.DoneMap[a.Id1] && a.DoneMap[a.Id2] {
				a.ChanCrash <- struct{}{}
				go a.finishSelect()
			}

			return true
		}
		return false
	}
	a.SM.AcceptAction(GoCtx, handleAction)
}

func (a *ActiveChildCard) SelectEnd() {
	a.Lock()
	defer a.Unlock()
	if a.Completed.Load() {
		return
	}
	a.finishSelect()
}

func (a *ActiveChildCard) finishSelect() {
	if a.Completed.Load() {
		return
	}
	a.Completed.Store(true)

	selected := a.computeFinalSelection()

	a.activateChildCards(selected)

	result := map[string]any{"selected_temp_id_list": selected}
	a.SM.SendActionById(a.Id1, BattleDto.NewAction(BattleDto.ActiveChildCard, BattleDto.Finish, result))
	a.SM.SendActionById(a.Id2, BattleDto.NewAction(BattleDto.ActiveChildCard, BattleDto.Finish, result))

	go a.SM.finish("SelectWeather")
}

func (a *ActiveChildCard) computeFinalSelection() []int {
	selected1 := a.TaskMap[a.Id1]
	selected2 := a.TaskMap[a.Id2]
	intersection := intersectionIntSlice(selected1, selected2)
	if len(intersection) > 5 {
		intersection = intersection[:5]
	}
	if len(intersection) == 5 {
		return intersection
	}

	validIds := a.validChildTempIds()
	rest := excludeIntSlice(validIds, intersection)
	need := 5 - len(intersection)
	randoms := Util.GetRandomElements(rest, need) //这是需要随机的,补充用户选不够的情况

	return append(intersection, randoms...)
}

// 获得所有的TempIds
func (a *ActiveChildCard) validChildTempIds() []int {
	result := make([]int, 0)
	a.SM.c.ChildList.Do(func(data *[]CardAbstract.Card) {
		for _, card := range *data {
			result = append(result, card.GetTempId())
		}
	})
	return result
}

func (a *ActiveChildCard) activateChildCards(selected []int) {
	a.SM.c.ChildList.Do(func(data *[]CardAbstract.Card) {
		for _, card := range *data {
			if containsInt(selected, card.GetTempId()) {
				card.GetInfo()["ChildState"] = BattleData.Active
			}
		}
	})
}

// 返回值: 不包含重复元素的切片，且保持原有的相对顺序
func uniqueIntSlice(ids []int) []int {
	result := make([]int, 0, len(ids))
	seen := make(map[int]struct{}, len(ids))
	for _, id := range ids {
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		result = append(result, id)
	}
	return result
}

// intersectionIntSlice 获取两个整数切片的交集
// a1, a2: 两个待比较的切片
// 返回值: 同时存在于 a1 和 a2 中的元素组成的切片（已去重）
func intersectionIntSlice(a1, a2 []int) []int {
	set := make(map[int]struct{}, len(a1))
	for _, v := range a1 {
		set[v] = struct{}{}
	}
	result := make([]int, 0)
	for _, v := range a2 {
		if _, ok := set[v]; ok {
			result = append(result, v)
		}
	}
	return uniqueIntSlice(result)
}

// excludeIntSlice 从 list 中排除掉存在于 excludes 中的元素（差集操作）
// list: 原始数据切片
// excludes: 需要排除的黑名单切片
// 返回值: 存在于 list 但不在 excludes 中的元素
func excludeIntSlice(list, excludes []int) []int {
	excludeSet := make(map[int]struct{}, len(excludes))
	for _, v := range excludes {
		excludeSet[v] = struct{}{}
	}
	result := make([]int, 0, len(list))
	for _, v := range list {
		if _, ok := excludeSet[v]; !ok {
			result = append(result, v)
		}
	}
	return result
}

// containsInt 判断切片中是否包含某个特定的整数
// list: 查找范围
// target: 目标整数
// 返回值: 找到返回 true，否则返回 false
func containsInt(list []int, target int) bool {
	for _, v := range list {
		if v == target {
			return true
		}
	}
	return false
}

//#endregion
//#region State:SelectWeather

type SelectWeather struct {
	StateTemplate
	StopChan         chan struct{}
	CrashChan        chan struct{}
	WeatherList      []int
	WeatherListMutex sync.Mutex
}

type SelectWeatherDto struct {
	Weather protocol.Weather `json:"weather" mapstructure:"weather"`
}

// ChangeWeather !!天气变化主函数
func ChangeWeather(w protocol.Weather, SM *StateMachine) {
	SM.WeatherLasting = 6
	SM.c.Weather.Store(int64(w))

}

// ToSelect 随机出3个天气,加上宁静
func (s *SelectWeather) ToSelect() []int {

	num := protocol.WeatherCanSelectNum //几个里随机
	k := 3                              //要几个
	nums := rand.Perm(num)
	res := make([]int, k)
	for i := 0; i < k; i++ {
		res[i] = nums[i] + 1
	}
	res = append(res, 0)
	return res
}

func (s *SelectWeather) process(GoCtx context.Context) {
	handleAction := func(id int, action BattleDto.Action, ResponseChan chan<- BattleDto.Action) bool {
		if action.ActionCode == BattleDto.SelectWeather && action.Predicates == BattleDto.Result {
			var data SelectWeatherDto
			s.SM.DataDecode(action, &data, id)
			s.WeatherListMutex.Lock() //上锁处理WeatherList
			defer s.WeatherListMutex.Unlock()

			fmt.Println("天气选择前端结果", data)
			var res int
			if data.Weather == -1 { //如果传-1,随机
				res = Util.GetRandomElements[int](s.WeatherList, 1)[0]
			} else {
				Is_have := false
				//检查他是不是给定选择内的
				for _, Value := range s.WeatherList {
					if Value == int(data.Weather) {
						Is_have = true
					}
				}
				if !Is_have {
					s.SM.SendActionById(id, BattleDto.NewErrAction(global.ResponseDataNotFound))
					return false
				} else { //有的话,就把传的给res
					res = int(data.Weather)
				}
			}

			s.SM.SendActionById(id, BattleDto.NewAction(BattleDto.SelectWeather, BattleDto.Succeed, data))
			s.SM.SendActionById(s.c.GetOpponentId(id), BattleDto.NewAction(BattleDto.SelectWeather, BattleDto.Succeed, data))
			s.CrashChan <- struct{}{}
			ChangeWeather(protocol.Weather(res), s.SM) //改res

			go s.SM.finish("SelectSkillCard")
		}

		return false
	}
	s.SM.AcceptAction(GoCtx, handleAction)
}

func (s *SelectWeather) timeEnding() {

	s.WeatherListMutex.Lock() //上锁处理WeatherList
	res := Util.GetRandomElements[int](s.WeatherList, 1)[0]
	s.WeatherListMutex.Unlock()

	ChangeWeather(protocol.Weather(res), s.SM)
	s.SM.SendActionById(s.SM.Id1, BattleDto.NewAction(BattleDto.SelectWeather, BattleDto.Succeed, SelectWeatherDto{Weather: protocol.Weather(res)}))
	s.SM.SendActionById(s.SM.Id2, BattleDto.NewAction(BattleDto.SelectWeather, BattleDto.Succeed, SelectWeatherDto{Weather: protocol.Weather(res)}))
	go s.SM.finish("SelectSkillCard")
}

func (s *SelectWeather) enter() {

	//----通知谁的钱多----
	opponentId := s.c.GetOpponentId(s.SM.GoldMoreUserId)
	s.SM.SendActionById(s.SM.GoldMoreUserId, BattleDto.NewAction(
		BattleDto.SelectWeather,
		BattleDto.Notify,
		map[string]bool{"is_more": true},
	))
	s.SM.SendActionById(opponentId, BattleDto.NewAction(
		BattleDto.SelectWeather,
		BattleDto.Notify,
		map[string]bool{"is_more": false},
	))
	//----通知谁的钱多----

	//-----让钱多的人选天气-----
	queryMap := make(map[string]any)
	queryMap["state_wait_time"] = Util.SendTime(global.SelectWeatherTime * time.Second)

	s.WeatherListMutex.Lock() //上锁处理WeatherList
	s.WeatherList = s.ToSelect()
	queryMap["weather_list"] = s.WeatherList
	s.WeatherListMutex.Unlock()

	s.SM.SendActionById(s.SM.GoldMoreUserId, BattleDto.NewAction(BattleDto.SelectWeather, BattleDto.Query, queryMap))
	//-----让钱多的人选天气-----

	//设置定时
	StopChan, CrashChan := Util.CreateTimer(global.SelectWeatherTime*time.Second, s.timeEnding)
	s.StopChan = StopChan
	s.CrashChan = CrashChan
}

func (s *SelectWeather) exit() {

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
	s.Mutex.Lock()
	defer s.Mutex.Unlock() //先不用上锁，毕竟没有race操作
	s.SM.SendActionById(s.Id1, BattleDto.NewAction(BattleDto.DeployCard, BattleDto.Finish, "技能牌全部选择完毕"))
	s.SM.SendActionById(s.Id2, BattleDto.NewAction(BattleDto.DeployCard, BattleDto.Finish, "技能牌全部选择完毕"))
	go s.SM.finish("Judge")
}

func (s *SelectSkillCard) enter() {

	s.SM.RoundNum.Store(s.SM.RoundNum.Add(1))
	chanStop, chanCrash := Util.CreateTimer(time.Second*global.SelectSkillCardTime, s.SelectEnd)
	s.ChanCrash = chanCrash
	s.ChanStop = chanStop

	s.SM.SendActionById(s.Id1, BattleDto.NewAction(BattleDto.DeployCard, BattleDto.Query, map[string]any{"state_wait_time": Util.SendTime(time.Second * global.SelectSkillCardTime), "where": BattleData.SkillCard}))
	s.SM.SendActionById(s.Id2, BattleDto.NewAction(BattleDto.DeployCard, BattleDto.Query, map[string]any{"state_wait_time": Util.SendTime(time.Second * global.SelectSkillCardTime), "where": BattleData.SkillCard}))

}
func (s *SelectSkillCard) exit() {
	s.SM.Mutex.Lock()
	defer s.SM.Mutex.Unlock()
	s.StateTemplate.exit()
	s.TaskMap[s.Id1] = false
	s.TaskMap[s.Id2] = false
	s.ChanCrash = nil
	s.ChanStop = nil
}
func (s *SelectSkillCard) process(GoCtx context.Context) {
	handleAction := func(id int, action BattleDto.Action, ResponseChan chan<- BattleDto.Action) bool {
		s.SM.Mutex.Lock()
		if s.TaskMap[id] {
			s.SM.SendActionById(s.Id1, BattleDto.NewErrAction(global.ResponseRepeatRequest))
			s.SM.Mutex.Unlock()
			return true
		}
		if action.ActionCode == BattleDto.DeployCard && action.Predicates == BattleDto.Result {

			var data BattleData.SelectCard
			err := mapstructure.Decode(action.ActionData, &data)
			if err != nil {
				fmt.Println(err)
				s.SM.SendActionById(id, BattleDto.NewErrAction(global.ResponseInvalidReqParams))
				s.SM.Mutex.Unlock()
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
					s.SM.Mutex.Unlock()
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
						s.SM.Mutex.Unlock()
						return true
					} else {
						s.SM.SendActionById(id, BattleDto.NewErrAction(global.BattleCardCategoryError))
						s.SM.Mutex.Unlock()
						return true
					}
				} else {
					s.SM.SendActionById(id, BattleDto.NewErrAction(global.BattleCardNotFound))
					s.SM.Mutex.Unlock()
					return true
				}

			} else {
				s.SM.SendActionById(id, BattleDto.NewErrAction(global.BattleCardCategoryError))
				s.SM.Mutex.Unlock()
				return true
			}

		}
		s.SM.Mutex.Unlock()
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
	IsTie             atomic.Bool
	WaitAnimationPlay atomic.Bool
}

type JudgeData struct {
	JudgeData int `json:"judge_data" mapstructure:"judge_data"`
}

func (J *Judge) SpecialInit() {
	J.TaskMap = make(map[int]int)
	J.IsTie.Store(false)
	J.WaitAnimationPlay.Store(false)
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
		J.IsTie.Store(true)
	} else {
		J.SM.Winner = J.Id1
		J.SM.Loser = J.Id2
		if JudgeWin(J.TaskMap[J.Id1], J.TaskMap[J.Id2]) == -1 {
			J.SM.Winner = J.Id2
			J.SM.Loser = J.Id1
		}

		J.SM.LoseMarkMap[J.SM.Loser] += 1
		J.SM.LoseMarkMap[J.SM.Winner] = 0

		J.c.PlayerDataMap[J.SM.Winner].UpdateEnergy(2)

	}

	J.WaitAnimationPlay.Store(true)

}

func (J *Judge) enter() {
	J.Mutex.Lock()
	defer J.Mutex.Unlock()
	J.TaskMap[J.Id1] = 3 //设为一个不可能值作为检查是否返回了
	J.TaskMap[J.Id2] = 3

	chanStop, chanCrash := Util.CreateTimer(time.Second*global.JudgeWaitTime, J.EndJudge)
	J.ChanCrash = chanCrash
	J.ChanStop = chanStop

	J.SM.SendActionById(J.Id1, BattleDto.NewAction(BattleDto.Judge, BattleDto.Query, NewStateWaitTime(global.JudgeWaitTime)))
	J.SM.SendActionById(J.Id2, BattleDto.NewAction(BattleDto.Judge, BattleDto.Query, NewStateWaitTime(global.JudgeWaitTime)))
}
func (J *Judge) exit() {
	J.Mutex.Lock()
	J.TaskMap[J.Id1] = 3
	J.TaskMap[J.Id1] = 3
	J.Mutex.Unlock()

	J.IsTie.Store(false)

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

		if J.WaitAnimationPlay.Load() && action.ActionCode == BattleDto.AnimationPlayEnd && action.Predicates == BattleDto.Notify {
			J.WaitAnimationPlay.Store(false)
			if !J.IsTie.Load() {
				go J.SM.finish("Combat")
			} else {
				go J.SM.finish("Judge")
			}
			return true
		}

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
				J.Mutex.Unlock()
				return true
			}
			if J.TaskMap[id] != 3 {
				J.SM.SendActionById(id, BattleDto.NewErrAction(global.ResponseRepeatRequest))
				J.Mutex.Unlock()
				return true
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
	CombatMap map[string][]BattleData.CombatDto
	WaitNum   int
}

func (c *Combat) SpecialInit() {
	c.CombatMap = make(map[string][]BattleData.CombatDto)

	c.WaitNum = 0
}

func (c *Combat) enter() {

	WaitTime := global.CombatWaitTime * time.Second

	c.SM.SendActionById(c.SM.Winner, BattleDto.NewAction(BattleDto.Combat, BattleDto.Query, NewStateWaitTime(WaitTime)))
	c.SM.SendActionById(c.SM.Loser, BattleDto.NewAction(BattleDto.Combat, BattleDto.Notify, NewStateWaitTime(WaitTime)))
	c.ChanStop, c.ChanCrash = Util.CreateTimer(WaitTime, c.CombatEnd)
}
func (c *Combat) CombatEnd() {

	c.SM.Mutex.Lock()
	defer c.SM.Mutex.Unlock()

	if _, ok := c.CombatMap["Winner"]; !ok {
		c.CombatMap["Winner"] = make([]BattleData.CombatDto, 0)
	}
	if _, ok := c.CombatMap["Loser"]; !ok {
		c.CombatMap["Loser"] = make([]BattleData.CombatDto, 0)
	}
	c.SM.CombatDataChan <- c.CombatMap
}

func (c *Combat) exit() {
	c.CombatMap = make(map[string][]BattleData.CombatDto)
	c.WaitNum = 0
}
func (c *Combat) process(GoCtx context.Context) {
	handleAction := func(id int, action BattleDto.Action, ResponseChan chan<- BattleDto.Action) bool {

		if action.ActionCode == BattleDto.Combat && action.Predicates == BattleDto.Result { //传递结果
			c.SM.Mutex.Lock()
			defer c.SM.Mutex.Unlock()
			var data []BattleData.CombatDto
			err := mapstructure.Decode(action.ActionData, &data)
			if err != nil {
				fmt.Println(err)
				c.SM.SendActionById(id, BattleDto.NewErrAction(global.ResponseInvalidReqParams))
				return true
			}
			fmt.Println("combat阶段前端发送的data:", data)
			//-----------把结果存入map,并且防止反复提交------------
			if id == c.SM.Winner {
				if _, ok := c.CombatMap["Winner"]; ok {
					c.SM.SendActionById(id, BattleDto.NewErrAction(global.ResponseRepeatRequest))
					return true
				}
				c.CombatMap["Winner"] = data
			} else {
				if _, ok := c.CombatMap["Loser"]; ok {
					c.SM.SendActionById(id, BattleDto.NewErrAction(global.ResponseRepeatRequest))
					return true
				}
				c.CombatMap["Loser"] = data
			}
			//-----------把结果存入map,并且防止反复提交------------

			c.WaitNum += 1
			c.SM.SendActionById(id, BattleDto.NewAction(BattleDto.Combat, BattleDto.Succeed, ""))
			//如果到2了,就转阶段了
			if c.WaitNum == 2 {
				c.SM.CombatDataChan <- c.CombatMap
				c.ChanCrash <- struct{}{}
				go c.SM.finish("CardCalc")
				return true
			}
			return true
		}

		return false
	}
	c.SM.AcceptAction(GoCtx, handleAction)
}

//endregion
//region State:CardCalc

type CardCalc struct {
	StateTemplate
	HaveDone atomic.Bool
}

func (s *CardCalc) SpecialInit() {
	s.HaveDone.Store(false)
}

func (s *CardCalc) enter() {
	s.SM.SendActionById(s.SM.Id2, BattleDto.NewAction(BattleDto.CardCalc, BattleDto.Notify, ""))
	s.SM.SendActionById(s.SM.Id1, BattleDto.NewAction(BattleDto.CardCalc, BattleDto.Notify, ""))
}

func (s *CardCalc) CalcBtCry() { //光环的效果
	WinParCard := s.c.PlayerDataMap[s.SM.Winner].GetBt(BattleData.ParentCard)
	LoserParCard := s.c.PlayerDataMap[s.SM.Loser].GetBt(BattleData.ParentCard)
	WinChiCard := s.c.PlayerDataMap[s.SM.Winner].GetBt(BattleData.ChildCard)
	LoserChiCard := s.c.PlayerDataMap[s.SM.Loser].GetBt(BattleData.ChildCard)
	Extc := func(ExtcCard CardAbstract.Card) {
		if ExtcCard != nil {
			ExtcCard.(CardAbstract.Character).BtCry()
		}
	}
	Extc(WinParCard) //按顺序执行四个多战吼
	Extc(LoserParCard)
	Extc(WinChiCard)
	Extc(LoserChiCard)
}

func (s *CardCalc) Switch(_data BattleData.CombatDto, UserId int) bool {
	data := _data.SelectCard

	if data.Where != BattleData.SkillCard {
		cardTempId := data.CardTempId
		if card, ok := s.c.PlayerDataMap[UserId].CardInHand[cardTempId]; ok { //手牌里有不有
			if _, ok := card.(CardAbstract.SkillCard); !ok { //不是技能牌

				playerData := s.c.PlayerDataMap[UserId]

				if _data.SelectCard.Where == BattleData.InHand { //检查是不是最后一张战斗牌
					if !playerData.CheckIs2Bt() {
						s.SM.SendActionById(UserId, BattleDto.NewErrAction(global.ResponseForbidden))
						return false
					}
				}
				if !playerData.IsCanUpdateEnergy(-1) { //减少能量
					s.SM.SendActionById(UserId, BattleDto.NewErrAction(global.BattleEnergyNotEnough))
					return false
				}
				playerData.UpdateEnergy(-1)

				//-----------------------换牌正文-----------------------

				playerData.SwitchCard(data.Where, card)
				return true
				//-----------------------换牌正文-----------------------

			} else {
				s.SM.SendActionById(UserId, BattleDto.NewErrAction(global.BattleCardCategoryError))
				return false
			}
		} else {
			s.SM.SendActionById(UserId, BattleDto.NewErrAction(global.BattleCardNotFound))
			return false
		}
	}
	return false
}

// 返回值是,是否使用技能或者攻击
func (s *CardCalc) CalcNotSwitch(data BattleData.CombatDto, UserId int) bool {
	opponentCardId := s.c.GetCardBt(UserId, data.OpponentWhere).GetTempId()
	res := false
	if data.Behavior == BattleData.Attack { //执行前端传过来的行为
		res = s.c.GetCardBt(UserId, data.SelfWhere).(CardAbstract.Character).Attack(opponentCardId)
	} else if data.Behavior == BattleData.Skill {
		res = s.c.GetCardBt(UserId, data.SelfWhere).(CardAbstract.Character).Skill(opponentCardId)
	}

	s.c.StackSettle() //执行效果堆栈
	return res
}

// SkillCalc 双方的法术牌结算都在这了//有执行
func (s *CardCalc) SkillCalc() {
	if s.c.PlayerDataMap[s.SM.Winner].SkillCardBT != nil {
		s.c.PlayerDataMap[s.SM.Winner].SkillCardBT.(CardAbstract.SkillCard).PlayMagic() //触发法术，然后，在法术这个函数里面，用和ctx的协议，把通知前端的action传出来
		s.c.StackSettle()
		s.c.ChildCardCheck()
	}

	if s.c.PlayerDataMap[s.SM.Loser].SkillCardBT != nil {
		s.c.PlayerDataMap[s.SM.Loser].SkillCardBT.(CardAbstract.SkillCard).PlayMagic()
		s.c.StackSettle() //执行效果堆栈
		s.c.ChildCardCheck()
	}
}

func (s *CardCalc) main() {
CalcLoop:
	for {

		select {
		case data := <-s.SM.CombatDataChan:
			fmt.Println("卡牌结算字典:", data)

			//-------------给连续输的人增加免伤-------------
			LoserBtCardList := s.c.GetBtAll(s.SM.Loser)
			if s.SM.LoseMarkMap[s.SM.Loser] >= 2 { //输两次了,有免伤
				for _, card := range LoserBtCardList {
					for range s.SM.LoseMarkMap[s.SM.Loser] - 1 {
						TempId := card.GetTempId()
						s.c.ProtoColPush(protocol.NewGiveBuff(&TempId, *protocol.NewBuffBase(protocol.DamageImmunity, 1, 0.28, s.c.CreateTempId())))
					}
				}
			}
			s.c.StackSettle()
			//-------------给连续输的人增加免伤-------------

			//回合开始结算天气//通知在里面了
			if protocol.WeatherExecPositionMap[protocol.Weather(s.c.Weather.Load())] == protocol.RoundStart {
				s.ExecWeather()
				s.c.StackSettle()
				s.c.ChildCardCheck()
			}

			//--------结算换牌start--------

			//-----记录换牌过程-----
			SwitchCardMap := make(map[int][]BattleData.SelectCard)
			SwitchCardMap[s.SM.Id2] = make([]BattleData.SelectCard, 0)
			SwitchCardMap[s.SM.Id1] = make([]BattleData.SelectCard, 0)
			//-----记录换牌过程-----

			//----执行换牌----
			SwitchCard := func(User string, UserId int, SwitchCardMap map[int][]BattleData.SelectCard) {
				DtoList := data[User]
				for _, Dto := range DtoList {
					if Dto.Behavior != BattleData.SwitchCard {
						break
					}
					if s.Switch(Dto, UserId) { //这里执行交换
						SwitchCardMap[UserId] = append(SwitchCardMap[UserId], Dto.SelectCard) //然后交换成功的记录
					}
				}
			}
			SwitchCard("Winner", s.SM.Winner, SwitchCardMap)
			s.c.ChildCardCheck()
			SwitchCard("Loser", s.SM.Loser, SwitchCardMap)
			s.c.ChildCardCheck()
			//----执行换牌----

			//通知换牌结算完成
			s.SM.SendActionById(s.SM.Id2, BattleDto.NewAction(BattleDto.DeployCard, BattleDto.Notify, transformToSelfAndOpponent[BattleData.SelectCard, []BattleData.SelectCard](SwitchCardMap, s.SM.Id2)))
			s.SM.SendActionById(s.SM.Id1, BattleDto.NewAction(BattleDto.DeployCard, BattleDto.Notify, transformToSelfAndOpponent[BattleData.SelectCard, []BattleData.SelectCard](SwitchCardMap, s.SM.Id1)))

			//-----------结算攻击或者技能-----------------

			//输出值是是否使用过攻击和技能
			Calc := func(User string, UserId int) bool {
				DtoList := data[User]
				res := false
				for _, Dto := range DtoList {
					if Dto.Behavior == BattleData.SwitchCard {
						s.SM.SendActionById(UserId, BattleDto.NewErrAction(global.BattleCantSwitch))
						continue
					}
					if s.CalcNotSwitch(Dto, UserId) {
						res = true
					}
				}
				return res
			}
			WinnerIsAction := Calc("Winner", s.SM.Winner) //获取赢的人是否使用过攻击和技能
			s.SM.WinnerIsAction.Store(WinnerIsAction)
			s.c.ChildCardCheck()
			Calc("Loser", s.SM.Loser)
			s.c.ChildCardCheck()
			//-----------结算攻击或者技能-----------------

			//结算法术(s.c.ChildCardCheck()在里面了)

			s.SM.SendActionById(s.SM.Id2, BattleDto.NewAction(BattleDto.SkillCardNotify, BattleDto.Notify, ""))
			s.SM.SendActionById(s.SM.Id1, BattleDto.NewAction(BattleDto.SkillCardNotify, BattleDto.Notify, ""))
			s.SkillCalc()
			s.c.ChildCardCheck()

			//回合结束结算天气//通知在里面了
			if protocol.WeatherExecPositionMap[protocol.Weather(s.c.Weather.Load())] == protocol.RoundEnd {
				s.ExecWeather()
				s.c.StackSettle()
				s.c.ChildCardCheck()
			}

			//结算buff
			s.SM.SendActionById(s.SM.Id2, BattleDto.NewAction(BattleDto.BuffCalcNotify, BattleDto.Notify, ""))
			s.SM.SendActionById(s.SM.Id1, BattleDto.NewAction(BattleDto.BuffCalcNotify, BattleDto.Notify, ""))
			CharacterCardList := s.c.GetCharacter()
			for _, Card := range CharacterCardList {
				Card.BuffRoundEnd(s.c)
				s.c.StackSettle()
			}
			s.c.ChildCardCheck()

			//全部结算完成,发个通知
			s.SM.SendActionById(s.SM.Id2, BattleDto.NewAction(BattleDto.CardCalc, BattleDto.Finish, ""))
			s.SM.SendActionById(s.SM.Id1, BattleDto.NewAction(BattleDto.CardCalc, BattleDto.Finish, ""))
			s.HaveDone.Store(true)
		default:
			break CalcLoop
		}
	}
}

// 结算天气//通知也在这个函数里面了的
func (s *CardCalc) ExecWeather() {
	s.SM.SendActionById(s.SM.Id2, BattleDto.NewAction(BattleDto.WeatherNotify, BattleDto.Notify, ""))
	s.SM.SendActionById(s.SM.Id1, BattleDto.NewAction(BattleDto.WeatherNotify, BattleDto.Notify, ""))
	cardBts := s.c.GetBtAll(-1)
	buffNeeds := make([]protocol.BuffNeed, len(cardBts)) //转化数组类型,数组没法隐式转化的
	for i, card := range cardBts {
		buffNeeds[i] = card
	}
	protocol.WeatherFuncMap[protocol.Weather(s.c.Weather.Load())](s.c, buffNeeds) //从天气执行函数map储存种取出来执行

}

func (s *CardCalc) exit() {
	s.HaveDone.Store(false)
	s.SM.ReduceWeatherLasting()
	s.SM.WinnerIsAction.Store(false)
}

func (s *CardCalc) process(GoCtx context.Context) {
	fmt.Println("进入cardcal的process了")
	handleAction := func(id int, action BattleDto.Action, ResponseChan chan<- BattleDto.Action) bool {
		if action.ActionCode == BattleDto.AnimationPlayEnd && action.Predicates == BattleDto.Notify && s.HaveDone.Load() {
			go s.SM.finish("SelectSkillCard")
			return true
		}
		return false
	}
	s.SM.AcceptAction(GoCtx, handleAction)
}

//endregion

// T是数据类型,转化传输的是[]T这种类型
func transformToSelfAndOpponent[T any, L []T](DataMap map[int]L, UserId int) map[string]L {
	res := make(map[string]L)
	res["self"] = make(L, 0)
	res["opponent"] = make(L, 0)
	for id, data := range DataMap {
		if id == UserId {
			res["self"] = data
		} else {
			res["opponent"] = data
		}
	}
	return res
}
