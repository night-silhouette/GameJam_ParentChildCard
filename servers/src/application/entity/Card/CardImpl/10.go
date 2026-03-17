package CardImpl

type Card10 struct {
	BaseCard
}

func NewCard10() *Card10 {
	return &Card10{}
}

func (c *Card10) Attack() {

}
func (c *Card10) Hurt() {
}

func (c *Card10) GetID() int {
	return 10
}
