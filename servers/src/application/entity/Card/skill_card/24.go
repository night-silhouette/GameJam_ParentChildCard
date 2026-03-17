package skill_card

type Card24 struct {
	SkillCardTemplate
}

func NewCard24() *Card24 {
	return &Card24{}
}

func (c *Card24) GetID() int {
	return 24
}
