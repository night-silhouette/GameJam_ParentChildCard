package battleservice

import (
	"context"
	"fmt"
	"pcc_card/Util"
	"pcc_card/application/entity/BattleData"
	"pcc_card/application/entity/Card/CardAbstract"
	"pcc_card/application/entity/CardMeta"
	"pcc_card/application/entity/protocol"
	"pcc_card/global"
	"pcc_card/presentation/handler/battlehandler/BattleDto"
	"sync"
	"sync/atomic"
	"time"
)

type Ctx struct {
	ParentContext context.Context

	StateMachine   *StateMachine
	entityCounter  int //自增字段，用来生成tempId
	EffectsStack   *EffectStack
	CardObserver   *CardObserver
	PlayerDataMap  map[int]*PlayerData
	CtxStateNotify *CtxStateNotify

	DisCardPool *Util.SafeContainer[CardAbstract.Card]
	TempIdCalc  *atomic.Int32 //计算tempid//现在记录的值,是用过的值

	Weather             atomic.Int64
	ChildList           *Util.SafeContainer[CardAbstract.Card]
	NeedInterrupt       atomic.Bool
	InterruptChan       chan struct{}
	InterruptListenFunc atomic.Value //func(id int, action BattleDto.Action, ResponseChan chan<- BattleDto.Action) bool

}

// cList是子牌堆,CardList是手牌堆
func InitCtx(c *Ctx, idA int, idB int, ParentContext context.Context, CardList map[int]map[int]CardAbstract.Card, TempIdCalc *atomic.Int32, cList []CardAbstract.Card) *Ctx {
	c.EffectsStack = NewEffectStack()
	c.entityCounter = 1
	c.ParentContext = ParentContext

	c.PlayerDataMap = make(map[int]*PlayerData, 2)
	c.PlayerDataMap[idA] = NewPlayerData(idA, CardList[idA])

	c.PlayerDataMap[idB] = NewPlayerData(idB, CardList[idB])
	c.CtxStateNotify = NewCtxStateNotify()
	//c.CardObserver = NewCardObserver(ParentContext, c)//弃用哨兵模式
	c.NeedInterrupt.Store(false)
	c.InterruptListenFunc.Store(func(id int, action BattleDto.Action, ResponseChan chan<- BattleDto.Action) bool { return false })
	c.DisCardPool = Util.NewSafeContainer[CardAbstract.Card](ParentContext, 8)
	c.InterruptChan = make(chan struct{}, 1)
	c.TempIdCalc = TempIdCalc
	c.ChildList = Util.NewSafeContainer[CardAbstract.Card](ParentContext, 8)
	c.ChildList.Do(func(data *[]CardAbstract.Card) {
		*data = cList
		for _, card := range *data {
			card.SetTempId(int(TempIdCalc.Add(1)))
		}
	})
	return c
}

// 把手牌里的和野生的标定数组里被激活的加起来
func (c *Ctx) GetNeedCheckChildCard() []CardAbstract.ChildCard {
	res := make([]CardAbstract.ChildCard, 0)
	c.ChildList.Do(func(data *[]CardAbstract.Card) {
		for _, card := range *data {
			if card.GetInfo()["ChildState"] == BattleData.Active {
				childCard := card.(CardAbstract.ChildCard)
				res = append(res, childCard)
			}
		}

	})
	for _, playerData := range c.PlayerDataMap {
		for _, Card := range playerData.CardInHand {
			if Card.GetInfo()["is_parent"] == false {
				childCard := Card.(CardAbstract.ChildCard)
				res = append(res, childCard)
			}
		}
	}
	return res
}

// 聚合了,卡到手牌的过程,和通知前端
func (c *Ctx) ChildCatch(card CardAbstract.ChildCard, UserId int) {
	OpponentId := c.GetOpponentId(UserId)
	playerData := c.PlayerDataMap[UserId]
	OpPlayerData := c.PlayerDataMap[OpponentId]
	playerData.dataMutex.Lock()
	defer playerData.dataMutex.Unlock()
	OpPlayerData.dataMutex.Lock()
	defer OpPlayerData.dataMutex.Unlock()

	Notify := func(Origin int) {
		DtoSelf := BattleData.ChildCardCatchDto{ //通知
			Origin:  Origin,
			Object:  UserId,
			DataAll: c.GetDataAll(UserId),
		}
		DtoOp := BattleData.ChildCardCatchDto{ //通知
			Origin:  Origin,
			Object:  UserId,
			DataAll: c.GetDataAll(OpponentId),
		}
		c.StateMachine.SendActionById(UserId, BattleDto.NewAction(BattleDto.ChildBelongChange, BattleDto.Result, DtoSelf))
		c.StateMachine.SendActionById(OpponentId, BattleDto.NewAction(BattleDto.ChildBelongChange, BattleDto.Result, DtoOp))
	}
	if card.GetInfo()["ChildState"] == BattleData.Active {
		playerData.CardInHand[card.GetTempId()] = card     //底层数据改牌
		card.GetInfo()["ChildState"] = BattleData.HasCatch //改状态
		Notify(-1)
		return

	} else if card.GetInfo()["ChildState"] == BattleData.HasCatch {
		delete(OpPlayerData.CardInHand, card.GetTempId())
		playerData.CardInHand[card.GetTempId()] = card
		Notify(OpponentId)
		return
	}

}

// 循环需要检查的数组,然后触发
func (c *Ctx) ChildCardCheck() {
	NeedCheckChildCard := c.GetNeedCheckChildCard()
	for _, card := range NeedCheckChildCard {
		cFunc := card.Check(c)
		flag, UserId := cFunc.Exec(c)
		if flag {
			c.ChildCatch(card, UserId)
			card.Trigger(c, UserId)
			c.StackSettle()
		}
	}
}

//__________________________________________EffectsStack______________________________________________

//todo

func (c *Ctx) StackSettle() { //执行函数
	c.resolveAllChains()
}

func (c *Ctx) resolveAllChains() {
	for {
		if c.NeedInterrupt.Load() {
			<-c.InterruptChan
			fmt.Println("中断停止了")
		}

		if c.EffectsStack.IsEmpty() {
			// 出口

			break

		}

		effect := c.EffectsStack.Pop()
		effect.Execute(c)
		fmt.Print(effect)
		fmt.Println("执行了")

	}
}

//__________________________________________CardObserver______________________________________________

type CardObserver struct {
	ctx           *Ctx
	ParentContext context.Context
	Collector     chan MetaCardState
}

func (O *CardObserver) DrainCollector() {
	for {
		select {
		case e := <-O.Collector:
			// 只要管道里有东西，就一直 Push 到栈里
			O.ctx.EffectsStack.Push(e.Effect)
		default:
			return
		}
	}
}

type MetaCardState struct {
	CardId int
	Effect protocol.Effect
}

//__________________________________________CtxStateNotify______________________________________________

// MetaCardBTChange 用于在CtxStateNotify传输状态的元数据
type MetaCardBTChange struct {
	Old CardAbstract.Card
	New CardAbstract.Card
}

// CtxStateNotify 内嵌到ctx里面，监听数据变化
type CtxStateNotify struct {
	ParentCardChange  chan MetaCardBTChange
	ChildCardBTChange chan MetaCardBTChange
	SkillCardBTChange chan MetaCardBTChange
}

func NewCtxStateNotify() *CtxStateNotify {
	N := &CtxStateNotify{}
	N.ParentCardChange = make(chan MetaCardBTChange, 4)
	N.ChildCardBTChange = make(chan MetaCardBTChange, 4)
	N.SkillCardBTChange = make(chan MetaCardBTChange, 4)
	return N

}

func (c *Ctx) SetParentCardBT(id int, new CardAbstract.Card) {
	pData := c.PlayerDataMap[id]
	old := pData.ParentCardBT
	pData.ParentCardBT = new
	c.CtxStateNotify.ParentCardChange <- MetaCardBTChange{old, new}
}
func (c *Ctx) SetChildCardBT(id int, new CardAbstract.Card) {
	pData := c.PlayerDataMap[id]
	old := pData.ChildCardBT
	pData.ChildCardBT = new
	c.CtxStateNotify.ParentCardChange <- MetaCardBTChange{old, new}
}
func (c *Ctx) SetSkillCardBT(id int, new CardAbstract.Card) {
	pData := c.PlayerDataMap[id]
	old := pData.SkillCardBT
	pData.SkillCardBT = new
	c.CtxStateNotify.ParentCardChange <- MetaCardBTChange{old, new}
}

//__________________________________________PlayerData______________________________________________

type PlayerData struct {
	ID           int
	CardInHand   map[int]CardAbstract.Card
	ParentCardBT CardAbstract.Card
	ChildCardBT  CardAbstract.Card
	SkillCardBT  CardAbstract.Card
	Energy       int

	dataMutex sync.RWMutex
}

func (p *PlayerData) GetEnergy() int {
	p.dataMutex.RLock()
	defer p.dataMutex.RUnlock()
	return p.Energy
}

// 是不是有两张出战卡,返回bool
func (p *PlayerData) CheckIs2Bt() bool {

	p.dataMutex.Lock()
	defer p.dataMutex.Unlock()
	flag := true

	if p.ParentCardBT == nil {
		flag = false
	}
	if p.ChildCardBT == nil {
		flag = false
	}
	return flag
}
func (p *PlayerData) UpdateEnergy(offset int) {
	p.dataMutex.Lock()
	defer p.dataMutex.Unlock()
	if p.isCanUpdateEnergy(offset) {
		temp := p.Energy + offset
		if temp > 5 {
			temp = 5
		}
		p.Energy = temp
	}
}

func (p *PlayerData) isCanUpdateEnergy(offset int) bool {
	if (p.Energy + offset) < 0 {
		return false
	}
	return true
}

func (p *PlayerData) IsCanUpdateEnergy(offset int) bool {
	p.dataMutex.RLock()
	defer p.dataMutex.RUnlock()
	return p.isCanUpdateEnergy(offset)
}

// 这样就和字典一样,可以直接取要的位置的卡,
func (p *PlayerData) GetBt(where BattleData.Where) CardAbstract.Card {
	p.dataMutex.Lock()
	defer p.dataMutex.Unlock()

	switch where {
	case BattleData.SkillCard:
		return p.SkillCardBT
	case BattleData.ChildCard:
		return p.ChildCardBT
	case BattleData.ParentCard:
		return p.ParentCardBT
	default:
		return nil
	}
}

// SwitchCard 没有做skillCard的鉴定//鉴定了原来的那个牌是不是空的,空的才会把他换下来//还有,纯下牌的情况也写进去了
func (p *PlayerData) SwitchCard(where BattleData.Where, card CardAbstract.Card) {
	p.dataMutex.Lock()
	defer p.dataMutex.Unlock()
	var SwitchedCard CardAbstract.Card
	if where == BattleData.ChildCard {
		SwitchedCard = p.ChildCardBT
		p.ChildCardBT = card
	}
	if where == BattleData.ParentCard {
		SwitchedCard = p.ParentCardBT
		p.ParentCardBT = card
	}
	if where == BattleData.InHand { //下牌
		SwitchedCard = card
	}
	if SwitchedCard != nil {
		p.CardInHand[SwitchedCard.GetTempId()] = SwitchedCard
	}
}

func NewPlayerData(ID int, CardInHand map[int]CardAbstract.Card) *PlayerData {
	p := &PlayerData{}
	p.CardInHand = CardInHand
	p.ID = ID
	p.Energy = 1 //能量初始值为一
	return p
}

//__________________________________________对外访问数据口(dto)______________________________________________

func (c *Ctx) GetCardInHard(id_self int) *BattleData.CardInHand {
	var id_opponent int
	for key := range c.PlayerDataMap {
		if key != id_self {
			id_opponent = key
			break
		}
	}
	res := &BattleData.CardInHand{}

	mapSelf := c.PlayerDataMap[id_self].CardInHand
	mapOpponent := c.PlayerDataMap[id_opponent].CardInHand

	res.Self = make([]BattleData.CardDto, 0, len(mapSelf))
	res.Opponent = make([]BattleData.CardDto, 0, len(mapOpponent))
	for _, card := range mapSelf {
		res.Self = append(res.Self, CardAbstract.GetCardDto(card))
	}
	for _, card := range mapOpponent {
		res.Opponent = append(res.Opponent, CardAbstract.GetCardDto(card))
	}
	return res
}

func (c *Ctx) GetBtCardInfo(id int) BattleData.BtCardInfo {
	GetDtoDefault := func(card CardAbstract.Card) BattleData.CardDto {
		if card == nil {
			res := BattleData.CardDto{}
			res.Id = -1
			res.BuffDtoList = make([]BattleData.BuffDto, 0)
			return res
		}
		return CardAbstract.GetCardDto(card)
	}

	var BtCardInfo BattleData.BtCardInfo
	for key, value := range c.PlayerDataMap {
		if key == id {
			BtCardInfo.Self.ChildCardBt = GetDtoDefault(value.ChildCardBT)
			BtCardInfo.Self.SkillCardBt = GetDtoDefault(value.SkillCardBT)
			BtCardInfo.Self.ParentCardBt = GetDtoDefault(value.ParentCardBT)
		} else {
			BtCardInfo.Opponent.ChildCardBt = GetDtoDefault(value.ChildCardBT)
			BtCardInfo.Opponent.SkillCardBt = GetDtoDefault(value.SkillCardBT)
			BtCardInfo.Opponent.ParentCardBt = GetDtoDefault(value.ParentCardBT)
		}
	}

	return BtCardInfo
}

// todo
// region protocol
func (c *Ctx) CreateTempId() int {
	return int(c.TempIdCalc.Add(1))
}

func (c *Ctx) ProtoColPush(e protocol.Effect) {
	c.EffectsStack.Push(e)
}

func (c *Ctx) ProtoColMoveDisCardPool(UserId int, TempId int) {
	c.MoveDisCardPool(UserId, TempId)
}

func (c *Ctx) ProtoColSetCardBt(UserId int, TempId int) {
	card := c.GetCardInHardByCardTempId(UserId, TempId)
	c.SetCardBt(UserId, card) //这个接口已经上锁了
}

// Notify 传-1，表全部,这是动画dto的通知
func (c *Ctx) Notify(AnimationDto BattleData.AnimationDto, UserId int) {
	if UserId == -1 {
		c.StateMachine.SendActionById(c.StateMachine.Id2, BattleDto.NewAction(BattleDto.AnimationNotify, BattleDto.Notify, AnimationDto))
		c.StateMachine.SendActionById(c.StateMachine.Id1, BattleDto.NewAction(BattleDto.AnimationNotify, BattleDto.Notify, AnimationDto))
		return
	} else {
		c.StateMachine.SendActionById(UserId, BattleDto.NewAction(BattleDto.AnimationNotify, BattleDto.Notify, AnimationDto))
		return
	}

}

// 不走卡牌hurt虚函数的攻击(走NoSourceHurt),也可以叫做无主攻击
func (c *Ctx) ProtoColAttackNoHurt(CardTempId int, Value int, Category BattleData.ValueChange) {
	Card := c.FindCard(CardTempId)
	Card.(CardAbstract.Character).NoSourceHurt(float64(Value), Category)
}
func (c *Ctx) GetWeather() protocol.Weather {
	return protocol.Weather(c.Weather.Load())
}

func (c *Ctx) ProtoColInterrupt(UserId int, InterruptDto *BattleData.InterruptDto, res chan []int, InterruptWaitTime time.Duration) {
	c.NeedInterrupt.Store(true)
	c.StateMachine.SendActionById(UserId, BattleDto.NewAction(BattleDto.Interrupt, BattleDto.Query, InterruptDto))
	c.StateMachine.SendActionById(c.GetOpponentId(UserId), BattleDto.NewAction(BattleDto.Interrupt, BattleDto.Notify, InterruptDto))
	var data BattleData.InterruptSelect
	var dataMutex sync.Mutex
	var DataIsOK atomic.Bool
	DataIsOK.Store(false)

	TimeEnding := func() { //结束回调
		if !DataIsOK.Load() { //随机取
			dataMutex.Lock()
			data.TempIdList = Util.GetRandomElements(InterruptDto.TempIdList, InterruptDto.SelectNum)
			dataMutex.Unlock()
		}
		c.StateMachine.SendActionById(UserId, BattleDto.NewAction(BattleDto.Interrupt, BattleDto.Succeed, map[string][]int{"temp_id_list": data.TempIdList}))
		res <- data.TempIdList
		c.NeedInterrupt.Store(false)
		c.InterruptListenFunc.Store(func(id int, action BattleDto.Action, ResponseChan chan<- BattleDto.Action) bool { return false })
	}
	StopChan, _ := Util.CreateTimer(InterruptWaitTime, TimeEnding) //定时
	c.InterruptListenFunc.Store(func(id int, action BattleDto.Action, ResponseChan chan<- BattleDto.Action) bool {

		if action.ActionCode == BattleDto.Interrupt && action.Predicates == BattleDto.Result && id == UserId {
			dataMutex.Lock()
			defer dataMutex.Unlock()
			fmt.Println(action.ActionData)
			if !c.StateMachine.DataDecode(action, &data, id) {
				return true
			}

			if len(data.TempIdList) != InterruptDto.SelectNum {
				c.StateMachine.SendActionById(id, BattleDto.NewErrAction(global.BattleCardNumErr))
				return true
			}
			if !Util.VerifyIncludes(InterruptDto.TempIdList, data.TempIdList) {
				c.StateMachine.SendActionById(id, BattleDto.NewErrAction(global.BattleCardNotFound))
				return true
			}
			DataIsOK.Store(true)
			StopChan <- struct{}{}
			return true
		}
		return false
	})
}

func (c *Ctx) ProtoColUpdateEnergy(UserId int, offset int) {
	playerData := c.PlayerDataMap[UserId]
	playerData.UpdateEnergy(offset)
}

func (c *Ctx) ProtoColCardBtAttack(SendTempId int, UserId int, TargetTempId int, AtkHp float64, Category BattleData.ValueChange) {
	var card CardAbstract.Character
	var ok bool

	if card, ok = c.FindCard(TargetTempId).(CardAbstract.Character); !ok {
		fmt.Println("转化失败bug")
		return
	}
	card.Hurt(SendTempId, AtkHp, Category)
}

func (c *Ctx) ProtoColCancelInterrupt() {
	c.InterruptChan <- struct{}{}
}

func (c *Ctx) GetTempIdByWhere(where BattleData.Where, userId int) int {
	playerdata := c.PlayerDataMap[userId]
	var res int
	switch where {
	case BattleData.ParentCard:
		if playerdata.ParentCardBT != nil {
			res = playerdata.ParentCardBT.GetTempId()
		}
	case BattleData.ChildCard:
		if playerdata.ChildCardBT != nil {
			res = playerdata.ChildCardBT.GetTempId()
		}
	case BattleData.SkillCard:
		if playerdata.SkillCardBT != nil {
			res = playerdata.SkillCardBT.GetTempId()
		}
	}
	return res
}
func (c *Ctx) ProtoColSetMaxHp(TargetTempId int, MaxHp float64) {
	Card := c.FindCard(TargetTempId)
	Card.GetInfo()["maxHp"] = MaxHp
}

// 扣血的底层方法
func (c *Ctx) ProtoColReduceCardBtHp(SendTempId int, TargetTempId int, ReduceHp float64) { //最后的底层方法
	var card CardAbstract.Character
	var ok bool
	if card, ok = c.FindCard(TargetTempId).(CardAbstract.Character); !ok {
		fmt.Println("转化失败bug")
		return
	}

	NowHp := card.GetHpNow()
	if NowHp-ReduceHp <= 0 {
		card.SetHpNow(0)
		card.Death(SendTempId)
		return
	}
	card.SetHpNow(NowHp - ReduceHp)

}
func (c *Ctx) ProtoColHealCardBt(TargetTempId int, HealHp float64) {
	var card CardAbstract.Character
	var ok bool
	if card, ok = c.FindCard(TargetTempId).(CardAbstract.Character); !ok {
		fmt.Println("转化失败bug")
		return
	}
	NowHp := card.GetHpNow()
	MaxHp := card.GetInfo()["maxHp"].(float64)
	if NowHp >= MaxHp {
		return
	}

	if NowHp+HealHp > MaxHp {
		card.SetHpNow(MaxHp)
		return
	}
	card.SetHpNow(NowHp + HealHp)
}
func (c *Ctx) Broad(v *CardMeta.BroadInfo) {
	AllCardList := c.GetAllCard()
	for _, card := range AllCardList {
		card.PutBroadInfo(v)
	}
}

func (c *Ctx) ProtoColGetCharacterCard(UserId int) []int {
	return c.GetCharacterCardinCardInHand(UserId)
}

func (c *Ctx) ProtoColCanUpdateEnergy(UserId int, offset int) bool {
	playData := c.PlayerDataMap[UserId]
	return playData.isCanUpdateEnergy(offset)
}

func (c *Ctx) ProtoSendAction(UserId int, action BattleDto.Action) {
	c.StateMachine.SendActionById(UserId, action)
}

func (c *Ctx) ProtoColSetDamageCardBt(UserId int, TargetTempId int, NewDamage float64) {
	var card CardAbstract.Character
	var ok bool
	if card, ok = c.FindCard(TargetTempId).(CardAbstract.Character); !ok {
		fmt.Println("转化失败bug")
		return
	}
	if NewDamage < 0 {
		NewDamage = 0
	}
	card.SetAtkNow(NewDamage)
}
func (c *Ctx) GetCharacterId() []int {
	list := c.GetCharacter()
	result := make([]int, len(list))
	for _, e := range list {
		result = append(result, e.GetTempId())
	}
	return result
}

func (c *Ctx) ProtoNotifyValue(Category BattleData.ValueChange, Value float64, TempId int, IsMiss bool) {
	c.StateMachine.SendActionById(c.StateMachine.Id1, BattleDto.NewAction(BattleDto.HpChange, BattleDto.Result, BattleData.CardCalcValueDto{
		TempId:   TempId,
		Category: Category,
		IsMiss:   IsMiss,
		Value:    Value,
		DataAll:  c.GetDataAll(c.StateMachine.Id1),
	}))
	c.StateMachine.SendActionById(c.StateMachine.Id2, BattleDto.NewAction(BattleDto.HpChange, BattleDto.Result, BattleData.CardCalcValueDto{
		TempId:   TempId,
		Category: Category,
		IsMiss:   IsMiss,
		Value:    Value,
		DataAll:  c.GetDataAll(c.StateMachine.Id2),
	}))
}

func (c *Ctx) ProtoNotifyCardMove(Object BattleData.Where, TempId int) {
	c.StateMachine.SendActionById(c.StateMachine.Id1, BattleDto.NewAction(BattleDto.PositionChange, BattleDto.Result, BattleData.CardCalcCardMoveDto{
		Object:  Object,
		TempId:  TempId,
		DataAll: c.GetDataAll(c.StateMachine.Id1),
	}))
	c.StateMachine.SendActionById(c.StateMachine.Id2, BattleDto.NewAction(BattleDto.PositionChange, BattleDto.Result, BattleData.CardCalcCardMoveDto{
		Object:  Object,
		TempId:  TempId,
		DataAll: c.GetDataAll(c.StateMachine.Id2),
	}))
}

//endregion

//todo
//region 对卡牌数据的操作算法

func (c *Ctx) GiveBuff(TempId int, buff *protocol.Buff) {
	card := c.FindCard(TempId)
	card.AddBuff(buff, c)
}

func (c *Ctx) CheckBuff(tempId int, buffId protocol.BuffId) (bool, *protocol.Buff) {
	card := c.FindCard(tempId)
	for _, b := range *card.GetBuffList() {
		if b.BuffId == buffId {
			return true, b
		}
	}
	return false, nil
}

// FindCard 这是一个总和方法,他会在两个人的手牌和出战斗牌里根据tempId找牌
func (c *Ctx) FindCard(tempId int) CardAbstract.Card {
	var card CardAbstract.Card
	res := c.GetCardInCardBtByCardTempId(c.StateMachine.Id1, tempId)
	if res != nil {
		card = res
	}
	res = c.GetCardInCardBtByCardTempId(c.StateMachine.Id2, tempId)
	if res != nil {
		card = res
	}
	res = c.GetCardInHardByCardTempId(c.StateMachine.Id1, tempId)
	if res != nil {
		card = res
	}
	res = c.GetCardInHardByCardTempId(c.StateMachine.Id2, tempId)
	if res != nil {
		card = res
	}
	return card
}

func (c *Ctx) FindCardForUserId(UserId int, tempId int) CardAbstract.Card {
	var card CardAbstract.Card
	res := c.GetCardInCardBtByCardTempId(UserId, tempId)
	if res != nil {
		card = res
	}
	res = c.GetCardInHardByCardTempId(UserId, tempId)
	if res != nil {
		card = res
	}
	return card
}

func (c *Ctx) GetOpponentId(selfId int) int {
	res := c.StateMachine.Id1
	if res == selfId {
		res = c.StateMachine.Id2
	}
	return res
}

func (c *Ctx) GetCardInCardBtByCardTempId(UserId int, TargetTempId int) CardAbstract.Card {
	player := c.PlayerDataMap[UserId]
	if player.ParentCardBT != nil {
		if player.ParentCardBT.GetTempId() == TargetTempId {
			return player.ParentCardBT
		}
	}
	if player.ChildCardBT != nil {
		if player.ChildCardBT.GetTempId() == TargetTempId {
			return player.ChildCardBT
		}
	}
	if player.SkillCardBT != nil {
		if player.SkillCardBT.GetTempId() == TargetTempId {
			return player.SkillCardBT
		}
	}

	return nil
}

func (c *Ctx) GetCardInHardByCardTempId(UserId int, CardTempId int) CardAbstract.Card {
	player := c.PlayerDataMap[UserId]
	for _, card := range player.CardInHand {
		if card.GetTempId() == CardTempId {
			return card
		}
	}
	return nil
}

func (c *Ctx) RandomSelectCard(id int) BattleData.Where { //这个是个不安全方法，就是你要保证他是一定还有牌可以上的（就是，没死掉的）
	playerData := c.PlayerDataMap[id]
	cardInHard := playerData.CardInHand
	for cardId, card := range cardInHard {
		if _, ok := card.(CardAbstract.Character); !ok {
			continue
		}
		fmt.Println(card.GetInfo())
		if card.GetInfo()["is_parent"].(bool) {
			c.SetParentCardBT(id, card)
			delete(cardInHard, cardId)
			return BattleData.ParentCard
		} else {
			c.SetChildCardBT(id, card)
			delete(cardInHard, cardId)
			return BattleData.ChildCard
		}
	}
	return -1
}

// CheckCardByWhere 判断这个地方是否有牌
func (c *Ctx) CheckCardByWhere(id int, where BattleData.Where) bool {
	playerData := c.PlayerDataMap[id]
	if where == BattleData.SkillCard && playerData.SkillCardBT == nil {
		return false
	}
	if where == BattleData.ParentCard && playerData.ParentCardBT == nil {
		return false
	}
	if where == BattleData.ChildCard && playerData.ChildCardBT == nil {
		return false
	}
	return true
}

// CheckCard 检查是否有角色牌出战
func (c *Ctx) CheckCard(id int) bool {
	playerData := c.PlayerDataMap[id]
	playerData.dataMutex.Lock()
	defer playerData.dataMutex.Unlock()
	flag := false //没有牌
	if playerData.ParentCardBT != nil {
		flag = true
	}
	if playerData.ChildCardBT != nil {
		flag = true
	}
	return flag
}

// SetCardBt 设置cardBT,已考虑子母Bt问题,从手牌里删除,判断了是否有牌，没牌才可以上
func (c *Ctx) SetCardBt(id int, card CardAbstract.Card) {
	playerData := c.PlayerDataMap[id]
	playerData.dataMutex.Lock()
	defer playerData.dataMutex.Unlock()
	if _, ok := card.(CardAbstract.SkillCard); ok {
		c.SetSkillCardBT(id, card)
		delete(playerData.CardInHand, card.GetTempId())
		return
	}
	if card.GetInfo()["is_parent"].(bool) && !c.CheckCardByWhere(id, BattleData.ParentCard) {
		c.SetParentCardBT(id, card)
		delete(playerData.CardInHand, card.GetTempId())
		return
	}
	if !card.GetInfo()["is_parent"].(bool) && !c.CheckCardByWhere(id, BattleData.ChildCard) {
		c.SetChildCardBT(id, card)
		delete(playerData.CardInHand, card.GetTempId())
		return
	}

}

// GetCardBt 根据where获取卡牌对象本体
func (c *Ctx) GetCardBt(id int, where BattleData.Where) CardAbstract.Card {
	playerData := c.PlayerDataMap[id]
	playerData.dataMutex.Lock()
	defer playerData.dataMutex.Unlock()
	if where == BattleData.SkillCard {
		return playerData.SkillCardBT
	}
	if where == BattleData.ParentCard {
		return playerData.ParentCardBT
	}
	if where == BattleData.ChildCard {
		return playerData.ChildCardBT
	}
	return nil
}

func (c *Ctx) GetDisCardDto() []BattleData.CardDto {
	result := make([]BattleData.CardDto, 0, 20)
	c.DisCardPool.Do(func(data *[]CardAbstract.Card) {
		for _, card := range *data {
			result = append(result, CardAbstract.GetCardDto(card))
		}
	})
	return result
}

// GetCharacterCardinCardInHand 获取在手牌中的所有角色卡
func (c *Ctx) GetCharacterCardinCardInHand(UserId int) []int {
	player := c.PlayerDataMap[UserId]
	TempIdList := make([]int, 0, 10)
	for _, card := range player.CardInHand {
		if characterCard, ok := card.(CardAbstract.Character); ok {
			TempIdList = append(TempIdList, characterCard.GetTempId())
		}
	}
	return TempIdList
}

// GetWhereByTempId 根据tempId 返回where，找不到返回-1
func (c *Ctx) GetWhereByTempId(UserId int, tempId int) BattleData.Where {
	player := c.PlayerDataMap[UserId]
	if player.ParentCardBT.GetTempId() == tempId {
		return BattleData.ParentCard
	}
	if player.ChildCardBT.GetTempId() == tempId {
		return BattleData.ChildCard
	}
	if player.SkillCardBT.GetTempId() == tempId {
		return BattleData.SkillCard
	}
	return -1
}

// DeleteCardBtByWhere 根据where删除bt
func (c *Ctx) DeleteCardBtByWhere(UserId int, where BattleData.Where) {
	player := c.PlayerDataMap[UserId]
	if where == BattleData.SkillCard {
		player.ParentCardBT = nil
		return
	}
	if where == BattleData.ParentCard {
		player.ParentCardBT = nil
		return
	}
	if where == BattleData.ChildCard {
		player.ChildCardBT = nil
		return
	}
}

// MoveCardBtToDisCardPool 把bt上的卡删掉，并且移动到discardpool
func (c *Ctx) MoveCardBtToDisCardPool(UserId int, tempId int) {
	var where BattleData.Where
	if where = c.GetWhereByTempId(UserId, tempId); where == -1 {
		return
	}
	card := c.GetCardInCardBtByCardTempId(UserId, tempId)
	c.DeleteCardBtByWhere(UserId, where)
	c.DisCardPool.Push(card)
}

func (c *Ctx) MoveDisCardPool(UserId int, tempId int) {
	player := c.PlayerDataMap[UserId]
	var where BattleData.Where
	if where = c.GetWhereByTempId(UserId, tempId); where != -1 {
		card := c.GetCardInCardBtByCardTempId(UserId, tempId)
		c.DeleteCardBtByWhere(UserId, where)
		card.SetOwnerId(-1)
		card.ReInitialize()
		c.DisCardPool.Push(card)
		return
	} else {
		card := c.GetCardInCardBtByCardTempId(UserId, tempId)
		delete(player.CardInHand, tempId)
		card.SetOwnerId(-1)
		card.ReInitialize()
		c.DisCardPool.Push(card)
		return
	}

}

func (c *Ctx) GetChildCardDto() []*BattleData.ChildCardDto {
	result := make([]*BattleData.ChildCardDto, 0, 10)

	c.ChildList.Do(func(data *[]CardAbstract.Card) {

		for _, card := range *data {
			el := BattleData.NewChildCardDto(CardAbstract.GetCardDto(card), card.GetInfo()["ChildState"].(BattleData.ChildState))

			result = append(result, el)
		}
	})
	return result
}

func (c *Ctx) GetDataAll(UseId int) *BattleData.DataAll {
	EnergyDto := BattleData.EnergyDto{
		Self:     c.PlayerDataMap[UseId].GetEnergy(),
		Opponent: c.PlayerDataMap[c.GetOpponentId(UseId)].GetEnergy(),
	}
	BtInfo := c.GetBtCardInfo(UseId)
	res := BattleData.NewDataAll(
		&BtInfo,
		c.GetCardInHard(UseId),
		&EnergyDto,
		c.GetDisCardDto(),
		c.GetChildCardDto(),
	)
	return res
}

// 获取全部的卡
func (c *Ctx) GetAllCard() []CardAbstract.Card {
	res := make([]CardAbstract.Card, 0)
	appendPlayer := func(p *PlayerData) {
		if p == nil {
			return
		}
		if p.ParentCardBT != nil {
			res = append(res, p.ParentCardBT)
		}
		if p.ChildCardBT != nil {
			res = append(res, p.ChildCardBT)
		}
		if p.SkillCardBT != nil {
			res = append(res, p.SkillCardBT)
		}
	}
	appendCardInHand := func(UserId int) {
		player := c.PlayerDataMap[UserId]
		for _, card := range player.CardInHand {
			res = append(res, card)
		}
	}
	appendCardInHand(c.StateMachine.Id1)
	appendCardInHand(c.StateMachine.Id2)
	appendPlayer(c.PlayerDataMap[c.StateMachine.Id1])
	appendPlayer(c.PlayerDataMap[c.StateMachine.Id2])
	return res
}

// 获取全部的角色卡
func (c *Ctx) GetCharacter() []CardAbstract.Card {

	res := make([]CardAbstract.Card, 0)

	appendPlayer := func(p *PlayerData) {
		if p == nil {
			return
		}
		if p.ParentCardBT != nil {
			res = append(res, p.ParentCardBT)
		}
		if p.ChildCardBT != nil {
			res = append(res, p.ChildCardBT)
		}
	}
	appendCardInHand := func(UserId int) {
		player := c.PlayerDataMap[UserId]
		for _, card := range player.CardInHand {
			if characterCard, ok := card.(CardAbstract.Character); ok {
				res = append(res, characterCard)
			}
		}
	}

	appendCardInHand(c.StateMachine.Id1)
	appendCardInHand(c.StateMachine.Id2)
	appendPlayer(c.PlayerDataMap[c.StateMachine.Id1])
	appendPlayer(c.PlayerDataMap[c.StateMachine.Id2])
	return res
}

// GetBtAll 获取场上所有的出战卡,如果输入-1,就是所有双方的,如果是有id,就是单人的
func (c *Ctx) GetBtAll(UserId int) []CardAbstract.Card {
	res := make([]CardAbstract.Card, 0)
	DiffUserId := func(id int) {
		player := c.PlayerDataMap[id]
		player.dataMutex.Lock()
		if player.ParentCardBT != nil {
			res = append(res, player.ParentCardBT)
		}
		if player.ChildCardBT != nil {
			res = append(res, player.ChildCardBT)
		}
		player.dataMutex.Unlock()
	}
	if UserId == -1 {
		DiffUserId(c.StateMachine.Id1)
		DiffUserId(c.StateMachine.Id2)
	} else {
		DiffUserId(UserId)
	}
	return res
}
func (c *Ctx) GetIds() []int {
	res := make([]int, 0)
	res = append(res, c.StateMachine.Id1)
	res = append(res, c.StateMachine.Id2)
	return res
}

func (c *Ctx) GetWinnerIsAction() bool {
	return c.StateMachine.WinnerIsAction.Load()
}
func (c *Ctx) GetWinnerId() int {
	return c.StateMachine.Winner
}

//endregion
