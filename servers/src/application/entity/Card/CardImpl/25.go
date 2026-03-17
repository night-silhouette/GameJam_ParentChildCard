package CardImpl

import "pcc_card/application/entity/Card/skill_card"

type Card25 struct {
	skill_card.SkillCardTemplate
}

func NewCard25() *Card25 {
	return &Card25{}
}

func (c *Card25) GetID() int {
	return 25
}
