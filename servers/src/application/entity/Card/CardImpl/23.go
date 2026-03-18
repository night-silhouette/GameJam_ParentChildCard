package CardImpl

type Card23 struct {
	BaseCard
}

func NewCard23() *Card23 {
	return &Card23{}
}

func (c *Card23) GetID() int {
	return 23
}
