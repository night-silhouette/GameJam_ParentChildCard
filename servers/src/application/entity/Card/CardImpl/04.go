package CardImpl

import "pcc_card/application/entity/Card/skill_card"

type Card04 struct {
	skill_card.SkillCardTemplate
}

func NewCard04() *Card04 {
	return &Card04{}
}

func (c *Card04) GetID() int {
	return 4
}
