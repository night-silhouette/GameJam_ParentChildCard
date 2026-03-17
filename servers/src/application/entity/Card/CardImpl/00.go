package CardImpl

type Card00 struct {
	skill_card.SkillCardTemplate
}

func NewCard00() *Card00 {
	return &Card00{}
}

func (c *Card00) GetID() int {
	return 0
}
