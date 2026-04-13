package battleservice

import (
	"context"
	"pcc_card/application/entity/BattleData"
	"pcc_card/application/entity/Card/CardAbstract"
)

type Ctx struct {
	ParentContext context.Context

	CardObserver   *CardObserver
	PlayerDataMap  map[int]*PlayerData
	CtxStateNotify *CtxStateNotify
	CardPool       *[]CardAbstract.Card
}

func NewCtx(idA int, idB int, CardPool *[]CardAbstract.Card, ParentContext context.Context) *Ctx {
	c := &Ctx{}
	c.ParentContext = ParentContext
	c.CardPool = CardPool
	c.PlayerDataMap = make(map[int]*PlayerData, 2)
	c.PlayerDataMap[idA] = NewPlayerData(idA)
	c.PlayerDataMap[idB] = NewPlayerData(idB)
	c.CtxStateNotify = NewCtxStateNotify()
	c.CardObserver = NewCardObserver(ParentContext, c)
	return c
}

//__________________________________________CardObserver______________________________________________

type CardObserver struct {
	ctx           *Ctx
	ParentContext context.Context
	collector     chan MetaCardStateCode
}

type MetaCardStateCode struct {
	CardId    int
	StateCode CardAbstract.StateCode
}

func NewCardObserver(ParentContext context.Context, ctx *Ctx) *CardObserver {
	o := &CardObserver{}
	o.ParentContext = ParentContext
	o.collector = make(chan MetaCardStateCode, 16)
	o.ctx = ctx
	CardPool := *o.ctx.CardPool

	for _, card := range CardPool {

		go func(card CardAbstract.Card) { //给每一张卡开一个哨兵
			var CardChan <-chan CardAbstract.StateCode
			CardChan = card.GetStateCodeChan()
			for {
				select {
				case code := <-CardChan:
					Meta := MetaCardStateCode{card.GetID(), code}
					o.collector <- Meta
				case <-o.ParentContext.Done():
					return
				}
			}
		}(card)

	}

	go o.CardResponse(ParentContext)
	return o
}
func (o *CardObserver) CardResponse(ParentContext context.Context) {
	for {
		select {
		case Meta := <-o.collector:
			switch Meta.StateCode {
			case CardAbstract.Died:
				
			}
		}
	}
}

//__________________________________________CtxStateNotify______________________________________________

// MetaCardBTChange 用于在CtxStateNotify传输状态的元数据
type MetaCardBTChange struct {
	Old *CardAbstract.Card
	New *CardAbstract.Card
}

// CtxStateNotify 内嵌到ctx里面，监听数据变化
type CtxStateNotify struct {
	ParentCardChange  chan MetaCardBTChange
	ChildCardBTChange chan MetaCardBTChange
	SkillCardBTChange chan MetaCardBTChange
}

func NewCtxStateNotify() *CtxStateNotify {
	N := &CtxStateNotify{}
	N.ParentCardChange = make(chan MetaCardBTChange)
	return N

}

func (c *Ctx) SetParentCardBT(id int, new *CardAbstract.Card) {
	pData := c.PlayerDataMap[id]
	old := pData.ParentCardBT
	pData.ParentCardBT = new
	c.CtxStateNotify.ParentCardChange <- MetaCardBTChange{old, new}
}
func (c *Ctx) SetChildCardBT(id int, new *CardAbstract.Card) {
	pData := c.PlayerDataMap[id]
	old := pData.ChildCardBT
	pData.ChildCardBT = new
	c.CtxStateNotify.ParentCardChange <- MetaCardBTChange{old, new}
}
func (c *Ctx) SetSkillCardBT(id int, new *CardAbstract.Card) {
	pData := c.PlayerDataMap[id]
	old := pData.SkillCardBT
	pData.SkillCardBT = new
	c.CtxStateNotify.ParentCardChange <- MetaCardBTChange{old, new}
}

//__________________________________________PlayerData______________________________________________

type PlayerData struct {
	ID           int
	CardInHand   *[]CardAbstract.Card
	ParentCardBT *CardAbstract.Card
	ChildCardBT  *CardAbstract.Card
	SkillCardBT  *CardAbstract.Card
}

func NewPlayerData(ID int) *PlayerData {
	p := &PlayerData{}
	p.CardInHand = &[]CardAbstract.Card{}
	p.ID = ID
	return p
}

//__________________________________________杂项______________________________________________

func (c *Ctx) GetCardInHard(id_self int) *BattleData.CardInHand {
	var id_opponent int
	for key := range c.PlayerDataMap {
		if key != id_self {
			id_opponent = key
		}
	}
	res := BattleData.CardInHand{}
	res.Self = make([]BattleData.CardDto, 0, 20)
	res.Opponent = make([]BattleData.CardDto, 0, 20)
	list_self := c.PlayerDataMap[id_self].CardInHand
	list_opponent := c.PlayerDataMap[id_opponent].CardInHand
	for i := 0; i < len(*list_self); i++ {
		card := (*list_self)[i]
		res.Self = append(res.Self, CardAbstract.GetCardDto(card))
	}
	for i := 0; i < len(*list_opponent); i++ {
		card := (*list_opponent)[i]
		res.Opponent = append(res.Opponent, CardAbstract.GetCardDto(card))
	}
	return &res
}
