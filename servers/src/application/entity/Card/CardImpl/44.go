package CardImpl

import "pcc_card/application/entity/Card/skill_card"

type Card44 struct {
	skill_card.SkillCardTemplate
}

func NewCard44() *Card44 {
	return &Card44{}
}

func (c *Card44) GetID() int {
	return 44
}
