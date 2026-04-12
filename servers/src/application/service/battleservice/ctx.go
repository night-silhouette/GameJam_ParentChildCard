package battleservice

import (
	"pcc_card/application/entity/BattleData"
	"pcc_card/application/entity/Card/CardAbstract"
)

// MetaCardChange 用于在CtxStateNotify传输状态的元数据
type MetaCardChange struct {
	Old *CardAbstract.Card
	New *CardAbstract.Card
}

// CtxStateNotify 内嵌到ctx里面，监听数据变化
type CtxStateNotify struct {
	ParentCardChange  chan MetaCardChange
	ChildCardBTChange chan MetaCardChange
	SkillCardBTChange chan MetaCardChange
}

func NewCtxStateNotify() *CtxStateNotify {
	N := &CtxStateNotify{}
	N.ParentCardChange = make(chan MetaCardChange)

	return N

}

type Ctx struct {
	PlayerDataMap  map[int]*PlayerData
	CtxStateNotify *CtxStateNotify
}

func (c *Ctx) SetParentCardBT(id int, new *CardAbstract.Card) {
	pData := c.PlayerDataMap[id]
	old := pData.ParentCardBT
	pData.ParentCardBT = new
	c.CtxStateNotify.ParentCardChange <- MetaCardChange{old, new}
}
func (c *Ctx) SetChildCardBT(id int, new *CardAbstract.Card) {
	pData := c.PlayerDataMap[id]
	old := pData.ChildCardBT
	pData.ChildCardBT = new
	c.CtxStateNotify.ParentCardChange <- MetaCardChange{old, new}
}
func (c *Ctx) SetSkillCardBT(id int, new *CardAbstract.Card) {
	pData := c.PlayerDataMap[id]
	old := pData.SkillCardBT
	pData.SkillCardBT = new
	c.CtxStateNotify.ParentCardChange <- MetaCardChange{old, new}
}

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

func NewCtx(idA int, idB int) *Ctx {
	c := &Ctx{}
	c.PlayerDataMap = make(map[int]*PlayerData, 2)
	c.PlayerDataMap[idA] = NewPlayerData(idA)
	c.PlayerDataMap[idB] = NewPlayerData(idB)
	c.CtxStateNotify = NewCtxStateNotify()
	return c
}

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
