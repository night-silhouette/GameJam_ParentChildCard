package CardImpl

import (
	"pcc_card/application/entity/Card/CardAbstract"
	"pcc_card/application/entity/protocol"
)

type Card0006 struct {
	CharacterBaseCard
}

func NewCard0006() *Card0006 {
	res := &Card0006{}
	res.CharacterBaseCard.Card = res
	return res
}

func (c *Card0006) GetID() int {
	return 6
}

func (c *Card0006) Clone() CardAbstract.Card {
	newCard := *c
	return &newCard
}

func (c *Card0006) Skill(TargetId int) bool {
	if !c.CharacterBaseCard.Skill(TargetId) {
		return false
	}

	c.GiveBuff(&TargetId, *protocol.NewBuffBase(protocol.HealingDecay, 4, 0.5, c.BtCtx.CreateTempId()))

	return true
}
