package battleservice

import (
	"context"
	"fmt"
	"log"
	"pcc_card/Util"
	"pcc_card/application/entity/BattleData"
	"pcc_card/application/entity/Card/CardAbstract"
	"pcc_card/application/entity/protocol"
	"pcc_card/global"
	"pcc_card/presentation/handler/battlehandler/BattleDto"
	"runtime"
)

type Ctx struct {
	ParentContext context.Context

	StateMachine   *StateMachine
	entityCounter  int //自增字段，用来生成tempId
	EffectsStack   *EffectStack
	CardObserver   *CardObserver
	PlayerDataMap  map[int]*PlayerData
	CtxStateNotify *CtxStateNotify
	CardPool       *[]CardAbstract.Card //和手牌数组里的是同一份对象，都是从总体复制出来的
	DisCardPool    *Util.SafeContainer[CardAbstract.Card]

	NeedInterrupt       bool
	InterruptChan       chan struct{}
	InterruptListenFunc func(id int, action BattleDto.Action, ResponseChan chan<- BattleDto.Action) bool
}

func NewCtx(idA int, idB int, CardPool *[]CardAbstract.Card, ParentContext context.Context) *Ctx {
	c := &Ctx{}
	c.EffectsStack = NewEffectStack()
	c.entityCounter = 1
	c.ParentContext = ParentContext
	c.CardPool = CardPool
	c.PlayerDataMap = make(map[int]*PlayerData, 2)
	c.PlayerDataMap[idA] = NewPlayerData(idA)
	c.PlayerDataMap[idB] = NewPlayerData(idB)
	c.CtxStateNotify = NewCtxStateNotify()
	c.CardObserver = NewCardObserver(ParentContext, c)
	c.NeedInterrupt = false
	c.InterruptListenFunc = func(id int, action BattleDto.Action, ResponseChan chan<- BattleDto.Action) bool { return false }
	c.DisCardPool = Util.NewSafeContainer[CardAbstract.Card](ParentContext, 8)
	return c
}

//__________________________________________EffectsStack______________________________________________

func (c *Ctx) StackSettle(action BattleDto.Action) { //执行函数

	select {
	case firstEffect := <-c.CardObserver.Collector:
		c.EffectsStack.Push(firstEffect.Effect)
		c.resolveAllChains(action)
	}
}

func (c *Ctx) resolveAllChains(action BattleDto.Action) {
	emptyCount := 0
	for {
		if c.NeedInterrupt {
			<-c.InterruptChan
		}

		c.CardObserver.DrainCollector()
		if c.EffectsStack.IsEmpty() {
			// 如果栈空了，不急着退出，先尝试让出执行权
			runtime.Gosched()
			c.CardObserver.DrainCollector()
			if c.EffectsStack.IsEmpty() {
				emptyCount++
				if emptyCount < 3 { // 连续三次检查都是空的，才真的算结束
					continue
				}
				// 真正的出口
				c.StateMachine.SendActionById(c.StateMachine.Id2, action)
				c.StateMachine.SendActionById(c.StateMachine.Id1, action)
				break
			}
		}
		emptyCount = 0 // 只要处理了效果，计数器清零

		effect := c.EffectsStack.Pop()
		effect.Execute(c)
		fmt.Println("执行")
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

func NewCardObserver(ParentContext context.Context, ctx *Ctx) *CardObserver {
	o := &CardObserver{}
	o.ParentContext = ParentContext
	o.Collector = make(chan MetaCardState, 128)
	o.ctx = ctx
	CardPool := *o.ctx.CardPool

	for _, card := range CardPool {

		go func(card CardAbstract.Card) { //给每一张卡开一个哨兵
			var CardChan <-chan protocol.Effect
			CardChan = card.GetStateCodeChan()
			for {
				select {
				case code := <-CardChan:
					Meta := MetaCardState{card.GetID(), code}
					select {
					case o.Collector <- Meta:
					default:

						ctx.StateMachine.SendActionById(ctx.StateMachine.Id1, BattleDto.NewErrAction(global.BattleEffectStackOverflow))
						ctx.StateMachine.SendActionById(ctx.StateMachine.Id2, BattleDto.NewErrAction(global.BattleEffectStackOverflow))
						log.Println("collector channel full")
					}
				case <-o.ParentContext.Done():
					return
				}
			}
		}(card)

	}
	return o
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
}

func NewPlayerData(ID int) *PlayerData {
	p := &PlayerData{}
	p.CardInHand = make(map[int]CardAbstract.Card)
	p.ID = ID
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
			res.BuffId = -1
			res.Id = -1
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

//todo
//region protocol

func (c *Ctx) ProtoColSetCardBt(UserId int, TempId int) {
	card := c.GetCardInHardByCardTempId(UserId, TempId)
	c.SetCardBt(UserId, card)
}

func (c *Ctx) Notify(AnimationDto BattleData.AnimationDto) {
	c.StateMachine.SendActionById(c.StateMachine.Id2, BattleDto.NewAction(BattleDto.CardCalc, BattleDto.Notify, AnimationDto))
	c.StateMachine.SendActionById(c.StateMachine.Id1, BattleDto.NewAction(BattleDto.CardCalc, BattleDto.Notify, AnimationDto))
}

func (c *Ctx) ProtoColInterrupt(UserId int, InterruptDto *BattleData.InterruptDto, res chan []int) {
	c.NeedInterrupt = true
	c.StateMachine.SendActionById(UserId, BattleDto.NewAction(BattleDto.Interrupt, BattleDto.Query, InterruptDto))
	c.StateMachine.SendActionById(c.GetOpponentId(UserId), BattleDto.NewAction(BattleDto.Interrupt, BattleDto.Notify, InterruptDto))
	c.InterruptListenFunc = func(id int, action BattleDto.Action, ResponseChan chan<- BattleDto.Action) bool {
		if action.ActionCode == BattleDto.Interrupt && action.Predicates == BattleDto.Result && id == UserId {
			var data BattleData.InterruptSelect
			if !c.StateMachine.DataDecode(action, data, id) {
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
			res <- data.TempIdList
			c.NeedInterrupt = false
			c.InterruptListenFunc = func(id int, action BattleDto.Action, ResponseChan chan<- BattleDto.Action) bool { return false }
			c.InterruptChan <- struct{}{}
			return true
		}
		return false
	}
}

func (c *Ctx) ProtoColCardBtAttack(SendTempId int, UserId int, TargetTempId int, AtkHp float64) {
	var card CardAbstract.Character
	var ok bool

	if card, ok = c.FindCard(TargetTempId).(CardAbstract.Character); !ok {
		fmt.Println("转化失败bug")
		return
	}
	card.Hurt(SendTempId, AtkHp)
}

func (c *Ctx) ProtoColReduceCardBtHp(SendTempId int, UserId int, TargetTempId int, ReduceHp float64) { //最后的底层方法
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
func (c *Ctx) ProtoColHealCardBt(UserId int, TargetTempId int, HealHp float64) {
	var card CardAbstract.Character
	var ok bool
	if card, ok = c.FindCard(TargetTempId).(CardAbstract.Character); !ok {
		fmt.Println("转化失败bug")
		return
	}
	NowHp := card.GetHpNow()
	MaxHp := card.GetInfo()["hp"].(float64)
	if NowHp >= MaxHp {
		return
	}

	if NowHp+HealHp > MaxHp {
		card.SetHpNow(MaxHp)
		return
	}
	card.SetHpNow(NowHp + HealHp)
}

func (c *Ctx) ProtoColGetCharacterCard(UserId int) []int {
	return c.GetCharacterCardinCardInHand(UserId)
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

//endregion

//region 对卡牌数据的操作算法

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

//endregion
