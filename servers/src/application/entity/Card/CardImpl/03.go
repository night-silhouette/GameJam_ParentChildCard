package CardImpl

import "pcc_card/application/entity/Card/skill_card"

type Card03 struct {
	skill_card.SkillCardTemplate
}

func NewCard03() *Card03 {
	return &Card03{}
}

func (c *Card03) GetID() int {
	return 3
}
