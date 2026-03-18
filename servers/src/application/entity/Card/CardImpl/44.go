package CardImpl

type Card44 struct {
	BaseCard
}

func NewCard44() *Card44 {
	return &Card44{}
}

func (c *Card44) GetID() int {
	return 44
}
