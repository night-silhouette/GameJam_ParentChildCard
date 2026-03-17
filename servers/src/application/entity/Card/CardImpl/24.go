package CardImpl

import "pcc_card/application/entity/Card/skill_card"

type Card24 struct {
	skill_card.SkillCardTemplate
}

func NewCard24() *Card24 {
	return &Card24{}
}

func (c *Card24) GetID() int {
	return 24
}
