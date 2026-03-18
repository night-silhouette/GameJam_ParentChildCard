package CardImpl

type Card25 struct {
	BaseCard
}

func NewCard25() *Card25 {
	return &Card25{}
}

func (c *Card25) GetID() int {
	return 25
}
