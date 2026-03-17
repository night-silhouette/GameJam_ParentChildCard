package skill_card

type Card02 struct {
	SkillCardTemplate
}

func NewCard02() *Card02 {
	return &Card02{}
}

func (c *Card02) GetID() int {
	return 2
}
