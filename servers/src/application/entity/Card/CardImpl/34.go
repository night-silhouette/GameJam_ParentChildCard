package CardImpl

type Card34 struct {
	BaseCard
}

func NewCard34() *Card34 {
	return &Card34{}
}

func (c *Card34) GetID() int {
	return 34
}
