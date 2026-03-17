package CardImpl

import "pcc_card/application/entity/Card/skill_card"

type Card33 struct {
	skill_card.SkillCardTemplate
}

func NewCard33() *Card33 {
	return &Card33{}
}

func (c *Card33) GetID() int {
	return 33
}
