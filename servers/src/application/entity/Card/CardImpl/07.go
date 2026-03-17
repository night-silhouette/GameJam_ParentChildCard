package CardImpl

type Card07 struct {
	BaseCard
}

func NewCard07() *Card07 {
	return &Card07{}
}

func (c *Card07) Attack() {

}
func (c *Card07) Hurt() {
}

func (c *Card07) GetID() int {
	return 7
}
