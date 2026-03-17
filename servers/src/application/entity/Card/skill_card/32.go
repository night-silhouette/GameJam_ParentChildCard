package skill_card

type Card32 struct {
	SkillCardTemplate
}

func NewCard32() *Card32 {
	return &Card32{}
}

func (c *Card32) GetID() int {
	return 32
}
