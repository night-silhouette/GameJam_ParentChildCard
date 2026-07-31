package CardImpl

import (
	"pcc_card/application/entity/Card/CardAbstract"
)

type Card0008 struct {
	CharacterBaseCard
	IsSkill bool
}

func NewCard0008() *Card0008 {
	res := &Card0008{}
	res.CharacterBaseCard.Card = res
	res.IsSkill = false
	return res
}

func (c *Card0008) GetID() int {
	return 8
}

func (c *Card0008) Clone() CardAbstract.Card {
	newCard := *c
	return &newCard
}

func (c *Card0008) NextRound() {
	if c.IsSkill {
		c.IsSkill = false
		c.EffectUpdateEnergy(3)
	}
}

func (c *Card0008) Skill(TargetId int) bool {
	FinalId := c.CheckGuard(TargetId)
	if !c.ShareSkill(FinalId) {
		return false
	}
	c.IsSkill = true
	return true
}
