package CardImpl

type Card01 struct {
	BaseCard
}

func NewCard01() *Card01 {
	return &Card01{}
}

func (c *Card01) GetID() int {
	return 1
}
