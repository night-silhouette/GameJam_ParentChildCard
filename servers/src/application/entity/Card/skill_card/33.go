package skill_card

type Card33 struct {
	SkillCardTemplate
}

func NewCard33() *Card33 {
	return &Card33{}
}

func (c *Card33) GetID() int {
	return 33
}
