package CardImpl

import "pcc_card/application/entity/Card/skill_card"

type Card02 struct {
	skill_card.SkillCardTemplate
}

func NewCard02() *Card02 {
	return &Card02{}
}

func (c *Card02) GetID() int {
	return 2
}
