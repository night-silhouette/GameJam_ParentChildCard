package skill_card

type Card44 struct {
	SkillCardTemplate
}

func NewCard44() *Card44 {
	return &Card44{}
}

func (c *Card44) GetID() int {
	return 44
}
