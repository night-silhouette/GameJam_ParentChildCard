package CardImpl

import (
	"pcc_card/application/entity/BattleData"
	"pcc_card/application/entity/Card/CardAbstract"
)

type Card0001 struct {
	CharacterBaseCard
}

func NewCard0001() *Card0001 {
	res := &Card0001{}
	res.CharacterBaseCard.Card = res
	return res
}

func (c *Card0001) GetID() int {
	return 1
}

func (c *Card0001) Clone() CardAbstract.Card {
	newCard := *c
	return &newCard
}

func (c *Card0001) Skill(TargetId int) bool {
	if !c.CharacterBaseCard.Skill(TargetId) {
		return false
	}
	c.EffectHeal(c.GetTempId(), 1)
	c.EffectAttack(TargetId, 1, BattleData.TrueDamage)
	return true
}
