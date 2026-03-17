package skill_card

type Card00 struct {
	SkillCardTemplate
}

func NewCard00() *Card00 {
	return &Card00{}
}

func (c *Card00) GetID() int {
	return 0
}
