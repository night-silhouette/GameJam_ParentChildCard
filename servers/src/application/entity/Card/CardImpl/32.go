package CardImpl

type Card32 struct {
	BaseCard
}

func NewCard32() *Card32 {
	return &Card32{}
}

func (c *Card32) GetID() int {
	return 32
}
