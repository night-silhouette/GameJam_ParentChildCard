package CardImpl

import "pcc_card/application/entity/Card/skill_card"

type Card05 struct {
	skill_card.SkillCardTemplate
}

func NewCard05() *Card05 {
	return &Card05{}
}

func (c *Card05) GetID() int {
	return 5
}
