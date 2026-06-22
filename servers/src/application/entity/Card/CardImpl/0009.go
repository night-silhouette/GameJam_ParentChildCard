package CardImpl

import (
	"pcc_card/application/entity/Card/CardAbstract"
	"pcc_card/application/entity/protocol"
)

type Card0009 struct {
	CharacterBaseCard
}

func NewCard0009() *Card0009 {
	res := &Card0009{}
	res.CharacterBaseCard.Card = res
	return res
}

func (c *Card0009) GetID() int {
	return 9
}

func (c *Card0009) Clone() CardAbstract.Card {
	newCard := *c
	return &newCard
}
func (c *Card0009) Skill(TargetId int) bool {
	FinalId := c.CheckGuard(TargetId)
	if !c.ShareSkill(FinalId) {
		return false
	}
	TempId := c.GetTempId()
	c.GiveBuff(&TempId, *protocol.NewBuffBase(protocol.Retaliate, 2, 0.5, c.BtCtx.CreateTempId()))
	return true
}
