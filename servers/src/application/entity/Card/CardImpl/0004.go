package CardImpl

import (
	"pcc_card/application/entity/BattleData"
	"pcc_card/application/entity/Card/CardAbstract"
)

type Card0004 struct {
	CharacterBaseCard
	IsSkill bool
}

func NewCard0004() *Card0004 {
	res := &Card0004{}
	res.CharacterBaseCard.Card = res
	res.IsSkill = false
	return res
}

func (c *Card0004) GetID() int {
	return 4
}

func (c *Card0004) Clone() CardAbstract.Card {
	newCard := *c
	return &newCard
}
func (c *Card0004) Skill(TargetId int) bool {
	FinalId := c.CheckGuard(TargetId)
	if !c.ShareSkill(FinalId) {
		return false
	}
	c.IsSkill = true
	return true
}

func (c *Card0004) RoundEnd() {
	c.BaseCard.RoundEnd()
	c.IsSkill = false
}

func (c *Card0004) Hurt(AttackId int, HurtHp float64, category BattleData.ValueChange) {
	c.Retaliate(AttackId, HurtHp)
	c.Notify(BattleData.AnHurt, -1, AttackId, c.GetTempId())

	if !c.IsSkill {
		c.EffectHurt(AttackId, HurtHp, category)
	} else {
		c.EffectHeal(c.GetTempId(), HurtHp)
	}

}

func (c *Card0004) NoSourceHurt(HurtHp float64, category BattleData.ValueChange) {
	if !c.IsSkill {
		c.EffectHurt(-1, HurtHp, category)
	} else {
		c.EffectHeal(c.GetTempId(), HurtHp)
	}
}
