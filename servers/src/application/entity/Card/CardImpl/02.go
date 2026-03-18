package CardImpl

type Card02 struct {
	BaseCard
}

func NewCard02() *Card02 {
	return &Card02{}
}

func (c *Card02) GetID() int {
	return 2
}
