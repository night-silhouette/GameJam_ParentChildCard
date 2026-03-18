package CardImpl

type Card04 struct {
	BaseCard
}

func NewCard04() *Card04 {
	return &Card04{}
}

func (c *Card04) GetID() int {
	return 4
}
