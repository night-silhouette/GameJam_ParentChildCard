package CardImpl

import "pcc_card/application/entity/Card/skill_card"

type Card23 struct {
	skill_card.SkillCardTemplate
}

func NewCard23() *Card23 {
	return &Card23{}
}

func (c *Card23) GetID() int {
	return 23
}
