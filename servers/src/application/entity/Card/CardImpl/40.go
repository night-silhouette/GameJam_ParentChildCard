package CardImpl

type Card40 struct {
	BaseCard
}

func NewCard40() *Card40 {
	return &Card40{}
}

func (c *Card40) Attack() {

}
func (c *Card40) Hurt() {
}

func (c *Card40) GetID() int {
	return 40
}
