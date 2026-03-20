package battleservice

import (
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
