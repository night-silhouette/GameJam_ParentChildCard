package skill_card

type Card35 struct {
	SkillCardTemplate
}

func NewCard35() *Card35 {
	return &Card35{}
}

func (c *Card35) GetID() int {
	return 35
}
