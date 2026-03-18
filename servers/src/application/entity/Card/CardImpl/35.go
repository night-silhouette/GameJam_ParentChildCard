package CardImpl

type Card35 struct {
	BaseCard
}

func NewCard35() *Card35 {
	return &Card35{}
}

func (c *Card35) GetID() int {
	return 35
}
