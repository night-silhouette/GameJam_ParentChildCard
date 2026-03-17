package CardImpl

import "pcc_card/application/entity/Card/skill_card"

type Card35 struct {
	skill_card.SkillCardTemplate
}

func NewCard35() *Card35 {
	return &Card35{}
}

func (c *Card35) GetID() int {
	return 35
}
