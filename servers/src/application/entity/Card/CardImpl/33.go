package CardImpl

type Card33 struct {
	BaseCard
}

func NewCard33() *Card33 {
	return &Card33{}
}

func (c *Card33) GetID() int {
	return 33
}
