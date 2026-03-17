package CardImpl

import "pcc_card/application/entity/Card/skill_card"

type Card32 struct {
	skill_card.SkillCardTemplate
}

func NewCard32() *Card32 {
	return &Card32{}
}

func (c *Card32) GetID() int {
	return 32
}
