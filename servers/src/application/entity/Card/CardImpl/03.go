package CardImpl

type Card03 struct {
	BaseCard
}

func NewCard03() *Card03 {
	return &Card03{}
}

func (c *Card03) GetID() int {
	return 3
}
