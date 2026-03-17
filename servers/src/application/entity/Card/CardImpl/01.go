package CardImpl

import "pcc_card/application/entity/Card/skill_card"

type Card01 struct {
	skill_card.SkillCardTemplate
}

func NewCard01() *Card01 {
	return &Card01{}
}

func (c *Card01) GetID() int {
	return 1
}
