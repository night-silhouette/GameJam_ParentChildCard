package battleservice

import (
	"pcc_card/application/entity/BattleData"
	"pcc_card/application/entity/Card/CardAbstract"
)

type Ctx struct {
	PlayerDataMap map[int]*PlayerData
}

type PlayerData struct {
	ID           int
	CardInHand   *[]CardAbstract.Card
	ParentCardBT *CardAbstract.Card
	ChildCardBT  *CardAbstract.Card
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
		res.Self = append(res.Self, card.GetCardDto())
	}
	for i := 0; i < len(*list_opponent); i++ {
		card := (*list_opponent)[i]
		res.Opponent = append(res.Opponent, card.GetCardDto())
	}
	return &res
}
