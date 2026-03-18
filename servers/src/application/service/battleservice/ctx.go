package battleservice

import (
	"pcc_card/application/entity/Card/CardAbstract"
)

type Ctx struct {
	DataA *PlayerData
	DataB *PlayerData
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
	c.DataA = NewPlayerData(idA)
	c.DataB = NewPlayerData(idB)
	return c
}
