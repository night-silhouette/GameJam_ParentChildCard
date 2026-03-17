package skill_card

type Card05 struct {
	SkillCardTemplate
}

func NewCard05() *Card05 {
	return &Card05{}
}

func (c *Card05) GetID() int {
	return 5
}
