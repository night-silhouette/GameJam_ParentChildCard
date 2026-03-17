package CardImpl

import "pcc_card/application/entity/Card/skill_card"

type Card34 struct {
	skill_card.SkillCardTemplate
}

func NewCard34() *Card34 {
	return &Card34{}
}

func (c *Card34) GetID() int {
	return 34
}
