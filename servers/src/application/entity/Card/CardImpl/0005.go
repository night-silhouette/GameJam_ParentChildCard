package CardImpl

import (
	"pcc_card/application/entity/Card/CardAbstract"
	"pcc_card/application/entity/CardMeta"
	"pcc_card/application/entity/protocol"
)

type Card0005 struct {
	CharacterBaseCard
}

func NewCard0005() *Card0005 {
	res := &Card0005{}
	res.CharacterBaseCard.Card = res
	return res
}

func (c *Card0005) GetID() int {
	return 5
}

func (c *Card0005) Clone() CardAbstract.Card {
	newCard := *c
	return &newCard
}

func (c *Card0005) Skill(TargetId int) bool {
	if !c.CharacterBaseCard.Skill(TargetId) {
		return false
	}
	return true
}

func (c *Card0005) BroadCallBack(v *CardMeta.BroadInfo) {
	if v.ControlSignal != CardMeta.Wound {
		return
	}
	TempId := c.GetTempId()
	c.GiveBuff(&TempId, *protocol.NewBuffBase(protocol.Binding, 1, 0, c.BtCtx.CreateTempId()))
}
